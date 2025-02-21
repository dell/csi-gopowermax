/*
 Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package inttest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	pmax "github.com/dell/gopowermax/v2"
	types "github.com/dell/gopowermax/v2/types/v100"
)

const (
	SleepTime = 10 * time.Second
)

var (
	client pmax.Pmax
	// Following values will be read from env variables set in user.env
	endpoint = "https://0.0.0.0:8443"
	// username should match an existing user in Unisphere
	username = "username"
	// password should be the value for the corresponding user in Unisphere
	password                  = "password"
	apiVersion                = "XX"
	symmetrixID               = "000000000001"
	remoteSymmetrixID         = "000000000001"
	localRDFGrpNo             = "0"
	remoteRDFGrpNo            = "0"
	defaultRepMode            = "ASYNC"
	defaultFCPortGroup        = "fc-pg"
	defaultiSCSIPortGroup     = "iscsi-pg"
	defaultNVMePortGroup      = "nvme-pg"
	defaultFCInitiator        = "FA-1D:0:00000000abcd000e"
	defaultFCInitiatorID      string
	fcInitiator1              = "00000000abcd000f"
	fcInitiator2              = "00000000abcd000g"
	fcInitiator3              = "10000000c9959b8e"
	defaultiSCSIInitiator     = "SE-1E:000:iqn.1993-08.org.debian:01:012a34b5cd6"
	defaultiSCSIInitiatorID   string
	iscsiInitiator1           = "iqn.1993-08.org.centos:01:012a34b5cd7"
	iscsiInitiator2           = "iqn.1993-08.org.centos:01:012a34b5cd8"
	defaultNVMeTCPInitiator   = "OR-1C:1:nqn.2014-08.org.nvmexpress:uuid:csi_k8_nvme:76B04D56EAB26A2E1509A7E98D3DFDB6"
	defaultNVMeTCPInitiatorID string
	nvmetcpInitiator1         = "nqn.2014-08.org.nvmexpress:uuid:csi_k8_nvme:76B04D56EAB26A2E1509A7E98D3DFDB6"
	nvmetcpInitiator2         = "nqn.2014-08.org.nvmexpress:uuid:csi_k8_nvme:76B04D56EAB26A2E1509A7E98D3DFDB6"
	defaultSRP                = "storage-pool"
	defaultServiceLevel       = "Diamond"
	volumePrefix              = "xx"
	sgPrefix                  = "zz"
	snapshotPrefix            = "snap"
	csiPrefix                 = "csi"
	// the test run will create these for the run and clean up in the end
	defaultStorageGroup          = "csi-Integration-Test"
	defaultSnapshotName          = "integration-test-snapshot"
	defaultProtectedStorageGroup = "csi-Integration-Test-Protected-SG"
	nonFASTManagedSG             = "csi-Integration-No-FAST"
	defaultSGWithSnapshotPolicy  = "csi-Integration-Test-With-Snapshot-Policy"
	defaultSnapshotPolicy        = "DailyDefault"
	defaultFCHost                = "IntegrationFCHost"
	defaultiSCSIHost             = "IntegrationiSCSIHost"
	defaultNVMeTCPHost           = "IntegrationNVMeTCPHost"
	defaultFCdirname             = "FA-1D"
	defaultFCportName            = "4"
	defaultiscsidirName          = "SE-1E"
	defaultiscsiportName         = "0"
	defaultFCDirectorID          = "OR-1C"
	defaultFCPortID              = "0"
	localVol, remoteVol          *types.Volume
)

func setDefaultVariables() {
	endpoint = setenvVariable("Endpoint", endpoint)
	username = setenvVariable("Username", username)
	password = setenvVariable("Password", password)
	apiVersion = strings.TrimSpace(setenvVariable("APIVersion", ""))
	symmetrixID = setenvVariable("SymmetrixID", symmetrixID)
	remoteSymmetrixID = setenvVariable("RemoteSymmetrixID", remoteSymmetrixID)
	defaultRepMode = setenvVariable("DefaultRepMode", defaultRepMode)
	localRDFGrpNo = setenvVariable("LocalRDFGrpNo", localRDFGrpNo)
	remoteRDFGrpNo = setenvVariable("RemoteRDFGrpInfo", remoteRDFGrpNo)
	defaultFCPortGroup = setenvVariable("DefaultFCPortGroup", defaultFCPortGroup)
	defaultiSCSIPortGroup = setenvVariable("DefaultiSCSIPortGroup", defaultiSCSIPortGroup)
	defaultFCInitiator = setenvVariable("DefaultFCInitiator", defaultFCInitiator)
	defaultFCInitiatorID = strings.Split(defaultFCInitiator, ":")[2]
	fcInitiator1 = setenvVariable("FCInitiator1", fcInitiator1)
	fcInitiator2 = setenvVariable("FCInitiator2", fcInitiator2)
	defaultiSCSIInitiator = setenvVariable("DefaultiSCSIInitiator", defaultiSCSIInitiator)
	defaultiSCSIInitiatorID = strings.Join(strings.Split(defaultiSCSIInitiator, ":")[2:], ":")
	iscsiInitiator1 = setenvVariable("ISCSIInitiator1", iscsiInitiator1)
	iscsiInitiator2 = setenvVariable("ISCSIInitiator2", iscsiInitiator2)
	defaultNVMeTCPInitiator = setenvVariable("DefaultNVMeInitiator", defaultNVMeTCPInitiator)
	defaultNVMeTCPInitiatorID = strings.Join(strings.Split(defaultNVMeTCPInitiator, ":")[2:], ":")
	defaultNVMePortGroup = setenvVariable("DefaultNVMeTCPPortGroup", defaultNVMePortGroup)
	nvmetcpInitiator1 = setenvVariable("NVMeInitiator1", nvmetcpInitiator1)
	nvmetcpInitiator2 = setenvVariable("NVMeInitiator2", nvmetcpInitiator2)
	volumePrefix = setenvVariable("VolumePrefix", volumePrefix)
	defaultSRP = setenvVariable("DefaultStoragePool", defaultSRP)
	defaultServiceLevel = setenvVariable("DefaultServiceLevel", defaultServiceLevel)
	sgPrefix = setenvVariable("SGPrefix", sgPrefix)
	snapshotPrefix = setenvVariable("SnapPrefix", snapshotPrefix)
	defaultFCdirname = setenvVariable("DefaultFCDirName", defaultFCdirname)
	defaultFCportName = setenvVariable("DefaultFCPortName", defaultFCportName)
	defaultiscsidirName = setenvVariable("DefaultISCSIDirName", defaultiscsidirName)
	defaultiscsiportName = setenvVariable("DefaultISCSIPortName", defaultiscsiportName)
	defaultFCDirectorID = setenvVariable("DefaultFCDirectorID", defaultFCDirectorID)
	defaultFCPortID = setenvVariable("DefaultFCPortID", defaultFCPortID)
	defaultProtectedStorageGroup = defaultProtectedStorageGroup + "-" + localRDFGrpNo + "-" + defaultRepMode
}

func TestMain(m *testing.M) {
	status := 0
	// Process environment variables
	setDefaultVariables()

	err := createDefaultSGAndHost() // Creates default storage group and host for the test
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = createRDFSetup() // Creates RDF setup for the test
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if st := m.Run(); st > status {
		status = st
	}
	fmt.Printf("status %d\n", status)
	doCleanUp := setenvVariable("Cleanup", "true")
	cleanupTests := []testing.InternalTest{}
	if doCleanUp != "false" {
		fmt.Println("========= CLEANUP ==========")
		cleanupTests = append(cleanupTests, testing.InternalTest{
			Name: "cleanupDefaultSGAndHOST",
			F:    cleanupDefaultSGAndHOST,
		})
	}
	// Always clean up the resources used in replication
	cleanupTests = append(cleanupTests, testing.InternalTest{
		Name: "cleanupRDFSetup",
		F:    cleanupRDFSetup,
	})
	afterRun(cleanupTests) // Cleans up the volumes and snapshots created for replication testing purposes.
}

func setenvVariable(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if key != "Username" && key != "Password" {
			fmt.Printf("%s=%s\n", key, defaultValue)
		}
		return defaultValue
	}
	if key != "Username" && key != "Password" {
		fmt.Printf("%s=%s\n", key, value)
	}
	return value
}

func createDefaultSGAndHost() error {
	fmt.Println("Creating default SG and host...")
	// Create default SG with SRP
	_, err := createStorageGroup(symmetrixID, defaultStorageGroup, defaultSRP, defaultServiceLevel, false, nil)
	if err != nil {
		return fmt.Errorf("failed to create SG: (%s)", err.Error())
	}

	// Create default SG without srp
	_, err = createStorageGroup(symmetrixID, nonFASTManagedSG, "none", "none", false, nil)
	if err != nil {
		return fmt.Errorf("failed to create non fast SG: (%s)", err.Error())
	}

	// Create default SG with snapshot policy
	optionalPayload := make(map[string]interface{})
	optionalPayload["snapshotPolicies"] = []string{defaultSnapshotPolicy}
	_, err = client.CreateStorageGroup(context.TODO(), symmetrixID, defaultSGWithSnapshotPolicy, defaultSRP, defaultServiceLevel, false, optionalPayload)
	if err != nil {
		return fmt.Errorf("failed to create SG with snapshot policy: (%s)", err.Error())
	}

	// Create default FC Host
	initiators := []string{defaultFCInitiatorID}
	_, err = createHost(symmetrixID, defaultFCHost, initiators, nil)
	if err != nil {
		return fmt.Errorf("failed to create default FC host: (%s)", err.Error())
	}
	// Create default iSCSI Host
	initiators = []string{defaultiSCSIInitiatorID}
	_, err = createHost(symmetrixID, defaultiSCSIHost, initiators, nil)
	if err != nil {
		return fmt.Errorf("failed to create default iSCSI host: (%s)", err.Error())
	}

	// Create default NVMeTCP Host
	initiators = []string{defaultNVMeTCPInitiatorID}
	_, err = createHost(symmetrixID, defaultNVMeTCPHost, initiators, nil)
	if err != nil {
		return fmt.Errorf("failed to create default NVMeTCP host: (%s)", err.Error())
	}
	return err
}

func createRDFSetup() error {
	fmt.Printf("Creating RDF Setup.....")

	// Creating default Protected SG

	_, err := createStorageGroup(symmetrixID, defaultProtectedStorageGroup, defaultSRP, defaultServiceLevel, false, nil)
	if err != nil {
		return fmt.Errorf("failed to create SG: (%s)", err.Error())
	}
	now := time.Now()

	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())

	// Creating source volume
	volOpts := make(map[string]interface{})
	localVol, err = client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultProtectedStorageGroup, volumeName, 50, volOpts)
	if err != nil {
		return fmt.Errorf("failed to create volume : (%s)", err.Error())
	}
	fmt.Printf("volume:\n%#v\n", localVol)

	// Creating SG Replica

	SGRDFInfo, err := client.CreateSGReplica(context.TODO(), symmetrixID, remoteSymmetrixID, defaultRepMode, localRDFGrpNo, defaultProtectedStorageGroup, defaultProtectedStorageGroup, defaultServiceLevel, false)
	fmt.Printf("SG info :\n%#v\n", SGRDFInfo)
	if err != nil {
		return fmt.Errorf("Error Creating SGReplica: %s", err.Error())
	}

	// Retrieving remote volume information for cleanup

	rdfPair, err := client.GetRDFDevicePairInfo(context.TODO(), symmetrixID, localRDFGrpNo, localVol.VolumeID)
	if err != nil {
		return fmt.Errorf("Error retrieving RDF device pair information: %s", err.Error())
	}
	tgtVolID := rdfPair.RemoteVolumeName
	remoteVol, err = client.GetVolumeByID(context.TODO(), remoteSymmetrixID, tgtVolID)
	if err != nil {
		return fmt.Errorf("Error retrieving volume information: %s", err.Error())
	}
	return err
}

// NVME
func cleanupDefaultSGAndHOST(t *testing.T) {
	fmt.Println("Cleaning up SG and host...")
	// delete default SG
	err := deleteStorageGroup(symmetrixID, defaultStorageGroup)
	if err != nil {
		t.Errorf("failed to delete default SG (%s) : (%s)", defaultStorageGroup, err.Error())
	}
	// delete default non fast SG
	err = deleteStorageGroup(symmetrixID, nonFASTManagedSG)
	if err != nil {
		t.Errorf("failed to delete default non fast SG (%s) : (%s)", nonFASTManagedSG, err.Error())
	}

	// delete default FC host
	err = deleteHost(symmetrixID, defaultFCHost)
	if err != nil {
		t.Errorf("failed to delete default FC host (%s) : (%s)", defaultFCHost, err.Error())
	}
	// delete default iSCSI host
	err = deleteHost(symmetrixID, defaultiSCSIHost)
	if err != nil {
		t.Errorf("failed to delete default iSCSI host (%s) : (%s)", defaultFCHost, err.Error())
	}
}

func cleanupRDFSetup(t *testing.T) {
	fmt.Println("Cleaning up RDF Setup...")

	// Terminating the Pair and removing the volumes from local SG and remote SG

	_, err := client.RemoveVolumesFromProtectedStorageGroup(context.TODO(), symmetrixID, defaultProtectedStorageGroup, remoteSymmetrixID, defaultProtectedStorageGroup, true, localVol.VolumeID)
	if err != nil {
		t.Errorf("failed to remove volumes from default Protected SG (%s) : (%s)", defaultProtectedStorageGroup, err.Error())
	}
	// Deleting local volume
	err = client.DeleteVolume(context.TODO(), symmetrixID, localVol.VolumeID)
	if err != nil {
		t.Error("DeleteVolume failed: " + err.Error())
	}
	// Test deletion of the volume again... should return an error
	err = client.DeleteVolume(context.TODO(), symmetrixID, localVol.VolumeID)
	if err == nil {
		t.Error("Expected an error saying volume was not found, but no error")
	}
	fmt.Printf("Received expected error: %s\n", err.Error())

	// Deleting remote volume
	err = client.DeleteVolume(context.TODO(), remoteSymmetrixID, remoteVol.VolumeID)
	if err != nil {
		t.Error("DeleteVolume failed: " + err.Error())
	}
	// Test deletion of the volume again... should return an error
	err = client.DeleteVolume(context.TODO(), remoteSymmetrixID, remoteVol.VolumeID)
	if err == nil {
		t.Error("Expected an error saying volume was not found, but no error")
	}
	fmt.Printf("Received expected error: %s\n", err.Error())

	// Deleting local SG
	err = deleteStorageGroup(symmetrixID, defaultProtectedStorageGroup)
	if err != nil {
		t.Errorf("failed to delete default SG (%s) : (%s)", defaultProtectedStorageGroup, err.Error())
	}

	// Deleting remote SG
	err = deleteStorageGroup(remoteSymmetrixID, defaultProtectedStorageGroup)
	if err != nil {
		t.Errorf("failed to delete default SG (%s) : (%s)", defaultProtectedStorageGroup, err.Error())
	}
}

func getClient() error {
	var err error
	client, err = pmax.NewClientWithArgs(endpoint, "CSI Driver for Dell EMC PowerMax v1.0",
		true, false, "")
	if err != nil {
		return err
	}
	err = client.Authenticate(context.TODO(), &pmax.ConfigConnect{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func TestAuthentication(_ *testing.T) {
	_ = getClient()
}

func TestGetSymmetrixIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	symIDList, err := client.GetSymmetrixIDList(context.TODO())
	if err != nil || symIDList == nil {
		t.Error("cannot get SymmetrixIDList: ", err.Error())
		return
	}
	if len(symIDList.SymmetrixIDs) == 0 {
		t.Error("expected at least one Symmetrix ID in list")
		return
	}
	for _, id := range symIDList.SymmetrixIDs {
		fmt.Printf("symmetrix ID: %s\n", id)
	}
}

func TestGetSymmetrix(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	symmetrix, err := client.GetSymmetrixByID(context.TODO(), symmetrixID)
	if err != nil || symmetrix == nil {
		t.Error("cannot get Symmetrix id "+symmetrixID, err.Error())
		return
	}
	fmt.Printf("Symmetrix %s: %#v\n", symmetrixID, symmetrix)
}

func TestGetVolumeIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	volumeIDList, err := client.GetVolumeIDList(context.TODO(), symmetrixID, "", false)
	if err != nil || volumeIDList == nil {
		t.Error("cannot get volumeIDList: ", err.Error())
		return
	}
	fmt.Printf("%d volume IDs\n", len(volumeIDList))
	// Make sure no duplicates
	dupMap := make(map[string]bool)
	for i := 0; i < len(volumeIDList); i++ {
		if volumeIDList[i] == "" {
			t.Error("Got an empty volume ID")
		}
		id := volumeIDList[i]
		if dupMap[id] == true {
			t.Error("Got duplicate ID:" + id)
		}
		dupMap[id] = true
	}

	volumeIDList, err = client.GetVolumeIDList(context.TODO(), symmetrixID, "csi", true)
	if err != nil || volumeIDList == nil {
		t.Error("cannot get volumeIDList: ", err.Error())
		return
	}
	fmt.Printf("%d CSI volume IDs\n", len(volumeIDList))
	for _, id := range volumeIDList {
		fmt.Printf("CSI volume: %s\n", id)
	}

	volumeIDList, err = client.GetVolumeIDList(context.TODO(), symmetrixID, "ce9072c0", true)
	if err != nil || volumeIDList == nil {
		t.Error("cannot get volumeIDList: ", err.Error())
		return
	}
	fmt.Printf("%d CSI volume IDs\n", len(volumeIDList))
}

func TestGetVolumeIDsWithParams(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	params := map[string]string{
		"tdev": "true",
	}
	volumeIDList, err := client.GetVolumeIDListWithParams(context.TODO(), symmetrixID, params)
	if err != nil || volumeIDList == nil {
		t.Error("cannot get volumeIDList with params", err.Error())
		return
	}
	fmt.Printf("%d Volume IDs with params\n", len(volumeIDList))
}

func TestGetVolume(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	// Get some CSI volumes
	volumeIDList, err := client.GetVolumeIDList(context.TODO(), symmetrixID, "csi", true)
	if err != nil || volumeIDList == nil {
		t.Error("cannot get CSI volumeIDList: ", err.Error())
		return
	}
	if len(volumeIDList) == 0 {
		t.Error("no CSI volumes")
		return
	}
	for i, id := range volumeIDList {
		if i >= 3 {
			break
		}
		volume, err := client.GetVolumeByID(context.TODO(), symmetrixID, id)
		if err != nil {
			t.Error("cannot retrieve Volume: " + err.Error())
		} else {
			fmt.Printf("Volume %#v\n", volume)
		}

	}
}

func TestGetNonExistentVolume(t *testing.T) {
	volume, err := client.GetVolumeByID(context.TODO(), symmetrixID, "88888")
	if err != nil {
		fmt.Printf("TestGetNonExistentVolume: %s\n", err.Error())
	} else {
		fmt.Printf("%#v\n", volume)
		t.Error("Expected volume 88888 to be non-existent")
	}
}

func TestGetStorageGroupIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	sgIDList, err := client.GetStorageGroupIDList(context.TODO(), symmetrixID, "", false)
	if err != nil || sgIDList == nil {
		t.Error("cannot get StorageGroupIDList: ", err.Error())
		return
	}
	if len(sgIDList.StorageGroupIDs) == 0 {
		t.Error("expected at least one StorageGroup ID in list")
		return
	}
	for _, id := range sgIDList.StorageGroupIDs {
		fmt.Printf("StorageGroup ID: %s\n", id)
	}
}

func TestGetStorageGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	storageGroup, err := client.GetStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup)
	if err != nil || storageGroup == nil {
		t.Error("Expected to find " + defaultStorageGroup + " but didn't")
		return
	}
	fmt.Printf("%#v\n", storageGroup)
}

func TestGetStorageGroupSnapshotPolicy(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	storageGroupSnapshotPolicy, err := client.GetStorageGroupSnapshotPolicy(context.TODO(), symmetrixID, defaultSnapshotPolicy, defaultSGWithSnapshotPolicy)
	if err != nil || storageGroupSnapshotPolicy == nil {
		t.Error("Expected to find " + defaultSGWithSnapshotPolicy + " but didn't")
		return
	}
	fmt.Printf("%#v\n", storageGroupSnapshotPolicy)

	// Cleanup -  remove snapshot policy from SG
	optionalPayload := make(map[string]interface{})
	optionalPayload["editStorageGroupActionParam"] = types.EditStorageGroupActionParam{
		EditSnapshotPoliciesParam: &types.EditSnapshotPoliciesParam{
			DisassociateSnapshotPolicyParam: &types.SnapshotPolicies{
				SnapshotPolicies: []string{defaultSnapshotPolicy},
			},
		},
	}
	_, err = client.UpdateStorageGroup(context.TODO(), symmetrixID, defaultSGWithSnapshotPolicy, optionalPayload)
	if err != nil {
		t.Errorf("failed to remove snapshot policy from SG (%s) : (%s)", defaultSGWithSnapshotPolicy, err.Error())
	} else {
		// Cleanup - delete SG after snapshot policy is removed
		err = deleteStorageGroup(symmetrixID, defaultSGWithSnapshotPolicy)
		if err != nil {
			t.Errorf("failed to delete default SG with snapshot policy (%s) : (%s)", defaultSGWithSnapshotPolicy, err.Error())
		}
	}
}

func TestGetStoragePool(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	storagePool, err := client.GetStoragePool(context.TODO(), symmetrixID, defaultSRP)
	if err != nil || storagePool == nil {
		t.Error("Expected to find " + defaultSRP + " but didn't")
		return
	}
	fmt.Printf("%#v\n", storagePool)
}

func createStorageGroup(symmetrixID, storageGroupID, srp, serviceLevel string, isThick bool, hostLimits *types.SetHostIOLimitsParam) (*types.StorageGroup, error) {
	if client == nil {
		err := getClient()
		if err != nil {
			return nil, err
		}
	}
	// Check if the SG exists on array
	storageGroup, err := client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	// Storage Group already exist, returning old one
	if storageGroup != nil && err == nil {
		return storageGroup, err
	}
	// Create a new storege group
	fmt.Println("Creating a new storage group...")
	optionalPayload := make(map[string]interface{})
	optionalPayload["hostLimits"] = hostLimits
	return client.CreateStorageGroup(context.TODO(), symmetrixID, storageGroupID, srp, serviceLevel, isThick, optionalPayload)
}

func deleteStorageGroup(symmetrixID, storageGroupID string) error {
	return client.DeleteStorageGroup(context.TODO(), symmetrixID, storageGroupID)
}

func TestCreateStorageGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	storageGroupID := fmt.Sprintf("csi-%s-Int%d-SG", sgPrefix, now.Nanosecond())
	storageGroup, err := createStorageGroup(symmetrixID, storageGroupID,
		defaultSRP, defaultServiceLevel, false, nil)
	if err != nil || storageGroup == nil {
		t.Error("Failed to create " + storageGroupID + " " + err.Error())
		return
	}
	fmt.Println("Fetching the newly create storage group from array")
	// Check if the SG exists on array
	storageGroup, err = client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err != nil || storageGroup == nil {
		t.Error("Expected to find " + storageGroupID + " but didn't")
		return
	}
	fmt.Printf("%#v\n", storageGroup)
	fmt.Println("Cleaning up the storage group: " + storageGroupID)
	err = deleteStorageGroup(symmetrixID, storageGroupID)
	if err != nil {
		t.Error("Failed to delete " + storageGroupID)
		return
	}
	// Check if the SG exists on array
	storageGroup, err = client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err == nil || storageGroup != nil {
		t.Error("Expected a failure in fetching " + storageGroupID + " but didn't")
		return
	}
	fmt.Println(fmt.Sprintf("Error received while fetching %s: %s", storageGroupID, err.Error()))
}

func TestCreateStorageGroupWithHostIOLimits(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	limits := types.SetHostIOLimitsParam{
		HostIOLimitMBSec:    "1",
		HostIOLimitIOSec:    "100",
		DynamicDistribution: "Never",
	}
	storageGroupID := fmt.Sprintf("csi-%s-Int%d-SG", sgPrefix, now.Nanosecond())
	storageGroup, err := createStorageGroup(symmetrixID, storageGroupID,
		defaultSRP, defaultServiceLevel, false, &limits)
	if err != nil || storageGroup == nil {
		t.Error("Failed to create " + storageGroupID + " " + err.Error())
		return
	}
	fmt.Println("Fetching the newly create storage group from array")
	// Check if the SG exists on array
	storageGroup, err = client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err != nil || storageGroup == nil {
		t.Error("Expected to find " + storageGroupID + " but didn't")
		return
	}
	fmt.Printf("%#v\n", storageGroup)
	fmt.Println("Cleaning up the storage group: " + storageGroupID)
	err = deleteStorageGroup(symmetrixID, storageGroupID)
	if err != nil {
		t.Error("Failed to delete " + storageGroupID)
		return
	}
	// Check if the SG exists on array
	storageGroup, err = client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err == nil || storageGroup != nil {
		t.Error("Expected a failure in fetching " + storageGroupID + " but didn't")
		return
	}
	fmt.Println(fmt.Sprintf("Error received while fetching %s: %s", storageGroupID, err.Error()))
}

func TestCreateStorageGroupNonFASTManaged(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	storageGroupID := fmt.Sprintf("csi-%s-Int%d-SG-No-FAST", sgPrefix, now.Nanosecond())
	storageGroup, err := createStorageGroup(symmetrixID, storageGroupID,
		"None", "None", false, nil)
	if err != nil || storageGroup == nil {
		t.Error("Failed to create " + storageGroupID)
		return
	}
	fmt.Printf("%#v\n", storageGroup)
	if storageGroup.SRP != "" {
		t.Error("Expected no SRP but received: " + storageGroup.SRP)
	}
	fmt.Println("Cleaning up the storage group: " + storageGroupID)
	err = client.DeleteStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err != nil {
		t.Error("Failed to delete " + storageGroupID)
		return
	}
	// Check if the SG exists on array
	storageGroup, err = client.GetStorageGroup(context.TODO(), symmetrixID, storageGroupID)
	if err == nil || storageGroup != nil {
		t.Error("Expected a failure in fetching " + storageGroupID + " but didn't")
		return
	}
	fmt.Println(fmt.Sprintf("Error received while fetching %s: %s", storageGroupID, err.Error()))
}

func TestGetJobs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	jobIDList, err := client.GetJobIDList(context.TODO(), symmetrixID, "")
	if err != nil {
		t.Error("failed to get Job ID LIst")
		return
	}
	for i, id := range jobIDList {
		if i >= 10 {
			break
		}
		job, err := client.GetJobByID(context.TODO(), symmetrixID, id)
		if err != nil {
			t.Error("failed to get job: " + id)
			return
		}
		fmt.Printf("%s\n", client.JobToString(job))
	}

	jobIDList, err = client.GetJobIDList(context.TODO(), symmetrixID, types.JobStatusRunning)
	if err != nil {
		t.Error("failed to get Job ID LIst")
		return
	}
	for i, id := range jobIDList {
		if i >= 10 {
			break
		}
		job, err := client.GetJobByID(context.TODO(), symmetrixID, id)
		if err != nil {
			t.Error("failed to get job: " + id)
			return
		}
		fmt.Printf("%s\n", client.JobToString(job))
		if job.Status != types.JobStatusRunning {
			t.Error("Expected Running job: " + client.JobToString(job))
		}
	}
}

func TestGetStoragePoolList(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	spList, err := client.GetStoragePoolList(context.TODO(), symmetrixID)
	if err != nil {
		t.Error("Failed to get StoragePoolList: " + err.Error())
		return
	}
	for _, value := range spList.StoragePoolIDs {
		fmt.Printf("Storage Resource Pool: %s\n", value)
	}
}

func TestGetMaskingViews(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	mvList, err := client.GetMaskingViewList(context.TODO(), symmetrixID)
	if err != nil {
		t.Error("Failed to get MaskingViewList: " + err.Error())
		return
	}
	for _, mvID := range mvList.MaskingViewIDs {
		fmt.Printf("Masking View: %s\n", mvID)
		mv, err := client.GetMaskingViewByID(context.TODO(), symmetrixID, mvID)
		if err != nil {
			t.Error("Failed to GetMaskingViewByID: ", err.Error())
			return
		}
		fmt.Printf("%#v\n", mv)
		conns, err := client.GetMaskingViewConnections(context.TODO(), symmetrixID, mvID, "")
		if err != nil {
			t.Error("Failed to GetMaskingViewConnections: ", err.Error())
			return
		}
		for _, conn := range conns {
			fmt.Printf("mv connection VolumeID %s HostLUNAddress %s InitiatorID %s DirectorPort %s\n",
				conn.VolumeID, conn.HostLUNAddress, conn.InitiatorID, conn.DirectorPort)
		}
	}
}

func TestCreateVolumeInStorageGroup1(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	payload := client.GetCreateVolInSGPayload(1, "CYL", volumeName, false, false, "", "", nil)

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		t.Error("Encoding error on json")
	}
	fmt.Printf("payload: %s\n", string(payloadBytes))

	job, err := client.UpdateStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, payload)
	if err != nil {
		t.Error("Error returned from UpdateStorageGroup")
		return
	}
	jobID := job.JobID
	job, err = client.WaitOnJobCompletion(context.TODO(), symmetrixID, jobID)
	if err == nil {
		idlist, err := client.GetVolumeIDList(context.TODO(), symmetrixID, volumeName, false)
		if err != nil {
			t.Error("Error getting volume IDs: " + err.Error())
		}
		for _, id := range idlist {
			cleanupVolume(id, volumeName, defaultStorageGroup, t)
		}
	}
}

func TestCreateVolumeInStorageGroup2(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroupS(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func TestModifyMobilityForVolume(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroupS(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	vol, err = client.ModifyMobilityForVolume(context.TODO(), symmetrixID, vol.VolumeID, true)
	if err != nil {
		t.Error(err)
		return
	}
	if !vol.MobilityIDEnabled {
		t.Errorf("Failed to modify mobilityID")
	}
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func TestCreateVolumeInStorageGroup2withUnit(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	capUnit := "TB"
	fmt.Printf("volumeName: %s\n", volumeName)
	volopts := make(map[string]interface{})
	volopts["capacityUnit"] = capUnit
	vol, err := client.CreateVolumeInStorageGroupS(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, "1", volopts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func TestAddVolumesInStorageGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	err = client.AddVolumesToStorageGroup(context.TODO(), symmetrixID, nonFASTManagedSG, true, vol.VolumeID)
	if err != nil {
		t.Error(err)
		return
	}
	sg, err := client.GetStorageGroup(context.TODO(), symmetrixID, nonFASTManagedSG)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("SG after adding volume: %#v\n", sg)
	// Remove the volume from SG as part of cleanup
	sg, err = client.RemoveVolumesFromStorageGroup(context.TODO(), symmetrixID, nonFASTManagedSG, true, vol.VolumeID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("SG after removing volume: %#v\n", sg)
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func cleanupVolume(volumeID string, volumeName string, storageGroup string, t *testing.T) {
	if volumeName != "" {
		vol, err := client.RenameVolume(context.TODO(), symmetrixID, volumeID, "_DEL"+volumeName)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("volume Renamed: %s\n", vol.VolumeIdentifier)
	}
	sg, err := client.RemoveVolumesFromStorageGroup(context.TODO(), symmetrixID, storageGroup, true, volumeID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("SG after removing volume: %#v\n", sg)
	pmax.Debug = true
	err = client.DeleteVolume(context.TODO(), symmetrixID, volumeID)
	if err != nil {
		t.Error("DeleteVolume failed: " + err.Error())
	}
	// Test deletion of the volume again... should return an error
	err = client.DeleteVolume(context.TODO(), symmetrixID, volumeID)
	if err == nil {
		t.Error("Expected an error saying volume was not found, but no error")
	}
	fmt.Printf("Received expected error: %s\n", err.Error())
}

func TestCreateVolumeInStorageGroupInParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping this test in short mode")
	}
	// make sure we have a client
	getClient()
	CreateVolumesInParallel(5, t)
}

func CreateVolumesInParallel(nVols int, t *testing.T) {
	fmt.Printf("testing CreateVolumeInStorageGroup with %d parallel requests\n", nVols)
	volIDList := make([]string, nVols)
	volOpts := make(map[string]interface{})
	// make channels for communication
	idchan := make(chan string, nVols)
	errchan := make(chan error, nVols)
	t0 := time.Now()

	// create a temporary storage group
	now := time.Now()
	storageGroupName := fmt.Sprintf("pmax-%s-Int%d-SG", sgPrefix, now.Nanosecond())
	_, err := createStorageGroup(symmetrixID, storageGroupName,
		defaultSRP, defaultServiceLevel, false, nil)
	if err != nil {
		t.Errorf("Unable to create temporary Storage Group: %s", storageGroupName)
	}
	// Send requests
	for i := 0; i < nVols; i++ {
		name := fmt.Sprintf("pmax-Int%d-Scale%d", now.Nanosecond(), i)
		go func(volumeName string, idchan chan string, errchan chan error) {
			var err error
			resp, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, storageGroupName, volumeName, 1, volOpts)
			if resp != nil {
				fmt.Printf("ID %s Name %s\n%#v\n", resp.VolumeID, volumeName, resp)
				idchan <- resp.VolumeID
			} else {
				idchan <- ""
			}
			errchan <- err
		}(name, idchan, errchan)
	}
	// Wait on complete, collecting ids and errors
	nerrors := 0
	for i := 0; i < nVols; i++ {
		var id string
		var err error
		id = <-idchan
		if id != "" {
			volIDList[i] = id
		}
		err = <-errchan
		if err != nil {
			err = fmt.Errorf("create volume received error: %s", err.Error())
			t.Error(err.Error())
			nerrors++
		}
	}
	t1 := time.Now()
	fmt.Printf("Create volume time for %d volumes %d errors: %v %v\n", nVols, nerrors, t1.Sub(t0).Seconds(), t1.Sub(t0).Seconds()/float64(nVols))
	fmt.Printf("%v\n", volIDList)
	time.Sleep(SleepTime)
	// Cleanup the volumes
	for _, id := range volIDList {
		cleanupVolume(id, "", storageGroupName, t)
	}
	// remove the temporary storage group
	err = client.DeleteStorageGroup(context.TODO(), symmetrixID, storageGroupName)
	if err != nil {
		t.Errorf("Unable to delete temporary Storage Group: %s", storageGroupName)
	}
}

func TestGetPorts(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	dirName := defaultFCdirname
	portName := defaultFCportName
	port, err := client.GetPort(context.TODO(), symmetrixID, dirName, portName)
	if err != nil {
		t.Errorf("Unable to read FC storage port %s %s: %s", dirName, portName, err)
		return
	}
	fmt.Printf("port %s:%s %#v\n", dirName, portName, port)
	dirName = defaultiscsidirName
	portName = defaultiscsiportName
	port, err = client.GetPort(context.TODO(), symmetrixID, dirName, portName)
	if err != nil {
		t.Errorf("Unable to read iSCSI storage port %s %s: %s", dirName, portName, err)
		return
	}
	fmt.Printf("port %s:%s %#v\n", dirName, portName, port)
}

func TestGetPortGroupIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	pgList, err := client.GetPortGroupList(context.TODO(), symmetrixID, "")
	if err != nil || pgList == nil {
		t.Error("cannot get PortGroupList: ", err.Error())
		return
	}
	if len(pgList.PortGroupIDs) == 0 {
		t.Error("expected at least one PortGroup ID in list")
		return
	}
	pgList, err = client.GetPortGroupList(context.TODO(), symmetrixID, "fibre")
	if err != nil || pgList == nil {
		t.Error("cannot get FC PortGroupList: ", err.Error())
		return
	}
	if len(pgList.PortGroupIDs) == 0 {
		t.Error("expected at least one FC PortGroup ID in list")
		return
	}
	pgList, err = client.GetPortGroupList(context.TODO(), symmetrixID, "iscsi")
	if err != nil || pgList == nil {
		t.Error("cannot get iSCSI PortGroupList: ", err.Error())
		return
	}
	if len(pgList.PortGroupIDs) == 0 {
		t.Error("expected at least one iSCSI PortGroup ID in list")
		return
	}
}

func TestGetPortGroupByFCID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	portGroup, err := client.GetPortGroupByID(context.TODO(), symmetrixID, defaultFCPortGroup)
	if err != nil || portGroup == nil {
		t.Error("Expected to find " + defaultFCPortGroup + " but didn't")
		return
	}
}

func TestGetPortGroupByiSCSIID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	portGroup, err := client.GetPortGroupByID(context.TODO(), symmetrixID, defaultiSCSIPortGroup)
	if err != nil || portGroup == nil {
		t.Error("Expected to find " + defaultiSCSIPortGroup + " but didn't")
		return
	}
	fmt.Printf("PortGroup: %#v\n", portGroup)
}

func TestGetInitiatorIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	initList, err := client.GetInitiatorList(context.TODO(), symmetrixID, "", true, false)
	if err != nil || initList == nil {
		t.Error("cannot get Initiator List: ", err.Error())
		return
	}
	if len(initList.InitiatorIDs) == 0 {
		t.Error("expected at least one Initiator ID in list")
		return
	}
	// Get the FC initiator list for the default initiator HBA
	initList, err = client.GetInitiatorList(context.TODO(), symmetrixID, defaultFCInitiatorID, false, true)
	if err != nil {
		t.Error("Receieved error : ", err.Error())
	}
	if len(initList.InitiatorIDs) != 0 {
		fmt.Println(initList.InitiatorIDs)
	} else {
		fmt.Println("Received an empty FC list")
		t.Error("Expected to find atleast one FC initiator")
	}
	// Get the iSCSI initiator list for the default IQN
	initList, err = client.GetInitiatorList(context.TODO(), symmetrixID, defaultiSCSIInitiatorID, true, true)
	if err != nil {
		t.Error("Receieved error : ", err.Error())
	}
	if len(initList.InitiatorIDs) != 0 {
		fmt.Println(initList.InitiatorIDs)
	} else {
		fmt.Println("Received an empty iSCSI list")
		t.Error("Expected to find atleast one iSCSI initiator")
	}

	// Get the initiator list for an IQN not on the array
	initList, err = client.GetInitiatorList(context.TODO(), symmetrixID, "iqn.1993-08.org.desian:01:5ae293b352a2", true, true)
	if err != nil {
		t.Error("Received error : ", err.Error())
	}
	if len(initList.InitiatorIDs) != 0 {
		fmt.Println(initList.InitiatorIDs)
	} else {
		fmt.Println("Received an empty list as expected for unknown IQN")
	}
}

func TestGetInitiatorByFCID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	initiator, err := client.GetInitiatorByID(context.TODO(), symmetrixID, defaultFCInitiator)
	if err != nil || initiator == nil {
		t.Error("Expected to find " + defaultFCInitiator + " but didn't")
		return
	}
	fmt.Printf("defaultFCInitator: %#v\n", initiator)
}

func TestGetInitiatorByiSCSIID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	initiator, err := client.GetInitiatorByID(context.TODO(), symmetrixID, defaultiSCSIInitiator)
	if err != nil || initiator == nil {
		t.Error("Expected to find " + defaultiSCSIInitiator + " but didn't")
		return
	}
	fmt.Printf("defaultiSCSIInitator: %#v\n", initiator)
}

// NVMe
func TestNVMeTCPGetInitiatorIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	initList, err := client.GetInitiatorList(context.TODO(), symmetrixID, "", false, false)
	if err != nil || initList == nil {
		t.Error("cannot get Initiator List: ", err.Error())
		return
	}
	if len(initList.InitiatorIDs) == 0 {
		t.Error("expected at least one Initiator ID in list")
		return
	}
	// Get the NVMeTCP initiator list for the default initiator ID
	initList, err = client.GetInitiatorList(context.TODO(), symmetrixID, defaultNVMeTCPInitiatorID, false, true)
	if err != nil {
		t.Error("Receieved error : ", err.Error())
	}
	if len(initList.InitiatorIDs) != 0 {
		fmt.Println(initList.InitiatorIDs)
	} else {
		fmt.Println("Received an empty NVMeTCP list")
		t.Error("Expected to find atleast one NVMeTCP initiator")
	}

	// Get the NVMeTCP initiator list for an NQN not on the array
	initList, err = client.GetInitiatorList(context.TODO(), symmetrixID, "nqn.2014-08.org.nvmexpress:uuid:csi_nvme:76b04d56e", false, true)
	if err != nil {
		t.Error("Received error : ", err.Error())
	}
	if len(initList.InitiatorIDs) != 0 {
		fmt.Println(initList.InitiatorIDs)
	} else {
		fmt.Println("Received an empty list as expected for unknown IQN")
	}
}

func TestGetInitiatorByNVMeTCPID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	initiator, err := client.GetInitiatorByID(context.TODO(), symmetrixID, defaultNVMeTCPInitiator)
	if err != nil || initiator == nil {
		t.Error("Expected to find " + defaultNVMeTCPInitiator + " but didn't")
		return
	}
	fmt.Printf("defaultNVMeTCPInitator: %#v\n", initiator)
}

func TestFCGetInitiators(t *testing.T) {
	// Get all the initiators and print the FC ones
	initList, err := client.GetInitiatorList(context.TODO(), symmetrixID, "", false, false)
	if err != nil || initList == nil {
		t.Error("cannot get Initiator List: ", err.Error())
		return
	}

	// Read our FC initiators from the /sys/class/fc_host directory
	cmd := exec.Command("/bin/sh", "-c", "cd /sys/class/fc_host; cat */port_name")
	bytes, err := cmd.Output()
	if err != nil {
		return
	}

	// Look for our initiators on the array
	ourInits := strings.Split(string(bytes), "\n")
	if len(ourInits) == 0 {
		// We have any initiators that we know
		return
	}
	fcInitiators := make([]string, 0)
	for _, ourInit := range ourInits {
		ourInit := strings.TrimSpace(strings.Replace(ourInit, "0x", "", 1))
		if ourInit == "" {
			continue
		}
		// fmt.Printf("ourInit: %s\n", ourInit)
		for _, init := range initList.InitiatorIDs {
			if strings.HasSuffix(init, ourInit) {
				fmt.Printf("initiator: %s\n", init)
				fcInitiators = append(fcInitiators, init)
			}
		}
	}

	// Print the matching initiator structures from the Symmetrix
	for _, fcInit := range fcInitiators {
		initiator, err := client.GetInitiatorByID(context.TODO(), symmetrixID, fcInit)
		if err != nil || initiator == nil {
			t.Errorf("Unable to read FC initiator: %s", err)
			return
		}
		fmt.Printf("%#v\n", initiator)
	}
}

func TestGetHostIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostList, err := client.GetHostList(context.TODO(), symmetrixID)
	for _, id := range hostList.HostIDs {
		fmt.Printf("Host ID: %s\n", id)
	}
	if err != nil || hostList == nil {
		t.Error("cannot get Host List: ", err.Error())
		return
	}
	if len(hostList.HostIDs) == 0 {
		t.Error("expected at least one Host ID in list")
		return
	}
}

func TestGetHostByFCID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	host, err := client.GetHostByID(context.TODO(), symmetrixID, defaultFCHost)
	if err != nil || host == nil {
		t.Error("Expected to find FC Host" + defaultFCHost + " but didn't")
		return
	}
	fmt.Printf("defaultHost: %#v\n", host)
}

func TestGetHostByiSCSIID(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	host, err := client.GetHostByID(context.TODO(), symmetrixID, defaultiSCSIHost)
	if err != nil || host == nil {
		t.Error("Expected to find " + defaultiSCSIHost + " but didn't")
		return
	}
	fmt.Printf("defaultHost: %#v\n", host)
}

func createHost(symmetrixID, hostID string, initiatorKeys []string, hostFlag *types.HostFlags) (*types.Host, error) {
	if client == nil {
		err := getClient()
		if err != nil {
			return nil, err
		}
	}
	// Check and return if the host exist on the array
	host, err := client.GetHostByID(context.TODO(), symmetrixID, hostID)
	if host != nil && err == nil {
		return host, nil
	}
	// host not found, create the host
	fmt.Println("Creating a new host...")
	return client.CreateHost(context.TODO(), symmetrixID, hostID, initiatorKeys, hostFlag)
}

