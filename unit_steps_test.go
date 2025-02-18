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
package pmax

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/dell/gopowermax/v2/mock"
	types "github.com/dell/gopowermax/v2/types/v100"
)

const (
	defaultUsername         = "username"
	defaultPassword         = "password"
	symID                   = "000197900046"
	remoteSymID             = ""
	srdfMode                = "ASYNC"
	testPortGroup           = "12se0042-iscsi-PG"
	testInitiator           = "SE-1E:000:iqn.1993-08.org.debian:01:5ae293b352a2"
	testInitiatorIQN        = "iqn.1993-08.org.debian:01:5ae293b352a2"
	testUpdateInitiatorIQN  = "iqn.1993-08.org.debian:01:5ae293b352a3"
	testUpdateInitiator     = "SE-1E:000:iqn.1993-08.org.debian:01:5ae293b352a3"
	testNVMETCPInitiatorNQN = "nqn.2014-08.org.nvmexpress:uuid:csi_k8_nvme:76b04d56eab26a2e1509a7e98d3dfdb6"
	testNVMETCPInitiator    = "5e07d33b-d1ee-497c-97d8-22c0337ed8b8"
	testHost                = "l2se0042_iscsi_ig"
	testHostGroup           = "l2se0042_43_iscsi_ig"
	testSG                  = "l2se0042_sg"
	mvID                    = "12se0042_mv"
	testFCInitiatorWWN      = "10000090fa66060a"
	testFCInitiator         = "FA-1D:4:10000090fa66060a"
	protocol                = "SCSI_FC"
	queryNASServerID        = "nas_server_id"
	queryName               = "name"
)

type uMV struct {
	maskingViewID  string
	hostID         string
	hostGroupID    string
	storageGroupID string
	portGroupID    string
}

// client91 makes a client with version 91
// flag91 is set when a valid connection of client91 is established
// client91 is used when flag91 is set to have successful API calls of version 91 even if APIVersion is 90
type unitContext struct {
	nGoRoutines int
	client      Pmax
	client91    Pmax
	err         error // First error observed
	flag91      bool

	symIDList                  *types.SymmetrixIDList
	sym                        *types.Symmetrix
	directorIDList             *types.DirectorIDList
	vol                        *types.Volume
	volList                    []string
	storageGroup               *types.StorageGroup
	storageGroupSnapshotPolicy *types.StorageGroupSnapshotPolicy
	storageGroupIDList         *types.StorageGroupIDList
	jobIDList                  []string
	job                        *types.Job
	storagePoolList            *types.StoragePoolList
	portGroupList              *types.PortGroupList
	portGroup                  *types.PortGroup
	initiatorList              *types.InitiatorList
	initiator                  *types.Initiator
	hostList                   *types.HostList
	host                       *types.Host
	hostGroup                  *types.HostGroup
	hostGroupList              *types.HostGroupList
	maskingViewList            *types.MaskingViewList
	maskingView                *types.MaskingView
	uMaskingView               *uMV
	addressList                []string
	targetList                 []ISCSITarget
	nvmeTCPTargetList          []NVMeTCPTarget
	storagePool                *types.StoragePool
	volIDList                  []string
	hostID                     string
	hostGroupID                string
	sgID                       string

	symRepCapabilities    *types.SymReplicationCapabilities
	sourceVolumeList      []types.VolumeList
	symVolumeList         *types.SymVolumeList
	volSnapList           *types.SnapshotVolumeGeneration
	volumeSnapshot        *types.VolumeSnapshot
	volSnapGenerationList *types.VolumeSnapshotGenerations
	volSnapGenerationInfo *types.VolumeSnapshotGeneration
	volResultPrivate      *types.VolumeResultPrivate
	storageGroupMetrics   *types.StorageGroupMetricsIterator
	volumesMetrics        *types.VolumeMetricsIterator
	fileSystemMetrics     *types.FileSystemMetricsIterator

	sgSnapshot              *types.StorageGroupSnapshot
	storageGroupSnapSetting *types.CreateStorageGroupSnapshot
	storageGroupSnap        *types.StorageGroupSnap
	storageGroupSnapIDs     *types.SnapID

	storageGroupPerfKeys *types.StorageGroupKeysResult
	arrayPerfKeys        *types.ArrayKeysResult

	snapshotPolicy            *types.SnapshotPolicy
	createSnapshotPolicy      *types.CreateSnapshotPolicyParam
	modifySnapshotPolicyParam *types.ModifySnapshotPolicyParam
	snapshotPolicyList        *types.SnapshotPolicyList

	inducedErrors struct {
		badCredentials bool
		badPort        bool
		badIP          bool
	}

	fileSystemList *types.FileSystemIterator
	nasServerList  *types.NASServerIterator
	nfsExportList  *types.NFSExportIterator
	fileSystem     *types.FileSystem
	nfsExport      *types.NFSExport
	nasServer      *types.NASServer
	fileInterface  *types.FileInterface
}

func (c *unitContext) reset() {
	Debug = true
	c.flag91 = false
	c.err = nil
	c.symIDList = nil
	c.sym = nil
	c.vol = nil
	c.volList = make([]string, 0)
	c.storageGroup = nil
	c.storageGroupIDList = nil
	c.portGroupList = nil
	c.portGroup = nil
	c.initiatorList = nil
	c.initiator = nil
	c.hostList = nil
	c.host = nil
	c.jobIDList = nil
	c.job = nil
	c.storagePoolList = nil
	c.maskingViewList = nil
	c.uMaskingView = nil
	c.maskingView = nil
	c.storagePool = nil
	MAXJobRetryCount = 5
	c.volIDList = make([]string, 0)
	c.hostID = ""
	c.hostGroupID = ""
	c.sgID = ""
	c.symRepCapabilities = nil
	c.sourceVolumeList = make([]types.VolumeList, 0)
	c.symVolumeList = nil
	c.volSnapList = nil
	c.volumeSnapshot = nil
	c.volSnapGenerationList = nil
	c.volSnapGenerationInfo = nil
	c.volResultPrivate = nil
	c.fileSystemList = nil
	c.fileSystem = nil
	c.nfsExport = nil
	c.nasServer = nil
	c.fileInterface = nil
	c.volumesMetrics = nil
	c.fileSystemMetrics = nil
}

func (c *unitContext) iInduceError(errorType string) error {
	mock.InducedErrors.InvalidJSON = false
	mock.InducedErrors.BadHTTPStatus = 0
	mock.InducedErrors.GetSymmetrixError = false
	mock.InducedErrors.GetVolumeIteratorError = false
	mock.InducedErrors.GetVolumeError = false
	mock.InducedErrors.UpdateVolumeError = false
	mock.InducedErrors.DeleteVolumeError = false
	mock.InducedErrors.DeviceInSGError = false
	mock.InducedErrors.GetStorageGroupError = false
	mock.InducedErrors.GetStorageGroupSnapshotPolicyError = false
	mock.InducedErrors.DeleteStorageGroupSnapshotError = false
	mock.InducedErrors.InvalidResponse = false
	mock.InducedErrors.UpdateStorageGroupError = false
	mock.InducedErrors.GetJobError = false
	mock.InducedErrors.JobFailedError = false
	mock.InducedErrors.VolumeNotCreatedError = false
	mock.InducedErrors.GetJobCannotFindRoleForUser = false
	mock.InducedErrors.CreateStorageGroupError = false
	mock.InducedErrors.StorageGroupAlreadyExists = false
	mock.InducedErrors.DeleteStorageGroupError = false
	mock.InducedErrors.GetStoragePoolListError = false
	mock.InducedErrors.GetMaskingViewError = false
	mock.InducedErrors.GetPortGroupError = false
	mock.InducedErrors.GetInitiatorError = false
	mock.InducedErrors.GetHostError = false
	mock.InducedErrors.MaskingViewAlreadyExists = false
	mock.InducedErrors.DeleteMaskingViewError = false
	mock.InducedErrors.CreateMaskingViewError = false
	mock.InducedErrors.UpdateMaskingViewError = false
	mock.InducedErrors.PortGroupNotFoundError = false
	mock.InducedErrors.InitiatorGroupNotFoundError = false
	mock.InducedErrors.StorageGroupNotFoundError = false
	mock.InducedErrors.CreateHostError = false
	mock.InducedErrors.DeleteHostError = false
	mock.InducedErrors.VolumeNotAddedError = false
	mock.InducedErrors.UpdateHostError = false
	mock.InducedErrors.GetPortError = false
	mock.InducedErrors.GetSpecificPortError = false
	mock.InducedErrors.GetPortISCSITargetError = false
	mock.InducedErrors.GetPortGigEError = false
	mock.InducedErrors.GetDirectorError = false
	mock.InducedErrors.GetStoragePoolError = false
	mock.InducedErrors.ExpandVolumeError = false
	mock.InducedErrors.UpdatePortGroupError = false
	mock.InducedErrors.ModifyMobilityError = false
	mock.InducedErrors.CreateHostGroupError = false
	mock.InducedErrors.GetHostGroupError = false
	mock.InducedErrors.UpdateHostGroupError = false
	mock.InducedErrors.DeleteHostGroupError = false
	mock.InducedErrors.GetFreeRDFGError = false
	mock.InducedErrors.GetLocalOnlineRDFDirsError = false
	mock.InducedErrors.GetRemoteRDFPortOnSANError = false
	mock.InducedErrors.GetLocalOnlineRDFPortsError = false
	mock.InducedErrors.GetLocalRDFPortDetailsError = false
	mock.InducedErrors.CreateRDFGroupError = false
	mock.InducedErrors.GetRDFGroupError = false
	mock.InducedErrors.GetFileSystemListError = false
	mock.InducedErrors.GetNFSExportListError = false
	mock.InducedErrors.GetNASServerListError = false
	mock.InducedErrors.GetFileSystemError = false
	mock.InducedErrors.CreateFileSystemError = false
	mock.InducedErrors.UpdateFileSystemError = false
	mock.InducedErrors.DeleteFileSystemError = false
	mock.InducedErrors.GetNASServerError = false
	mock.InducedErrors.UpdateNASServerError = false
	mock.InducedErrors.DeleteNASServerError = false
	mock.InducedErrors.GetNFSExportError = false
	mock.InducedErrors.CreateNFSExportError = false
	mock.InducedErrors.UpdateNFSExportError = false
	mock.InducedErrors.DeleteNFSExportError = false
	mock.InducedErrors.GetFileInterfaceError = false
	mock.InducedErrors.ExecuteActionError = false
	mock.InducedErrors.CreateSnapshotPolicyError = false
	mock.InducedErrors.GetStorageGroupSnapshotError = false
	mock.InducedErrors.GetStorageGroupSnapshotSnapError = false
	mock.InducedErrors.GetStorageGroupSnapshotSnapDetailError = false
	mock.InducedErrors.GetStorageGroupSnapshotSnapModifyError = false
	mock.InducedErrors.GetSnapshotPolicyError = false
	mock.InducedErrors.GetSnapshotPolicyListError = false
	mock.InducedErrors.CreateSnapshotPolicyError = false
	mock.InducedErrors.ModifySnapshotPolicyError = false
	mock.InducedErrors.DeleteSnapshotPolicyError = false
	switch errorType {
	case "InvalidJSON":
		mock.InducedErrors.InvalidJSON = true
	case "httpStatus500":
		mock.InducedErrors.BadHTTPStatus = 500
	case "GetSymmetrixError":
		mock.InducedErrors.GetSymmetrixError = true
	case "GetVolumeIteratorError":
		mock.InducedErrors.GetVolumeIteratorError = true
	case "GetVolumeError":
		mock.InducedErrors.GetVolumeError = true
	case "UpdateVolumeError":
		mock.InducedErrors.UpdateVolumeError = true
	case "DeleteVolumeError":
		mock.InducedErrors.DeleteVolumeError = true
	case "DeviceInSGError":
		mock.InducedErrors.DeviceInSGError = true
	case "GetStorageGroupError":
		mock.InducedErrors.GetStorageGroupError = true
	case "GetStorageGroupSnapshotPolicyError":
		mock.InducedErrors.GetStorageGroupSnapshotPolicyError = true
	case "DeleteStorageGroupSnapshotError":
		mock.InducedErrors.DeleteStorageGroupSnapshotError = true
	case "InvalidResponse":
		mock.InducedErrors.InvalidResponse = true
	case "UpdateStorageGroupError":
		mock.InducedErrors.UpdateStorageGroupError = true
	case "GetJobError":
		mock.InducedErrors.GetJobError = true
	case "JobFailedError":
		mock.InducedErrors.JobFailedError = true
	case "VolumeNotCreatedError":
		mock.InducedErrors.VolumeNotCreatedError = true
	case "GetJobCannotFindRoleForUser":
		mock.InducedErrors.GetJobCannotFindRoleForUser = true
	case "CreateStorageGroupError":
		mock.InducedErrors.CreateStorageGroupError = true
	case "StorageGroupAlreadyExists":
		mock.InducedErrors.StorageGroupAlreadyExists = true
	case "DeleteStorageGroupError":
		mock.InducedErrors.DeleteStorageGroupError = true
	case "GetStoragePoolListError":
		mock.InducedErrors.GetStoragePoolListError = true
	case "GetMaskingViewError":
		mock.InducedErrors.GetMaskingViewError = true
	case "GetPortGroupError":
		mock.InducedErrors.GetPortGroupError = true
	case "GetInitiatorError":
		mock.InducedErrors.GetInitiatorError = true
	case "GetHostError":
		mock.InducedErrors.GetHostError = true
	case "CreateMaskingViewError":
		mock.InducedErrors.CreateMaskingViewError = true
	case "UpdateMaskingViewError":
		mock.InducedErrors.UpdateMaskingViewError = true
	case "MaskingViewAlreadyExists":
		mock.InducedErrors.MaskingViewAlreadyExists = true
	case "DeleteMaskingViewError":
		mock.InducedErrors.DeleteMaskingViewError = true
	case "PortGroupNotFoundError":
		mock.InducedErrors.PortGroupNotFoundError = true
	case "InitiatorGroupNotFoundError":
		mock.InducedErrors.InitiatorGroupNotFoundError = true
	case "StorageGroupNotFoundError":
		mock.InducedErrors.StorageGroupNotFoundError = true
	case "CreateHostError":
		mock.InducedErrors.CreateHostError = true
	case "DeleteHostError":
		mock.InducedErrors.DeleteHostError = true
	case "VolumeNotAddedError":
		mock.InducedErrors.VolumeNotAddedError = true
	case "UpdateHostError":
		mock.InducedErrors.UpdateHostError = true
	case "GetPortError":
		mock.InducedErrors.GetPortError = true
	case "GetSpecificPortError":
		mock.InducedErrors.GetSpecificPortError = true
	case "GetPortGigEError":
		mock.InducedErrors.GetPortGigEError = true
	case "GetPortISCSITargetError":
		mock.InducedErrors.GetPortISCSITargetError = true
	case "GetPortNVMeTCPTargetError":
		mock.InducedErrors.GetPortNVMeTCPTargetError = true
	case "GetDirectorError":
		mock.InducedErrors.GetDirectorError = true
	case "GetStoragePoolError":
		mock.InducedErrors.GetStoragePoolError = true
	case "GetSymVolumeError":
		mock.InducedErrors.GetSymVolumeError = true
	case "DeleteSnapshotError":
		mock.InducedErrors.DeleteSnapshotError = true
	case "GetGenerationError":
		mock.InducedErrors.GetGenerationError = true
	case "GetPrivateVolumeIterator":
		mock.InducedErrors.GetPrivateVolumeIterator = true
	case "GetVolSnapsError":
		mock.InducedErrors.GetVolSnapsError = true
	case "GetPrivVolumeByIDError":
		mock.InducedErrors.GetPrivVolumeByIDError = true
	case "GetStorageGroupSnapshotError":
		mock.InducedErrors.GetStorageGroupSnapshotError = true
	case "GetStorageGroupSnapshotSnapError":
		mock.InducedErrors.GetStorageGroupSnapshotSnapError = true
	case "GetStorageGroupSnapshotSnapDetailError":
		mock.InducedErrors.GetStorageGroupSnapshotSnapDetailError = true
	case "GetStorageGroupSnapshotSnapModifyError":
		mock.InducedErrors.GetStorageGroupSnapshotSnapModifyError = true
	case "CreatePortGroupError":
		mock.InducedErrors.CreatePortGroupError = true
	case "UpdatePortGroupError":
		mock.InducedErrors.UpdatePortGroupError = true
	case "DeletePortGroupError":
		mock.InducedErrors.DeletePortGroupError = true
	case "ExpandVolumeError":
		mock.InducedErrors.ExpandVolumeError = true
	case "ModifyMobilityError":
		mock.InducedErrors.ModifyMobilityError = true
	case "CreateHostGroupError":
		mock.InducedErrors.CreateHostGroupError = true
	case "GetHostGroupError":
		mock.InducedErrors.GetHostGroupError = true
	case "UpdateHostGroupError":
		mock.InducedErrors.UpdateHostGroupError = true
	case "DeleteHostGroupError":
		mock.InducedErrors.DeleteHostGroupError = true
	case "GetHostGroupListError":
		mock.InducedErrors.GetHostGroupListError = true
	case "GetStorageGroupMetricsError":
		mock.InducedErrors.GetStorageGroupMetricsError = true
	case "GetVolumesMetricsError":
		mock.InducedErrors.GetVolumesMetricsError = true
	case "GetFileSysMetricsError":
		mock.InducedErrors.GetFileSysMetricsError = true
	case "GetStorageGroupPerfKeyError":
		mock.InducedErrors.GetStorageGroupPerfKeyError = true
	case "GetArrayPerfKeyError":
		mock.InducedErrors.GetArrayPerfKeyError = true
	case "GetFreeRDFGError":
		mock.InducedErrors.GetFreeRDFGError = true
	case "GetLocalOnlineRDFDirsError":
		mock.InducedErrors.GetLocalOnlineRDFDirsError = true
	case "GetRemoteRDFPortOnSANError":
		mock.InducedErrors.GetRemoteRDFPortOnSANError = true
	case "GetLocalOnlineRDFPortsError":
		mock.InducedErrors.GetLocalOnlineRDFPortsError = true
	case "GetLocalRDFPortDetailsError":
		mock.InducedErrors.GetLocalRDFPortDetailsError = true
	case "CreateRDFGroupError":
		mock.InducedErrors.CreateRDFGroupError = true
	case "GetRDFGroupError":
		mock.InducedErrors.GetRDFGroupError = true
	case "GetSnapshotPolicyError":
		mock.InducedErrors.GetSnapshotPolicyError = true
	case "CreateSnapshotPolicyError":
		mock.InducedErrors.CreateSnapshotPolicyError = true
	case "GetSnapshotPolicyListError":
		mock.InducedErrors.GetSnapshotPolicyListError = true
	case "ModifySnapshotPolicyError":
		mock.InducedErrors.ModifySnapshotPolicyError = true
	case "DeleteSnapshotPolicyError":
		mock.InducedErrors.DeleteSnapshotPolicyError = true
	case "GetFileSystemListError":
		mock.InducedErrors.GetFileSystemListError = true
	case "GetNFSExportListError":
		mock.InducedErrors.GetNFSExportListError = true
	case "GetNASServerListError":
		mock.InducedErrors.GetNASServerListError = true
	case "GetFileSystemError":
		mock.InducedErrors.GetFileSystemError = true
	case "CreateFileSystemError":
		mock.InducedErrors.CreateFileSystemError = true
	case "UpdateFileSystemError":
		mock.InducedErrors.UpdateFileSystemError = true
	case "DeleteFileSystemError":
		mock.InducedErrors.DeleteFileSystemError = true
	case "GetNASServerError":
		mock.InducedErrors.GetNASServerError = true
	case "UpdateNASServerError":
		mock.InducedErrors.UpdateNASServerError = true
	case "DeleteNASServerError":
		mock.InducedErrors.DeleteNASServerError = true
	case "GetNFSExportError":
		mock.InducedErrors.GetNFSExportError = true
	case "CreateNFSExportError":
		mock.InducedErrors.CreateNFSExportError = true
	case "UpdateNFSExportError":
		mock.InducedErrors.UpdateNFSExportError = true
	case "DeleteNFSExportError":
		mock.InducedErrors.DeleteNFSExportError = true
	case "GetFileInterfaceError":
		mock.InducedErrors.GetFileInterfaceError = true
	case "ExecuteActionError":
		mock.InducedErrors.ExecuteActionError = true
	case "none":
	default:
		return fmt.Errorf("unknown errorType: %s", errorType)
	}
	return nil
}

func (c *unitContext) aValidConnection() error {
	c.reset()
	mock.Reset()
	if c.client == nil {
		apiVersion := strings.TrimSpace(os.Getenv("APIVersion"))
		err := c.iCallAuthenticateWithEndpointCredentials("", "", apiVersion)
		if err != nil {
			return err
		}
	}
	c.checkGoRoutines("aValidConnection")
	c.client.SetAllowedArrays([]string{})
	return nil
}

// Make a client with apiversion 91
func (c *unitContext) aValidv91Connection(_ int) error {
	c.reset()
	mock.Reset()
	// set the flag to insure client91 is used while making functions calls
	c.flag91 = true
	if c.client91 == nil {
		apiVersion := APIVersion91
		err := c.iCallAuthenticateWithEndpointCredentials("", "", apiVersion)
		if err != nil {
			return err
		}
	}
	c.checkGoRoutines("aValidV91Connection")
	c.client91.SetAllowedArrays([]string{})
	return nil
}

func (c *unitContext) checkGoRoutines(tag string) {
	goroutines := runtime.NumGoroutine()
	fmt.Printf("goroutines %s new %d old groutines %d\n", tag, goroutines, c.nGoRoutines)
	c.nGoRoutines = goroutines
}

func (c *unitContext) iCallAuthenticateWithEndpointCredentials(endpoint, credentials, apiVersion string) error {
	URL := mockServer.URL
	switch endpoint {
	case "badurl":
		URL = "https://127.0.0.99:2222"
	case "nilurl":
		URL = ""
	}
	fmt.Printf("apiVersion: %s\n", apiVersion)
	client, err := NewClientWithArgs(URL, "", true, false, "")
	if err != nil {
		c.err = err
		return nil
	}
	password := defaultPassword
	if credentials == "bad" {
		password = "xxx"
	}
	err = client.Authenticate(context.TODO(), &ConfigConnect{
		Endpoint: endpoint,
		Username: defaultUsername,
		Password: password,
	})
	if err == nil {
		if apiVersion == APIVersion91 {
			c.client91 = client
		} else {
			c.client = client
		}
	}
	c.err = err
	return nil
}

func (c *unitContext) theErrorMessageContains(expected string) error {
	if expected == "none" {
		if c.err == nil {
			return nil
		}
		return fmt.Errorf("Unexpected error: %s", c.err)
	}
	// We expected an error message
	if c.err == nil {
		return fmt.Errorf("Expected error message %s but no error was recorded", expected)
	}
	if strings.Contains(c.err.Error(), expected) {
		return nil
	}
	return fmt.Errorf("Expected error message to contain: %s but the error message was: %s", expected, c.err)
}

func (c *unitContext) iCallGetSymmetrixIDList() error {
	c.symIDList, c.err = c.client.GetSymmetrixIDList(context.TODO())
	return nil
}

func (c *unitContext) iGetAValidSymmetrixIDListIfNoError() error {
	if c.err == nil {
		if c.symIDList == nil {
			return fmt.Errorf("SymmetrixIDList nil")
		}
		if len(c.symIDList.SymmetrixIDs) == 0 {
			return fmt.Errorf("No IDs in SymmetrixIDList")
		}
	}
	return nil
}

func (c *unitContext) iCallGetSymmetrixByID(id string) error {
	c.sym, c.err = c.client.GetSymmetrixByID(context.TODO(), id)
	return nil
}

func (c *unitContext) iGetAValidSymmetrixObjectIfNoError() error {
	if c.err == nil {
		if c.sym == nil {
			return fmt.Errorf("Symmetrix nil")
		}
		fmt.Printf("Symmetrix: %#v", c.sym)
		if c.sym.SymmetrixID == "" || c.sym.Ucode == "" || c.sym.Model == "" || c.sym.DiskCount <= 0 {
			return fmt.Errorf("Problem with Symmetrix fields SymmetrixID Ucode Model or DiskCount")
		}
	}
	return nil
}

func (c *unitContext) iHaveVolumes(number int) error {
	for i := 1; i <= number; i++ {
		id := fmt.Sprintf("%05d", i)
		size := 7
		volumeIdentifier := "Vol" + id
		mock.AddNewVolume(id, volumeIdentifier, size, mock.DefaultStorageGroup)
		c.volIDList = append(c.volIDList, id)
		// mock.Data.VolumeIDToIdentifier[id] = fmt.Sprintf("Vol%05d", i)
		// mock.Data.VolumeIDToSGList[id] = make([]string, 0)
	}
	return nil
}

func (c *unitContext) iCallGetVolumeByID(volID string) error {
	c.vol, c.err = c.client.GetVolumeByID(context.TODO(), symID, volID)
	return nil
}

func (c *unitContext) iGetAValidVolumeObjectIfNoError(id string) error {
	if c.err != nil {
		return nil
	}
	if c.vol.VolumeID != id {
		return fmt.Errorf("Expected volume %s but got %s", id, c.vol.VolumeID)
	}
	return nil
}

func (c *unitContext) iCallGetVolumeIDList(volumeIdentifier string) error {
	var like bool
	if strings.Contains(volumeIdentifier, "<like>") {
		volumeIdentifier = strings.TrimPrefix(volumeIdentifier, "<like>")
		like = true
	}
	c.volList, c.err = c.client.GetVolumeIDList(context.TODO(), symID, volumeIdentifier, like)
	return nil
}

func (c *unitContext) iCallGetVolumeIDListWithParams() error {
	param := map[string]string{
		"tdev":   "true",
		"status": "Ready,<like>Read",
		"cap_gb": ">10.0",
		"cap_tb": "=1.0",
	}
	c.volList, c.err = c.client.GetVolumeIDListWithParams(context.TODO(), symID, param)
	return nil
}

func (c *unitContext) iCallGetVolumeIDListInStorageGroup(sgID string) error {
	c.volList, c.err = c.client.GetVolumeIDListInStorageGroup(context.TODO(), symID, sgID)
	return nil
}

func (c *unitContext) iGetAValidVolumeIDListWithIfNoError(nvols int) error {
	if c.err != nil {
		return nil
	}
	if len(c.volList) != nvols {
		return fmt.Errorf("Expected %d volumes but got %d", nvols, len(c.volList))
	}
	return nil
}

func (c *unitContext) iExpandVolumeToSize(volumeID string, sizeStr string) error {
	if c.err != nil {
		return nil
	}

	if size, err := strconv.Atoi(sizeStr); err == nil {
		c.vol, c.err = c.client.ExpandVolume(context.TODO(), symID, volumeID, 0, size)
	} else {
		return err
	}

	return nil
}

func (c *unitContext) iExpandVolumeToSizeWithUnit(volumeID string, sizeStr string, capUnits string) error {
	if c.err != nil {
		return nil
	}

	if size, err := strconv.Atoi(sizeStr); err == nil {
		c.vol, c.err = c.client.ExpandVolume(context.TODO(), symID, volumeID, 0, size, capUnits)
	} else {
		return err
	}

	return nil
}

func (c *unitContext) iCallModifyMobilityForVolume(volumeID string, mobility string) error {
	if c.err != nil {
		return nil
	}
	mobilityBool, _ := strconv.ParseBool(mobility)
	c.vol, c.err = c.client.ModifyMobilityForVolume(context.TODO(), symID, volumeID, mobilityBool)
	return nil
}

func (c *unitContext) iValidateVolumeSize(volumeID string, sizeStr string) error {
	if c.err != nil {
		return nil
	}

	c.vol, c.err = c.client.GetVolumeByID(context.TODO(), symID, volumeID)
	size, err := strconv.Atoi(sizeStr)
	if err == nil && float64(size) != c.vol.CapacityGB {
		return fmt.Errorf("Expected volume %s to be size %s, but was %d", volumeID, sizeStr, size)
	} else if err != nil {
		return err
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroupIDList() error {
	c.storageGroupIDList, c.err = c.client.GetStorageGroupIDList(context.TODO(), symID, "", false)
	return nil
}

func (c *unitContext) iGetAValidStorageGroupIDListIfNoErrors() error {
	if c.err != nil {
		return nil
	}
	if len(c.storageGroupIDList.StorageGroupIDs) == 0 {
		return fmt.Errorf("Expected storage group IDs to be returned but there were none")
	}
	for _, id := range c.storageGroupIDList.StorageGroupIDs {
		fmt.Printf("StorageGroup: %s\n", id)
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroup(sgID string) error {
	c.storageGroup, c.err = c.client.GetStorageGroup(context.TODO(), symID, sgID)
	return nil
}

func (c *unitContext) iCallGetStorageGroupSnapshotPolicy(symID, snapshotPolicyID, storageGroupID string) error {
	c.storageGroupSnapshotPolicy, c.err = c.client.GetStorageGroupSnapshotPolicy(context.TODO(), symID, snapshotPolicyID, storageGroupID)
	return nil
}

func (c *unitContext) iGetAValidStorageGroupIfNoErrors() error {
	if c.err != nil {
		return nil
	}
	if c.storageGroup.StorageGroupID == "" || c.storageGroup.Type == "" {
		return fmt.Errorf("Expected StorageGroup to have StorageGroupID and Type but didn't")
	}
	return nil
}

func (c *unitContext) iGetAValidStorageGroupSnapshotPolicyObjectIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.storageGroupSnapshotPolicy.StorageGroupID == "IntSGB" {
		return fmt.Errorf("Storage Group [IntSGB] on Symmetrix [000000000002] cannot be found")
	}
	return nil
}

func (c *unitContext) iCallGetStoragePool(srpID string) error {
	c.storagePool, c.err = c.client.GetStoragePool(context.TODO(), symID, srpID)
	return nil
}

func (c *unitContext) iGetAValidGetStoragePoolIfNoErrors() error {
	if c.err != nil {
		return nil
	}
	if c.storagePool.StoragePoolID == "" {
		return fmt.Errorf("Expected StoragePool to have StoragePoolID and Type but didn't")
	}
	return nil
}

func (c *unitContext) iHaveJobs(numberOfJobs int) error {
	for i := 1; i <= numberOfJobs; i++ {
		jobID := fmt.Sprintf("job%d", i)
		mock.NewMockJob(jobID, "RUNNING", "SUCCEEDED", "")
	}
	return nil
}

func (c *unitContext) iCallGetJobIDListWith(statusQuery string) error {
	c.jobIDList, c.err = c.client.GetJobIDList(context.TODO(), symID, statusQuery)
	return nil
}

func (c *unitContext) iGetAValidJobsIDListWithIfNoErrors(numberOfEntries int) error {
	if c.err != nil {
		return nil
	}
	if len(c.jobIDList) != numberOfEntries {
		return fmt.Errorf("Expected %d jobs ids to be returned but got %d", numberOfEntries, len(c.jobIDList))
	}
	return nil
}

func (c *unitContext) iCreateAJobWithInitialStateAndFinalState(initialState, finalState string) error {
	mock.NewMockJob("myjob", initialState, finalState, "")
	return nil
}

func (c *unitContext) iCallGetJobByID() error {
	c.job, c.err = c.client.GetJobByID(context.TODO(), symID, "myjob")
	return nil
}

func (c *unitContext) iGetAValidJobWithStateIfNoError(expectedState string) error {
	if c.err != nil {
		return nil
	}
	if c.job.Status != expectedState {
		return fmt.Errorf("Expected job state to be %s but instead it was %s", expectedState, c.job.Status)
	}
	return nil
}

func (c *unitContext) iCallWaitOnJobCompletion() error {
	c.job, c.err = c.client.WaitOnJobCompletion(context.TODO(), symID, "myjob")
	return nil
}

func (c *unitContext) iCallCreateVolumeInStorageGroupWithNameAndSize(volumeName string, sizeInCylinders int) error {
	volOpts := make(map[string]interface{})
	if !c.flag91 {
		c.vol, c.err = c.client.CreateVolumeInStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	} else {
		c.vol, c.err = c.client91.CreateVolumeInStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	}
	return nil
}

func (c *unitContext) iCallCreateVolumeInStorageGroupWithNameAndSizeAndUnit(volumeName string, sizeInCylinders int, capUnit string) error {
	volOpts := make(map[string]interface{})
	volOpts["capacityUnit"] = capUnit
	volOpts["enableMobility"] = false
	if !c.flag91 {
		c.vol, c.err = c.client.CreateVolumeInStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	} else {
		c.vol, c.err = c.client91.CreateVolumeInStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	}
	return nil
}

func (c *unitContext) iCallCreateVolumeInStorageGroupSWithNameAndSize(volumeName string, sizeInCylinders int) error {
	volOpts := make(map[string]interface{})
	if !c.flag91 {
		c.vol, c.err = c.client.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	} else {
		c.vol, c.err = c.client91.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts)
	}
	return nil
}

func (c *unitContext) iCallCreateVolumeInStorageGroupSWithNameAndSizeAndUnit(volumeName string, sizeInCylinders int, capUnit string) error {
	var size interface{}
	volOpts := make(map[string]interface{})
	volOpts["capacityUnit"] = capUnit
	volOpts["enableMobility"] = false
	if capUnit != "CYL" {
		size = strconv.Itoa(sizeInCylinders)
	} else {
		size = sizeInCylinders
	}

	if !c.flag91 {
		c.vol, c.err = c.client.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, size, volOpts)
	} else {
		c.vol, c.err = c.client91.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, size, volOpts)
	}
	return nil
}