func deleteHost(symmetrixID, hostID string) error {
	return client.DeleteHost(context.TODO(), symmetrixID, hostID)
}

func TestCreateFCHost(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			return
		}
	}
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, fcInitiator1)
	host, err := createHost(symmetrixID, "IntTestFCHost", initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create FC host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, FC host", host)
	err = deleteHost(symmetrixID, "IntTestFCHost")
	if err != nil {
		t.Error("Could not delete FC Host: " + err.Error())
	}
}

func TestUpdateHostFlags(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			return
		}
	}

	hostID := "IntTestHostFlags"
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, fcInitiator1)
	host, err := createHost(symmetrixID, hostID, initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create FC host but didn't: " + err.Error())
		return
	}

	hostFlags := &types.HostFlags{
		VolumeSetAddressing: &types.HostFlag{
			Enabled:  true,
			Override: true,
		},
		DisableQResetOnUA:   &types.HostFlag{},
		EnvironSet:          &types.HostFlag{},
		AvoidResetBroadcast: &types.HostFlag{},
		OpenVMS: &types.HostFlag{
			Override: true,
		},
		SCSI3:               &types.HostFlag{},
		Spc2ProtocolVersion: &types.HostFlag{},
		SCSISupport1:        &types.HostFlag{},
	}

	host, err = client.UpdateHostFlags(context.TODO(), symmetrixID, hostID, hostFlags)
	if err != nil || host == nil {
		t.Error("Failed to Update FC hostflags " + err.Error())
		return
	}

	if !strings.Contains(host.EnabledFlags, "Volume_Set_Addressing") {
		t.Error("Expected to Update FC hostflags but didn't: " + err.Error())
		return
	}

	fmt.Printf("%#v\n, FC host", host)
	err = deleteHost(symmetrixID, hostID)
	if err != nil {
		t.Error("Could not delete FC Host: " + err.Error())
	}
}