func (c *unitContext) iCallCreateVolumeInStorageGroupSWithNameAndSizeWithMetaDataHeaders(volumeName string, sizeInCylinders int) error {
	metadata := make(http.Header)
	metadata.Set("x-csi-pv-name", "testPVName")
	metadata.Set("x-csi-pv-claimname", "testPVClaimName")
	metadata.Set("x-csi-pv-namespace", "testPVNamespace")
	volOpts := make(map[string]interface{})
	if !c.flag91 {
		c.vol, c.err = c.client.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts, metadata)
	} else {
		c.vol, c.err = c.client91.CreateVolumeInStorageGroupS(context.TODO(), symID, mock.DefaultStorageGroup, volumeName, sizeInCylinders, volOpts, metadata)
	}
	return nil
}

func (c *unitContext) iGetAValidVolumeWithNameIfNoError(volumeName string) error {
	if c.err != nil {
		return nil
	}
	if c.vol.VolumeIdentifier != volumeName {
		return fmt.Errorf("Expected volume named %s but got %s", volumeName, c.vol.VolumeIdentifier)
	}
	return nil
}

func (c *unitContext) iGetAValidVolumeWithMobilityModified(mobility string) error {
	if c.err != nil {
		return nil
	}
	mobilityBool, _ := strconv.ParseBool(mobility)
	if c.vol.MobilityIDEnabled != mobilityBool {
		return fmt.Errorf("Expected volume mobility-enabled: %v but %v ", mobilityBool, c.vol.MobilityIDEnabled)
	}
	return nil
}

func (c *unitContext) iCallCreateStorageGroupWithNameAndSrpAndSl(sgName, srp, serviceLevel string) error {
	if !c.flag91 {
		c.storageGroup, c.err = c.client.CreateStorageGroup(context.TODO(), symID, sgName, srp, serviceLevel, false, nil)
	} else {
		c.storageGroup, c.err = c.client91.CreateStorageGroup(context.TODO(), symID, sgName, srp, serviceLevel, false, nil)
	}
	return nil
}

func (c *unitContext) iCallCreateStorageGroupWithNameAndSrpAndSlAndHostLimits(sgName, srp, serviceLevel string, hl string) error {
	limits := convertStringSliceOfHostLimitsToHostLimitParams(hl)
	if !c.flag91 {
		c.storageGroup, c.err = c.client.CreateStorageGroup(context.TODO(), symID, sgName, srp, serviceLevel, false, *limits)
	} else {
		c.storageGroup, c.err = c.client91.CreateStorageGroup(context.TODO(), symID, sgName, srp, serviceLevel, false, *limits)
	}
	return nil
}

func (c *unitContext) iGetAValidStorageGroupWithNameIfNoError(sgName string) error {
	if c.err != nil {
		return nil
	}
	if c.storageGroup.StorageGroupID != sgName {
		return fmt.Errorf("Expected StorageGroup to have name %s", sgName)
	}
	return nil
}

func (c *unitContext) iCallDeleteStorageGroup(sgID string) error {
	c.err = c.client.DeleteStorageGroup(context.TODO(), symID, sgID)
	return nil
}

func (c *unitContext) iCallGetStoragePoolList() error {
	c.storagePoolList, c.err = c.client.GetStoragePoolList(context.TODO(), symID)
	return nil
}

func (c *unitContext) iGetAValidStoragePoolListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.storagePoolList == nil || len(c.storagePoolList.StoragePoolIDs) <= 0 || c.storagePoolList.StoragePoolIDs[0] != "SRP_1" {
		return fmt.Errorf("Expected StoragePoolList to have SRP_1 but it didn't")
	}
	return nil
}

func (c *unitContext) iCallRemoveVolumeFromStorageGroup() error {
	if !c.flag91 {
		c.storageGroup, c.err = c.client.RemoveVolumesFromStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, true, c.vol.VolumeID)
	} else {
		c.storageGroup, c.err = c.client91.RemoveVolumesFromStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, true, c.vol.VolumeID)
	}
	return nil
}

func (c *unitContext) theVolumeIsNoLongerAMemberOfTheStorageGroupIfNoError() error {
	if c.err != nil {
		return nil
	}
	sgIDList := mock.Data.VolumeIDToSGList[c.vol.VolumeID]
	if len(sgIDList) == 0 {
		return nil
	}
	for _, sgid := range sgIDList {
		if sgid == mock.DefaultStorageGroup {
			return fmt.Errorf("Volume contained the Storage Group %s which was not expected", mock.DefaultStorageGroup)
		}
	}
	return nil
}

func (c *unitContext) iCallRenameVolumeWith(newName string) error {
	c.vol, c.err = c.client.RenameVolume(context.TODO(), symID, c.vol.VolumeID, newName)
	return nil
}

func (c *unitContext) iCallInitiateDeallocationOfTracksFromVolume() error {
	c.job, c.err = c.client.InitiateDeallocationOfTracksFromVolume(context.TODO(), symID, c.vol.VolumeID)
	return nil
}

func (c *unitContext) iCallDeleteVolume() error {
	c.err = c.client.DeleteVolume(context.TODO(), symID, c.vol.VolumeID)
	return nil
}

func (c *unitContext) iHaveAMaskingView(maskingViewID string) error {
	sgID := maskingViewID + "-sg"
	pgID := maskingViewID + "-pg"
	hostID := maskingViewID + "-host"
	localMaskingView := &uMV{
		maskingViewID:  maskingViewID,
		hostID:         hostID,
		hostGroupID:    "",
		storageGroupID: sgID,
		portGroupID:    pgID,
	}
	initiators := []string{testInitiatorIQN}
	mock.AddInitiator(testInitiator, testInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	mock.AddHost(hostID, "iSCSI", initiators)
	mock.AddStorageGroup(sgID, "SRP_1", "Diamond")
	mock.AddMaskingView(maskingViewID, sgID, hostID, pgID)
	c.uMaskingView = localMaskingView
	return nil
}

func (c *unitContext) iCallGetMaskingViewList() error {
	c.maskingViewList, c.err = c.client.GetMaskingViewList(context.TODO(), symID)
	return nil
}

func (c *unitContext) iGetAValidMaskingViewListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.maskingViewList == nil || len(c.maskingViewList.MaskingViewIDs) == 0 {
		return fmt.Errorf("Expected item in MaskingViewList but got none")
	}
	found := false
	for _, id := range c.maskingViewList.MaskingViewIDs {
		fmt.Printf("MaskingView: %s\n", id)
		if c.uMaskingView.maskingViewID == id {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("Expected to find %s in MaskingViewList but didn't", c.uMaskingView.maskingViewID)
	}
	return nil
}

func (c *unitContext) iCallGetMaskingViewByID(mvID string) error {
	c.maskingView, c.err = c.client.GetMaskingViewByID(context.TODO(), symID, mvID)
	return nil
}

func (c *unitContext) iGetAValidMaskingViewIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.maskingView == nil {
		return fmt.Errorf("Expecting a masking view but received none")
	}
	if c.hostGroupID == "" {
		if c.maskingView.HostID != c.uMaskingView.hostID {
			return fmt.Errorf("Expecting host %s but got %s", c.uMaskingView.hostID, c.maskingView.HostID)
		}
	} else {
		if c.maskingView.HostID != c.uMaskingView.hostGroupID {
			return fmt.Errorf("Expecting hostgroup %s but got %s", c.uMaskingView.hostGroupID, c.maskingView.HostID)
		}
	}
	if c.maskingView.PortGroupID != c.uMaskingView.portGroupID {
		return fmt.Errorf("Expecting portgroup %s but got %s", c.uMaskingView.portGroupID, c.maskingView.PortGroupID)
	}
	if c.maskingView.StorageGroupID != c.uMaskingView.storageGroupID {
		return fmt.Errorf("Expecting storagegroup %s but got %s", c.uMaskingView.storageGroupID, c.maskingView.StorageGroupID)
	}
	return nil
}

func (c *unitContext) iCallDeleteMaskingView() error {
	c.err = c.client.DeleteMaskingView(context.TODO(), symID, c.uMaskingView.maskingViewID)
	return nil
}

func (c *unitContext) iCallRenameMaskingViewWith(newName string) error {
	c.maskingView, c.err = c.client.RenameMaskingView(context.TODO(), symID, c.uMaskingView.maskingViewID, newName)
	return nil
}

func (c *unitContext) iHaveAPortGroup() error {
	mock.AddPortGroupWithPortID(testPortGroup, "ISCSI", []string{"SE-1E:000"})
	return nil
}

func (c *unitContext) iCallGetPortGroupList() error {
	c.portGroupList, c.err = c.client.GetPortGroupList(context.TODO(), symID, "")
	return nil
}

func (c *unitContext) iGetAValidPortGroupListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.portGroupList == nil || len(c.portGroupList.PortGroupIDs) == 0 {
		return fmt.Errorf("Expected item in PortGroupList but got none")
	}
	return nil
}

func (c *unitContext) iCallGetPortGroupByID() error {
	c.portGroup, c.err = c.client.GetPortGroupByID(context.TODO(), symID, testPortGroup)
	return nil
}

func (c *unitContext) iCallCreatePortGroup(groupName string, strSliceOfPorts string) error {
	if c.err != nil {
		return nil
	}

	initialPorts := convertStringSliceOfPortsToPortKeys(strSliceOfPorts)
	c.portGroup, c.err = c.client.CreatePortGroup(context.TODO(), symID, groupName, initialPorts, protocol)
	return nil
}

func (c *unitContext) iCallRenamePortGroupWith(newName string) error {
	c.portGroup, c.err = c.client.RenamePortGroup(context.TODO(), symID, c.portGroup.PortGroupID, newName)
	return nil
}

func (c *unitContext) iCallUpdatePortGroup(groupName string, strUpdatePorts string) error {
	if c.err != nil {
		return nil
	}

	updatedPorts := convertStringSliceOfPortsToPortKeys(strUpdatePorts)
	c.portGroup, c.err = c.client.UpdatePortGroup(context.TODO(), symID, groupName, updatedPorts)
	return nil
}

func (c *unitContext) iExpectedThesePortsInPortGroup(strSliceOfPorts string) error {
	if c.err != nil {
		return nil
	}

	expectedPorts := convertStringSliceOfPortsToPortKeys(strSliceOfPorts)
	if c.portGroup == nil {
		return errors.New("could not find any portGroup. Make sure test was set up with a PortGroup")
	}

	expectedPortsLen := len(expectedPorts)
	portGroupLen := len(c.portGroup.SymmetrixPortKey)
	if expectedPortsLen != portGroupLen {
		return fmt.Errorf("expected number of ports does not match. Expected %d, but portGroup %s has %d", expectedPortsLen, c.portGroup.PortGroupID, portGroupLen)
	}

	portKeySlice := make([]string, 0)
	portsInPortGroup := make(map[string]bool)
	for _, its := range c.portGroup.SymmetrixPortKey {
		thisKey := fmt.Sprintf("%s:%s", its.DirectorID, its.PortID)
		if !portsInPortGroup[thisKey] {
			portsInPortGroup[thisKey] = true
			portKeySlice = append(portKeySlice, thisKey)
		}
	}

	for _, its := range expectedPorts {
		thisKey := fmt.Sprintf("%s:%s", its.DirectorID, its.PortID)
		if !portsInPortGroup[thisKey] {
			return fmt.Errorf("list of ports in PortGroup do not match expected list. Expected %s, but got %s", strSliceOfPorts, strings.Join(portKeySlice, ","))
		}
	}

	return nil
}

func (c *unitContext) iCallDeletePortGroup(groupName string) error {
	if c.err != nil {
		return nil
	}

	c.err = c.client.DeletePortGroup(context.TODO(), symID, groupName)

	return nil
}

func (c *unitContext) thePortGroupShouldNotExist(groupName string) error {
	if c.err != nil {
		return nil
	}

	c.portGroupList, c.err = c.client.GetPortGroupList(context.TODO(), symID, "")
	for _, id := range c.portGroupList.PortGroupIDs {
		if id == groupName {
			return fmt.Errorf("PortGroup %s was not expected, but is in PortGroup list", groupName)
		}
	}
	return nil
}

func (c *unitContext) iGetPortGroupIfNoError(groupName string) error {
	if c.err != nil {
		return nil
	}

	if c.portGroup.PortGroupID != groupName {
		return fmt.Errorf("Expected to get Port Group %s, but received %s",
			c.portGroup.PortGroupID, groupName)
	}
	return nil
}

func (c *unitContext) iGetAValidPortGroupIfNoError() error {
	if c.err != nil {
		return nil
	}

	if c.portGroup.PortGroupID != testPortGroup {
		return fmt.Errorf("Expected to get Port Group %s, but received %s",
			c.portGroup.PortGroupID, testPortGroup)
	}
	return nil
}

func (c *unitContext) iHaveAISCSIHost(hostName string) error {
	initiators := []string{testInitiatorIQN}
	mock.AddInitiator(testInitiator, testInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	c.hostID = hostName
	c.host, c.err = mock.AddHost(c.hostID, "iSCSI", initiators)
	return nil
}

func (c *unitContext) iHaveAFCHost(hostName string) error {
	initiators := []string{testFCInitiatorWWN}
	mock.AddInitiator(testFCInitiator, testFCInitiatorWWN, "Fibre", []string{"FA-1D:4"}, "")
	c.hostID = hostName
	c.host, c.err = mock.AddHost(c.hostID, "Fibre", initiators)
	return nil
}

func (c *unitContext) iHaveANVMETCPHost(hostName string) error {
	initiators := []string{testNVMETCPInitiatorNQN}
	mock.AddInitiator(testNVMETCPInitiator, testNVMETCPInitiatorNQN, "GigE", []string{"OR-1C:001"}, "")
	c.hostID = hostName
	c.host, c.err = mock.AddHost(c.hostID, "NVMe", initiators)
	return nil
}

func (c *unitContext) iCallGetHostList() error {
	c.hostList, c.err = c.client.GetHostList(context.TODO(), symID)
	return nil
}

func (c *unitContext) iGetAValidHostListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.hostList == nil || len(c.hostList.HostIDs) == 0 {
		return fmt.Errorf("Expected item in HostList but got none")
	}
	fmt.Println(c.hostList)
	return nil
}

func (c *unitContext) iCallGetHostByID(hostID string) error {
	c.host, c.err = c.client.GetHostByID(context.TODO(), symID, hostID)
	return nil
}

func (c *unitContext) iGetAValidHostIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.host.HostID != c.hostID {
		return fmt.Errorf("Expected to get Host %s, but received %s",
			c.host.HostID, c.hostID)
	}
	return nil
}

func (c *unitContext) iCallCreateHost(hostName string) error {
	initiatorList := make([]string, 1)
	c.hostID = hostName
	initiatorList[0] = testInitiatorIQN
	mock.AddInitiator(testInitiator, testInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	c.host, c.err = c.client.CreateHost(context.TODO(), symID, hostName, initiatorList, nil)
	return nil
}

func (c *unitContext) iCallUpdateHost() error {
	initiatorList := make([]string, 1)
	initiatorList[0] = testUpdateInitiatorIQN
	mock.AddInitiator(testUpdateInitiator, testUpdateInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	c.host, c.err = c.client.UpdateHostInitiators(context.TODO(), symID, c.host, initiatorList)
	return nil
}

func (c *unitContext) iCallUpdateHostFlags() error {
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
	c.host, c.err = c.client.UpdateHostFlags(context.TODO(), symID, c.hostID, hostFlags)
	return nil
}

func (c *unitContext) iCallDeleteHost(hostName string) error {
	c.err = c.client.DeleteHost(context.TODO(), symID, hostName)
	return nil
}

func (c *unitContext) iHaveAInitiator() error {
	return nil
}

func (c *unitContext) iCallGetInitiatorList() error {
	c.initiatorList, c.err = c.client.GetInitiatorList(context.TODO(), symID, "", false, false)
	return nil
}

func (c *unitContext) iCallGetInitiatorListWithFilters() error {
	c.initiatorList, c.err = c.client.GetInitiatorList(context.TODO(), symID, testInitiatorIQN, true, true)
	return nil
}

func (c *unitContext) iGetAValidInitiatorListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.initiatorList == nil || len(c.initiatorList.InitiatorIDs) == 0 {
		return fmt.Errorf("Expected item in InitiatorList but got none")
	}
	return nil
}

func (c *unitContext) iCallGetInitiatorByID() error {
	mock.AddInitiator(testInitiator, testInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	c.initiator, c.err = c.client.GetInitiatorByID(context.TODO(), symID, testInitiator)
	return nil
}

func (c *unitContext) iGetAValidInitiatorIfNoError() error {
	if c.err != nil {
		return nil
	}

	if c.initiator.InitiatorID != testInitiatorIQN {
		return fmt.Errorf("Expected to get initiator %s, but received %s",
			c.initiator.InitiatorID, testInitiator)
	}
	return nil
}

func (c *unitContext) iHaveAStorageGroup(sgID string) error {
	c.sgID = sgID
	mock.AddStorageGroup(sgID, "SRP_1", "Diamond")
	return nil
}

func (c *unitContext) iCallCreateMaskingViewWithHost(mvID string) error {
	localMaskingView := &uMV{
		maskingViewID:  mvID,
		hostID:         c.hostID,
		hostGroupID:    "",
		storageGroupID: c.sgID,
		portGroupID:    testPortGroup,
	}
	c.uMaskingView = localMaskingView
	c.maskingView, c.err = c.client.CreateMaskingView(context.TODO(), symID, mvID, c.sgID, c.hostID, true, testPortGroup)
	return nil
}

func (c *unitContext) iCallCreateMaskingViewWithHostGroup(mvID string) error {
	localMaskingView := &uMV{
		maskingViewID:  mvID,
		hostID:         "",
		hostGroupID:    c.hostGroupID,
		storageGroupID: c.sgID,
		portGroupID:    testPortGroup,
	}
	c.uMaskingView = localMaskingView
	c.maskingView, c.err = c.client.CreateMaskingView(context.TODO(), symID, mvID, c.sgID, c.hostGroupID, false, testPortGroup)
	return nil
}

func (c *unitContext) iHaveAHostGroup(hostGroupID string) error {
	// Create a host instead of host group
	c.hostGroupID = hostGroupID
	initiators := []string{testInitiatorIQN}
	mock.AddInitiator(testInitiator, testInitiatorIQN, "GigE", []string{"SE-1E:000"}, "")
	mock.AddHost(hostGroupID, "iSCSI", initiators)
	return nil
}

func (c *unitContext) iCallAddVolumesToStorageGroup(sgID string) error {
	if !c.flag91 {
		c.err = c.client.AddVolumesToStorageGroup(context.TODO(), symID, sgID, true, c.volIDList...)
	} else {
		c.err = c.client91.AddVolumesToStorageGroup(context.TODO(), symID, sgID, true, c.volIDList...)
	}
	return nil
}

func (c *unitContext) iCallAddVolumesToStorageGroupS(sgID string) error {
	if !c.flag91 {
		c.err = c.client.AddVolumesToStorageGroupS(context.TODO(), symID, sgID, true, c.volIDList...)
	} else {
		c.err = c.client91.AddVolumesToStorageGroupS(context.TODO(), symID, sgID, true, c.volIDList...)
	}
	return nil
}

func (c *unitContext) thenTheVolumesArePartOfStorageGroupIfNoError() error {
	if c.err != nil {
		return nil
	}
	sgList := mock.Data.VolumeIDToSGList["00001"]
	fmt.Printf("%v", sgList)
	for volumeID := range mock.Data.VolumeIDToIdentifier {
		fmt.Println(volumeID)
		sgList := mock.Data.VolumeIDToSGList[volumeID]
		fmt.Printf("%v", sgList)
		volumeFound := false
		for _, sg := range sgList {
			if sg == mock.DefaultStorageGroup {
				volumeFound = true
				break
			}
		}
		if !volumeFound {
			return fmt.Errorf("Couldn't find volume in storage group")
		}
	}
	return nil
}

func (c *unitContext) iCallGetListOfTargetAddresses() error {
	c.addressList, c.err = c.client.GetListOfTargetAddresses(context.TODO(), symID)
	return nil
}

func (c *unitContext) iRecieveIPAddresses(count int) error {
	if len(c.addressList) != count {
		return fmt.Errorf("Expected to get %d addresses but recieved %d", count, len(c.addressList))
	}
	return nil
}

func (c *unitContext) iHaveAnAllowedListOf(listOfAllowedArrays string) error {
	// turn the string into a slice
	results := convertStringToSlice(listOfAllowedArrays)

	// set the list of allowed arrays
	if !c.flag91 {
		c.client.SetAllowedArrays(results)
	} else {
		c.client91.SetAllowedArrays(results)
	}
	return nil
}

func (c *unitContext) itContainsArrays(count int) error {
	allowed := c.client.GetAllowedArrays()
	if len(allowed) != count {
		return fmt.Errorf("Received the wrong number of arrays in the allowed list. Expected %d but have a allowed list of %v", count, allowed)
	}
	return nil
}

func (c *unitContext) shouldInclude(include string) error {
	// turn the list of specified arrays into a slice
	results := convertStringToSlice(include)

	// make sure each one is in the allowed list of arrays
	for _, a := range results {
		if ok, _ := c.client.IsAllowedArray(a); ok == false {
			return fmt.Errorf("Expected array (%s) to be in the allowed list but it was not found", a)
		}
	}
	return nil
}

func (c *unitContext) shouldNotInclude(exclude string) error {
	// turn the list of specified arrays into a slice
	results := convertStringToSlice(exclude)

	// make sure each one is not in the allowed list of arrays
	for _, a := range results {
		if ok, _ := c.client.IsAllowedArray(a); ok == true {
			return fmt.Errorf("Expected array (%s) to not be in the allowd list but it was", a)
		}
	}
	return nil
}

func (c *unitContext) iGetAValidSymmetrixIDListThatContainsAndDoesNotContains(included string, excluded string) error {
	includedArrays := convertStringToSlice(included)
	// make sure all the included arrays exist in the response
	for _, array := range includedArrays {
		found := false
		for _, expectedArray := range c.symIDList.SymmetrixIDs {
			if array == expectedArray {
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("Expected array %s to be included in %v, but it was not", array, c.symIDList.SymmetrixIDs)
		}
	}

	excludedArrays := convertStringToSlice(excluded)
	// make sure all the excluded arrays do NOT exist in the response
	for _, array := range excludedArrays {
		found := false
		for _, expectedArray := range c.symIDList.SymmetrixIDs {
			if array == expectedArray {
				found = true
			}
		}
		if found != false {
			return fmt.Errorf("Expected array %s to be excluded in %v, but it was not", array, c.symIDList.SymmetrixIDs)
		}
	}

	return nil
}

func convertStringToSlice(input string) []string {
	results := make([]string, 0)
	st := strings.Split(input, ",")
	for i := range st {
		t := strings.TrimSpace(st[i])
		if t != "" {
			results = append(results, t)
		}
	}
	return results
}

func (c *unitContext) iExcuteTheCapabilitiesOnTheSymmetrixArray() error {
	c.symRepCapabilities, c.err = c.client.GetReplicationCapabilities(context.TODO())
	return nil
}

func (c *unitContext) iCallGetSnapVolumeListWithAnd(queryKey, queryValue string) error {
	if queryKey != "" {
		if queryValue == "true" {
			c.symVolumeList, c.err = c.client.GetSnapVolumeList(context.TODO(), symID, types.QueryParams{
				queryKey: true,
			})
		}
	} else {
		c.symVolumeList, c.err = c.client.GetSnapVolumeList(context.TODO(), symID, nil)
	}
	return nil
}

func (c *unitContext) iShouldGetListOfVolumesHavingSnapshots() error {
	if c.err != nil {
		return nil
	}
	if len(c.symVolumeList.Name) == 0 {
		return fmt.Errorf("No Volumes with Snapshot found")
	}
	return nil
}

func (c *unitContext) iCallGetVolumeSnapInfoWithVolume(volID string) error {
	c.volSnapList, c.err = c.client.GetVolumeSnapInfo(context.TODO(), symID, volID)
	return nil
}

func (c *unitContext) iShouldGetAListOfSnapshotsIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.volSnapList == nil {
		return fmt.Errorf("There is no Snapshot for provided source volume")
	}
	return nil
}

func (c *unitContext) iCallCreateSnapshotWithAndSnapshotOnIt(volIDs, snapID string) error {
	c.sourceVolumeList = c.createVolumeList(volIDs)
	c.err = c.client.CreateSnapshot(context.TODO(), symID, snapID, c.sourceVolumeList, 0)
	return nil
}

func (c *unitContext) iGetAValidSnapshotObjectIfNoError() error {
	if c.err != nil {
		return nil
	}
	for i := range c.sourceVolumeList {
		sourceVol := c.sourceVolumeList[i]
		if mock.Data.VolIDToSnapshots[sourceVol.Name] == nil {
			return fmt.Errorf("The snaphshot does not exist for source volume %s", sourceVol.Name)
		}
	}
	return nil
}

func (c *unitContext) iCallGetSnapshotInfoWithAndSnapshotNameOnIt(volID, SnapID string) error {
	c.volumeSnapshot, c.err = c.client.GetSnapshotInfo(context.TODO(), symID, volID, SnapID)
	return nil
}

func (c *unitContext) iShouldGetTheSnapshotDetailsIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.volumeSnapshot == nil {
		return fmt.Errorf("The Snapshot does not exist")
	}
	return nil
}

func (c *unitContext) iCallGetSnapshotGenerationsWithAndSnapshotOnIt(volID, SnapID string) error {
	c.volSnapGenerationList, c.err = c.client.GetSnapshotGenerations(context.TODO(), symID, volID, SnapID)
	return nil
}

func (c *unitContext) iShouldGetTheGenerationListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.volSnapGenerationList == nil {
		return fmt.Errorf("The Generation List for the snaphshot does not exist")
	}
	return nil
}

func (c *unitContext) iCallGetSnapshotGenerationWithSnapshotAndOnIt(volID string, snapID string, genID int64) error {
	c.volSnapGenerationInfo, c.err = c.client.GetSnapshotGenerationInfo(context.TODO(), symID, volID, snapID, genID)
	return nil
}

func (c *unitContext) iShouldGetAGenerationInfoIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.volSnapGenerationInfo == nil {
		return fmt.Errorf("The Generation for the snaphshot does not exist")
	}
	return nil
}

func (c *unitContext) iCallModifySnapshotWithAnd(sourceVols, targetVols, SnapID, newSnapID string, genID int64, action string) error {
	sourceVolumeList := c.createVolumeList(sourceVols)
	targetVolumeList := c.createVolumeList(targetVols)
	c.err = c.client.ModifySnapshot(context.TODO(), symID, sourceVolumeList, targetVolumeList, SnapID, action, newSnapID, genID, false)

	return nil
}

func (c *unitContext) iCallModifySnapshotSWithAnd(sourceVols, targetVols, SnapID, newSnapID string, genID int64, action string) error {
	sourceVolumeList := c.createVolumeList(sourceVols)
	targetVolumeList := c.createVolumeList(targetVols)
	c.err = c.client.ModifySnapshotS(context.TODO(), symID, sourceVolumeList, targetVolumeList, SnapID, action, newSnapID, genID, false)

	return nil
}

func (c *unitContext) iCallDeleteSnapshotWithSnapshotAndOnIt(sourceVols, SnapID string, genID int64) error {
	sourceVolumeList := c.createVolumeList(sourceVols)
	c.err = c.client.DeleteSnapshot(context.TODO(), symID, SnapID, sourceVolumeList, genID)
	return nil
}

func (c *unitContext) iCallDeleteSnapshotSWithSnapshotAndOnIt(sourceVols, SnapID string, genID int64) error {
	sourceVolumeList := c.createVolumeList(sourceVols)
	c.err = c.client.DeleteSnapshotS(context.TODO(), symID, SnapID, sourceVolumeList, genID)
	return nil
}

func (c *unitContext) iCallGetPrivVolumeByIDWith(volID string) error {
	c.volResultPrivate, c.err = c.client.GetPrivVolumeByID(context.TODO(), symID, volID)
	return nil
}