func TestCreateiSCSIHost(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			return
		}
	}
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, iscsiInitiator1)
	host, err := createHost(symmetrixID, "IntTestiSCSIHost", initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	err = deleteHost(symmetrixID, "IntTestiSCSIHost")
	if err != nil {
		t.Error("Could not delete Host: " + err.Error())
	}
}

func TestCreateNVMeTCPHost(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			return
		}
	}
	// In case a prior test left the host in the system
	err := deleteHost(symmetrixID, "IntTestNVMeTCPHost")
	if err != nil {
		t.Error("Could not delete Host: " + err.Error())
	}
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, nvmetcpInitiator1)
	host, err := createHost(symmetrixID, "IntTestNVMeTCPHost", initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	err = deleteHost(symmetrixID, "IntTestNVMeTCPHost")
	if err != nil {
		t.Error("Could not delete Host: " + err.Error())
	}
}

func TestCreateNVMeTCPMaskingView(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping this test in short mode")
	}
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostID := "IntTestNVMeTCPMV-Host"
	// In case a prior test left the Host in the system...
	client.DeleteHost(context.TODO(), symmetrixID, hostID)

	// create a Host with some initiators
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, nvmetcpInitiator1)
	fmt.Println("Setting up a host before creation of masking view")
	host, err := createHost(symmetrixID, hostID, initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	// add a volume in defaultStorageGroup
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, time.Now().Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error("Expected to create a volume but didn't" + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	maskingViewID := "IntTestNVMeTCPMV"
	maskingView, err := client.CreateMaskingView(context.TODO(), symmetrixID, maskingViewID, defaultStorageGroup,
		hostID, true, defaultNVMePortGroup)
	if err != nil {
		t.Error("Expected to create MV with NVMeTCP initiator and port but didn't: " + err.Error())
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Println("Fetching the newly created masking view from array")
	// Check if the MV exists on array
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err != nil || maskingView == nil {
		t.Error("Expected to find " + maskingViewID + " but didn't")
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Printf("%#v\n", maskingView)
	fmt.Println("Cleaning up the masking view")
	err = client.DeleteMaskingView(context.TODO(), symmetrixID, maskingViewID)
	if err != nil {
		t.Error("Failed to delete " + maskingViewID)
		return
	}
	fmt.Println("Sleeping for 20 seconds")
	time.Sleep(20 * time.Second)
	// Confirm if the masking view got deleted
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err == nil {
		t.Error("Expected a failure in fetching MV: " + maskingViewID + "but didn't")
		fmt.Printf("%#v\n", maskingView)
		return
	}
	fmt.Println(fmt.Sprintf("Error in fetching %s: %s", maskingViewID, err.Error()))
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
	cleanupHost(symmetrixID, hostID, t)
}

func TestCreateFCMaskingView(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping this test in short mode")
	}
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostID := "IntTestFCMV-Host"
	// In case a prior test left the Host in the system...
	client.DeleteHost(context.TODO(), symmetrixID, hostID)

	// create a Host with some initiators
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, fcInitiator2)
	fmt.Println("Setting up a host before creation of masking view")
	host, err := createHost(symmetrixID, hostID, initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	// add a volume in defaultStorageGroup
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, time.Now().Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error("Expected to create a volume but didn't" + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	maskingViewID := "IntTestFCMV"
	maskingView, err := client.CreateMaskingView(context.TODO(), symmetrixID, maskingViewID, defaultStorageGroup,
		hostID, true, defaultFCPortGroup)
	if err != nil {
		t.Error("Expected to create MV with FC initiator and port but didn't: " + err.Error())
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Println("Fetching the newly created masking view from array")
	// Check if the MV exists on array
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err != nil || maskingView == nil {
		t.Error("Expected to find " + maskingViewID + " but didn't")
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Printf("%#v\n", maskingView)
	fmt.Println("Cleaning up the masking view")
	err = client.DeleteMaskingView(context.TODO(), symmetrixID, maskingViewID)
	if err != nil {
		t.Error("Failed to delete " + maskingViewID)
		return
	}
	fmt.Println("Sleeping for 20 seconds")
	time.Sleep(20 * time.Second)
	// Confirm if the masking view got deleted
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err == nil {
		t.Error("Expected a failure in fetching MV: " + maskingViewID + "but didn't")
		fmt.Printf("%#v\n", maskingView)
		return
	}
	fmt.Println(fmt.Sprintf("Error in fetching %s: %s", maskingViewID, err.Error()))
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
	cleanupHost(symmetrixID, hostID, t)
}

func TestRenameMaskingView(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping this test in short mode")
	}
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostID := "IntTestFCMV-Host"
	// In case a prior test left the Host in the system...
	client.DeleteHost(context.TODO(), symmetrixID, hostID)

	// create a Host with some initiators
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, fcInitiator2)
	fmt.Println("Setting up a host before creation of masking view")
	host, err := createHost(symmetrixID, hostID, initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	// add a volume in defaultStorageGroup
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, time.Now().Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error("Expected to create a volume but didn't" + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	maskingViewID := "IntTestFCMV"
	maskingView, err := client.CreateMaskingView(context.TODO(), symmetrixID, maskingViewID, defaultStorageGroup,
		hostID, true, defaultFCPortGroup)
	if err != nil {
		t.Error("Expected to create MV with FC initiator and port but didn't: " + err.Error())
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Println("Fetching the newly created masking view from array")
	// Check if the MV exists on array
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err != nil || maskingView == nil {
		t.Error("Expected to find " + maskingViewID + " but didn't")
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Printf("%#v\n", maskingView)
	newMaskingViewID := "IntTestFCMV_new"
	fmt.Println("Renaming Masking view")
	renamedMaskingView, err1 := client.RenameMaskingView(context.TODO(), symmetrixID, maskingViewID, newMaskingViewID)
	if err1 != nil || renamedMaskingView == nil {
		t.Error("Error While Renaming Masking View")
		cleanupHost(symmetrixID, hostID, t)
	} else {
		fmt.Println("Successfully renamed Masking View!")
		fmt.Printf("%#v\n", renamedMaskingView)
	}
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err == nil {
		t.Error("Expected a failure in fetching MV: " + maskingViewID + "but didn't")
		fmt.Printf("%#v\n", maskingView)
		return
	}
	fmt.Println("Cleaning up the masking view")
	err = client.DeleteMaskingView(context.TODO(), symmetrixID, newMaskingViewID)
	if err != nil {
		t.Error("Failed to delete " + maskingViewID)
		return
	}
	fmt.Println("Sleeping for 20 seconds")
	time.Sleep(20 * time.Second)
	// Confirm if the masking view got deleted
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, newMaskingViewID)
	if err == nil {
		t.Error("Expected a failure in fetching MV: " + newMaskingViewID + "but didn't")
		fmt.Printf("%#v\n", maskingView)
		return
	}
	fmt.Println(fmt.Sprintf("Error in fetching %s: %s", maskingViewID, err.Error()))
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
	cleanupHost(symmetrixID, hostID, t)
}

func TestCreatePortGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	portGroupID := "IntTestPG"
	portKeys := make([]types.PortKey, 0)
	portKey := types.PortKey{
		DirectorID: defaultFCDirectorID,
		PortID:     defaultFCPortID,
	}
	portKeys = append(portKeys, portKey)
	portGroup, err := client.CreatePortGroup(context.TODO(), symmetrixID, portGroupID, portKeys, "SCSI_FC")
	if err != nil {
		t.Error("Couldn't create port group")
		return
	}
	fmt.Println(portGroup.PortGroupID)
	err = client.DeletePortGroup(context.TODO(), symmetrixID, portGroup.PortGroupID)
	if err != nil {
		t.Error("Couldn't delete port group")
		return
	}
}

func TestUpdatePortGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	portGroupID := "IntTestPG"
	portKeys := make([]types.PortKey, 0)
	portKey := types.PortKey{
		DirectorID: defaultFCDirectorID,
		PortID:     defaultFCPortID,
	}
	portKeys = append(portKeys, portKey)
	portGroup, err := client.CreatePortGroup(context.TODO(), symmetrixID, portGroupID, portKeys, "SCSI_FC")
	if err != nil {
		t.Error("Couldn't create port group")
		return
	}
	newPortGroupID := portGroup.PortGroupID
	portKeysUpdated := make([]types.PortKey, 0)
	portKey = types.PortKey{
		DirectorID: "OR-2C",
		PortID:     defaultFCPortID,
	}
	portKeysUpdated = append(portKeysUpdated, portKey)
	portGroup, err = client.UpdatePortGroup(context.TODO(), symmetrixID, newPortGroupID, portKeysUpdated)
	if err != nil {
		t.Errorf("Couldn't update port group: %s ", err.Error())
		return
	}
	portKeys = portGroup.SymmetrixPortKey
	// #nosec G602
	if portKeys[0].DirectorID != "OR-2C" {
		t.Errorf("Couldnt modify port group details: %s", err.Error())
		return
	}

	err = client.DeletePortGroup(context.TODO(), symmetrixID, newPortGroupID)
	if err != nil {
		t.Error("Couldn't delete port group")
		return
	}
}

func TestRenamePortGroup(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}

	portGroupID := "IntTestPG1"
	newPortGroupID := "test_rename_portgroup"
	portKeys := make([]types.PortKey, 0)
	portKey := types.PortKey{
		DirectorID: defaultFCDirectorID,
		PortID:     defaultFCPortID,
	}
	portKeys = append(portKeys, portKey)
	_, err := client.CreatePortGroup(context.TODO(), symmetrixID, portGroupID, portKeys, "SCSI_FC")
	if err != nil {
		t.Error("Couldn't create port group")
		return
	}

	_, err = client.GetPortGroupByID(context.TODO(), symmetrixID, portGroupID)
	if err != nil {
		t.Error("Couldn't get port group")
		return
	}

	portGroup, err := client.RenamePortGroup(context.TODO(), symmetrixID, portGroupID, newPortGroupID)
	if err != nil {
		t.Error("Couldn't rename port group")
		return
	}

	if portGroup.PortGroupID != newPortGroupID {
		t.Error("Couldn't rename port group")
		return
	}

	err = client.DeletePortGroup(context.TODO(), symmetrixID, newPortGroupID)
	if err != nil {
		t.Error("Couldn't delete port group")
		return
	}
}

func TestCreateiSCSIMaskingView(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping this test in short mode")
	}
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostID := "IntTestiSCSIMV-Host"
	// In case a prior test left the Host in the system...
	client.DeleteHost(context.TODO(), symmetrixID, hostID)

	// create a Host with some initiators
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, iscsiInitiator2)
	fmt.Println("Setting up a host before creation of masking view")
	host, err := createHost(symmetrixID, hostID, initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create iscsi host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)
	// add a volume in defaultStorageGroup
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, time.Now().Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 1, volOpts)
	if err != nil {
		t.Error("Expected to create a volume but didn't" + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	// create the masking view
	maskingViewID := "IntTestiSCSIMV"
	maskingView, err := client.CreateMaskingView(context.TODO(), symmetrixID, maskingViewID, defaultStorageGroup,
		hostID, true, defaultiSCSIPortGroup)
	if err != nil {
		t.Error("Expected to create MV with iscsi initiator host and port but didn't: " + err.Error())
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Println("Fetching the newly created masking view from array")
	// Check if the MV exists on array
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err != nil || maskingView == nil {
		t.Error("Expected to find " + maskingViewID + " but didn't")
		cleanupHost(symmetrixID, hostID, t)
		return
	}
	fmt.Printf("%#v\n", maskingView)
	fmt.Println("Cleaning up the masking view")
	err = client.DeleteMaskingView(context.TODO(), symmetrixID, maskingViewID)
	if err != nil {
		t.Error("Failed to delete " + maskingViewID)
		return
	}
	fmt.Println("Sleeping for 20 seconds")
	time.Sleep(20 * time.Second)
	// Confirm if the masking view got deleted
	maskingView, err = client.GetMaskingViewByID(context.TODO(), symmetrixID, maskingViewID)
	if err == nil {
		t.Error("Expected a failure in fetching MV: " + maskingViewID + "but didn't")
		fmt.Printf("%#v\n", maskingView)
		return
	}
	fmt.Println(fmt.Sprintf("Error in fetching %s: %s", maskingViewID, err.Error()))
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
	cleanupHost(symmetrixID, hostID, t)
}

func cleanupHost(symmetrixID string, hostID string, t *testing.T) {
	fmt.Println("Cleaning up the host")
	err := client.DeleteHost(context.TODO(), symmetrixID, hostID)
	if err != nil {
		t.Error("Failed to delete " + hostID)
	}
	return
}

func TestUpdateHostInitiators(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	// In case a prior test left the Host in the system...
	client.DeleteHost(context.TODO(), symmetrixID, "IntTestHost")

	// create a Host with some initiators
	initiatorKeys := make([]string, 0)
	initiatorKeys = append(initiatorKeys, iscsiInitiator1)
	host, err := client.CreateHost(context.TODO(), symmetrixID, "IntTestHost", initiatorKeys, nil)
	if err != nil || host == nil {
		t.Error("Expected to create host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)

	// change the list of initiators and update the host
	updatedInitiators := make([]string, 0)
	updatedInitiators = append(updatedInitiators, iscsiInitiator1, iscsiInitiator2)
	host, err = client.UpdateHostInitiators(context.TODO(), symmetrixID, host, updatedInitiators)
	if err != nil || host == nil {
		t.Error("Expected to update host but didn't: " + err.Error())
		return
	}
	fmt.Printf("%#v\n, host", host)

	// validate that we have the right number of intiators
	if len(host.Initiators) != len(updatedInitiators) {
		msg := fmt.Sprintf("Expected %d initiators but received %d", len(updatedInitiators), len(host.Initiators))
		t.Error(msg)
		return
	}

	// delete the host
	err = client.DeleteHost(context.TODO(), symmetrixID, "IntTestHost")
	if err != nil {
		t.Error("Could not delete Host: " + err.Error())
	}
}

func TestGetTargetAddresses(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	addresses, err := client.GetListOfTargetAddresses(context.TODO(), symmetrixID)
	if err != nil {
		t.Error("Error calling GetListOfTargetAddresses " + err.Error())
		return
	}
	fmt.Printf("Addresses: %v\n", addresses)
}

func TestGetISCSITargets(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	targets, err := client.GetISCSITargets(context.TODO(), symmetrixID)
	if err != nil {
		t.Error("Error calling GetISCSITargets " + err.Error())
		return
	}
	fmt.Printf("Targets: %v\n", targets)
}

func TestGetNVMeTCPTargets(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	targets, err := client.GetNVMeTCPTargets(context.TODO(), symmetrixID)
	if err != nil {
		t.Error("Error calling GetNVMeTCPTargets " + err.Error())
		return
	}
	fmt.Printf("Targets: %v\n", targets)
}

func TestExpandVolume(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	// create a volume
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 26, volOpts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	// expand Volume
	expandedSize := 30
	fmt.Println("Doing VolumeExpansion")
	expandedVol, err := client.ExpandVolume(context.TODO(), symmetrixID, vol.VolumeID, 0, expandedSize)
	if err != nil {
		t.Error("Error in Volume Expansion: " + err.Error())
		return
	}
	// check expand size
	if expandedVol.CapacityCYL != expandedSize {
		t.Error("Size mismatch after Expansion: " + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", expandedVol)
	fmt.Printf("Expanded Volume Size:\n%d\n", expandedVol.CapacityCYL)
	// all ok delete Volume
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func TestExpandVolumeWithUnit(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	// create a volume
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 26, volOpts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("volume:\n%#v\n", vol)
	// expand Volume
	expandedSize := 30
	capUnit := "GB"
	fmt.Println("Doing VolumeExpansion")
	expandedVol, err := client.ExpandVolume(context.TODO(), symmetrixID, vol.VolumeID, 0, expandedSize, capUnit)
	if err != nil {
		t.Error("Error in Volume Expansion: " + err.Error())
		return
	}
	// check expand size
	if expandedVol.CapacityGB != float64(expandedSize) {
		t.Error("Size mismatch after Expansion: " + err.Error())
		return
	}
	fmt.Printf("volume:\n%#v\n", expandedVol)
	fmt.Printf("Expanded Volume Size:\n%d\n", expandedVol.CapacityCYL)
	// all ok delete Volume
	cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
}

func TestHostGroup_CRUDOperation(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}
	hostGroupID := "IntTestHostGroup"
	hostIDs := []string{defaultFCHost}
	hostFlags := &types.HostFlags{
		VolumeSetAddressing: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		DisableQResetOnUA: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		EnvironSet: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		AvoidResetBroadcast: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		OpenVMS: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		SCSI3: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		Spc2ProtocolVersion: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		SCSISupport1: &types.HostFlag{
			Enabled:  false,
			Override: false,
		},
		ConsistentLUN: false,
	}
	_, err := client.CreateHostGroup(context.TODO(), symmetrixID, hostGroupID, hostIDs, hostFlags)
	if err != nil {
		t.Error("Couldn't create Host group")
		return
	}
	hostGroup, err := client.GetHostGroupByID(context.Background(), symmetrixID, hostGroupID)
	if err != nil || hostGroup.HostGroupID != hostGroupID {
		t.Error("Couldn't fetch the created Host group")
		return
	}

	hostFlags = &types.HostFlags{
		VolumeSetAddressing: &types.HostFlag{},
		DisableQResetOnUA:   &types.HostFlag{},
		EnvironSet:          &types.HostFlag{},
		AvoidResetBroadcast: &types.HostFlag{},
		OpenVMS:             &types.HostFlag{},
		SCSI3:               &types.HostFlag{},
		Spc2ProtocolVersion: &types.HostFlag{
			Enabled:  true,
			Override: true,
		},
		SCSISupport1:  &types.HostFlag{},
		ConsistentLUN: false,
	}

	hostGroup, err = client.UpdateHostGroupFlags(context.TODO(), symmetrixID, hostGroupID, hostFlags)
	if err != nil || hostGroup == nil {
		t.Error("Failed to Update FC hostflags " + err.Error())
		return
	}

	if !strings.Contains(hostGroup.EnabledFlags, "SPC2_Protocol_Version") {
		t.Error("Expected to Update FC hostflags but didn't: " + err.Error())
		return
	}

	newHostGroupID := hostGroupID + "-updated"

	hostGroup, err = client.UpdateHostGroupName(context.TODO(), symmetrixID, hostGroupID, newHostGroupID)
	if err != nil || hostGroup == nil {
		t.Error("Failed to rename hostgroup " + err.Error())
		return
	}

	if hostGroup.HostGroupID != newHostGroupID {
		t.Error("Expected to Update FC hostGroup name but didn't: " + err.Error())
		return
	}

	hostGroup, err = client.UpdateHostGroupHosts(context.TODO(), symmetrixID, hostGroup.HostGroupID, []string{})
	if err != nil {
		t.Error("Failed to Update hosts for the hostgroup " + err.Error())
		return
	}

	if len(hostGroup.Hosts) != 0 {
		t.Error("Expected to Remove FC hostGroup hosts but didn't: " + err.Error())
		return
	}

	hostGroup, err = client.UpdateHostGroupHosts(context.TODO(), symmetrixID, hostGroup.HostGroupID, hostIDs)
	if err != nil {
		t.Error("Failed to Update hosts for the hostgroup " + err.Error())
		return
	}

	if len(hostGroup.Hosts) == 0 {
		t.Error("Expected to Add FC hostGroup hosts but didn't: " + err.Error())
		return
	}

	fmt.Printf("%#v\n, FC hostGroup", hostGroup)

	err = client.DeleteHostGroup(context.TODO(), symmetrixID, hostGroup.HostGroupID)
	if err != nil {
		t.Error("Couldn't delete host group")
		return
	}
}

func TestGetHostGroupIDs(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Errorf("Unable to get/create pmax client: (%s)", err.Error())
			return
		}
	}

	// Creating a test hostgroup and asserting atleast one hostgroup is returned by the API.
	hostGroupID := "IntTestHostGroup"
	hostIDs := []string{defaultFCHost}
	hostFlags := &types.HostFlags{
		VolumeSetAddressing: &types.HostFlag{},
		DisableQResetOnUA:   &types.HostFlag{},
		EnvironSet:          &types.HostFlag{},
		AvoidResetBroadcast: &types.HostFlag{},
		OpenVMS:             &types.HostFlag{},
		SCSI3:               &types.HostFlag{},
		Spc2ProtocolVersion: &types.HostFlag{},
		SCSISupport1:        &types.HostFlag{},
		ConsistentLUN:       false,
	}
	_, err := client.CreateHostGroup(context.TODO(), symmetrixID, hostGroupID, hostIDs, hostFlags)
	if err != nil {
		t.Error("Couldn't create Host group")
		return
	}

	hostGroupList, err := client.GetHostGroupList(context.TODO(), symmetrixID)
	for _, id := range hostGroupList.HostGroupIDs {
		fmt.Printf("HostGroup ID: %s\n", id)
	}
	if err != nil || hostGroupList == nil {
		t.Error("cannot get HostGroup List: ", err.Error())
		return
	}
	if len(hostGroupList.HostGroupIDs) == 0 {
		t.Error("expected at least one HostGroup ID in list")
		return
	}

	err = client.DeleteHostGroup(context.TODO(), symmetrixID, hostGroupID)
	if err != nil {
		t.Error("Couldn't delete host group")
		return
	}
}

func TestGetStorageGroupPerfKeysAndMetrics(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	queryParams := []string{"HostMBReads"}
	perfKeys, err := client.GetStorageGroupPerfKeys(context.TODO(), symmetrixID)
	if err != nil {
		t.Errorf("Failed to call perf keys API")
		return
	}
	lastTime := int64(0)
	firstTime := int64(0)
	for _, storagePerfKey := range perfKeys.StorageGroupInfos {
		if storagePerfKey.StorageGroupID == defaultStorageGroup {
			lastTime = storagePerfKey.LastAvailableDate
			firstTime = storagePerfKey.LastAvailableDate
		}
	}
	if lastTime == int64(0) || firstTime == int64(0) {
		t.Errorf("Failed to get storage %s perf keys", defaultStorageGroup)
		return
	}
	metrics, err := client.GetStorageGroupMetrics(context.TODO(), symmetrixID, defaultStorageGroup, queryParams, firstTime, lastTime)
	if err != nil {
		t.Errorf("Failed to get storage group %s metrics", defaultStorageGroup)
		return
	}
	fmt.Printf("HostMBReads of storage group %s: %f \n", defaultStorageGroup, metrics.ResultList.Result[0].HostMBReads)
}

func TestGetVolumesMetrics(t *testing.T) {
	if client == nil {
		err := getClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	queryParams := []string{"MBRead"}
	// Ensure that default SG has at least 1 vol
	now := time.Now()
	volumeName := fmt.Sprintf("csi%s-Int%d", volumePrefix, now.Nanosecond())
	fmt.Printf("volumeName: %s\n", volumeName)
	volOpts := make(map[string]interface{})
	vol, err := client.CreateVolumeInStorageGroup(context.TODO(), symmetrixID, defaultStorageGroup, volumeName, 10, volOpts)
	defer cleanupVolume(vol.VolumeID, volumeName, defaultStorageGroup, t)
	if err != nil {
		t.Error(err)
		return
	}

	perfKeys, err := client.GetArrayPerfKeys(context.TODO())
	if err != nil {
		t.Errorf("Failed to call perf keys API")
		return
	}
	lastTime := int64(0)
	firstTime := int64(0)
	for _, storagePerfKey := range perfKeys.ArrayInfos {
		if storagePerfKey.SymmetrixID == symmetrixID {
			lastTime = storagePerfKey.LastAvailableDate
			firstTime = storagePerfKey.LastAvailableDate
		}
	}
	if lastTime == int64(0) || firstTime == int64(0) {
		t.Errorf("Failed to get array %s perf keys", symmetrixID)
		return
	}
	fmt.Println("... getting VolumeMetrics ...")
	metrics, err := client.GetVolumesMetrics(context.TODO(), symmetrixID, defaultStorageGroup, queryParams, firstTime, lastTime)
	if err != nil {
		t.Errorf("Failed to get volume in storage group %s Metrics", defaultStorageGroup)
		return
	}
	fmt.Printf("MBRead of volume %s: %f \n", metrics.ResultList.Result[0].VolumeID, metrics.ResultList.Result[0].VolumeResult[0].MBRead)

	fmt.Println("... getting VolumeMetricsByID ...")

	metrics, err = client.GetVolumesMetricsByID(context.TODO(), symmetrixID, vol.VolumeID, queryParams, firstTime, lastTime)
	if err != nil {
		t.Errorf("Failed to get volume in storage group %s Metrics", defaultStorageGroup)
		return
	}
	fmt.Printf("MBRead of volume %s: %f \n", metrics.ResultList.Result[0].VolumeID, metrics.ResultList.Result[0].VolumeResult[0].MBRead)
}