func (c *unitContext) iCallGetStorageGroupSnapshotsWith(storageGroupID string) error {
	c.sgSnapshot, c.err = c.client.GetStorageGroupSnapshots(context.TODO(), symID, storageGroupID, false, false)
	return nil
}

func (c *unitContext) iCallGetStorageGroupSnapshotsWithAndParam(storageGroupID, params string) error {
	var exludeManualSnaps bool
	var exludeSlSnaps bool
	param := strings.Split(params, ",")
	for _, p := range param {
		if p == "exludeManualSnaps" {
			exludeManualSnaps = true
		}
		if p == "exludeSlSnaps" {
			exludeSlSnaps = true
		}
	}
	c.sgSnapshot, c.err = c.client.GetStorageGroupSnapshots(context.TODO(), symID, storageGroupID, exludeManualSnaps, exludeSlSnaps)
	return nil
}

func (c *unitContext) iCallCreateStorageGroupSnapshotWith(storageGroupID string) error {
	c.storageGroupSnap, c.err = c.client.CreateStorageGroupSnapshot(context.TODO(), symID, storageGroupID, c.storageGroupSnapSetting)
	return nil
}

func (c *unitContext) iShouldGetStorageGroupSnapshotInformationIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.sgSnapshot == nil {
		return fmt.Errorf("The storage group snapshot does not exist")
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroupSnapshotSnapIDsWithAnd(storageGroupID string, snapshotID string) error {
	c.storageGroupSnapIDs, c.err = c.client.GetStorageGroupSnapshotSnapIDs(context.TODO(), symID, storageGroupID, snapshotID)
	return nil
}

func (c *unitContext) iShouldGetStorageGroupSnapshotSnapIDsIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.storageGroupSnapIDs == nil {
		return fmt.Errorf("The storage group snapshot snap does not exist")
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroupSnapshotSnapWithAndAnd(storageGroupID string, snapshotID string, snapID string) error {
	c.storageGroupSnap, c.err = c.client.GetStorageGroupSnapshotSnap(context.TODO(), symID, storageGroupID, snapshotID, snapID)
	return nil
}

func (c *unitContext) iShouldGetStorageGroupSnapshotSnapDetailInformationIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.storageGroupSnap == nil {
		return fmt.Errorf("The storage group snapshot snap details does not exist")
	}
	return nil
}

func (c *unitContext) iCallModifyStorageGroupSnapshotWithAndAndAndAction(storageGroupID string, snapshotID string, snapID string, action string) error {
	var payload *types.ModifyStorageGroupSnapshot
	switch action {
	case "rename":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "Rename",
			Rename: types.RenameSnapshotAction{
				NewStorageGroupSnapshotName: "sg_1_snap_2",
			},
		}
	case "restore":
		payload = &types.ModifyStorageGroupSnapshot{
			Action:  "Restore",
			Restore: types.RestoreSnapshotAction{},
		}
	case "link":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "Link",
			Link: types.LinkSnapshotAction{
				StorageGroupName: "sg_1_2",
			},
		}
	case "relink":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "Relink",
			Relink: types.RelinkSnapshotAction{
				StorageGroupName: "sg_1_2",
			},
		}
	case "unlink":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "Unlink",
			Unlink: types.UnlinkSnapshotAction{
				StorageGroupName: "sg_1_2",
			},
		}
	case "setmode":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "SetMode",
			SetMode: types.SetModeSnapshotAction{
				StorageGroupName: "sg_1_snap_2",
			},
		}
	case "timeToLive":
		payload = &types.ModifyStorageGroupSnapshot{
			Action:     "SetTimeToLive",
			TimeToLive: types.TimeToLiveSnapshotAction{},
		}
	case "secure":
		payload = &types.ModifyStorageGroupSnapshot{
			Action: "SetSecure",
			Secure: types.SecureSnapshotAction{},
		}
	case "persist":
		payload = &types.ModifyStorageGroupSnapshot{
			Action:  "Persist",
			Persist: types.PresistSnapshotAction{},
		}
	}
	c.storageGroupSnap, c.err = c.client.ModifyStorageGroupSnapshot(context.TODO(), symID, storageGroupID, snapshotID, snapID, payload)
	return nil
}

func (c *unitContext) iShouldModifyStorageGroupSnapshotSnapIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.storageGroupSnap == nil {
		return fmt.Errorf("The storage group snapshot snap details does not exist")
	}
	return nil
}

func (c *unitContext) iCallDeleteStorageGroupSnapshotWithAndAnd(storageGroupID string, snapshotID string, snapID string) error {
	c.err = c.client.DeleteStorageGroupSnapshot(context.TODO(), symID, storageGroupID, snapshotID, snapID)
	return nil
}

func (c *unitContext) iShouldGetAPrivateVolumeInformationIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.volResultPrivate == nil {
		return fmt.Errorf("The private information for the volume does not exist")
	}
	return nil
}

func (c *unitContext) iShouldGetAValidResponseIfNoError() error {
	if c.err != nil {
		return nil
	}
	return nil
}

func (c *unitContext) thereShouldBeNoErrors() error {
	return c.err
}

// createVolumeList will extract all the volumes and will return a list of type VolumeList
func (c *unitContext) createVolumeList(volIDs string) []types.VolumeList {
	var VolumeList []types.VolumeList
	volNames := strings.Split(volIDs, ",")
	for i := 0; i < len(volNames); i++ {
		VolumeList = append(VolumeList, types.VolumeList{Name: volNames[i]})
	}
	return VolumeList
}

// convertStringSliceOfPortsToPortKeys - Given a comma delimited string of ports in
// this format "<DIRECTOR>:<PORT>", produce a slice of types.PortKey values
func convertStringSliceOfPortsToPortKeys(strListOfPorts string) []types.PortKey {
	initialPorts := make([]types.PortKey, 0)
	for _, it := range convertStringToSlice(strListOfPorts) {
		dirAndPort := strings.Split(it, ":")
		port := types.PortKey{
			DirectorID: dirAndPort[0],
			PortID:     dirAndPort[1],
		}
		initialPorts = append(initialPorts, port)
	}
	return initialPorts
}

func convertStringSliceOfHostLimitsToHostLimitParams(strListOfHostLimts string) *map[string]interface{} {
	hostLimitValues := strings.Split(strListOfHostLimts, ":")
	hostLimits := types.SetHostIOLimitsParam{
		HostIOLimitMBSec:    hostLimitValues[0],
		HostIOLimitIOSec:    hostLimitValues[1],
		DynamicDistribution: hostLimitValues[2],
	}

	optionalPayload := make(map[string]interface{})
	optionalPayload["hostLimits"] = &hostLimits
	return &optionalPayload
}

func (c *unitContext) iCallGetISCSITargets() error {
	c.targetList, c.err = c.client.GetISCSITargets(context.TODO(), symID)
	return nil
}

func (c *unitContext) iCallGetNVMeTCPTargets() error {
	c.nvmeTCPTargetList, c.err = c.client.GetNVMeTCPTargets(context.TODO(), symID)
	return nil
}

func (c *unitContext) iCallRefreshSymmetrix(id string) error {
	c.err = c.client.RefreshSymmetrix(context.TODO(), id)
	return nil
}

func (c *unitContext) iRecieveTargets(count int) error {
	if len(c.targetList) != count {
		return fmt.Errorf("expected to get %d targets but recieved %d", count, len(c.targetList))
	}
	return nil
}

func (c *unitContext) iCallUpdateHostName(newName string) error {
	c.host, c.err = c.client.UpdateHostName(context.TODO(), symID, c.hostID, newName)
	return nil
}

func (c *unitContext) iCallCreateSGReplica(mode string) error {
	localSG, remoteSG := mock.DefaultStorageGroup, mock.DefaultStorageGroup // Using same names for local and remote storage groups
	remoteServiceLevel := mock.DefaultServiceLevel                          // Using the same service level as local
	var rdfgNumber string
	if mode == "METRO" {
		rdfgNumber = fmt.Sprintf("%d", mock.DefaultMetroRDFGNo)
	} else {
		rdfgNumber = fmt.Sprintf("%d", mock.DefaultAsyncRDFGNo)
	}
	_, c.err = c.client.CreateSGReplica(context.TODO(), symID, mock.DefaultRemoteSymID, mode, rdfgNumber, localSG, remoteSG, remoteServiceLevel, false)
	return nil
}

func (c *unitContext) thenSGShouldBeReplicated() error {
	if c.err != nil {
		return nil
	}
	if mock.Data.StorageGroupIDToRDFStorageGroup[mock.DefaultStorageGroup] == nil ||
		mock.Data.StorageGroupIDToStorageGroup[mock.DefaultStorageGroup].Unprotected {
		return fmt.Errorf("storage group not protected")
	}
	return nil
}

func (c *unitContext) checkReplication(volume *types.Volume, compliment string) bool {
	isReplicated := strings.Contains(volume.Type, "RDF") || len(volume.RDFGroupIDList) > 0
	if compliment == "not" {
		isReplicated = !isReplicated
	}
	return isReplicated
}

func (c *unitContext) theVolumesShouldBeReplicated(compliment string) error {
	if c.err != nil {
		return nil
	}
	for _, volumeID := range mock.Data.StorageGroupIDToVolumes[mock.DefaultStorageGroup] {
		volume := mock.Data.VolumeIDToVolume[volumeID]
		if !c.checkReplication(volume, compliment) {
			return fmt.Errorf("volumes[%s] should %s be replicated", volumeID, compliment)
		}
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroupRDFInfo() error {
	var sgrdf *types.StorageGroupRDFG
	sgrdf, c.err = c.client.GetStorageGroupRDFInfo(context.TODO(), symID, mock.DefaultStorageGroup, fmt.Sprintf("%d", mock.DefaultAsyncRemoteRDFGNo))
	if c.err == nil {
		if sgrdf.SymmetrixID != symID ||
			sgrdf.StorageGroupName != mock.DefaultStorageGroup ||
			sgrdf.RdfGroupNumber != mock.DefaultAsyncRDFGNo {
			c.err = fmt.Errorf("the returned storage group doesn't contain proper details")
		}
	}
	return nil
}

func (c *unitContext) iCallGetRDFDevicePairInfo() error {
	var devicePairInfo *types.RDFDevicePair
	localVolID := c.volIDList[0]
	remoteVolID := c.volIDList[0]

	devicePairInfo, c.err = c.client.GetRDFDevicePairInfo(context.TODO(), symID, fmt.Sprintf("%d", mock.DefaultAsyncRemoteRDFGNo), localVolID)
	if c.err == nil {
		if devicePairInfo.LocalSymmID != symID || devicePairInfo.RemoteSymmID != mock.DefaultRemoteSymID ||
			devicePairInfo.LocalVolumeName != localVolID || devicePairInfo.RemoteVolumeName != remoteVolID {
			c.err = fmt.Errorf("incorrect rdf-pair info returned")
		}
	}
	return nil
}

func (c *unitContext) iCallGetProtectedStorageGroup() error {
	var protectedSG *types.RDFStorageGroup
	protectedSG, c.err = c.client.GetProtectedStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup)
	if c.err == nil {
		if protectedSG.SymmetrixID != symID || protectedSG.Name != mock.DefaultStorageGroup {
			c.err = fmt.Errorf("protected sg with incorrect details returned")
		}
		if _, ok := mock.Data.StorageGroupIDToRDFStorageGroup[mock.DefaultStorageGroup]; !ok && protectedSG.Rdf {
			c.err = fmt.Errorf("storage group is not protected, but the rdf is set to true")
		}
	}
	return nil
}

func (c *unitContext) iCallGetRDFGroup() error {
	var rdfGroup *types.RDFGroup
	rdfGroup, c.err = c.client.GetRDFGroupByID(context.TODO(), symID, fmt.Sprintf("%d", mock.DefaultAsyncRemoteRDFGNo))
	if c.err == nil {
		if !rdfGroup.Async || rdfGroup.RemoteSymmetrix != mock.DefaultRemoteSymID || rdfGroup.RdfgNumber != mock.DefaultAsyncRemoteRDFGNo {
			c.err = fmt.Errorf("rdf group with incorrect details returned")
		}
	}
	return nil
}

func (c *unitContext) iCallAddVolumesToProtectedStorageGroup(mode string) error {
	if mode == "ASYNC" {
		c.err = c.client.AddVolumesToProtectedStorageGroup(context.TODO(), symID, mock.DefaultASYNCProtectedSG, mock.DefaultRemoteSymID, mock.DefaultASYNCProtectedSG, false, c.volIDList...)
	} else {
		c.err = c.client.AddVolumesToProtectedStorageGroup(context.TODO(), symID, mock.DefaultMETROProtectedSG, mock.DefaultRemoteSymID, mock.DefaultMETROProtectedSG, false, c.volIDList...)
	}
	return nil
}

func (c *unitContext) iCallCreateVolumeInProtectedStorageGroupSWithNameAndSize(volumeName string, sizeInCylinders int) error {
	volOpts := make(map[string]interface{})
	volOpts["capacityUnit"] = "CYL"
	volOpts["enableMobility"] = "false"
	c.vol, c.err = c.client.CreateVolumeInProtectedStorageGroupS(context.TODO(), symID, mock.DefaultRemoteSymID, mock.DefaultASYNCProtectedSG, mock.DefaultASYNCProtectedSG, volumeName, sizeInCylinders, volOpts)
	return nil
}

func (c *unitContext) iCallRemoveVolumesFromProtectedStorageGroup() error {
	_, c.err = c.client.RemoveVolumesFromProtectedStorageGroup(context.TODO(), symID, mock.DefaultStorageGroup, mock.DefaultRemoteSymID, mock.DefaultStorageGroup, false, c.volIDList...)
	return nil
}

func (c *unitContext) iCallCreateRDFPair(mode string) error {
	if mode == "ASYNC" {
		_, c.err = c.client.CreateRDFPair(context.TODO(), symID, fmt.Sprintf("%d", mock.DefaultAsyncRDFGNo), c.volIDList[0], mode, "", false, false)
	} else {
		_, c.err = c.client.CreateRDFPair(context.TODO(), symID, fmt.Sprintf("%d", mock.DefaultMetroRDFGNo), c.volIDList[0], mode, "", false, false)
	}
	return nil
}

func (c *unitContext) iCallExecuteAction(action string) error {
	c.err = c.client.ExecuteReplicationActionOnSG(context.TODO(), symID, action, mock.DefaultASYNCProtectedSG, fmt.Sprintf("%d", mock.DefaultAsyncRDFGNo), false, false, true)
	return nil
}

func (c *unitContext) iCallGetFreeLocalAndRemoteRDFg() error {
	_, c.err = c.client.GetFreeLocalAndRemoteRDFg(context.TODO(), mock.DefaultSymmetrixID, mock.DefaultRemoteSymID)
	return nil
}

func (c *unitContext) iCallGetLocalOnlineRDFDirs() error {
	_, c.err = c.client.GetLocalOnlineRDFDirs(context.TODO(), mock.DefaultSymmetrixID)
	return nil
}

func (c *unitContext) iCallGetLocalOnlineRDFPorts() error {
	_, c.err = c.client.GetLocalOnlineRDFPorts(context.TODO(), mock.DefaultRDFDir, mock.DefaultRemoteSymID)
	return nil
}

func (c *unitContext) iCallGetLocalRDFPortDetails() error {
	_, c.err = c.client.GetLocalRDFPortDetails(context.TODO(), mock.DefaultSymmetrixID, mock.DefaultRDFDir, mock.DefaultRDFPort)
	return nil
}

func (c *unitContext) iCallGetRDFGroupListWithQuery(query string) error {
	if query != "" {
		switch query {
		case "remote_symmetrix_id":
			_, c.err = c.client.GetRDFGroupList(context.TODO(), mock.DefaultSymmetrixID, types.QueryParams{query: mock.DefaultRemoteSymID})
		case "volume_count":
			_, c.err = c.client.GetRDFGroupList(context.TODO(), mock.DefaultSymmetrixID, types.QueryParams{query: 1})
		}
	} else {
		_, c.err = c.client.GetRDFGroupList(context.TODO(), mock.DefaultSymmetrixID, nil)
	}
	return nil
}

func (c *unitContext) iCallGetRemoteRDFPortOnSAN() error {
	_, c.err = c.client.GetRemoteRDFPortOnSAN(context.TODO(), mock.DefaultSymmetrixID, mock.DefaultRDFDir, fmt.Sprintf("%d", mock.DefaultRDFPort))
	return nil
}

func (c *unitContext) iCallExecuteCreateRDFGroup() error {
	createRDFgPayload := new(types.RDFGroupCreate)
	createRDFgPayload.LocalPorts = []types.RDFPortDetails{
		{
			SymmID:     mock.DefaultSymmetrixID,
			DirNum:     33,
			DirID:      mock.DefaultRDFDir,
			PortNum:    mock.DefaultRDFPort,
			PortOnline: true,
			PortWWN:    "5000097200007003",
		},
	}
	createRDFgPayload.RemotePorts = []types.RDFPortDetails{
		{
			SymmID:     mock.DefaultRemoteSymID,
			DirNum:     33,
			DirID:      mock.DefaultRDFDir,
			PortNum:    mock.DefaultRDFPort,
			PortOnline: true,
			PortWWN:    "5000097200007003",
		},
	}
	createRDFgPayload.Label = mock.DefaultAsyncRDFLabel
	createRDFgPayload.LocalRDFNum = mock.DefaultAsyncRDFGNo
	createRDFgPayload.RemoteRDFNum = mock.DefaultAsyncRemoteRDFGNo
	c.err = c.client.ExecuteCreateRDFGroup(context.TODO(), mock.DefaultSymmetrixID, createRDFgPayload)
	return nil
}

func (c *unitContext) iCallCreateHostGroupWithFlags(hostGroupID string, setHostFlags string) error {
	hostIDs := make([]string, 1)
	c.hostGroupID = hostGroupID
	hostIDs[0] = testHost
	if setHostFlags == "true" {
		hostFlags := &types.HostFlags{
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
		c.hostGroup, c.err = c.client.CreateHostGroup(context.TODO(), symID, hostGroupID, hostIDs, hostFlags)
	} else {
		c.hostGroup, c.err = c.client.CreateHostGroup(context.TODO(), symID, hostGroupID, hostIDs, nil)
	}

	return nil
}

func (c *unitContext) iCallUpdateHostGroupFlagsWithFlags(hostGroupID string, updateHostFlags string) error {
	hostIDs := make([]string, 1)
	c.hostGroupID = hostGroupID
	hostIDs[0] = testHost
	if updateHostFlags == "true" {
		hostFlags := &types.HostFlags{
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
		c.hostGroup, c.err = c.client.UpdateHostGroupFlags(context.TODO(), symID, hostGroupID, hostFlags)
	} else {
		c.hostGroup, c.err = c.client.UpdateHostGroupFlags(context.TODO(), symID, hostGroupID, nil)
	}

	return nil
}

func (c *unitContext) iCallUpdateHostGroupHostsWithHosts(hostGroupID, hostID string) error {
	hostIDs := make([]string, 1)
	c.hostGroupID = hostGroupID
	hostIDs[0] = hostID
	c.hostGroup, c.err = c.client.UpdateHostGroupHosts(context.TODO(), symID, c.hostGroupID, hostIDs)
	return nil
}

func (c *unitContext) iCallUpdateHostGroupName(newName string) error {
	c.hostGroup, c.err = c.client.UpdateHostGroupName(context.TODO(), symID, c.hostGroupID, newName)
	c.hostGroupID = newName
	return nil
}

func (c *unitContext) iGetAValidHostGroupIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.hostGroup.HostGroupID != c.hostGroupID {
		return fmt.Errorf("Expected to get HostGroup %s, but received %s",
			c.hostGroup.HostGroupID, c.hostGroupID)
	}
	return nil
}

func (c *unitContext) iHaveAValidHostGroup(hostGroupname string) error {
	hostIDs := make([]string, 1)
	c.hostGroupID = hostGroupname
	hostIDs[0] = testHost
	hostFlags := &types.HostFlags{
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
	c.hostGroup, c.err = mock.AddHostGroup(c.hostGroupID, hostIDs, hostFlags)
	return nil
}

func (c *unitContext) iCallGetHostGroupByID(hostGroupID string) error {
	c.hostGroup, c.err = c.client.GetHostGroupByID(context.TODO(), symID, hostGroupID)
	return nil
}

func (c *unitContext) iCallDeleteHostGroup(hostGroupName string) error {
	c.err = c.client.DeleteHostGroup(context.TODO(), symID, hostGroupName)
	return nil
}

func (c *unitContext) iCallGetHostGroupList() error {
	c.hostGroupList, c.err = c.client.GetHostGroupList(context.TODO(), symID)
	return nil
}

func (c *unitContext) iGetAValidHostGroupListIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.hostGroupList == nil || len(c.hostGroupList.HostGroupIDs) == 0 {
		return fmt.Errorf("Expected item in HostGroupList but got none")
	}
	fmt.Println(c.hostGroupList)
	return nil
}

func (c *unitContext) iCallGetStorageGroupMetrics() error {
	var metrics *types.StorageGroupMetricsIterator
	metrics, c.err = c.client.GetStorageGroupMetrics(context.TODO(), symID, mock.DefaultStorageGroup, []string{"HostMBReads"}, 0, 0)
	c.storageGroupMetrics = metrics
	return nil
}

func (c *unitContext) iGetStorageGroupMetrics() error {
	if c.err == nil {
		if c.storageGroupMetrics == nil {
			return fmt.Errorf("StorageGroupMetrics nil")
		}
		if len(c.storageGroupMetrics.ResultList.Result) == 0 {
			return fmt.Errorf("no metric in StorageGroupMetrics")
		}
	}
	return nil
}

func (c *unitContext) iCallGetVolumesMetrics() error {
	var metrics *types.VolumeMetricsIterator
	metrics, c.err = c.client.GetVolumesMetrics(context.TODO(), symID, mock.DefaultStorageGroup, []string{"MBReads"}, 0, 0)
	c.volumesMetrics = metrics
	return nil
}

func (c *unitContext) iCallGetVolumesMetricsByIDFor(volID string) error {
	var metrics *types.VolumeMetricsIterator
	metrics, c.err = c.client.GetVolumesMetricsByID(context.TODO(), symID, volID, []string{"MBReads"}, 0, 0)
	c.volumesMetrics = metrics
	return nil
}

func (c *unitContext) iGetVolumesMetrics() error {
	if c.err == nil {
		if c.volumesMetrics == nil {
			return fmt.Errorf("volumesMetrics nil")
		}
		if len(c.volumesMetrics.ResultList.Result) == 0 {
			return fmt.Errorf("no metric in volumesMetrics")
		}
	}
	return nil
}

func (c *unitContext) iCallGetFileSystemMetricsByIDFor(fileSystemID string) error {
	var metrics *types.FileSystemMetricsIterator
	metrics, c.err = c.client.GetFileSystemMetricsByID(context.TODO(), symID, fileSystemID, []string{"MBReads"}, 0, 0)
	c.fileSystemMetrics = metrics
	return nil
}

func (c *unitContext) iGetFileMetrics() error {
	if c.err == nil {
		if c.fileSystemMetrics == nil {
			return fmt.Errorf("fileMetrics nil")
		}
		if len(c.fileSystemMetrics.ResultList.Result) == 0 {
			return fmt.Errorf("no metric in filemetrics")
		}
	}
	return nil
}

func (c *unitContext) iCallGetStorageGroupPerfKeys() error {
	var perfKeys *types.StorageGroupKeysResult
	perfKeys, c.err = c.client.GetStorageGroupPerfKeys(context.TODO(), symID)
	c.storageGroupPerfKeys = perfKeys
	return nil
}

func (c *unitContext) iGetStorageGroupPerfKeys() error {
	if c.err == nil {
		if c.storageGroupPerfKeys == nil {
			return fmt.Errorf("storage group performance keys nil")
		}
	}
	return nil
}

func (c *unitContext) iCallGetArrayPerfKeys() error {
	var perfKeys *types.ArrayKeysResult
	perfKeys, c.err = c.client.GetArrayPerfKeys(context.TODO())
	c.arrayPerfKeys = perfKeys
	return nil
}

func (c *unitContext) iGetArrayPerfKeys() error {
	if c.err == nil {
		if c.arrayPerfKeys == nil {
			return fmt.Errorf("array performance keys nil")
		}
	}
	return nil
}

func (c *unitContext) iCallGetSnapshotPolicyWith(snapshotPolicyID string) error {
	c.snapshotPolicy, c.err = c.client.GetSnapshotPolicy(context.TODO(), symID, snapshotPolicyID)
	return nil
}

func (c *unitContext) iCallCreateSnapshotPolicyWith(snapshotPolicyID string) error {
	c.snapshotPolicy, c.err = c.client.CreateSnapshotPolicy(context.TODO(), symID, snapshotPolicyID, "1 Hour", 10, 2, 2, nil)
	return nil
}

func (c *unitContext) iCallCreateSnapshotPolicyWithAndPayload(snapshotPolicyID, payload string) error {
	opPayload := make(map[string]interface{})
	if payload == "cloudSnapshotPolicyDetails" {
		opPayload["cloudSnapshotPolicyDetails"] = &types.CloudSnapshotPolicyDetails{
			CloudRetentionDays: 1,
			CloudProviderName:  "emc",
		}
	}
	if payload == "localSnapshotPolicyDetails" {
		opPayload["localSnapshotPolicyDetails"] = &types.LocalSnapshotPolicyDetails{
			Secure:        true,
			SnapshotCount: 1,
		}
	}
	c.snapshotPolicy, c.err = c.client.CreateSnapshotPolicy(context.TODO(), symID, snapshotPolicyID, "1 Hour", 10, 2, 2, opPayload)
	return nil
}

func (c *unitContext) iShouldGetSnapshotPolicytInformationIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.snapshotPolicy == nil {
		return fmt.Errorf("The snapshot policy does not exist")
	}
	return nil
}

func (c *unitContext) iCallModifySnapshotPolicyWithAndAnd(snapshotPolicyID string, action string, updatedName string) error {
	optionalPayload := make(map[string]interface{})
	modifySnapshotPolicyParam := &types.ModifySnapshotPolicyParam{
		SnapshotPolicyName: updatedName,
	}
	optionalPayload["modify"] = modifySnapshotPolicyParam

	c.err = c.client.UpdateSnapshotPolicy(context.TODO(), symID, action, snapshotPolicyID, optionalPayload)
	return nil
}

func (c *unitContext) iShouldModifySnapshotPolicyIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.snapshotPolicy == nil {
		return fmt.Errorf("The snapshot policy could not be modified")
	}
	return nil
}

func (c *unitContext) iCallAddRemoveStorageGrpFromSnapshotPolicyAndAnd(snapshotPolicyID string, action string, sgName string) error {
	optionalPayload := make(map[string]interface{})
	if action == "AssociateToStorageGroups" {
		associateStorageGroupParam := &types.AssociateStorageGroupParam{
			StorageGroupName: []string{sgName},
		}
		optionalPayload["associateStorageGroupParam"] = associateStorageGroupParam
	} else {
		disassociateStorageGroupParam := &types.DisassociateStorageGroupParam{
			StorageGroupName: []string{sgName},
		}
		optionalPayload["disassociateStorageGroupParam"] = disassociateStorageGroupParam
	}

	c.err = c.client.UpdateSnapshotPolicy(context.TODO(), symID, action, snapshotPolicyID, optionalPayload)
	return nil
}

func (c *unitContext) iShoulAddRemoveStorageGrpFromSnapshotPolicyIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.snapshotPolicy == nil {
		return fmt.Errorf("The snapshot policy could not be modified")
	}
	return nil
}

func (c *unitContext) iCallDeleteSnapshotPolicy(snapshotPolicyID string) error {
	c.err = c.client.DeleteSnapshotPolicy(context.TODO(), symID, snapshotPolicyID)
	return nil
}

func (c *unitContext) iCallGetSnapshotPolicyList() error {
	c.snapshotPolicyList, c.err = c.client.GetSnapshotPolicyList(context.TODO(), symID)
	return nil
}

func (c *unitContext) iShouldGetListOfSnapshotPoliciesIfNoError() error {
	if c.err != nil {
		return nil
	}
	if c.snapshotPolicyList == nil {
		return fmt.Errorf("Could not get the list of snapshot policies")
	}
	return nil
}

func (c *unitContext) iCallGetFileSystemList() error {
	c.fileSystemList, c.err = c.client.GetFileSystemList(context.TODO(), symID, nil)
	return nil
}

func (c *unitContext) iCallGetFileSystemListWithParam() error {
	query := types.QueryParams{queryName: mock.DefaultFSName, queryNASServerID: mock.DefaultNASServerID}
	c.fileSystemList, c.err = c.client.GetFileSystemList(context.TODO(), symID, query)
	return nil
}

func (c *unitContext) iGetAValidFileSystemIDListIfNoError() error {
	if c.err == nil {
		if c.fileSystemList == nil {
			return fmt.Errorf("fileSystemIDList nil")
		}
		if c.fileSystemList.Count == 0 {
			return fmt.Errorf("no IDs in fileSystemIDList")
		}
	}
	return nil
}

func (c *unitContext) iCallGetNASServerList() error {
	c.nasServerList, c.err = c.client.GetNASServerList(context.TODO(), symID, nil)
	return nil
}

func (c *unitContext) iCallGetNASServerListWithParam() error {
	query := types.QueryParams{queryName: mock.DefaultNASServerName}
	c.nasServerList, c.err = c.client.GetNASServerList(context.TODO(), symID, query)
	return nil
}

func (c *unitContext) iCallGetNFSExportListWithParam() error {
	query := types.QueryParams{queryName: "nfs-1"}
	c.nfsExportList, c.err = c.client.GetNFSExportList(context.TODO(), symID, query)
	return nil
}

func (c *unitContext) iCallGetNFSExportList() error {
	c.nfsExportList, c.err = c.client.GetNFSExportList(context.TODO(), symID, nil)
	return nil
}

func (c *unitContext) iGetAValidNASServerIDListIfNoError() error {
	if c.err == nil {
		if c.nasServerList == nil {
			return fmt.Errorf("nasServer List nil")
		}
		if len(c.nasServerList.Entries) == 0 {
			return fmt.Errorf("no IDs in nasServerList")
		}
	}
	return nil
}

func (c *unitContext) iGetAValidNFSExportIDListIfNoError() error {
	if c.err == nil {
		if c.nfsExportList == nil {
			return fmt.Errorf("nfsExportList nil")
		}
		if c.nfsExportList.Count == 0 {
			return fmt.Errorf("no IDs in nfsExportList")
		}
	}
	return nil
}

func (c *unitContext) iCallGetFileSystemByID(fsID string) error {
	c.fileSystem, c.err = c.client.GetFileSystemByID(context.TODO(), symID, fsID)
	return nil
}

func (c *unitContext) iGetAValidFileSystemObjectIfNoError() error {
	if c.err == nil {
		if c.fileSystem == nil {
			return fmt.Errorf("fileSystem nil")
		}
	}
	return nil
}

func (c *unitContext) iCallCreateFileSystem(fsName string) error {
	c.fileSystem, c.err = c.client.CreateFileSystem(context.TODO(), symID, fsName, "nas-1", "Diamond", 20000)
	return nil
}

func (c *unitContext) iCallModifyFileSystemOn(fsID string) error {
	payload := types.ModifyFileSystem{
		SizeTotal: 8000,
	}
	c.fileSystem, c.err = c.client.ModifyFileSystem(context.TODO(), symID, fsID, payload)
	return nil
}

func (c *unitContext) iCallDeleteFileSystem() error {
	var fs string
	if c.fileSystem == nil {
		// for disallowed array
		fs = "id3"
	} else {
		fs = c.fileSystem.ID
	}
	c.err = c.client.DeleteFileSystem(context.TODO(), symID, fs)
	return nil
}

func (c *unitContext) iCallDeleteNASServer(nasID string) error {
	c.err = c.client.DeleteNASServer(context.TODO(), symID, nasID)
	return nil
}

func (c *unitContext) iCallGetNASServerByID(nasID string) error {
	c.nasServer, c.err = c.client.GetNASServerByID(context.TODO(), symID, nasID)
	return nil
}

func (c *unitContext) iCallModifyNASServerOn(nasID string) error {
	payload := types.ModifyNASServer{Name: "new-name"}
	c.nasServer, c.err = c.client.ModifyNASServer(context.TODO(), symID, nasID, payload)
	return nil
}

func (c *unitContext) iGetAValidNASServerObjectIfNoError() error {
	if c.err == nil {
		if c.nasServer == nil {
			return fmt.Errorf("nasServer nil")
		}
	}
	return nil
}

func (c *unitContext) iCallCreateNFSExport(nfsName string) error {
	payload := types.CreateNFSExport{
		StorageResource: "id1",
		Path:            "/id1",
		Name:            nfsName,
		DefaultAccess:   "ReadWrite",
	}
	c.nfsExport, c.err = c.client.CreateNFSExport(context.TODO(), symID, payload)
	return nil
}

func (c *unitContext) iCallDeleteNFSExport(nfsID string) error {
	c.err = c.client.DeleteNFSExport(context.TODO(), symID, nfsID)
	return nil
}

func (c *unitContext) iCallGetNFSExportByID(nfsID string) error {
	c.nfsExport, c.err = c.client.GetNFSExportByID(context.TODO(), symID, nfsID)
	return nil
}

func (c *unitContext) iCallModifyNFSExportOn(nfsID string) error {
	payload := types.ModifyNFSExport{
		Name: "updated-name",
	}
	c.nfsExport, c.err = c.client.ModifyNFSExport(context.TODO(), symID, nfsID, payload)
	return nil
}

func (c *unitContext) iGetAValidNFSExportObjectIfNoError() error {
	if c.err == nil {
		if c.nfsExport == nil {
			return fmt.Errorf("nfs export nil")
		}
	}
	return nil
}

func (c *unitContext) iCallGetFileInterfaceByID(interfaceID string) error {
	c.fileInterface, c.err = c.client.GetFileInterfaceByID(context.TODO(), symID, interfaceID)
	return nil
}

func (c *unitContext) iGetAValidFileInterfaceObjectIfNoError() error {
	if c.err == nil {
		if c.fileInterface == nil {
			return fmt.Errorf("file interface nil")
		}
	}
	return nil
}

func (c *unitContext) iCallGetDirectorIDList() error {
	c.directorIDList, c.err = c.client.GetDirectorIDList(context.TODO(), symID)
	return nil
}

func UnitTestContext(s *godog.ScenarioContext) {
	c := &unitContext{}
	s.Step(`^I induce error "([^"]*)"$`, c.iInduceError)
	s.Step(`^I call authenticate with endpoint "([^"]*)" credentials "([^"]*)" apiversion "([^"]*)"$`, c.iCallAuthenticateWithEndpointCredentials)
	s.Step(`^the error message contains "([^"]*)"$`, c.theErrorMessageContains)
	s.Step(`^a valid connection$`, c.aValidConnection)
	s.Step(`^a valid v(\d+) connection$`, c.aValidv91Connection)
	s.Step(`^I call GetSymmetrixIDList$`, c.iCallGetSymmetrixIDList)
	s.Step(`^I get a valid Symmetrix ID List if no error$`, c.iGetAValidSymmetrixIDListIfNoError)
	s.Step(`^I call GetSymmetrixByID "([^"]*)"$`, c.iCallGetSymmetrixByID)
	s.Step(`^I get a valid Symmetrix Object if no error$`, c.iGetAValidSymmetrixObjectIfNoError)
	s.Step(`^I have (\d+) volumes$`, c.iHaveVolumes)
	s.Step(`^I call GetVolumeByID "([^"]*)"$`, c.iCallGetVolumeByID)
	s.Step(`^I get a valid Volume Object "([^"]*)" if no error$`, c.iGetAValidVolumeObjectIfNoError)
	s.Step(`^I call GetVolumeIDList "([^"]*)"$`, c.iCallGetVolumeIDList)
	s.Step(`^I get a valid VolumeIDList with (\d+) if no error$`, c.iGetAValidVolumeIDListWithIfNoError)
	s.Step(`^I call GetVolumeIDListWithParams`, c.iCallGetVolumeIDListWithParams)
	s.Step(`^I call GetStorageGroupIDList$`, c.iCallGetStorageGroupIDList)
	s.Step(`^I get a valid StorageGroupIDList if no errors$`, c.iGetAValidStorageGroupIDListIfNoErrors)
	s.Step(`^I call GetStorageGroup "([^"]*)"$`, c.iCallGetStorageGroup)
	s.Step(`^I call GetStorageGroupSnapshotPolicy with "([^"]*)" "([^"]*)" "([^"]*)"$`, c.iCallGetStorageGroupSnapshotPolicy)
	s.Step(`^I have a StorageGroup "([^"]*)"$`, c.iHaveAStorageGroup)
	s.Step(`^I get a valid StorageGroup if no errors$`, c.iGetAValidStorageGroupIfNoErrors)
	s.Step(`^I get a valid StorageGroupSnapshotPolicy Object if no error$`, c.iGetAValidStorageGroupSnapshotPolicyObjectIfNoError)
	s.Step(`^I have (\d+) jobs$`, c.iHaveJobs)
	s.Step(`^I call GetJobIDList with "([^"]*)"$`, c.iCallGetJobIDListWith)
	s.Step(`^I get a valid JobsIDList with (\d+) if no errors$`, c.iGetAValidJobsIDListWithIfNoErrors)
	s.Step(`^I create a job with initial state "([^"]*)" and final state "([^"]*)"$`, c.iCreateAJobWithInitialStateAndFinalState)
	s.Step(`^I call GetJobByID$`, c.iCallGetJobByID)
	s.Step(`^I get a valid Job with state "([^"]*)" if no error$`, c.iGetAValidJobWithStateIfNoError)
	s.Step(`^I call WaitOnJobCompletion$`, c.iCallWaitOnJobCompletion)
	s.Step(`^I call RefreshSymmetrix "([^"]*)"$`, c.iCallRefreshSymmetrix)
	// Volumes
	s.Step(`^I call CreateVolumeInStorageGroup with name "([^"]*)" and size (\d+)$`, c.iCallCreateVolumeInStorageGroupWithNameAndSize)
	s.Step(`^I call CreateVolumeInStorageGroup with name "([^"]*)" and size (\d+) and unit "([^"]*)"$`, c.iCallCreateVolumeInStorageGroupWithNameAndSizeAndUnit)
	s.Step(`^I call CreateVolumeInStorageGroupS with name "([^"]*)" and size (\d+)$`, c.iCallCreateVolumeInStorageGroupSWithNameAndSize)
	s.Step(`^I call CreateVolumeInStorageGroupS with name "([^"]*)" and size (\d+) and unit "([^"]*)"$`, c.iCallCreateVolumeInStorageGroupSWithNameAndSizeAndUnit)
	s.Step(`^I call CreateVolumeInStorageGroupSWithMetaDataHeaders with name "([^"]*)" and size (\d+)$`, c.iCallCreateVolumeInStorageGroupSWithNameAndSizeWithMetaDataHeaders)
	s.Step(`^I get a valid Volume with name "([^"]*)" if no error$`, c.iGetAValidVolumeWithNameIfNoError)
	s.Step(`^I validate that volume has mobility modified to "([^"]*)"$`, c.iGetAValidVolumeWithMobilityModified)
	s.Step(`^I call CreateStorageGroup with name "([^"]*)" and srp "([^"]*)" and sl "([^"]*)" and hostlimits "([^"]*)"$`, c.iCallCreateStorageGroupWithNameAndSrpAndSlAndHostLimits)
	s.Step(`^I call CreateStorageGroup with name "([^"]*)" and srp "([^"]*)" and sl "([^"]*)"$`, c.iCallCreateStorageGroupWithNameAndSrpAndSl)
	s.Step(`^I call DeleteStorageGroup "([^"]*)"$`, c.iCallDeleteStorageGroup)
	s.Step(`^I get a valid StorageGroup with name "([^"]*)" if no error$`, c.iGetAValidStorageGroupWithNameIfNoError)
	s.Step(`^I call GetStoragePoolList$`, c.iCallGetStoragePoolList)
	s.Step(`^I get a valid StoragePoolList if no error$`, c.iGetAValidStoragePoolListIfNoError)
	s.Step(`^I call RemoveVolumeFromStorageGroup$`, c.iCallRemoveVolumeFromStorageGroup)
	s.Step(`^the volume is no longer a member of the Storage Group if no error$`, c.theVolumeIsNoLongerAMemberOfTheStorageGroupIfNoError)
	s.Step(`^I call RenameVolume with "([^"]*)"$`, c.iCallRenameVolumeWith)
	s.Step(`^I call ModifyMobility for Volume with id "([^"]*)" to mobility "([^"]*)"$`, c.iCallModifyMobilityForVolume)
	s.Step(`^I call InitiateDeallocationOfTracksFromVolume$`, c.iCallInitiateDeallocationOfTracksFromVolume)
	s.Step(`^I call DeleteVolume$`, c.iCallDeleteVolume)
	s.Step(`^I expand volume "([^"]*)" to "([^"]*)" in GB$`, c.iExpandVolumeToSize)
	s.Step(`^I validate that volume "([^"]*)" has has size "([^"]*)" in GB$`, c.iValidateVolumeSize)
	s.Step(`^I expand volume "([^"]*)" to "([^"]*)" in "([^"]*)"$`, c.iExpandVolumeToSizeWithUnit)
	// Masking View
	s.Step(`^I have a MaskingView "([^"]*)"$`, c.iHaveAMaskingView)
	s.Step(`^I call GetMaskingViewList$`, c.iCallGetMaskingViewList)
	s.Step(`^I get a valid MaskingViewList if no error$`, c.iGetAValidMaskingViewListIfNoError)
	s.Step(`^I call GetMaskingViewByID "([^"]*)"$`, c.iCallGetMaskingViewByID)
	s.Step(`^I get a valid MaskingView if no error$`, c.iGetAValidMaskingViewIfNoError)
	s.Step(`^I call RenameMaskingView with "([^"]*)"$`, c.iCallRenameMaskingViewWith)
	s.Step(`^I call CreateMaskingViewWithHost "([^"]*)"$`, c.iCallCreateMaskingViewWithHost)
	s.Step(`^I call CreateMaskingViewWithHostGroup "([^"]*)"$`, c.iCallCreateMaskingViewWithHostGroup)
	s.Step(`^I call DeleteMaskingView$`, c.iCallDeleteMaskingView)
	// Port Group
	s.Step(`^I have a PortGroup$`, c.iHaveAPortGroup)
	s.Step(`^I call GetPortGroupList$`, c.iCallGetPortGroupList)
	s.Step(`^I get a valid PortGroupList if no error$`, c.iGetAValidPortGroupListIfNoError)
	s.Step(`^I call GetPortGroupByID$`, c.iCallGetPortGroupByID)
	s.Step(`^I get a valid PortGroup if no error$`, c.iGetAValidPortGroupIfNoError)
	s.Step(`^I get PortGroup "([^"]*)" if no error$`, c.iGetPortGroupIfNoError)
	s.Step(`^I call CreatePortGroup "([^"]*)" with ports "([^"]*)"$`, c.iCallCreatePortGroup)
	s.Step(`^I call RenamePortGroup with "([^"]*)"$`, c.iCallRenamePortGroupWith)
	s.Step(`^I call UpdatePortGroup "([^"]*)" with ports "([^"]*)"$`, c.iCallUpdatePortGroup)
	s.Step(`^I call DeletePortGroup "([^"]*)"$`, c.iCallDeletePortGroup)
	s.Step(`^I expect PortGroup to have these ports "([^"]*)"$`, c.iExpectedThesePortsInPortGroup)
	s.Step(`^the PortGroup "([^"]*)" should not exist`, c.thePortGroupShouldNotExist)
	// Host
	s.Step(`^I have a FC Host "([^"]*)"$`, c.iHaveAFCHost)
	s.Step(`^I have a ISCSI Host "([^"]*)"$`, c.iHaveAISCSIHost)
	s.Step(`^I have a NVMeTCP Host "([^"]*)"$`, c.iHaveANVMETCPHost)
	s.Step(`^I call GetHostList$`, c.iCallGetHostList)
	s.Step(`^I get a valid HostList if no error$`, c.iGetAValidHostListIfNoError)
	s.Step(`^I call GetHostByID "([^"]*)"$`, c.iCallGetHostByID)
	s.Step(`^I get a valid Host if no error$`, c.iGetAValidHostIfNoError)
	// Initiator
	s.Step(`^I have a Initiator$`, c.iHaveAInitiator)
	s.Step(`^I call GetInitiatorList$`, c.iCallGetInitiatorList)
	s.Step(`^I call GetInitiatorList with filters$`, c.iCallGetInitiatorListWithFilters)
	s.Step(`^I get a valid InitiatorList if no error$`, c.iGetAValidInitiatorListIfNoError)
	s.Step(`^I call GetInitiatorByID$`, c.iCallGetInitiatorByID)
	s.Step(`^I get a valid Initiator if no error$`, c.iGetAValidInitiatorIfNoError)
	// HostGroup/Host
	s.Step(`^I call CreateHostGroup "([^"]*)" with flags "([^"]*)"$`, c.iCallCreateHostGroupWithFlags)
	s.Step(`^I get a valid HostGroup if no error$`, c.iGetAValidHostGroupIfNoError)
	s.Step(`^I have a HostGroup "([^"]*)"$`, c.iHaveAHostGroup)
	s.Step(`^I have a valid HostGroup "([^"]*)"$`, c.iHaveAValidHostGroup)
	s.Step(`^I call GetHostGroupByID "([^"]*)"$`, c.iCallGetHostGroupByID)
	s.Step(`^I call UpdateHostGroupFlags "([^"]*)" with flags "([^"]*)"$`, c.iCallUpdateHostGroupFlagsWithFlags)
	s.Step(`^I call UpdateHostGroupHosts "([^"]*)" with hosts "([^"]*)"$`, c.iCallUpdateHostGroupHostsWithHosts)
	s.Step(`^I call UpdateHostGroupName "([^"]*)"$`, c.iCallUpdateHostGroupName)
	s.Step(`^I call DeleteHostGroup "([^"]*)"$`, c.iCallDeleteHostGroup)
	s.Step(`^I call GetHostGroupList$`, c.iCallGetHostGroupList)
	s.Step(`^I get a valid HostGroupList if no error$`, c.iGetAValidHostGroupListIfNoError)
	s.Step(`^I call CreateHost "([^"]*)"$`, c.iCallCreateHost)
	s.Step(`^I call DeleteHost "([^"]*)"$`, c.iCallDeleteHost)
	s.Step(`^I call AddVolumesToStorageGroup "([^"]*)"$`, c.iCallAddVolumesToStorageGroup)
	s.Step(`^I call AddVolumesToStorageGroupS "([^"]*)"$`, c.iCallAddVolumesToStorageGroupS)
	s.Step(`^then the Volumes are part of StorageGroup if no error$`, c.thenTheVolumesArePartOfStorageGroupIfNoError)
	s.Step(`^I call UpdateHost$`, c.iCallUpdateHost)
	s.Step(`^I call UpdateHostFlags$`, c.iCallUpdateHostFlags)
	s.Step(`^I call GetVolumeIDListInStorageGroup "([^"]*)"$`, c.iCallGetVolumeIDListInStorageGroup)
	// GetListOftargetAddresses
	s.Step(`^I call GetListOfTargetAddresses$`, c.iCallGetListOfTargetAddresses)
	s.Step(`^I recieve (\d+) IP addresses$`, c.iRecieveIPAddresses)
	s.Step(`^I call GetStoragePool "([^"]*)"$`, c.iCallGetStoragePool)
	s.Step(`^I get a valid GetStoragePool if no errors$`, c.iGetAValidGetStoragePoolIfNoErrors)
	// Allowed List of arrays
	s.Step(`^I have an allowed list of "([^"]*)"$`, c.iHaveAnAllowedListOf)
	s.Step(`^it contains (\d+) arrays$`, c.itContainsArrays)
	s.Step(`^should include "([^"]*)"$`, c.shouldInclude)
	s.Step(`^should not include "([^"]*)"$`, c.shouldNotInclude)

	s.Step(`^I get a valid Symmetrix ID List that contains "([^"]*)" and does not contains "([^"]*)"$`, c.iGetAValidSymmetrixIDListThatContainsAndDoesNotContains)

	// SG Snapshot
	s.Step(`^I call GetStorageGroupSnapshots with "([^"]*)"$`, c.iCallGetStorageGroupSnapshotsWith)
	s.Step(`^I call GetStorageGroupSnapshots with "([^"]*)" and param "([^"]*)"$`, c.iCallGetStorageGroupSnapshotsWithAndParam)
	s.Step(`^I should get storage group snapshot information if no error$`, c.iShouldGetStorageGroupSnapshotInformationIfNoError)
	s.Step(`^I call CreateStorageGroupSnapshot with "([^"]*)"$`, c.iCallCreateStorageGroupSnapshotWith)
	s.Step(`^I call GetStorageGroupSnapshotSnapIDs with "([^"]*)" and "([^"]*)"$`, c.iCallGetStorageGroupSnapshotSnapIDsWithAnd)
	s.Step(`^I should get storage group snapshot snap ids if no error$`, c.iShouldGetStorageGroupSnapshotSnapIDsIfNoError)
	s.Step(`^I call GetStorageGroupSnapshotSnap with "([^"]*)" and "([^"]*)" and "([^"]*)"$`, c.iCallGetStorageGroupSnapshotSnapWithAndAnd)
	s.Step(`^I should get storage group snapshot snap detail information if no error$`, c.iShouldGetStorageGroupSnapshotSnapDetailInformationIfNoError)
	s.Step(`^I call ModifyStorageGroupSnapshot with "([^"]*)" and "([^"]*)" and "([^"]*)" and action "([^"]*)"$`, c.iCallModifyStorageGroupSnapshotWithAndAndAndAction)
	s.Step(`^I call DeleteStorageGroupSnapshot with "([^"]*)" and "([^"]*)" and "([^"]*)"$`, c.iCallDeleteStorageGroupSnapshotWithAndAnd)
	s.Step(`^I should modify storage group snapshot snap if no error$`, c.iShouldModifyStorageGroupSnapshotSnapIfNoError)

	// Snapshot
	s.Step(`^I excute the capabilities on the symmetrix array$`, c.iExcuteTheCapabilitiesOnTheSymmetrixArray)
	s.Step(`^I call GetSnapVolumeList with "([^"]*)" and "([^"]*)"$`, c.iCallGetSnapVolumeListWithAnd)
	s.Step(`^I should get a list of volumes having snapshots if no error$`, c.iShouldGetListOfVolumesHavingSnapshots)
	s.Step(`^I call GetVolumeSnapInfo with volume "([^"]*)"$`, c.iCallGetVolumeSnapInfoWithVolume)
	s.Step(`^I should get a list of snapshots if no error$`, c.iShouldGetAListOfSnapshotsIfNoError)
	s.Step(`^I call CreateSnapshot with "([^"]*)" and snapshot "([^"]*)" on it$`, c.iCallCreateSnapshotWithAndSnapshotOnIt)
	s.Step(`^I get a valid Snapshot object if no error$`, c.iGetAValidSnapshotObjectIfNoError)
	s.Step(`^I call GetSnapshotInfo with "([^"]*)" and snapshot "([^"]*)" on it$`, c.iCallGetSnapshotInfoWithAndSnapshotNameOnIt)
	s.Step(`^I should get the snapshot details if no error$`, c.iShouldGetTheSnapshotDetailsIfNoError)
	s.Step(`^I call GetSnapshotGenerations with "([^"]*)" and snapshot "([^"]*)" on it$`, c.iCallGetSnapshotGenerationsWithAndSnapshotOnIt)
	s.Step(`^I should get the generation list if no error$`, c.iShouldGetTheGenerationListIfNoError)
	s.Step(`^I call GetSnapshotGeneration with "([^"]*)", snapshot "([^"]*)" and (\d+) on it$`, c.iCallGetSnapshotGenerationWithSnapshotAndOnIt)
	s.Step(`^I should get a generation Info if no error$`, c.iShouldGetAGenerationInfoIfNoError)
	s.Step(`^I call ModifySnapshot with "([^"]*)", "([^"]*)", "([^"]*)", "([^"]*)", (\d+) and "([^"]*)"$`, c.iCallModifySnapshotWithAnd)
	s.Step(`^I call ModifySnapshotS with "([^"]*)", "([^"]*)", "([^"]*)", "([^"]*)", (\d+) and "([^"]*)"$`, c.iCallModifySnapshotSWithAnd)
	s.Step(`^I should get a valid response if no error$`, c.iShouldGetAValidResponseIfNoError)
	s.Step(`^I call DeleteSnapshot with "([^"]*)", snapshot "([^"]*)" and (\d+)  on it$`, c.iCallDeleteSnapshotWithSnapshotAndOnIt)
	s.Step(`^I call DeleteSnapshotS with "([^"]*)", snapshot "([^"]*)" and (\d+)  on it$`, c.iCallDeleteSnapshotSWithSnapshotAndOnIt)
	s.Step(`^I call GetPrivVolumeByID with "([^"]*)"$`, c.iCallGetPrivVolumeByIDWith)
	s.Step(`^I should get a private volume information if no error$`, c.iShouldGetAPrivateVolumeInformationIfNoError)
	s.Step(`^I call GetISCSITargets$`, c.iCallGetISCSITargets)
	s.Step(`^I call GetNVMeTCPTargets$`, c.iCallGetNVMeTCPTargets)
	s.Step(`^I recieve (\d+) targets$`, c.iRecieveTargets)
	s.Step(`^there should be no errors$`, c.thereShouldBeNoErrors)
	s.Step(`^I call UpdateHostName "([^"]*)"$`, c.iCallUpdateHostName)

	// SRDF
	s.Step(`^I call CreateSGReplica with "([^"]*)"$`, c.iCallCreateSGReplica)
	s.Step(`^then SG should be replicated$`, c.thenSGShouldBeReplicated)
	s.Step(`^I call GetStorageGroupRDFInfo$`, c.iCallGetStorageGroupRDFInfo)
	s.Step(`^I call GetRDFDevicePairInfo$`, c.iCallGetRDFDevicePairInfo)
	s.Step(`^I call GetProtectedStorageGroup$`, c.iCallGetProtectedStorageGroup)
	s.Step(`^I call GetRDFGroup$`, c.iCallGetRDFGroup)
	s.Step(`^I call AddVolumesToProtectedStorageGroup with "([^"]*)"$`, c.iCallAddVolumesToProtectedStorageGroup)
	s.Step(`^I call CreateVolumeInProtectedStorageGroupS with name "([^"]*)" and size (\d+)$`, c.iCallCreateVolumeInProtectedStorageGroupSWithNameAndSize)
	s.Step(`^the volumes should "([^"]*)" be replicated$`, c.theVolumesShouldBeReplicated)
	s.Step(`^I call RemoveVolumesFromProtectedStorageGroup$`, c.iCallRemoveVolumesFromProtectedStorageGroup)
	s.Step(`^I call CreateRDFPair with "([^"]*)"$`, c.iCallCreateRDFPair)
	s.Step(`^I call ExecuteAction "([^"]*)"$`, c.iCallExecuteAction)

	// Performance Metrics
	s.Step(`^I call GetStorageGroupMetrics$`, c.iCallGetStorageGroupMetrics)
	s.Step(`^I get StorageGroupMetrics$`, c.iGetStorageGroupMetrics)
	s.Step(`^I call GetVolumesMetrics$`, c.iCallGetVolumesMetrics)
	s.Step(`^I get VolumesMetrics$`, c.iGetVolumesMetrics)
	s.Step(`^I call GetFileSystemMetricsByID for "([^"]*)"$`, c.iCallGetFileSystemMetricsByIDFor)
	s.Step(`^I call GetVolumesMetricsByID for "([^"]*)"$`, c.iCallGetVolumesMetricsByIDFor)
	s.Step(`^I get FileMetrics$`, c.iGetFileMetrics)

	// Performance keys
	s.Step(`^I call GetStorageGroupPerfKeys$`, c.iCallGetStorageGroupPerfKeys)
	s.Step(`^I get StorageGroupPerfKeys$`, c.iGetStorageGroupPerfKeys)
	s.Step(`^I call GetArrayPerfKeys$`, c.iCallGetArrayPerfKeys)
	s.Step(`^I get ArrayPerfKeys$`, c.iGetArrayPerfKeys)
	s.Step(`^I call GetFreeLocalAndRemoteRDFg$`, c.iCallGetFreeLocalAndRemoteRDFg)
	s.Step(`^I call ExecuteCreateRDFGroup$`, c.iCallExecuteCreateRDFGroup)
	s.Step(`^I call GetLocalOnlineRDFDirs$`, c.iCallGetLocalOnlineRDFDirs)
	s.Step(`^I call GetLocalOnlineRDFPorts$`, c.iCallGetLocalOnlineRDFPorts)
	s.Step(`^I call GetLocalRDFPortDetails$`, c.iCallGetLocalRDFPortDetails)
	s.Step(`^^I call GetRDFGroupList with query "([^"]*)"$`, c.iCallGetRDFGroupListWithQuery)
	s.Step(`^I call GetRemoteRDFPortOnSAN$`, c.iCallGetRemoteRDFPortOnSAN)

	// Snapshot Policy
	s.Step(`^I call CreateSnapshotPolicy with "([^"]*)"$`, c.iCallCreateSnapshotPolicyWith)
	s.Step(`^I call CreateSnapshotPolicy with "([^"]*)" and payload "([^"]*)"$`, c.iCallCreateSnapshotPolicyWithAndPayload)
	s.Step(`^I call GetSnapshotPolicy with "([^"]*)"$`, c.iCallGetSnapshotPolicyWith)
	s.Step(`^I should get snapshot policy information if no error$`, c.iShouldGetSnapshotPolicytInformationIfNoError)
	s.Step(`^I call ModifySnapshotPolicy with  "([^"]*)" and "([^"]*)" and "([^"]*)"$`, c.iCallModifySnapshotPolicyWithAndAnd)
	s.Step(`^I should modify snapshot policy if no error$`, c.iShouldModifySnapshotPolicyIfNoError)
	s.Step(`^I call AddRemoveStorageGrpFromSnapshotPolicy with  "([^"]*)" and "([^"]*)" and "([^"]*)"$`, c.iCallAddRemoveStorageGrpFromSnapshotPolicyAndAnd)
	s.Step(`^I call DeleteSnapshotPolicy "([^"]*)"$`, c.iCallDeleteSnapshotPolicy)
	s.Step(`^I call GetSnapshotPolicyList`, c.iCallGetSnapshotPolicyList)
	s.Step(`^I should get list of snapshot policies  if no error$`, c.iShouldGetListOfSnapshotPoliciesIfNoError)

	// File APIs
	s.Step(`^I call GetFileSystemList$`, c.iCallGetFileSystemList)
	s.Step(`^I get a valid FileSystem ID List if no error$`, c.iGetAValidFileSystemIDListIfNoError)
	s.Step(`^I call GetNASServerList$`, c.iCallGetNASServerList)
	s.Step(`^I call GetNFSExportList$`, c.iCallGetNFSExportList)
	s.Step(`^I get a valid NAS Server ID List if no error$`, c.iGetAValidNASServerIDListIfNoError)
	s.Step(`^I get a valid NFS Export ID List if no error$`, c.iGetAValidNFSExportIDListIfNoError)
	s.Step(`^I call GetFileSystemByID "([^"]*)"$`, c.iCallGetFileSystemByID)
	s.Step(`^I call CreateFileSystem "([^"]*)"$`, c.iCallCreateFileSystem)
	s.Step(`^I get a valid fileSystem Object if no error$`, c.iGetAValidFileSystemObjectIfNoError)
	s.Step(`^I call ModifyFileSystem on "([^"]*)"$`, c.iCallModifyFileSystemOn)
	s.Step(`^I call DeleteFileSystem$`, c.iCallDeleteFileSystem)
	s.Step(`^I call DeleteNASServer "([^"]*)"$`, c.iCallDeleteNASServer)
	s.Step(`^I call GetNASServerByID "([^"]*)"$`, c.iCallGetNASServerByID)
	s.Step(`^I call ModifyNASServer on "([^"]*)"$`, c.iCallModifyNASServerOn)
	s.Step(`^I get a valid nasServer Object if no error$`, c.iGetAValidNASServerObjectIfNoError)
	s.Step(`^I call CreateNFSExport "([^"]*)"$`, c.iCallCreateNFSExport)
	s.Step(`^I call DeleteNFSExport "([^"]*)"$`, c.iCallDeleteNFSExport)
	s.Step(`^I call GetNFSExportByID "([^"]*)"$`, c.iCallGetNFSExportByID)
	s.Step(`^I call ModifyNFSExport on "([^"]*)"$`, c.iCallModifyNFSExportOn)
	s.Step(`^I get a valid NFSExport object if no error$`, c.iGetAValidNFSExportObjectIfNoError)
	s.Step(`^I call GetFileSystemListWithParam$`, c.iCallGetFileSystemListWithParam)
	s.Step(`^I call GetNASServerListWithParam$`, c.iCallGetNASServerListWithParam)
	s.Step(`^I call GetNFSExportListWithParam`, c.iCallGetNFSExportListWithParam)
	s.Step(`^I call GetFileInterfaceByID "([^"]*)"$`, c.iCallGetFileInterfaceByID)
	s.Step(`^I get a valid fileInterface Object if no error$`, c.iGetAValidFileInterfaceObjectIfNoError)
}
