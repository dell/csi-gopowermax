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

package mock

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	types "github.com/dell/gopowermax/v2/types/v100"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
)

// constants
const (
	APIVersion               = "{apiversion}"
	PREFIX                   = "/univmax/restapi/" + APIVersion
	PREFIXNOVERSION          = "/univmax/restapi"
	PRIVATEPREFIX            = "/univmax/restapi/private/" + APIVersion
	INTERNALPREFIX           = "/univmax/restapi/internal/100"
	defaultUsername          = "username"
	defaultPassword          = "password"
	Debug                    = false
	DefaultStorageGroup      = "CSI-Test-SG-1"
	DefaultStorageGroup1     = "CSI-Test-SG-2"
	DefaultASYNCProtectedSG  = "csi-rep-sg-ns-test"
	DefaultMETROProtectedSG  = "csi-rep-sg-ns-test"
	DefaultSymmetrixID       = "000197900046"
	DefaultRemoteSymID       = "000000000013"
	DefaultRDFDir            = "OR-1C"
	DefaultRDFPort           = 3
	PostELMSRSymmetrixID     = "000197900047"
	DefaultStoragePool       = "SRP_1"
	DefaultServiceLevel      = "Optimized"
	DefaultFcStoragePortWWN  = "5000000000000001"
	DefaultAsyncRDFGNo       = 13
	DefaultAsyncRemoteRDFGNo = 13
	DefaultAsyncRDFLabel     = "csi-mock-async"
	DefaultMetroRDFGNo       = 14
	DefaultRemoteRDFGNo      = 14
	DefaultMetroRDFLabel     = "csi-mock-metro"
	RemoteArrayHeaderKey     = "RemoteArray"
	RemoteArrayHeaderValue   = "true"
	DefaultNASServerID       = "64xxx7a6-03b5-xxx-xxx-0zzzz8200209"
	DefaultNASServerName     = "nas-1"
	DefaultFSID              = "64xxx7a6-03b5-xxx-xxx-0zzzz8200208"
	DefaultFSName            = "fs-ds-1"
)

const (
	_ = 1 << (10 * iota)
	// KiB ...
	KiB
	// MiB ...
	MiB
	// GiB ...
	GiB
	// TiB ...
	TiB
	// PiB ...
	PiB
)

var mockCacheMutex sync.Mutex

// Data are internal tables the Mock Unisphere uses to provide functionality.
var Data struct {
	VolumeIDToIdentifier          map[string]string
	VolumeIDToSize                map[string]int
	VolumeIDIteratorList          []string
	VolumeIDToSGList              map[string][]string
	MaskingViewIDToHostID         map[string]string
	MaskingViewIDToHostGroupID    map[string]string
	MaskingViewIDToPortGroupID    map[string]string
	MaskingViewIDToStorageGroupID map[string]string
	StorageGroupIDToMaskingViewID map[string]string
	JobIDToMockJob                map[string]*JobInfo
	StorageGroupIDToNVolumes      map[string]int
	StorageGroupIDToStorageGroup  map[string]*types.StorageGroup
	StorageGroupIDToVolumes       map[string][]string
	MaskingViewIDToMaskingView    map[string]*types.MaskingView
	InitiatorIDToInitiator        map[string]*types.Initiator
	HostIDToHost                  map[string]*types.Host
	PortGroupIDToPortGroup        map[string]*types.PortGroup
	PortIDToSymmetrixPortType     map[string]*types.SymmetrixPortType
	VolumeIDToVolume              map[string]*types.Volume
	HostGroupIDToHostGroup        map[string]*types.HostGroup
	JSONDir                       string
	InitiatorHost                 string

	// Snapshots
	VolIDToSnapshots  map[string]map[string]*types.Snapshot
	SnapIDToLinkedVol map[string]map[string]*types.LinkedVolumes

	// SRDF
	StorageGroupIDToRDFStorageGroup map[string]*types.RDFStorageGroup
	AsyncRDFGroup                   *types.RDFGroup
	MetroRDFGroup                   *types.RDFGroup
	AsyncSGRDFInfo                  *types.SGRDFInfo
	MetroSGRDFInfo                  *types.SGRDFInfo

	// File
	FileSysIDToFileSystem    map[string]*types.FileSystem
	NFSExportIDToNFSExport   map[string]*types.NFSExport
	NASServerIDToNASServer   map[string]*types.NASServer
	FileIntIDtoFileInterface map[string]*types.FileInterface
}

var InducedErrors = new(inducedErrors)

// InducedErrors constants
type inducedErrors struct {
	NoConnection                           bool
	InvalidJSON                            bool
	BadHTTPStatus                          int
	GetSymmetrixError                      bool
	GetVolumeIteratorError                 bool
	GetVolumeError                         bool
	UpdateVolumeError                      bool
	DeleteVolumeError                      bool
	DeviceInSGError                        bool
	GetStorageGroupError                   bool
	GetStorageGroupSnapshotPolicyError     bool
	InvalidResponse                        bool
	GetStoragePoolError                    bool
	UpdateStorageGroupError                bool
	GetJobError                            bool
	JobFailedError                         bool
	VolumeNotCreatedError                  bool
	GetJobCannotFindRoleForUser            bool
	CreateStorageGroupError                bool
	StorageGroupAlreadyExists              bool
	DeleteStorageGroupError                bool
	GetStoragePoolListError                bool
	GetPortGroupError                      bool
	GetPortError                           bool
	GetSpecificPortError                   bool
	GetPortISCSITargetError                bool
	GetPortNVMeTCPTargetError              bool
	GetPortGigEError                       bool
	GetDirectorError                       bool
	GetInitiatorError                      bool
	GetInitiatorByIDError                  bool
	GetHostError                           bool
	CreateHostError                        bool
	DeleteHostError                        bool
	UpdateHostError                        bool
	GetMaskingViewError                    bool
	CreateMaskingViewError                 bool
	UpdateMaskingViewError                 bool
	MaskingViewAlreadyExists               bool
	DeleteMaskingViewError                 bool
	PortGroupNotFoundError                 bool
	InitiatorGroupNotFoundError            bool
	StorageGroupNotFoundError              bool
	VolumeNotAddedError                    bool
	GetMaskingViewConnectionsError         bool
	ResetAfterFirstError                   bool
	CreateSnapshotError                    bool
	DeleteSnapshotError                    bool
	LinkSnapshotError                      bool
	RenameSnapshotError                    bool
	GetSymVolumeError                      bool
	GetVolSnapsError                       bool
	GetGenerationError                     bool
	GetPrivateVolumeIterator               bool
	SnapshotNotLicensed                    bool
	UnisphereMismatchError                 bool
	TargetNotDefinedError                  bool
	SnapshotExpired                        bool
	InvalidSnapshotName                    bool
	GetPrivVolumeByIDError                 bool
	CreatePortGroupError                   bool
	UpdatePortGroupError                   bool
	DeletePortGroupError                   bool
	ExpandVolumeError                      bool
	MaxSnapSessionError                    bool
	GetSRDFInfoError                       bool
	VolumeRdfTypesError                    bool
	GetSRDFPairInfoError                   bool
	GetProtectedStorageGroupError          bool
	CreateSGReplicaError                   bool
	GetRDFGroupError                       bool
	GetSGOnRemote                          bool
	GetSGWithVolOnRemote                   bool
	RDFGroupHasPairError                   bool
	GetRemoteVolumeError                   bool
	InvalidLocalVolumeError                bool
	InvalidRemoteVolumeError               bool
	FetchResponseError                     bool
	RemoveVolumesFromSG                    bool
	ModifyMobilityError                    bool
	GetHostGroupError                      bool
	CreateHostGroupError                   bool
	DeleteHostGroupError                   bool
	UpdateHostGroupError                   bool
	GetHostGroupListError                  bool
	GetStorageGroupMetricsError            bool
	GetVolumesMetricsError                 bool
	GetFileSysMetricsError                 bool
	GetStorageGroupPerfKeyError            bool
	GetArrayPerfKeyError                   bool
	GetFreeRDFGError                       bool
	GetLocalOnlineRDFDirsError             bool
	GetRemoteRDFPortOnSANError             bool
	GetLocalOnlineRDFPortsError            bool
	GetLocalRDFPortDetailsError            bool
	CreateRDFGroupError                    bool
	GetStorageGroupSnapshotError           bool
	DeleteStorageGroupSnapshotError        bool
	GetStorageGroupSnapshotSnapError       bool
	GetStorageGroupSnapshotSnapDetailError bool
	GetStorageGroupSnapshotSnapModifyError bool
	GetSnapshotPolicyError                 bool
	GetSnapshotPolicyListError             bool
	CreateSnapshotPolicyError              bool
	ModifySnapshotPolicyError              bool
	DeleteSnapshotPolicyError              bool
	GetFileSystemListError                 bool
	GetNFSExportListError                  bool
	GetNASServerListError                  bool
	GetFileSystemError                     bool
	CreateFileSystemError                  bool
	UpdateFileSystemError                  bool
	DeleteFileSystemError                  bool
	GetNASServerError                      bool
	UpdateNASServerError                   bool
	DeleteNASServerError                   bool
	GetNFSExportError                      bool
	CreateNFSExportError                   bool
	UpdateNFSExportError                   bool
	DeleteNFSExportError                   bool
	GetFileInterfaceError                  bool
	ExecuteActionError                     bool
	GetFreshMetrics                        bool
}

// hasError checks to see if the specified error (via pointer)
// is set. If so it returns true, else false.
// Additionally, if ResetAfterFirstError is set, the first error
// condition will be reset to no longer be an error condition.
func hasError(errorType *bool) bool {
	if *errorType {
		if InducedErrors.ResetAfterFirstError {
			*errorType = false
			InducedErrors.ResetAfterFirstError = false
		}
		return true
	}
	return false
}

func SafeSetInducedError(inducedErrsPtr interface{}, errName string, value interface{}) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()

	v := reflect.ValueOf(inducedErrsPtr)

	// Check if it is a pointer
	if v.Kind() != reflect.Ptr {
		return errors.New("expected a pointer to struct")
	}

	// Dereference the pointer
	v = v.Elem()

	// Check if it is a struct
	if v.Kind() != reflect.Struct {
		return errors.New("expected a struct")
	}

	field := v.FieldByName(errName)
	if !field.IsValid() {
		return errors.New("invalid field name")
	}

	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	// Check if the type of the value matches the field type
	if field.Type() != reflect.TypeOf(value) {
		return errors.New("incompatible type")
	}

	field.Set(reflect.ValueOf(value))
	return nil
}

func SafeGetInducedError(inducedErrsPtr interface{}, errName string) (errValue interface{}, err error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()

	v := reflect.ValueOf(inducedErrsPtr)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("expected a pointer to struct")
	}

	v = v.Elem()

	// Check if it is a struct
	if v.Kind() != reflect.Struct {
		return nil, errors.New("expected a struct")
	}

	field := v.FieldByName(errName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", errName)
	}

	return field.Interface(), nil
}

// Reset : re-initializes the variables
func Reset() {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()

	InducedErrors.NoConnection = false
	InducedErrors.InvalidJSON = false
	InducedErrors.BadHTTPStatus = 0
	InducedErrors.GetSymmetrixError = false
	InducedErrors.GetVolumeIteratorError = false
	InducedErrors.GetVolumeError = false
	InducedErrors.UpdateVolumeError = false
	InducedErrors.DeleteVolumeError = false
	InducedErrors.DeviceInSGError = false
	InducedErrors.GetStorageGroupError = false
	InducedErrors.GetStorageGroupSnapshotPolicyError = false
	InducedErrors.InvalidResponse = false
	InducedErrors.UpdateStorageGroupError = false
	InducedErrors.ModifyMobilityError = false
	InducedErrors.GetJobError = false
	InducedErrors.JobFailedError = false
	InducedErrors.VolumeNotCreatedError = false
	InducedErrors.GetJobCannotFindRoleForUser = false
	InducedErrors.CreateStorageGroupError = false
	InducedErrors.StorageGroupAlreadyExists = false
	InducedErrors.DeleteStorageGroupError = false
	InducedErrors.GetStoragePoolListError = false
	InducedErrors.GetStoragePoolError = false
	InducedErrors.GetPortGroupError = false
	InducedErrors.GetPortError = false
	InducedErrors.GetSpecificPortError = false
	InducedErrors.GetPortISCSITargetError = false
	InducedErrors.GetPortNVMeTCPTargetError = false
	InducedErrors.GetPortGigEError = false
	InducedErrors.GetDirectorError = false
	InducedErrors.GetInitiatorError = false
	InducedErrors.GetInitiatorByIDError = false
	InducedErrors.GetHostError = false
	InducedErrors.CreateHostError = false
	InducedErrors.DeleteHostError = false
	InducedErrors.UpdateHostError = false
	InducedErrors.GetMaskingViewError = false
	InducedErrors.CreateMaskingViewError = false
	InducedErrors.UpdateMaskingViewError = false
	InducedErrors.MaskingViewAlreadyExists = false
	InducedErrors.DeleteMaskingViewError = false
	InducedErrors.PortGroupNotFoundError = false
	InducedErrors.InitiatorGroupNotFoundError = false
	InducedErrors.StorageGroupNotFoundError = false
	InducedErrors.VolumeNotAddedError = false
	InducedErrors.GetMaskingViewConnectionsError = false
	InducedErrors.ResetAfterFirstError = false
	InducedErrors.CreateSnapshotError = false
	InducedErrors.LinkSnapshotError = false
	InducedErrors.GetSymVolumeError = false
	InducedErrors.GetVolSnapsError = false
	InducedErrors.DeleteSnapshotError = false
	InducedErrors.RenameSnapshotError = false
	InducedErrors.GetGenerationError = false
	InducedErrors.GetPrivateVolumeIterator = false
	InducedErrors.SnapshotNotLicensed = false
	InducedErrors.UnisphereMismatchError = false
	InducedErrors.TargetNotDefinedError = false
	InducedErrors.SnapshotExpired = false
	InducedErrors.InvalidSnapshotName = false
	InducedErrors.GetPrivVolumeByIDError = false
	InducedErrors.CreatePortGroupError = false
	InducedErrors.UpdatePortGroupError = false
	InducedErrors.DeletePortGroupError = false
	InducedErrors.ExpandVolumeError = false
	InducedErrors.MaxSnapSessionError = false
	InducedErrors.GetSRDFInfoError = false
	InducedErrors.VolumeRdfTypesError = false
	InducedErrors.GetSRDFPairInfoError = false
	InducedErrors.GetProtectedStorageGroupError = false
	InducedErrors.CreateSGReplicaError = false
	InducedErrors.GetRDFGroupError = false
	InducedErrors.GetSGOnRemote = false
	InducedErrors.GetSGWithVolOnRemote = false
	InducedErrors.RDFGroupHasPairError = false
	InducedErrors.InvalidLocalVolumeError = false
	InducedErrors.InvalidRemoteVolumeError = false
	InducedErrors.GetRemoteVolumeError = false
	InducedErrors.FetchResponseError = false
	InducedErrors.RemoveVolumesFromSG = false
	InducedErrors.GetHostGroupError = false
	InducedErrors.CreateHostGroupError = false
	InducedErrors.DeleteHostGroupError = false
	InducedErrors.UpdateHostGroupError = false
	InducedErrors.GetHostGroupListError = false
	InducedErrors.GetStorageGroupMetricsError = false
	InducedErrors.GetVolumesMetricsError = false
	InducedErrors.GetFileSysMetricsError = false
	InducedErrors.GetStorageGroupPerfKeyError = false
	InducedErrors.GetArrayPerfKeyError = false
	InducedErrors.GetFreeRDFGError = false
	InducedErrors.GetLocalOnlineRDFDirsError = false
	InducedErrors.GetRemoteRDFPortOnSANError = false
	InducedErrors.GetLocalOnlineRDFPortsError = false
	InducedErrors.GetLocalRDFPortDetailsError = false
	InducedErrors.CreateRDFGroupError = false
	InducedErrors.GetStorageGroupSnapshotSnapDetailError = false
	InducedErrors.GetSnapshotPolicyError = false
	InducedErrors.GetSnapshotPolicyListError = false
	InducedErrors.CreateSnapshotPolicyError = false
	InducedErrors.ModifySnapshotPolicyError = false
	InducedErrors.DeleteSnapshotPolicyError = false
	InducedErrors.GetStorageGroupSnapshotError = false
	InducedErrors.CreateSnapshotPolicyError = false
	InducedErrors.GetStorageGroupSnapshotSnapError = false
	InducedErrors.GetStorageGroupSnapshotSnapModifyError = false
	InducedErrors.GetFileSystemListError = false
	InducedErrors.GetNFSExportListError = false
	InducedErrors.GetNASServerListError = false
	InducedErrors.GetFileSystemError = false
	InducedErrors.CreateFileSystemError = false
	InducedErrors.UpdateFileSystemError = false
	InducedErrors.DeleteFileSystemError = false
	InducedErrors.GetNASServerError = false
	InducedErrors.UpdateNASServerError = false
	InducedErrors.DeleteNASServerError = false
	InducedErrors.GetNFSExportError = false
	InducedErrors.CreateNFSExportError = false
	InducedErrors.UpdateNFSExportError = false
	InducedErrors.DeleteNFSExportError = false
	InducedErrors.GetFileInterfaceError = false
	InducedErrors.ExecuteActionError = false
	InducedErrors.GetFreshMetrics = false
	Data.JSONDir = "mock"
	Data.VolumeIDToIdentifier = make(map[string]string)
	Data.VolumeIDToSize = make(map[string]int)
	Data.VolumeIDIteratorList = make([]string, 0)
	Data.VolumeIDToSGList = make(map[string][]string)
	Data.MaskingViewIDToHostID = make(map[string]string)
	Data.MaskingViewIDToHostGroupID = make(map[string]string)
	Data.MaskingViewIDToPortGroupID = make(map[string]string)
	Data.MaskingViewIDToStorageGroupID = make(map[string]string)
	Data.StorageGroupIDToMaskingViewID = make(map[string]string)
	Data.JobIDToMockJob = make(map[string]*JobInfo)
	Data.StorageGroupIDToNVolumes = make(map[string]int)
	Data.StorageGroupIDToNVolumes[DefaultStorageGroup] = 0
	Data.StorageGroupIDToStorageGroup = make(map[string]*types.StorageGroup)
	Data.MaskingViewIDToMaskingView = make(map[string]*types.MaskingView)
	Data.InitiatorIDToInitiator = make(map[string]*types.Initiator)
	Data.HostIDToHost = make(map[string]*types.Host)
	Data.PortGroupIDToPortGroup = make(map[string]*types.PortGroup)
	Data.PortIDToSymmetrixPortType = make(map[string]*types.SymmetrixPortType)
	Data.VolumeIDToVolume = make(map[string]*types.Volume)
	Data.StorageGroupIDToVolumes = make(map[string][]string)
	Data.VolIDToSnapshots = make(map[string]map[string]*types.Snapshot)
	Data.SnapIDToLinkedVol = make(map[string]map[string]*types.LinkedVolumes)
	Data.StorageGroupIDToRDFStorageGroup = make(map[string]*types.RDFStorageGroup)
	Data.HostGroupIDToHostGroup = make(map[string]*types.HostGroup)
	Data.FileSysIDToFileSystem = make(map[string]*types.FileSystem)
	Data.NFSExportIDToNFSExport = make(map[string]*types.NFSExport)
	Data.NASServerIDToNASServer = make(map[string]*types.NASServer)
	Data.FileIntIDtoFileInterface = make(map[string]*types.FileInterface)
	Data.AsyncRDFGroup = &types.RDFGroup{
		RdfgNumber:          DefaultAsyncRDFGNo,
		Label:               DefaultAsyncRDFLabel,
		RemoteRdfgNumber:    DefaultAsyncRDFGNo,
		RemoteSymmetrix:     DefaultRemoteSymID,
		NumDevices:          0,
		TotalDeviceCapacity: 0.0,
		Modes:               []string{"Asynchronous"},
		Type:                "Dynamic",
		Async:               true,
	}
	Data.MetroRDFGroup = &types.RDFGroup{
		RdfgNumber:          DefaultMetroRDFGNo,
		Label:               DefaultMetroRDFLabel,
		RemoteRdfgNumber:    DefaultMetroRDFGNo,
		RemoteSymmetrix:     DefaultRemoteSymID,
		NumDevices:          0,
		TotalDeviceCapacity: 0.0,
		Modes:               []string{"Active"},
		Type:                "Metro",
		Metro:               true,
	}
	Data.AsyncSGRDFInfo = &types.SGRDFInfo{
		RdfGroupNumber: DefaultAsyncRDFGNo,
		VolumeRdfTypes: []string{"R1"},
		States:         []string{"Consistent"},
		Modes:          []string{"Asynchronous"},
		LargerRdfSides: []string{"Equal"},
	}
	Data.MetroSGRDFInfo = &types.SGRDFInfo{
		RdfGroupNumber: DefaultMetroRDFGNo,
		VolumeRdfTypes: []string{"R1"},
		States:         []string{"Consistent"},
		Modes:          []string{"Active"},
		LargerRdfSides: []string{"Equal"},
	}
	initMockCache()
}

func initMockCache() {
	// Initialize SGs
	addStorageGroup("CSI-Test-SG-1", "SRP_1", "Diamond")       // #nosec G20
	addStorageGroup("CSI-Test-SG-2", "SRP_1", "Diamond")       // #nosec G20
	addStorageGroup("CSI-Test-SG-3", "SRP_2", "Silver")        // #nosec G20
	addStorageGroup("CSI-Test-SG-4", "SRP_2", "Optimized")     // #nosec G20
	addStorageGroup("CSI-Test-SG-5", "SRP_2", "None")          // #nosec G20
	addStorageGroup("CSI-Test-SG-6", "None", "None")           // #nosec G20
	addStorageGroup("CSI-Test-Fake-Remote-SG", "None", "None") // #nosec G20
	// Initialize protected SG
	addStorageGroup(DefaultASYNCProtectedSG, "None", "None")        // #nosec G20
	addStorageGroup(DefaultMETROProtectedSG, "None", "None")        // #nosec G20
	addRDFStorageGroup(DefaultASYNCProtectedSG, DefaultRemoteSymID) // #nosec G20
	addRDFStorageGroup(DefaultMETROProtectedSG, DefaultRemoteSymID) // #nosec G20

	// ISCSI directors
	iscsiDir1 := "SE-1E"
	iscsidir1PortKey1 := iscsiDir1 + ":" + "4"
	// iscsiDir2 := "SE-2E"
	// FC directors
	fcDir1 := "FA-1D"
	fcDir2 := "FA-2D"
	fcDir1PortKey1 := fcDir1 + ":" + "5"
	fcDir2PortKey1 := fcDir2 + ":" + "1"
	// Add Port groups
	addPortGroupWithPortID("csi-pg", "Fibre", []string{fcDir1PortKey1, fcDir2PortKey1}) // #nosec G20
	// Initialize initiators
	// Initialize Hosts
	initNode1List := make([]string, 0)
	iqnNode1 := "iqn.1993-08.org.centos:01:5ae577b352a0"
	initNode1 := iscsidir1PortKey1 + ":" + iqnNode1
	initNode1List = append(initNode1List, iqnNode1)
	addInitiator(initNode1, iqnNode1, "GigE", []string{iscsidir1PortKey1}, "") // #nosec G20
	addHost("CSI-Test-Node-1", "iSCSI", initNode1List)                         // #nosec G20
	initNode2List := make([]string, 0)
	iqn1Node2 := "iqn.1993-08.org.centos:01:5ae577b352a1"
	iqn2Node2 := "iqn.1993-08.org.centos:01:5ae577b352a2"
	init1Node2 := iscsidir1PortKey1 + ":" + iqn1Node2
	init2Node2 := iscsidir1PortKey1 + ":" + iqn2Node2
	initNode2List = append(initNode2List, iqn1Node2)
	initNode2List = append(initNode2List, iqn2Node2)
	addInitiator(init1Node2, iqn1Node2, "GigE", []string{iscsidir1PortKey1}, "")       // #nosec G20
	addInitiator(init2Node2, iqn2Node2, "GigE", []string{iscsidir1PortKey1}, "")       // #nosec G20
	addHost("CSI-Test-Node-2", "iSCSI", initNode2List)                                 // #nosec G20
	addMaskingView("CSI-Test-MV-1", "CSI-Test-SG-1", "CSI-Test-Node-1", "iscsi_ports") // #nosec G20

	initNode3List := make([]string, 0)
	hba1Node3 := "20000090fa9278dd"
	hba2Node3 := "20000090fa9278dc"
	init1Node3 := fcDir1PortKey1 + ":" + hba1Node3
	init2Node3 := fcDir2PortKey1 + ":" + hba1Node3
	init3Node3 := fcDir1PortKey1 + ":" + hba2Node3
	init4Node3 := fcDir2PortKey1 + ":" + hba2Node3
	addInitiator(init1Node3, hba1Node3, "Fibre", []string{fcDir1PortKey1}, "") // #nosec G20
	addInitiator(init2Node3, hba1Node3, "Fibre", []string{fcDir2PortKey1}, "") // #nosec G20
	addInitiator(init3Node3, hba2Node3, "Fibre", []string{fcDir1PortKey1}, "") // #nosec G20
	addInitiator(init4Node3, hba2Node3, "Fibre", []string{fcDir2PortKey1}, "") // #nosec G20
	initNode3List = append(initNode3List, hba1Node3)
	initNode3List = append(initNode3List, hba2Node3)
	addHost("CSI-Test-Node-3-FC", "Fibre", initNode3List) // #nosec G20
	addTempSnapshots()
	addFileObjects()
}

func AddFileObjects() {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addFileObjects()
}

// AddFileObjects adds file objects for mock objects
func addFileObjects() {
	// Add a File System
	addNewFileSystem("id1", DefaultFSName, 4000)
	// Add a NFS Export
	addNewNFSExport("id1", "nfs-0")
	addNewNFSExport("id2", "nfs-del")
	// Add a NAS Server
	addNewNASServer("id1", "nas-1")
	addNewNASServer("id2", "nas-del")
	// Add a FileInterface
	addNewFileInterface("id1", "interface-1")
}

var mockRouter http.Handler

// GetHandler returns the http handler
func GetHandler() http.Handler {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if Debug {
				log.Printf("handler called: %s %s", r.Method, r.URL)
			}
			invalidJSONErr, err := SafeGetInducedError(InducedErrors, "InvalidJSON")
			if err != nil {
				writeError(w, "failed to get induced error for InvalidJSON", http.StatusRequestTimeout)
				return
			}
			noConnectionErr, err := SafeGetInducedError(InducedErrors, "NoConnection")
			if err != nil {
				writeError(w, "failed to get induced error for NoConnection", http.StatusRequestTimeout)
				return
			}
			badHTTPStatusErr, err := SafeGetInducedError(InducedErrors, "BadHTTPStatus")
			if err != nil {
				writeError(w, "failed to get induced error for BadHTTPStatus", http.StatusRequestTimeout)
				return
			}

			if invalidJSONErr.(bool) {
				w.Write([]byte(`this is not json`)) // #nosec G20
			} else if noConnectionErr.(bool) {
				writeError(w, "No Connection", http.StatusRequestTimeout)
			} else if badHTTPStatusErr.(int) != 0 {
				writeError(w, "Internal Error", InducedErrors.BadHTTPStatus)
			} else {
				if mockRouter != nil {
					mockRouter.ServeHTTP(w, r)
				} else {
					getRouter().ServeHTTP(w, r)
				}
			}
		})
	return handler
}

func getRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/host/{id}", HandleHost)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/host", HandleHost)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/hostgroup/{id}", HandleHostGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/hostgroup", HandleHostGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/initiator/{id}", HandleInitiator)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/initiator", HandleInitiator)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/portgroup/{id}", HandlePortGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/portgroup", HandlePortGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/storagegroup/{id}", HandleStorageGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/storagegroup", HandleStorageGroup)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/maskingview/{mvID}/connections", HandleMaskingViewConnections)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/maskingview/{mvID}", HandleMaskingView)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/maskingview", HandleMaskingView)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/srp/{id}", HandleStorageResourcePool)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/srp", HandleStorageResourcePool)
	router.HandleFunc(PREFIXNOVERSION+"/common/Iterator/{iterId}/page", HandleIterator)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/volume/{volID}", HandleVolume)
	router.HandleFunc(PREFIX+"/sloprovisioning/symmetrix/{symid}/volume", HandleVolume)
	router.HandleFunc(PRIVATEPREFIX+"/sloprovisioning/symmetrix/{symid}/volume", HandlePrivVolume)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/director/{director}/port/{id}", HandlePort)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/director/{director}/port", HandlePort)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/director/{id}", HandleDirector)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/director", HandleDirector)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/job/{jobID}", HandleJob)
	router.HandleFunc(PREFIX+"/system/symmetrix/{symid}/job", HandleJob)
	router.HandleFunc(PREFIX+"/system/symmetrix/{id}", HandleSymmetrix)
	router.HandleFunc(PREFIX+"/system/symmetrix", HandleSymmetrix)
	router.HandleFunc(PREFIX+"/system/version", HandleVersion)
	router.HandleFunc(PREFIX+"/version", HandleVersion)
	router.HandleFunc(PREFIXNOVERSION+"/version", HandleVersion)
	router.HandleFunc(PREFIX+"/system/symmetrix/{id}/refresh", HandleSymmetrix)
	router.HandleFunc("/", HandleNotFound)

	// StorageGroup Snapshots
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot", HandleGetStorageGroupSnapshots)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid", HandleGetStorageGroupSnapshotsSnapsIDs)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid/{snapID}", HandleGetStorageGroupSnapshotsSnapsDetails)

	// Snapshot
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/snapshot/{SnapID}", HandleSnapshot)
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/volume", HandleSymVolumes)
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/volume/{volID}/snapshot", HandleVolSnaps)
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}", HandleVolSnaps)
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}/generation", HandleGenerations)
	router.HandleFunc(PRIVATEPREFIX+"/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}/generation/{genID}", HandleGenerations)
	router.HandleFunc(PREFIX+"/replication/capabilities/symmetrix", HandleCapabilities)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symID}/snapshot_policy/{snapshotPolicyID}/storagegroup/{storageGroupID}", HandleStorageGroupSnapshotPolicy)

	// SRDF
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/rdf_group", HandleRDFGroup)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/rdf_group/{rdf_no}", HandleRDFGroup)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{id}", HandleRDFStorageGroup)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{id}/rdf_group", HandleRDFStorageGroup)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/storagegroup/{id}/rdf_group/{rdf_no}", HandleSGRDF)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/rdf_group/{rdf_no}/volume/{volume_id}", HandleRDFDevicePair)
	router.HandleFunc(INTERNALPREFIX+"/file/symmetrix/{symID}/rdf_group_numbers_free", HandleFreeRDF)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symID}/rdf_director", HandleRDFDirector)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symID}/rdf_director/{dir}/port", HandleRDFPort)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symID}/rdf_director/{dir}/port/{port}", HandleRDFPort)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symID}/rdf_director/{dir}/port/{port}/remote_port", HandleRDFRemotePort)

	// Performance Metrics
	router.HandleFunc(PREFIXNOVERSION+"/performance/StorageGroup/metrics", HandleStorageGroupMetrics)
	router.HandleFunc(PREFIXNOVERSION+"/performance/Volume/metrics", HandleVolumeMetrics)
	router.HandleFunc(PREFIXNOVERSION+"/performance/file/filesystem/metrics", HandleFileSysMetrics)

	// Performance Keys
	router.HandleFunc(PREFIXNOVERSION+"/performance/StorageGroup/keys", HandleStorageGroupPerfKeys)
	router.HandleFunc(PREFIXNOVERSION+"/performance/Array/keys", HandleArrayPerfKeys)

	// Snapshot Policy
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/snapshot_policy/{snapshotPolicyId}", HandleGetSnapshotPolicy)
	router.HandleFunc(PREFIX+"/replication/symmetrix/{symid}/snapshot_policy", HandleCreateSnapshotPolicy)

	// File APIs
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/file_system/{fsID}", HandleFileSystem)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/file_system", HandleFileSystem)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/nfs_export/{nfsID}", HandleNFSExport)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/nfs_export", HandleNFSExport)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/nas_server/{nasID}", HandleNASServer)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/nas_server", HandleNASServer)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/file_interface", HandleFileInterface)
	router.HandleFunc(PREFIX+"/file/symmetrix/{symid}/file_interface/{interfaceID}", HandleFileInterface)

	mockRouter = router
	return router
}

// GET /replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid/{snapID}
// PUT /replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid/{snapID}
// DELETE /replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid/{snapID}
func HandleGetStorageGroupSnapshotsSnapsDetails(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleGetStorageGroupSnapshotsSnapsDetails(w, r)
}

func handleGetStorageGroupSnapshotsSnapsDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPut && r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetStorageGroupSnapshotSnapDetailError {
		writeError(w, "Could not get StorageGroup Snapshots Snap Ids: induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.GetStorageGroupSnapshotSnapModifyError {
		writeError(w, "Could not get StorageGroup Snapshots Snap Ids: induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.DeleteStorageGroupSnapshotError {
		writeError(w, "Could not delete StorageGroup Snapshots Snap Ids: induced error", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		sgCreateSnap := &types.StorageGroupSnap{
			Name:       "sg_1_snap",
			Generation: 1,
			SnapID:     2,
			Timestamp:  "1234",
		}
		writeJSON(w, sgCreateSnap)
		return
	}
	if r.Method == http.MethodPut {
		sgCreateSnap := &types.StorageGroupSnap{
			Name:       "sg_1_snap_2",
			Generation: 1,
			SnapID:     2,
			Timestamp:  "1234",
		}
		writeJSON(w, sgCreateSnap)
		return
	}
}

// GET /replication/symmetrix/{symid}/storagegroup/{StorageGroupId}/snapshot/{snapshotId}/snapid
func HandleGetStorageGroupSnapshotsSnapsIDs(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleGetStorageGroupSnapshotsSnapsIDs(w, r)
}

func handleGetStorageGroupSnapshotsSnapsIDs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetStorageGroupSnapshotSnapError {
		writeError(w, "Could not get StorageGroup Snapshots Snap Ids: induced error", http.StatusBadRequest)
		return
	}
	snaps := make([]int64, 1)
	snaps = append(snaps, 1234)
	snapIDs := &types.SnapID{
		SnapIDs: snaps,
	}
	writeJSON(w, snapIDs)
}

// GET /replication/symmetrix/{symid}/storagegroup/{SnapID}/snapshot
// POST /replication/symmetrix/{symid}/storagegroup/{SnapID}/snapshot
func HandleGetStorageGroupSnapshots(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleGetStorageGroupSnapshots(w, r)
}

func handleGetStorageGroupSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetStorageGroupSnapshotError {
		writeError(w, "Could not get StorageGroup Snapshots: induced error", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		names := make([]string, 1)
		names = append(names, "sg_1_snap")
		namesAndCounts := make([]types.SnapshotNameAndCounts, 1)
		namesAndCounts = append(namesAndCounts, types.SnapshotNameAndCounts{
			Name:               "name",
			SnapshotCount:      1,
			NewestTimestampUtc: 123,
		})
		sgSnap := &types.StorageGroupSnapshot{
			Name:                   names,
			SlSnapshotName:         names,
			SnapshotNamesAndCounts: namesAndCounts,
		}
		writeJSON(w, sgSnap)
	} else {
		sgCreateSnap := &types.StorageGroupSnap{
			Name:       "sg_1_snap",
			Generation: 1,
			SnapID:     2,
			Timestamp:  "1234",
		}
		writeJSON(w, sgCreateSnap)
	}
}

// GET /replication/symmetrix/{symid}/snapshot_policy/{snapshotPolicyId}
// PUT /replication/symmetrix/{symid}/snapshot_policy/{snapshotPolicyId}
// DELETE /replication/symmetrix/{symid}/snapshot_policy/{snapshotPolicyId}
func HandleGetSnapshotPolicy(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleGetSnapshotPolicy(w, r)
}

func handleGetSnapshotPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPut && r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetSnapshotPolicyError {
		writeError(w, "Could not get Snapshot Policy : induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.ModifySnapshotPolicyError {
		writeError(w, "Could not update Snapshot Policy : induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.DeleteSnapshotPolicyError {
		writeError(w, "Could not delete Snapshot Policy : induced error", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {

		snapPolicy := &types.SnapshotPolicy{
			SymmetrixID:            "000197902572",
			SnapshotPolicyName:     "WeeklyDefault",
			SnapshotCount:          13,
			IntervalMinutes:        10080,
			OffsetMinutes:          10074,
			Suspended:              false,
			Secure:                 false,
			LastTimeUsed:           "23:53:15 Sun, 07 May 2023 +0000",
			StorageGroupCount:      1,
			ComplianceCountWarning: 10,
			Type:                   "local",
		}
		writeJSON(w, snapPolicy)
	}
	if r.Method == http.MethodPut {
		snapPolicy := &types.SnapshotPolicy{
			SymmetrixID:            "000197902572",
			SnapshotPolicyName:     "WeeklyDefault",
			SnapshotCount:          13,
			IntervalMinutes:        10080,
			OffsetMinutes:          10074,
			Suspended:              false,
			Secure:                 false,
			LastTimeUsed:           "23:53:15 Sun, 07 May 2023 +0000",
			StorageGroupCount:      1,
			ComplianceCountWarning: 10,
			Type:                   "local",
		}
		writeJSON(w, snapPolicy)
	}
}

// POST /replication/symmetrix/{symid}/snapshot_policy
// GET /replication/symmetrix/{symid}/snapshot_policy
func HandleCreateSnapshotPolicy(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleCreateSnapshotPolicy(w, r)
}

func handleCreateSnapshotPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.CreateSnapshotPolicyError {
		writeError(w, "Could not Create Snapshot Policy : induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.GetSnapshotPolicyListError {
		writeError(w, "Could not get Snapshot Policy List: induced error", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		ids := []string{"Test123", "Test345"}
		snapPolicyList := &types.SnapshotPolicyList{
			SnapshotPolicyIDs: ids,
		}
		writeJSON(w, snapPolicyList)
	}
	if r.Method == http.MethodPost {

		snapPolicy := &types.SnapshotPolicy{
			SymmetrixID:            "000197902572",
			SnapshotPolicyName:     "WeeklyDefault",
			SnapshotCount:          13,
			IntervalMinutes:        10080,
			OffsetMinutes:          10074,
			Suspended:              false,
			Secure:                 false,
			LastTimeUsed:           "23:53:15 Sun, 07 May 2023 +0000",
			StorageGroupCount:      1,
			ComplianceCountWarning: 10,
			Type:                   "local",
		}
		writeJSON(w, snapPolicy)
	}
}

// GET univmax/restapi/100/replication/symmetrix/{symID}/rdf_director/{dir}/port?online=true
// GET univmax/restapi/100/replication/symmetrix/{symID}/rdf_director/{dir}/port/{port}
func HandleRDFPort(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFPort(w, r)
}

func handleRDFPort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetLocalOnlineRDFPortsError {
		writeError(w, "Could not retrieve local online RDF ports: induced error", http.StatusBadRequest)
		return
	}
	if InducedErrors.GetLocalRDFPortDetailsError {
		writeError(w, "Could not retrieve local RDF port: induced error", http.StatusBadRequest)
		return
	}
	routeParams := mux.Vars(r)
	portID := routeParams["port"]
	if portID != "" {
		rdfPorts := &types.RDFPortDetails{
			SymmID:     DefaultSymmetrixID,
			DirNum:     33,
			DirID:      DefaultRDFDir,
			PortNum:    DefaultRDFPort,
			PortOnline: true,
			PortWWN:    "5000097200007003",
		}
		writeJSON(w, rdfPorts)
	} else {
		rdfPorts := &types.RDFPortList{RdfPorts: []string{"3"}}
		writeJSON(w, rdfPorts)
	}
}

// GET univmax/restapi/100/replication/symmetrix/{symID}/rdf_director/{dir}/port/{port}/remote_port
func HandleRDFRemotePort(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFRemotePort(w, r)
}

func handleRDFRemotePort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetRemoteRDFPortOnSANError {
		writeError(w, "Could not retrieve remote RDF port: induced error", http.StatusBadRequest)
		return
	}
	routeParams := mux.Vars(r)
	portID := routeParams["port"]
	if portID == "" {
		writeError(w, "portID is nil in request, can not retrieve remote port details", http.StatusBadRequest)
		return
	}
	remotePorts := &types.RemoteRDFPortDetails{
		RemotePorts: []types.RDFPortDetails{
			{
				SymmID:     DefaultRemoteSymID,
				DirNum:     33,
				DirID:      DefaultRDFDir,
				PortNum:    DefaultRDFPort,
				PortOnline: true,
				PortWWN:    "5000097200007003",
			},
		},
	}
	writeJSON(w, remotePorts)
}

// GET univmax/restapi/100/replication/symmetrix/{symID}/rdf_director?online=true
func HandleRDFDirector(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFDirector(w, r)
}

func handleRDFDirector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetLocalOnlineRDFDirsError {
		writeError(w, "Could not retrieve RDF director: induced error", http.StatusBadRequest)
		return
	}
	dir := &types.RDFDirList{
		RdfDirs: []string{DefaultRDFDir, "OR-2C"},
	}
	writeJSON(w, dir)
}

// GET univmax/restapi/internal/100/file/symmetrix/{symID}/rdf_group_numbers_free?remote_symmetrix_id={remoteSymID}
func HandleFreeRDF(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleFreeRDF(w, r)
}

func handleFreeRDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetFreeRDFGError {
		writeError(w, "Could not retrieve free RDF group: induced error", http.StatusBadRequest)
		return
	}
	nxtFreeRDFG := &types.NextFreeRDFGroup{
		LocalRdfGroup:  []int{DefaultAsyncRDFGNo},
		RemoteRdfGroup: []int{DefaultAsyncRDFGNo},
	}
	writeJSON(w, nxtFreeRDFG)
}

// NewVolume creates a new mock volume with the specified characteristics.
func NewVolume(volumeID, volumeIdentifier string, size int, sgList []string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	Data.VolumeIDToIdentifier[volumeID] = volumeIdentifier
	fmt.Printf("NewVolume: id %s name %s\n", volumeID, volumeIdentifier)
	Data.VolumeIDToSize[volumeID] = size
	Data.VolumeIDToSGList[volumeID] = sgList
}

// TO be used for the endpoints that don't have handlers yet
func handleTODO(w http.ResponseWriter, _ *http.Request) {
	writeError(w, "Endpoint not implemented yet", http.StatusNotImplemented)
}

// GET, POST /univmax/restapi/APIVersion/replication/symmetrix/{symID}/rdf_group/{rdf_no}/volume/{volume_id}
func HandleRDFDevicePair(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFDevicePair(w, r)
}

func handleRDFDevicePair(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRDFDevicePairInfo(w, r)
	case http.MethodPost:
		handleRDFDevicePairCreation(w, r)
	default:
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRDFDevicePairCreation(w http.ResponseWriter, r *http.Request) {
	// TODO: Update mock cache based on the request payload.
	routeParams := mux.Vars(r)
	rdfPairs := new(types.RDFDevicePairList)
	rdfPairs.RDFDevicePair = []types.RDFDevicePair{
		{
			RemoteVolumeName:     routeParams["volume_id"],
			LocalVolumeName:      routeParams["volume_id"],
			RemoteSymmID:         routeParams["symid"],
			LocalSymmID:          routeParams["symid"],
			LocalRdfGroupNumber:  DefaultAsyncRDFGNo,
			RemoteRdfGroupNumber: DefaultAsyncRemoteRDFGNo,
		},
	}
	writeJSON(w, rdfPairs)
}

// GET /univmax/restapi/APIVersion/replication/symmetrix/{symID}/rdf_group/{rdf_no}/volume/{volume_id}
func handleRDFDevicePairInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if InducedErrors.GetSRDFPairInfoError {
		writeError(w, "Could not retrieve pair info", http.StatusBadRequest)
		return
	}
	routeParams := mux.Vars(r)
	var volumeConfig string
	if routeParams["symid"] == Data.AsyncRDFGroup.RemoteSymmetrix {
		volumeConfig = "RDF2+TDEV"
	} else {
		volumeConfig = "RDF1+TDEV"
	}
	rdfDevicePairInfo := &types.RDFDevicePair{
		LocalRdfGroupNumber:  Data.AsyncRDFGroup.RdfgNumber,
		RemoteRdfGroupNumber: Data.AsyncRDFGroup.RdfgNumber,
		LocalSymmID:          routeParams["symid"],
		RemoteSymmID:         Data.AsyncRDFGroup.RemoteSymmetrix,
		LocalVolumeName:      routeParams["volume_id"],
		RemoteVolumeName:     routeParams["volume_id"],
		VolumeConfig:         volumeConfig,
		RdfMode:              Data.AsyncRDFGroup.Modes[0],
		RdfpairState:         "Consistent",
		LargerRdfSide:        "Equal",
	}
	writeJSON(w, rdfDevicePairInfo)
}

// GET, POST /univmax/restapi/APIVersion/replication/symmetrix/{symID}/rdf_group/
// GET /univmax/restapi/APIVersion/replication/symmetrix/{symID}/rdf_group/{rdf_no}
func HandleRDFGroup(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFGroup(w, r)
}

func handleRDFGroup(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.CreateRDFGroupError {
		writeError(w, "error creating RDF group: induced error", http.StatusNotFound)
		return
	}
	if InducedErrors.GetRDFGroupError {
		writeError(w, "the specified RA group does not exist: induced error", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		routeParams := mux.Vars(r)
		rdfGroupNumber := routeParams["rdf_no"]
		returnRDFGroup(w, rdfGroupNumber)
	case http.MethodPost:
		writeJSON(w, Data.AsyncRDFGroup)
	default:
		writeError(w, "Method["+r.Method+"] not allowed", http.StatusMethodNotAllowed)
	}
}

// ReturnRDFGroup - Returns RDF group information from mock cache
func ReturnRDFGroup(w http.ResponseWriter, rdfg string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnRDFGroup(w, rdfg)
}

func returnRDFGroup(w http.ResponseWriter, rdfg string) {
	if rdfg != "" {
		if rdfg != fmt.Sprintf("%d", Data.AsyncRDFGroup.RdfgNumber) && rdfg != fmt.Sprintf("%d", Data.MetroRDFGroup.RdfgNumber) {
			writeError(w, "The specified RA group is not valid", http.StatusNotFound)
		} else {
			if InducedErrors.RDFGroupHasPairError {
				Data.AsyncRDFGroup.NumDevices = 1
			}
			writeJSON(w, Data.AsyncRDFGroup)
		}
	} else {
		rdflist := &types.RDFGroupList{
			RDFGroupCount: 1,
			RDFGroupIDs: []types.RDFGroupIDL{
				{
					RDFGNumber:  1,
					Label:       "mock",
					RemoteSymID: DefaultRemoteSymID,
					GroupType:   "Dynamic",
				},
			},
		}
		writeJSON(w, rdflist)
	}
}

// GET /univmax/restapi/APIVersion/replication/symmetrix/{symid}/storagegroup/{id}
// POST /univmax/restapi/APIVersion/replication/symmetrix/{symid}/storagegroup/{id}/rdf_group
func HandleRDFStorageGroup(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleRDFStorageGroup(w, r)
}

func handleRDFStorageGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetProtectedStorageGroupError {
			writeError(w, "The requested storage group cannot be found: induced error", http.StatusNotFound)
			return
		}
		if InducedErrors.FetchResponseError {
			writeError(w, "Error fetching response", http.StatusBadRequest)
		}
		handleSGRDFFetch(w, r)
	case http.MethodPost:
		if InducedErrors.CreateSGReplicaError {
			writeError(w, "Failed to create SG replica: induced error", http.StatusNotFound)
			return
		}
		handleSGRDFCreation(w, r)
	default:
		writeError(w, "Method["+r.Method+"] not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSGRDFFetch(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	storageGroupID := routeParams["id"]
	symmetrixID := routeParams["symid"]
	var (
		rdfStorageGroup *types.RDFStorageGroup
		ok              bool
	)
	if _, ok = Data.StorageGroupIDToStorageGroup[storageGroupID]; !ok {
		writeError(w, "The requested storage group does not exist", http.StatusNotFound)
		return
	}
	if rdfStorageGroup, ok = Data.StorageGroupIDToRDFStorageGroup[storageGroupID]; !ok {
		rdfStorageGroup = &types.RDFStorageGroup{
			SymmetrixID: symmetrixID,
			Name:        storageGroupID,
			Rdf:         false,
		}
	}
	if InducedErrors.RDFGroupHasPairError {
		rdfStorageGroup.NumDevicesNonGk = 1
	}
	writeJSON(w, rdfStorageGroup)
}

func handleSGRDFCreation(w http.ResponseWriter, r *http.Request) {
	sgsrdf := new(types.CreateSGSRDF)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(sgsrdf); err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}
	routeParams := mux.Vars(r)
	mode := sgsrdf.ReplicationMode
	storageGroupName := routeParams["id"]
	symmetrixID := routeParams["symid"]
	if _, err := addRDFStorageGroup(storageGroupName, symmetrixID); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, volumeID := range Data.StorageGroupIDToVolumes[storageGroupName] {
		if _, ok := Data.VolumeIDToVolume[volumeID]; !ok {
			continue
		}
		volume := Data.VolumeIDToVolume[volumeID]
		volume.Type = "RDF1+TDEV"
		if strings.Compare(mode, "Active") == 0 {
			volume.RDFGroupIDList = []types.RDFGroupID{
				{RDFGroupNumber: Data.MetroRDFGroup.RdfgNumber},
			}
		} else {
			volume.RDFGroupIDList = []types.RDFGroupID{
				{RDFGroupNumber: Data.AsyncRDFGroup.RdfgNumber},
			}
		}
	}
	sgrdfInfo := new(types.SGRDFInfo)
	dataToCopy := new(types.SGRDFInfo)
	if mode == "Active" {
		dataToCopy = Data.MetroSGRDFInfo
	} else {
		dataToCopy = Data.AsyncSGRDFInfo
	}
	err := copier.Copy(sgrdfInfo, dataToCopy)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sgrdfInfo.SymmetrixID = symmetrixID
	sgrdfInfo.StorageGroupName = storageGroupName
	writeJSON(w, sgrdfInfo)
}

// GET, PUT /replication/symmetrix/{symid}/storagegroup/{id}/rdf_group/{rdf_no}
func HandleSGRDF(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleSGRDF(w, r)
}

func handleSGRDF(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleSGRDFInfo(w, r)
	case http.MethodPut:
		handleSGRDFAction(w, r)
	default:
		writeError(w, "Method["+r.Method+"] not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSGRDFInfo(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.GetSRDFInfoError {
		writeError(w, "Error retrieving SRDF Info: induced error", http.StatusRequestTimeout)
		return
	}
	routeParams := mux.Vars(r)
	rdfNo := routeParams["rdf_no"]
	if rdfNo != fmt.Sprintf("%d", Data.AsyncRDFGroup.RdfgNumber) && rdfNo != fmt.Sprintf("%d", Data.MetroRDFGroup.RdfgNumber) {
		writeError(w, "The specified RA group is not valid", http.StatusNotFound)
	} else {
		sgrdfInfo := new(types.SGRDFInfo)
		var err error
		if rdfNo == fmt.Sprintf("%d", Data.AsyncRDFGroup.RdfgNumber) {
			err = copier.Copy(sgrdfInfo, Data.AsyncSGRDFInfo)
		} else {
			err = copier.Copy(sgrdfInfo, Data.MetroSGRDFInfo)
		}
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if InducedErrors.VolumeRdfTypesError {
			sgrdfInfo.VolumeRdfTypes = []string{"invalid"}
		}
		sgrdfInfo.SymmetrixID = routeParams["symid"]
		sgrdfInfo.StorageGroupName = routeParams["id"]
		writeJSON(w, sgrdfInfo)
	}
}

func handleSGRDFAction(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.ExecuteActionError {
		writeError(w, "Failed to execute action on RDFG: induced error", http.StatusBadRequest)
		return
	}
	routeParams := mux.Vars(r)
	rdfNo := routeParams["rdf_no"]
	decoder := json.NewDecoder(r.Body)
	modifySRDFGParam := &types.ModifySGRDFGroup{}
	err := decoder.Decode(modifySRDFGParam)
	if err != nil {
		writeError(w, "problem decoding PUT ACTION payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	action := modifySRDFGParam.Action
	performActionOnRDFSG(w, rdfNo, action)
}

// PerformActionOnRDFSG updates rdfNo with given action
func PerformActionOnRDFSG(w http.ResponseWriter, rdfNo, action string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	performActionOnRDFSG(w, rdfNo, action)
}

func performActionOnRDFSG(w http.ResponseWriter, rdfNo, action string) {
	if rdfNo != fmt.Sprintf("%d", Data.AsyncRDFGroup.RdfgNumber) && rdfNo != fmt.Sprintf("%d", Data.MetroRDFGroup.RdfgNumber) {
		writeError(w, "The specified RA group is not valid", http.StatusNotFound)
	} else {
		// we only support actions on ASYNC
		switch action {
		case "Establish":
			Data.AsyncSGRDFInfo.States = []string{"Consistent"}
			return
		case "Suspend":
			Data.AsyncSGRDFInfo.States = []string{"Suspended"}
			return
		case "Resume":
			Data.AsyncSGRDFInfo.States = []string{"Consistent"}
			return
		case "Failback":
			Data.AsyncSGRDFInfo.States = []string{"Consistent"}
			return
		case "Failover":
			Data.AsyncSGRDFInfo.States = []string{"Failed Over"}
			return
		case "Swap":
			Data.AsyncSGRDFInfo.States = []string{"Consistent"}
			return
		}
	}
}

// GET /univmax/restapi/system/version
func HandleVersion(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleVersion(w, r)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	auth := defaultUsername + ":" + defaultPassword
	authExpected := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
	// Check for valid credentials
	authSupplied := r.Header.Get("Authorization")
	if authExpected != authSupplied {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	apiversion := vars["apiversion"]
	// check the apiversion
	switch apiversion {
	case "90":
		w.Write([]byte(`{ "version": "V9.0.1.6" }`)) // #nosec G20
		break
	case "": // for version 91, as URL does not have apiversion in V9.1
		w.Write([]byte(`{ "version": "V9.1.0.2" }`)) // #nosec G20
		break
	default:
		writeError(w, "Unsupport API version: "+apiversion, http.StatusServiceUnavailable)
	}
}

// GET /univmax/restapi/APIVersion/system/symmetrix/{id}"
// GET /univmax/restapi/APIVersion/system/symmetrix"
func HandleSymmetrix(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleSymmetrix(w, r)
}

func handleSymmetrix(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.GetSymmetrixError {
		writeError(w, "Error retrieving Symmetrix: induced error", http.StatusRequestTimeout)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		returnJSONFile(Data.JSONDir, "symmetrixList.json", w, nil)
	}
	if id != "000197900046" && id != "000197900047" && id != DefaultRemoteSymID {
		writeError(w, "Symmetrix not found", http.StatusNotFound)
		return
	}
	if id == "000197900046" {
		returnJSONFile(Data.JSONDir, "symmetrix46.json", w, nil)
	} else if id == "000197900047" {
		returnJSONFile(Data.JSONDir, "symmetrix47.json", w, nil)
	} else {
		returnJSONFile(Data.JSONDir, "symmetrix13.json", w, nil)
	}
}

func HandleStorageResourcePool(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleStorageResourcePool(w, r)
}

func handleStorageResourcePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srpID := vars["id"]
	if InducedErrors.GetStoragePoolListError {
		writeError(w, "Error retrieving StoragePools: induced error", http.StatusRequestTimeout)
		return
	}
	if InducedErrors.GetStoragePoolError {
		writeError(w, "Error retrieving Storage Pool(s): induced error", http.StatusRequestTimeout)
		return
	}
	if srpID == "" {
		returnJSONFile(Data.JSONDir, "storageResourcePool.json", w, nil)
	}
	replacements := make(map[string]string)
	replacements["__SRP_ID__"] = "SRP_1"
	returnJSONFile(Data.JSONDir, "storage_pool_template.json", w, replacements)
}

// GET /univmax/restapi/API_VERSION/sloprovisioning/symmetrix/{id}/volume/{id}
// GET /univmax/restapi/API_VERSION/sloprovisioning/symmetrix/{id}/volume
func HandleVolume(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleVolume(w, r)
}

func handleVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volID := vars["volID"]
	switch r.Method {
	case http.MethodGet:
		if volID == "" {
			if InducedErrors.GetVolumeIteratorError {
				writeError(w, "Error getting VolumeIterator: induced error", http.StatusRequestTimeout)
				return
			}
			// Here we want a volume iterator.
			var like bool
			queryParams := r.URL.Query()
			volumeIdentifier := queryParams.Get("volume_identifier")
			if strings.Contains(volumeIdentifier, "<like>") {
				like = true
				volumeIdentifier = strings.TrimPrefix(volumeIdentifier, "<like>")
			}
			// Copy data to Data.VolumeIDIteratorList, while checking for volumeIdentifier match if needed
			Data.VolumeIDIteratorList = make([]string, 0)
			for _, vol := range Data.VolumeIDToVolume {
				if volumeIdentifier != "" {
					if like {
						if !strings.Contains(vol.VolumeIdentifier, volumeIdentifier) {
							continue
						}
					} else {
						if vol.VolumeIdentifier != volumeIdentifier {
							continue
						}
					}
				}
				Data.VolumeIDIteratorList = append(Data.VolumeIDIteratorList, vol.VolumeID)
			}
			if Debug {
				fmt.Printf("Data.VolumeIDIteratorList %#v", Data.VolumeIDIteratorList)
			}
			iter := &types.VolumeIterator{
				Count:          len(Data.VolumeIDIteratorList),
				ID:             "Volume",
				MaxPageSize:    10,
				ExpirationTime: 0,
			}
			numberToDo := len(Data.VolumeIDIteratorList)
			if numberToDo > iter.MaxPageSize {
				numberToDo = iter.MaxPageSize
			}
			iter.ResultList.From = 1
			iter.ResultList.To = numberToDo
			for i := iter.ResultList.From - 1; i <= iter.ResultList.To-1; i++ {
				volIDList := types.VolumeIDList{VolumeIDs: Data.VolumeIDIteratorList[i]}
				iter.ResultList.VolumeList = append(iter.ResultList.VolumeList, volIDList)
			}
			if Debug {
				fmt.Printf("iter: %#v\n", iter)
			}
			encoder := json.NewEncoder(w)
			err := encoder.Encode(iter)
			if err != nil {
				writeError(w, "json encoding error", http.StatusInternalServerError)
			}
			return
		}
		if InducedErrors.GetVolumeError {
			writeError(w, "Error retrieving Volume: induced error", http.StatusRequestTimeout)
			return
		}
		if volID != "" {
			if vars["symid"] == Data.AsyncRDFGroup.RemoteSymmetrix {
				returnVolume(w, volID, true)
			} else {
				returnVolume(w, volID, false)
			}
		}

	case http.MethodPut:
		if InducedErrors.UpdateVolumeError {
			writeError(w, "Error updating Volume: induced error", http.StatusRequestTimeout)
			return
		}
		if volID == "" {
			writeError(w, "Volume ID must be supplied", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updateVolumePayload := &types.EditVolumeParam{}
		err := decoder.Decode(updateVolumePayload)
		if err != nil {
			writeError(w, "problem decoding PUT Volume payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("PUT volume payload: %#v\n", updateVolumePayload)
		executionOption := updateVolumePayload.ExecutionOption
		if updateVolumePayload.EditVolumeActionParam.FreeVolumeParam != nil {
			freeVolume(w, updateVolumePayload.EditVolumeActionParam.FreeVolumeParam, volID, executionOption)
			return
		}
		if updateVolumePayload.EditVolumeActionParam.EnableMobilityIDParam != nil {
			modifyMobility(w, updateVolumePayload.EditVolumeActionParam.EnableMobilityIDParam, volID, executionOption)
			return
		}
		if updateVolumePayload.EditVolumeActionParam.ModifyVolumeIdentifierParam != nil {
			if vars["symid"] == Data.AsyncRDFGroup.RemoteSymmetrix {
				renameVolume(w, updateVolumePayload.EditVolumeActionParam.ModifyVolumeIdentifierParam, volID, executionOption, true)
			} else {
				renameVolume(w, updateVolumePayload.EditVolumeActionParam.ModifyVolumeIdentifierParam, volID, executionOption, false)
			}
			return
		}
		if updateVolumePayload.EditVolumeActionParam.ExpandVolumeParam != nil {
			expandVolume(w, updateVolumePayload.EditVolumeActionParam.ExpandVolumeParam, volID, executionOption)
			return
		}
	case http.MethodDelete:
		if InducedErrors.DeleteVolumeError {
			writeError(w, "Error deleting Volume: induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.DeviceInSGError {
			writeError(w, "Error deleting Volume: induced error - device is a member of a storage group", http.StatusForbidden)
			return
		}
		err := deleteVolume(volID)
		if err != nil {
			writeError(w, "error deleteVolume: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// DeleteVolume - Deletes volume from cache
func DeleteVolume(volID string) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return deleteVolume(volID)
}

func deleteVolume(volID string) error {
	vol, ok := Data.VolumeIDToVolume[volID]
	if ok {
		if vol.NumberOfStorageGroups > 0 {
			return errors.New("Volume present in storage group. Can't be deleted")
		}
		Data.VolumeIDToVolume[volID] = nil
	} else {
		return errors.New("Could not find volume")
	}
	return nil
}

func returnVolume(w http.ResponseWriter, volID string, remote bool) {
	if volID != "" {
		if vol, ok := Data.VolumeIDToVolume[volID]; ok {
			newVol := new(types.Volume)
			err := copier.Copy(newVol, vol)
			if err != nil {
				writeError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("volume: %#v\n", vol)
			if InducedErrors.InvalidLocalVolumeError {
				newVol.StorageGroupIDList = nil
			}
			if remote {
				if InducedErrors.FetchResponseError {
					writeError(w, "Error fetching response", http.StatusBadRequest)
					return
				}
				if InducedErrors.GetRemoteVolumeError {
					writeError(w, "Could not find volume", http.StatusNotFound)
					return
				}
				if InducedErrors.InvalidRemoteVolumeError {
					newVol.StorageGroupIDList = nil
				}
				if !strings.Contains(vol.Type, "RDF") {
					writeError(w, "Could not find volume", http.StatusNotFound)
					return
				}
				newVol.Type = strings.ReplaceAll(newVol.Type, "RDF1", "RDF2")
				newVol.VolumeIdentifier = ""
			}
			writeJSON(w, newVol)
			return
		}
		writeError(w, "Could not find volume: "+volID, http.StatusNotFound)
	}
}

// FreeVolume - handler for free volume job
func FreeVolume(w http.ResponseWriter, param *types.FreeVolumeParam, volID string, executionOption string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	freeVolume(w, param, volID, executionOption)
}

// This returns a job for freeing space in a volume
func freeVolume(w http.ResponseWriter, _ *types.FreeVolumeParam, volID string, executionOption string) {
	if executionOption != types.ExecutionOptionAsynchronous {
		writeError(w, "expected ASYNCHRONOUS", http.StatusBadRequest)
		return
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("sloprovisioning/system/%s/volume/%s", DefaultSymmetrixID, volID)
	if InducedErrors.JobFailedError {
		newMockJob(volID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		newMockJob(volID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
	}
	returnJobByID(w, volID)
}

// RenameVolume - renames volume in cache
func RenameVolume(w http.ResponseWriter, param *types.ModifyVolumeIdentifierParam, volID string, executionOption string, remote bool) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	renameVolume(w, param, volID, executionOption, remote)
}

// This returns the volume itself after renaming
func renameVolume(w http.ResponseWriter, param *types.ModifyVolumeIdentifierParam, volID string, executionOption string, remote bool) {
	if executionOption != types.ExecutionOptionSynchronous {
		writeError(w, "expected SYNCHRONOUS", http.StatusBadRequest)
		return
	}
	Data.VolumeIDToVolume[volID].VolumeIdentifier = param.VolumeIdentifier.IdentifierName
	returnVolume(w, volID, remote)
}

// ModifyMobility modifes the mobility-id-enabled parameter of volume
func ModifyMobility(w http.ResponseWriter, param *types.EnableMobilityIDParam, volID string, executionOption string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	modifyMobility(w, param, volID, executionOption)
}

func modifyMobility(w http.ResponseWriter, param *types.EnableMobilityIDParam, volID string, _ string) {
	if InducedErrors.ModifyMobilityError {
		writeError(w, "Error modifying mobility for volume: induced error", http.StatusRequestTimeout)
		return
	}
	Data.VolumeIDToVolume[volID].MobilityIDEnabled = param.EnableMobilityID
	returnVolume(w, volID, false)
}

// ExpandVolume - Expands volume size in cache
func ExpandVolume(w http.ResponseWriter, param *types.ExpandVolumeParam, volID string, executionOption string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	expandVolume(w, param, volID, executionOption)
}

// This returns the volume itself after expanding the volume's size
func expandVolume(w http.ResponseWriter, param *types.ExpandVolumeParam, volID string, executionOption string) {
	if InducedErrors.ExpandVolumeError {
		writeError(w, "Error expanding volume: induced error", http.StatusRequestTimeout)
		return
	}
	if executionOption != types.ExecutionOptionSynchronous {
		writeError(w, "expected SYNCHRONOUS", http.StatusBadRequest)
		return
	}

	newSize, err := strconv.ParseFloat(param.VolumeAttribute.VolumeSize, 64)
	switch param.VolumeAttribute.CapacityUnit {
	case "MB":
		newSize = newSize * MiB / GiB
	case "TB":
		newSize = newSize * TiB / GiB
	case "PB":
		newSize = newSize * PiB / GiB
	case "GB":
	}

	if err == nil {
		Data.VolumeIDToVolume[volID].CapacityGB = newSize
	} else {
		writeError(w, fmt.Sprintf("Could not convert expand size parameter in request (%s)", param.VolumeAttribute.VolumeSize), http.StatusBadRequest)
		return
	}
	returnVolume(w, volID, false)
}

// JobInfo is used to simulate a job in Unisphere.
// The first call to read it returns Status as the InitialState.
// Subsequent calls return the Status as the FinalState.
type JobInfo struct {
	Job          types.Job
	InitialState string
	FinalState   string
}

// NewMockJob creates a JobInfo that can be queried
func NewMockJob(jobID string, initialState string, finalState string, resourceLink string) *JobInfo {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return newMockJob(jobID, initialState, finalState, resourceLink)
}

func newMockJob(jobID string, initialState string, finalState string, resourceLink string) *JobInfo {
	job := new(JobInfo)
	job.Job.JobID = jobID
	job.InitialState = initialState
	job.FinalState = finalState
	job.Job.Status = "SCHEDULED"
	job.Job.ResourceLink = resourceLink
	Data.JobIDToMockJob[jobID] = job
	return job
}

func HandleJob(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleJob(w, r)
}

func handleJob(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.GetJobError {
		writeError(w, "Error getting Job(s): induced error", http.StatusRequestTimeout)
		return
	}
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	if jobID == "" {
		queryParams := r.URL.Query()
		// Return a job id list
		jobIDList := new(types.JobIDList)
		jobIDList.JobIDs = make([]string, 0)
		for key := range Data.JobIDToMockJob {
			job := Data.JobIDToMockJob[key].Job
			if queryParams.Get("status") == "" || queryParams.Get("status") == job.Status {
				jobIDList.JobIDs = append(jobIDList.JobIDs, key)
			}
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(jobIDList) // #nosec G20
		return
	}
	// Return a specific job
	if InducedErrors.GetJobCannotFindRoleForUser {
		InducedErrors.GetJobCannotFindRoleForUser = false
		writeError(w, "Cannot find role for user", http.StatusInternalServerError)
		return
	}
	returnJobByID(w, jobID)
}

// ReturnJobByID - Returns job based on ID from mock cache
func ReturnJobByID(w http.ResponseWriter, jobID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnJobByID(w, jobID)
}

func returnJobByID(w http.ResponseWriter, jobID string) {
	job := Data.JobIDToMockJob[jobID]
	if job == nil {
		// Not found
		writeError(w, "Job not found: "+jobID, http.StatusNotFound)
		return
	}
	if job.Job.Status == job.InitialState {
		job.Job.Status = job.FinalState
		job.Job.CompletedDate = time.Now().String()
		job.Job.Result = "Mock job completed"
	} else {
		job.Job.Status = job.InitialState
		job.Job.Result = "Mock job in-progress"
	}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&job.Job)
	if err != nil {
		writeError(w, "json encoding error", http.StatusInternalServerError)
	}
}

// /unixvmax/restapi/common/Iterator/{iterID]/page}
func HandleIterator(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleIterator(w, r)
}

func handleIterator(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		queryParams := r.URL.Query()
		from := queryParams.Get("from")
		to := queryParams.Get("to")
		fmt.Printf("mux iterId %s from %s to %s\n", vars["iterId"], from, to)
		result := &types.VolumeResultList{}
		result.From, err = strconv.Atoi(from)
		if err != nil {
			writeError(w, "bad from query parameter", http.StatusBadRequest)
		}
		result.To, err = strconv.Atoi(to)
		if err != nil {
			writeError(w, "bad from query parameter", http.StatusBadRequest)
		}
		for i := result.From - 1; i < result.To-1; i++ {
			volIDList := types.VolumeIDList{VolumeIDs: Data.VolumeIDIteratorList[i]}
			result.VolumeList = append(result.VolumeList, volIDList)
		}
		if Debug {
			fmt.Printf("volumeResultList: %#v\n", result)
		}
		encoder := json.NewEncoder(w)
		err := encoder.Encode(result)
		if err != nil {
			writeError(w, "volumeResultList json encoding error", http.StatusInternalServerError)
		}
	case http.MethodDelete:
		// Nothing to do, will return
	}
}

func HandleStorageGroupSnapshotPolicy(w http.ResponseWriter, _ *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleStorageGroupSnapshotPolicy(w, nil)
}

func handleStorageGroupSnapshotPolicy(w http.ResponseWriter, _ *http.Request) {
	if InducedErrors.GetStorageGroupSnapshotPolicyError {
		writeError(w, "Error retrieving storage group snapshot policy: induced error", http.StatusRequestTimeout)
		return
	}
	returnJSONFile(Data.JSONDir, "storage_group_snapshot_policy.json", w, nil)
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/storagegroup/{id}
// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/storagegroup
// /univmax/restapi/91/sloprovisioning/symmetrix/{symid}/storagegroup/{id}
// /univmax/restapi/91/sloprovisioning/symmetrix/{symid}/storagegroup
func HandleStorageGroup(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleStorageGroup(w, r)
}

func handleStorageGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sgID := vars["id"]
	apiversion := vars["apiversion"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetStorageGroupError {
			writeError(w, "Error retrieving Storage Group(s): induced error", http.StatusRequestTimeout)
			return
		}
		if vars["symid"] == Data.AsyncRDFGroup.RemoteSymmetrix && strings.Contains(sgID, "rep") {
			returnStorageGroup(w, sgID, true)
		} else {
			returnStorageGroup(w, sgID, false)
		}

	case http.MethodPut:
		if InducedErrors.UpdateStorageGroupError {
			writeError(w, "Error updating Storage Group: induced error", http.StatusRequestTimeout)
			return
		}
		if sgID == "" {
			writeError(w, "storage group ID must be supplied", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		if apiversion == "90" {
			updateSGPayload := &types.UpdateStorageGroupPayload{}
			err := decoder.Decode(updateSGPayload)
			if err != nil {
				writeError(w, "problem decoding PUT StorageGroup payload: "+err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Printf("PUT StorageGroup payload: %#v\n", updateSGPayload)
			editPayload := updateSGPayload.EditStorageGroupActionParam
			if editPayload.ExpandStorageGroupParam != nil {
				expandPayload := editPayload.ExpandStorageGroupParam
				addVolumeParam := expandPayload.AddVolumeParam
				if addVolumeParam != nil {
					name := addVolumeParam.VolumeIdentifier.IdentifierName
					var size string
					if len(addVolumeParam.VolumeAttributes) > 0 {
						size = addVolumeParam.VolumeAttributes[0].VolumeSize
					}
					size = addVolumeParam.VolumeAttributes[0].VolumeSize
					addVolumeToStorageGroupTest(w, name, size, sgID)
				}
				addSpecificVolumeParam := expandPayload.AddSpecificVolumeParam
				if addSpecificVolumeParam != nil {
					addSpecificVolumeToStorageGroup(w, addSpecificVolumeParam.VolumeIDs, sgID)
				}
			}
			if editPayload.RemoveVolumeParam != nil {
				removeVolumeFromStorageGroup(w, editPayload.RemoveVolumeParam.VolumeIDs, sgID)
			}
		} else {
			// for apiVersion 91
			updateSGPayload := &types.UpdateStorageGroupPayload{}
			err := decoder.Decode(updateSGPayload)
			if err != nil {
				writeError(w, "problem decoding PUT StorageGroup payload: "+err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Printf("PUT StorageGroup payload: %#v\n", updateSGPayload)
			editPayload := updateSGPayload.EditStorageGroupActionParam
			if editPayload.ExpandStorageGroupParam != nil {
				expandPayload := editPayload.ExpandStorageGroupParam
				addVolumeParam := expandPayload.AddVolumeParam
				if addVolumeParam != nil {
					name := addVolumeParam.VolumeAttributes[0].VolumeIdentifier.IdentifierName
					size := addVolumeParam.VolumeAttributes[0].VolumeSize
					addVolumeToStorageGroupTest(w, name, size, sgID)
				}
				addSpecificVolumeParam := expandPayload.AddSpecificVolumeParam
				if addSpecificVolumeParam != nil {
					addSpecificVolumeToStorageGroup(w, addSpecificVolumeParam.VolumeIDs, sgID)
				}
			}
			if editPayload.RemoveVolumeParam != nil {
				removeVolumeFromStorageGroup(w, editPayload.RemoveVolumeParam.VolumeIDs, sgID)
			}
		}
	case http.MethodPost:
		if InducedErrors.CreateStorageGroupError {
			writeError(w, "Error creating Storage Group: induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.StorageGroupAlreadyExists {
			writeError(w, "The requested storage group resource already exists", http.StatusConflict)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createSGPayload := &types.CreateStorageGroupParam{}
		err := decoder.Decode(createSGPayload)
		if err != nil {
			writeError(w, "problem decoding POST StorageGroup payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("POST StorageGroup payload: %#v\n", createSGPayload)
		sgID := createSGPayload.StorageGroupID
		// Data.StorageGroupIDToNVolumes[sgID] = 0
		// fmt.Println("SG Name: ", sgID)
		addStorageGroupFromCreateParams(createSGPayload)
		if vars["symid"] == Data.AsyncRDFGroup.RemoteSymmetrix {
			returnStorageGroup(w, sgID, true)
		} else {
			returnStorageGroup(w, sgID, false)
		}

	case http.MethodDelete:
		if InducedErrors.DeleteStorageGroupError {
			writeError(w, "Error deleting storage group: induced error", http.StatusRequestTimeout)
			return
		}
		removeStorageGroup(w, sgID)

	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/maskingview/{id}/connections
func HandleMaskingViewConnections(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleMaskingViewConnections(w, r)
}

func handleMaskingViewConnections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		queryParams := r.URL.Query()
		volID := queryParams.Get("volume_id")
		if InducedErrors.GetMaskingViewConnectionsError {
			writeError(w, "Error retrieving Masking View Connections: induced error", http.StatusRequestTimeout)
			return
		}
		if volID == "" {
			// Return a response for all volumes
			index := 1
			result := &types.MaskingViewConnectionsResult{
				MaskingViewConnections: make([]*types.MaskingViewConnection, 0),
			}
			for id := range Data.VolumeIDToVolume {
				conn1 := &types.MaskingViewConnection{
					VolumeID:       id,
					HostLUNAddress: fmt.Sprintf("%4d", index),
					CapacityGB:     "0.1",
					InitiatorID:    "iqn.1993-08.org.debian:01:8f21cc8ad2a7",
					DirectorPort:   "SE-1E:000",
					LoggedIn:       false,
					OnFabric:       true,
				}
				result.MaskingViewConnections = append(result.MaskingViewConnections, conn1)
				conn2 := &types.MaskingViewConnection{
					VolumeID:       id,
					HostLUNAddress: fmt.Sprintf("%4d", index),
					CapacityGB:     "0.1",
					InitiatorID:    "iqn.1993-08.org.debian:01:8f21cc8ad2a7",
					DirectorPort:   "SE-2E:000",
					LoggedIn:       false,
					OnFabric:       true,
				}
				result.MaskingViewConnections = append(result.MaskingViewConnections, conn2)
				index++
			}
			writeJSON(w, result)
			return
		}
		replacements := make(map[string]string)
		replacements["__VOLUME_ID__"] = volID
		returnJSONFile(Data.JSONDir, "masking_view_connections_template.json", w, replacements)
	}
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/maskingview/{id}
// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/maskingview
func HandleMaskingView(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleMaskingView(w, r)
}

func handleMaskingView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mvID := vars["mvID"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetMaskingViewError {
			writeError(w, "Error retrieving Masking View(s): induced error", http.StatusRequestTimeout)
			return
		}
		returnMaskingView(w, mvID)

	case http.MethodPost:
		if InducedErrors.CreateMaskingViewError {
			writeError(w, "Failed to create masking view: induced error", http.StatusRequestTimeout)
			return
		} else if InducedErrors.MaskingViewAlreadyExists {
			writeError(w, "The requested masking view resource already exists", http.StatusConflict)
			return
		} else if InducedErrors.PortGroupNotFoundError {
			writeError(w, "Port Group on Symmetrix cannot be found", http.StatusInternalServerError)
		} else if InducedErrors.InitiatorGroupNotFoundError {
			writeError(w, "Initiator Group on Symmetrix cannot be found", http.StatusInternalServerError)
		} else if InducedErrors.StorageGroupNotFoundError {
			writeError(w, "Storage Group on Symmetrix cannot be found", http.StatusInternalServerError)
		}
		decoder := json.NewDecoder(r.Body)
		createMVPayload := &types.MaskingViewCreateParam{}
		err := decoder.Decode(createMVPayload)
		if err != nil {
			writeError(w, "problem decoding POST Masking View payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("POST MaskingView payload: %#v\n", createMVPayload)
		mvID := createMVPayload.MaskingViewID
		// Data.StorageGroupIDToNVolumes[sgID] = 0
		fmt.Println("MV Name: ", mvID)
		addMaskingViewFromCreateParams(createMVPayload)
		returnMaskingView(w, mvID)

	case http.MethodPut:
		if InducedErrors.UpdateMaskingViewError {
			writeError(w, "Error updating Masking View: induced error", http.StatusRequestTimeout)
			return
		}
		// if mvID == "" {
		// 	writeError(w, "Masking View ID must be supplied", http.StatusBadRequest)
		// 	return
		// }
		decoder := json.NewDecoder(r.Body)
		updateMaskingViewPayload := &types.EditMaskingViewParam{}
		err := decoder.Decode(updateMaskingViewPayload)
		if err != nil {
			writeError(w, "problem decoding PUT Masking View payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("PUT masking view payload: %#v\n", updateMaskingViewPayload)
		executionOption := updateMaskingViewPayload.ExecutionOption
		if &updateMaskingViewPayload.EditMaskingViewActionParam.RenameMaskingViewParam != nil {
			renameMaskingView(w, &updateMaskingViewPayload.EditMaskingViewActionParam.RenameMaskingViewParam, mvID, executionOption)
			return
		}

	case http.MethodDelete:
		if InducedErrors.DeleteMaskingViewError {
			writeError(w, "Error deleting Masking view: induced error", http.StatusRequestTimeout)
			return
		}
		removeMaskingView(w, mvID)

	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

func newStorageGroup(storageGroupID string, maskingViewID string, storageResourcePoolID string,
	serviceLevel string, numOfVolumes int,
) {
	numOfMaskingViews := 0
	if maskingViewID != "" {
		numOfMaskingViews = 1
	}
	childStorageGroups := []string{}
	maskingViews := make([]string, 0)
	if maskingViewID != "" {
		maskingViews = append(maskingViews, maskingViewID)
	}
	storageGroup := &types.StorageGroup{
		StorageGroupID:    storageGroupID,
		SLO:               serviceLevel,
		SRP:               storageResourcePoolID,
		Workload:          "None",
		SLOCompliance:     "STABLE",
		NumOfVolumes:      numOfVolumes,
		NumOfChildSGs:     0,
		NumOfParentSGs:    0,
		NumOfMaskingViews: numOfMaskingViews,
		NumOfSnapshots:    0,
		CapacityGB:        234.5,
		DeviceEmulation:   "FBA",
		Type:              "Standalone",
		Unprotected:       true,
		ChildStorageGroup: childStorageGroups,
		MaskingView:       maskingViews,
	}
	Data.StorageGroupIDToStorageGroup[storageGroupID] = storageGroup
	volumes := make([]string, 0)
	Data.StorageGroupIDToVolumes[storageGroupID] = volumes
}

func newMaskingView(maskingViewID string, storageGroupID string, hostID string, portGroupID string) {
	maskingView := &types.MaskingView{
		MaskingViewID:  maskingViewID,
		HostID:         hostID,
		HostGroupID:    "",
		PortGroupID:    portGroupID,
		StorageGroupID: storageGroupID,
	}
	Data.MaskingViewIDToMaskingView[maskingViewID] = maskingView
}

func AddStorageGroup(storageGroupID string, storageResourcePoolID string,
	serviceLevel string,
) (*types.StorageGroup, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addStorageGroup(storageGroupID, storageResourcePoolID, serviceLevel)
}

// AddStorageGroup - Adds a storage group to the mock data cache
func addStorageGroup(storageGroupID string, storageResourcePoolID string,
	serviceLevel string,
) (*types.StorageGroup, error) {
	if _, ok := Data.StorageGroupIDToStorageGroup[storageGroupID]; ok {
		return nil, errors.New("The requested storage group resource already exists")
	}
	newStorageGroup(storageGroupID, "", storageResourcePoolID, serviceLevel, 0)
	return Data.StorageGroupIDToStorageGroup[storageGroupID], nil
}

func AddRDFStorageGroup(storageGroupID, symmetrixID string) (*types.RDFStorageGroup, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addRDFStorageGroup(storageGroupID, symmetrixID)
}

// AddRDFStorageGroup ...
func addRDFStorageGroup(storageGroupID, symmetrixID string) (*types.RDFStorageGroup, error) {
	if _, ok := Data.StorageGroupIDToRDFStorageGroup[storageGroupID]; ok {
		return nil, fmt.Errorf("rdfStorageGroup already exists")
	}
	Data.StorageGroupIDToStorageGroup[storageGroupID].Unprotected = false
	rdfSG := &types.RDFStorageGroup{
		Name:        storageGroupID,
		SymmetrixID: symmetrixID,
		Rdf:         true,
	}
	Data.StorageGroupIDToRDFStorageGroup[storageGroupID] = rdfSG
	return rdfSG, nil
}

// RemoveStorageGroup - Removes a storage group from the mock data cache
func RemoveStorageGroup(w http.ResponseWriter, storageGroupID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeStorageGroup(w, storageGroupID)
}

func removeStorageGroup(w http.ResponseWriter, storageGroupID string) {
	if InducedErrors.GetSGOnRemote {
		storageGroupID = "CSI-Test-Fake-Remote-SG"
	}
	sg, ok := Data.StorageGroupIDToStorageGroup[storageGroupID]
	if !ok {
		fmt.Println("Storage Group " + storageGroupID + " doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if sg.NumOfMaskingViews != 0 {
		fmt.Println("Can't delete a storage group which is part of masking view")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	volumes := Data.StorageGroupIDToVolumes[storageGroupID]
	if InducedErrors.RemoveVolumesFromSG {
		volumes = nil
	}
	if len(volumes) > 0 {
		fmt.Println("Can't delete a storage group which has volumes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	delete(Data.StorageGroupIDToStorageGroup, storageGroupID)
	delete(Data.StorageGroupIDToStorageGroup, storageGroupID+"-remote")
	delete(Data.StorageGroupIDToRDFStorageGroup, storageGroupID)
	delete(Data.StorageGroupIDToRDFStorageGroup, storageGroupID+"-remote")
}

func addMaskingViewFromCreateParams(createParams *types.MaskingViewCreateParam) {
	mvID := createParams.MaskingViewID
	hostID := ""
	hostGroupID := ""
	if createParams.HostOrHostGroupSelection.UseExistingHostParam != nil {
		hostID = createParams.HostOrHostGroupSelection.UseExistingHostParam.HostID
	} else if createParams.HostOrHostGroupSelection.UseExistingHostGroupParam != nil {
		hostGroupID = createParams.HostOrHostGroupSelection.UseExistingHostGroupParam.HostGroupID
	}
	portGroupID := createParams.PortGroupSelection.UseExistingPortGroupParam.PortGroupID
	sgID := createParams.StorageGroupSelection.UseExistingStorageGroupParam.StorageGroupID
	if hostID != "" {
		addMaskingView(mvID, sgID, hostID, portGroupID) // #nosec G20
	} else if hostGroupID != "" {
		addMaskingView(mvID, sgID, hostGroupID, portGroupID) // #nosec G20
	}
}

// AddMaskingView - Adds a masking view to the mock data cache
func AddMaskingView(maskingViewID string, storageGroupID string, hostID string, portGroupID string) (*types.MaskingView, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addMaskingView(maskingViewID, storageGroupID, hostID, portGroupID)
}

func addMaskingView(maskingViewID string, storageGroupID string, hostID string, portGroupID string) (*types.MaskingView, error) {
	if _, ok := Data.MaskingViewIDToMaskingView[maskingViewID]; ok {
		return nil, errors.New("Error! Masking View already exists")
	}
	if _, ok := Data.StorageGroupIDToStorageGroup[storageGroupID]; !ok {
		return nil, errors.New("Storage Group doesn't exist")
	}
	/*if _, ok := Data.PortGroupIDToPortGroup[portGroupID]; !ok {
		return errors.New("Port Group doesn't exist")
	}*/
	if _, ok := Data.HostIDToHost[hostID]; !ok {
		return nil, errors.New("Host doesn't exist")
	}
	newMaskingView(maskingViewID, storageGroupID, hostID, portGroupID)
	// Update host
	Data.HostIDToHost[hostID].MaskingviewIDs = append(Data.HostIDToHost[hostID].MaskingviewIDs, maskingViewID)
	Data.HostIDToHost[hostID].NumberMaskingViews++
	// Update Storage Group
	currentMaskingViewIDs := Data.StorageGroupIDToStorageGroup[storageGroupID].MaskingView
	Data.StorageGroupIDToStorageGroup[storageGroupID].MaskingView = append(
		currentMaskingViewIDs, maskingViewID)
	Data.StorageGroupIDToStorageGroup[storageGroupID].NumOfMaskingViews++
	// Update the volume cache
	for _, volumeID := range Data.StorageGroupIDToVolumes[storageGroupID] {
		Data.VolumeIDToVolume[volumeID].NumberOfFrontEndPaths = 1
	}
	return Data.MaskingViewIDToMaskingView[maskingViewID], nil
}

// RenameMaskingView - Renames masking view
func RenameMaskingView(w http.ResponseWriter, param *types.RenameMaskingViewParam, maskingViewID string, executionOption string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	renameMaskingView(w, param, maskingViewID, executionOption)
}

func renameMaskingView(w http.ResponseWriter, param *types.RenameMaskingViewParam, maskingViewID string, executionOption string) {
	if executionOption != types.ExecutionOptionSynchronous {
		writeError(w, "expected SYNCHRONOUS", http.StatusBadRequest)
		return
	}
	Data.MaskingViewIDToMaskingView[maskingViewID].MaskingViewID = param.NewMaskingViewName
	returnMaskingView(w, maskingViewID)
}

// RemoveMaskingView - Removes a masking view from the mock data cache
func RemoveMaskingView(w http.ResponseWriter, maskingViewID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeMaskingView(w, maskingViewID)
}

func removeMaskingView(w http.ResponseWriter, maskingViewID string) {
	mv, ok := Data.MaskingViewIDToMaskingView[maskingViewID]
	if !ok {
		fmt.Println("Masking View " + maskingViewID + " doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Handle storage groups
	storageGroupID := mv.StorageGroupID
	Data.StorageGroupIDToStorageGroup[storageGroupID].NumOfMaskingViews--
	currentMaskingViewIDs := Data.StorageGroupIDToStorageGroup[storageGroupID].MaskingView
	newMaskingViewIDs := make([]string, 0)
	for _, mvID := range currentMaskingViewIDs {
		if mvID != maskingViewID {
			newMaskingViewIDs = append(newMaskingViewIDs, mvID)
		}
	}
	Data.StorageGroupIDToStorageGroup[storageGroupID].MaskingView = newMaskingViewIDs
	// Handle Hosts
	hostID := mv.HostID
	Data.HostIDToHost[hostID].NumberMaskingViews--
	currentMaskingViewIDs = Data.HostIDToHost[hostID].MaskingviewIDs
	newMaskingViewIDs = make([]string, 0)
	for _, mvID := range currentMaskingViewIDs {
		if mvID != maskingViewID {
			newMaskingViewIDs = append(newMaskingViewIDs, mvID)
		}
	}
	Data.HostIDToHost[hostID].MaskingviewIDs = newMaskingViewIDs
	// Check if we need to update the number of front end paths for volumes
	// Loop through volumes of this particular SG
	if volumeIDs, ok := Data.StorageGroupIDToVolumes[storageGroupID]; ok {
		// First construct a list of all SGs
		tempSGList := make([]string, 0)
		for _, volumeID := range volumeIDs {
			if vol, ok1 := Data.VolumeIDToVolume[volumeID]; ok1 {
				tempSGList = append(tempSGList, vol.StorageGroupIDList...)
			}
		}
		// Remove duplicates
		tempSGList = uniqueElements(tempSGList)
		// Filter out SGs in masking Views
		sgIDsInMaskingView := make([]string, 0)
		for _, sgID := range tempSGList {
			if sg, ok1 := Data.StorageGroupIDToStorageGroup[sgID]; ok1 {
				if sg.NumOfMaskingViews > 0 {
					sgIDsInMaskingView = append(sgIDsInMaskingView, sgID)
				}
			}
		}
		// Now Update the number of front end paths
		for _, volumeID := range volumeIDs {
			if vol, ok1 := Data.VolumeIDToVolume[volumeID]; ok1 {
				update := compareAndCheck(vol.StorageGroupIDList, sgIDsInMaskingView)
				if update {
					vol.NumberOfFrontEndPaths = 0
				}
			}
		}
	}
	delete(Data.StorageGroupIDToStorageGroup, maskingViewID)
}

// compareAndCheck - compares two string slices and returns true if the slices are equal or false if they aren't
func compareAndCheck(slice1 []string, slice2 []string) bool {
	for _, item := range slice1 {
		for _, item1 := range slice2 {
			if item == item1 {
				return false
			}
		}
	}
	return true
}

// uniqueElements - Removes duplicates from a string slice and returns a slice containing unique elements only
func uniqueElements(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// newVolume creates a new mock volume with the specified characteristics.
func newVolume(volumeID, volumeIdentifier string, size int, sgList []string) {
	volume := &types.Volume{
		VolumeID:              volumeID,
		Type:                  "TDEV",
		Emulation:             "FBA",
		SSID:                  "FFFFFFFF",
		AllocatedPercent:      0,
		CapacityGB:            0.0,
		FloatCapacityMB:       0.0,
		CapacityCYL:           size,
		Status:                "Ready",
		Reserved:              false,
		Pinned:                false,
		VolumeIdentifier:      volumeIdentifier,
		WWN:                   "600009700001979000465330303" + volumeID,
		EffectiveWWN:          "600009700001979000465330303" + volumeID,
		NGUID:                 "600009700001979000465330303" + volumeID,
		Encapsulated:          false,
		NumberOfStorageGroups: 1,
		NumberOfFrontEndPaths: 0,
		StorageGroupIDList:    sgList,
	}
	if _, ok := Data.StorageGroupIDToRDFStorageGroup[sgList[0]]; ok {
		volume.Type = "RDF1+TDEV"
		volume.RDFGroupIDList = []types.RDFGroupID{
			{RDFGroupNumber: Data.AsyncRDFGroup.RdfgNumber},
		}
	}
	Data.VolumeIDToVolume[volumeID] = volume
}

// AddNewVolume - Add a volume to the mock data cache
func AddNewVolume(volumeID, volumeIdentifier string, size int, storageGroupID string) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addNewVolume(volumeID, volumeIdentifier, size, storageGroupID)
}

func addNewVolume(volumeID, volumeIdentifier string, size int, storageGroupID string) error {
	if _, ok := Data.VolumeIDToVolume[volumeID]; ok {
		return errors.New("The requested volume already exists")
	}
	if _, ok := Data.StorageGroupIDToStorageGroup[storageGroupID]; !ok {
		return errors.New("The requested storage group resource doesn't exist")
	}
	sgList := []string{storageGroupID}
	newVolume(volumeID, volumeIdentifier, size, sgList)
	Data.StorageGroupIDToStorageGroup[storageGroupID].NumOfVolumes++
	currentVolumes := Data.StorageGroupIDToVolumes[storageGroupID]
	newVolumes := append(currentVolumes, volumeID)
	Data.StorageGroupIDToVolumes[storageGroupID] = newVolumes
	return nil
}

func newInitiator(initiatorID string, initiatorName string, initiatorType string, dirPortKeys []types.PortKey, hostID string) {
	// maskingViewIDs := []string{}
	initiator := &types.Initiator{
		InitiatorID:          initiatorName,
		SymmetrixPortKey:     dirPortKeys,
		InitiatorType:        initiatorType,
		FCID:                 "0",
		IPAddress:            "192.168.1.175",
		Host:                 hostID,
		HostID:               hostID,
		HostGroupIDs:         []string{},
		LoggedIn:             true,
		OnFabric:             true,
		FlagsInEffect:        "Common_Serial_Number(C), SCSI_3(SC3), SPC2_Protocol_Version(SPC2)",
		NumberVols:           1,
		NumberHostGroups:     0,
		NumberMaskingViews:   0,
		NumberPowerPathHosts: 0,
	}
	Data.InitiatorIDToInitiator[initiatorID] = initiator
}

func AddInitiator(initiatorID string, initiatorName string, initiatorType string, dirPortKeys []string, hostID string) (*types.Initiator, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addInitiator(initiatorID, initiatorName, initiatorType, dirPortKeys, hostID)
}

// AddInitiator - Adds an initiator to the mock data cache
func addInitiator(initiatorID string, initiatorName string, initiatorType string, dirPortKeys []string, hostID string) (*types.Initiator, error) {
	if _, ok := Data.InitiatorIDToInitiator[initiatorID]; ok {
		return nil, errors.New("Error! Initiator already exists")
	}
	// if host id is supplied, check for existence of host
	if hostID != "" {
		if _, ok := Data.HostIDToHost[hostID]; !ok {
			return nil, errors.New("Error! Host doesn't exist")
		}
	}
	portKeys := make([]types.PortKey, 0)
	for _, dirPortKey := range dirPortKeys {
		dirPortDetails := strings.Split(dirPortKey, ":")
		portKey := types.PortKey{
			DirectorID: dirPortDetails[0],
			PortID:     dirPortKey,
		}
		portKeys = append(portKeys, portKey)
	}
	newInitiator(initiatorID, initiatorName, initiatorType, portKeys, hostID)
	return Data.InitiatorIDToInitiator[initiatorID], nil
}

// ReturnInitiator - Returns initiator from mock cache based on initiator id
func ReturnInitiator(w http.ResponseWriter, initiatorID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnInitiator(w, initiatorID)
}

func returnInitiator(w http.ResponseWriter, initiatorID string) {
	if initiatorID != "" {
		if init, ok := Data.InitiatorIDToInitiator[initiatorID]; ok {
			writeJSON(w, init)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		initIDs := make([]string, 0)
		for k := range Data.InitiatorIDToInitiator {
			initIDs = append(initIDs, k)
		}
		initiatorIDList := &types.InitiatorList{
			InitiatorIDs: initIDs,
		}
		writeJSON(w, initiatorIDList)
	}
}

func newHost(hostID string, hostType string, initiatorIDs []string) {
	maskingViewIDs := []string{}
	host := &types.Host{
		HostID:             hostID,
		NumberMaskingViews: 0,
		NumberInitiators:   int64(len(initiatorIDs)),
		NumberHostGroups:   0,
		PortFlagsOverride:  false,
		ConsistentLun:      false,
		EnabledFlags:       "",
		DisabledFlags:      "",
		HostType:           hostType,
		Initiators:         initiatorIDs,
		MaskingviewIDs:     maskingViewIDs,
		NumPowerPathHosts:  0,
	}
	Data.HostIDToHost[hostID] = host
}

func AddHost(hostID string, hostType string, initiatorIDs []string) (*types.Host, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addHost(hostID, hostType, initiatorIDs)
}

// AddHost - Adds a host to the mock data cache
func addHost(hostID string, hostType string, initiatorIDs []string) (*types.Host, error) {
	if _, ok := Data.HostIDToHost[hostID]; ok {
		return nil, errors.New("Error! Host already exists")
	}
	validInitiators := false
	// Check if initiators exist
	for _, initID := range initiatorIDs {
		for _, v := range Data.InitiatorIDToInitiator {
			if v.InitiatorID == initID {
				if v.HostID == "" {
					validInitiators = true
					break
				}
			}
		}
		if !validInitiators {
			break
		}
	}
	if !validInitiators {
		errormsg := errors.New("error: Some initiators don't exist or are not valid")
		fmt.Println(errormsg)
		return nil, errormsg
	}
	newHost(hostID, hostType, initiatorIDs)
	// Update the initiators
	for _, initID := range initiatorIDs {
		for k, v := range Data.InitiatorIDToInitiator {
			if v.InitiatorID == initID {
				Data.InitiatorIDToInitiator[k].HostID = hostID
				Data.InitiatorIDToInitiator[k].Host = hostID
				break
			}
		}
	}
	fmt.Println(Data.HostIDToHost[hostID])
	return Data.HostIDToHost[hostID], nil
}

// RemoveHost - Removes host from mock cache
func RemoveHost(hostID string) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return removeHost(hostID)
}

// removeHost - Remove a host from the mock data cache
func removeHost(hostID string) error {
	host, ok := Data.HostIDToHost[hostID]
	if !ok {
		return errors.New("Error! Host doesn't exist")
	}
	if host.NumberMaskingViews > 0 {
		return errors.New("Error! Host is part of a masking view")
	}
	Data.HostIDToHost[hostID] = nil
	return nil
}

func newPortGroup(portGroupID string, portGroupType string, portKeys []types.PortKey) {
	portGroup := &types.PortGroup{
		PortGroupID:        portGroupID,
		SymmetrixPortKey:   portKeys,
		NumberPorts:        int64(len(portKeys)),
		NumberMaskingViews: 0,
		PortGroupType:      portGroupType,
	}
	Data.PortGroupIDToPortGroup[portGroupID] = portGroup
}

// addPortGroup - Adds a port group to the mock data cache
func addPortGroup(portGroupID string, portGroupType string, portKeys []types.PortKey) (*types.PortGroup, error) {
	if _, ok := Data.PortGroupIDToPortGroup[portGroupID]; ok {
		return nil, errors.New("Error! Port Group already exists")
	}
	newPortGroup(portGroupID, portGroupType, portKeys)
	return Data.PortGroupIDToPortGroup[portGroupID], nil
}

// updatePortGroup - Update PortGroup by ID by adding 'addKeys' and removing 'removeKeys'
func updatePortGroup(portGroupID string, editPayload *types.EditPortGroupActionParam) (*types.PortGroup, error) {
	pg, ok := Data.PortGroupIDToPortGroup[portGroupID]
	if !ok {
		return nil, fmt.Errorf("error! PortGroup %s does not exist", portGroupID)
	}

	// Collect the ports to add (if any)
	addKeys := make([]types.PortKey, 0)
	if editPayload.AddPortParam != nil {
		addKeys = convertToPortKeys(editPayload.AddPortParam.Ports)
	}

	// Collect the ports to remove (if any)
	removeKeys := make([]types.PortKey, 0)
	if editPayload.RemovePortParam != nil {
		removeKeys = convertToPortKeys(editPayload.RemovePortParam.Ports)
	}

	// Add to the list of ports
	pg.SymmetrixPortKey = append(pg.SymmetrixPortKey, addKeys...)

	// Remove from the list of ports in the PortGroup
	for _, key := range removeKeys {
		pg.SymmetrixPortKey = removePortKey(pg.SymmetrixPortKey, key)
	}

	if editPayload.RenamePortGroupParam != nil && editPayload.RenamePortGroupParam.NewPortGroupName != "" {
		portGroupID = editPayload.RenamePortGroupParam.NewPortGroupName
		pg.PortGroupID = portGroupID
	}

	// Update the PortGroup mapping with the update PortGroup
	Data.PortGroupIDToPortGroup[portGroupID] = pg
	return pg, nil
}

// convertToPortKeys - Convert a slice of types.SymmetrixPortKeyType to slice of types.PortKey
func convertToPortKeys(symmPorts []types.SymmetrixPortKeyType) []types.PortKey {
	if symmPorts == nil || len(symmPorts) == 0 {
		return make([]types.PortKey, 0)
	}

	out := make([]types.PortKey, len(symmPorts))
	for idx, it := range symmPorts {
		out[idx] = types.PortKey{
			DirectorID: it.DirectorID,
			PortID:     it.PortID,
		}
	}

	return out
}

// removePortKey - delete PortKey 'key' from the slice
func removePortKey(slice []types.PortKey, keyToRemove types.PortKey) []types.PortKey {
	index := -1
	// Find the index in the slice that has the match
	for it, thisKey := range slice {
		if thisKey.DirectorID == keyToRemove.DirectorID && thisKey.PortID == keyToRemove.PortID {
			index = it
			break
		}
	}
	if index != -1 {
		// Found the index with matching port
		copy(slice[index:], slice[index+1:])
		return slice[:len(slice)-1]
	}
	// No match was found, return unchanged slice
	return slice
}

func UpdatePortGroupFromParams(portGroupID string, updateParams *types.EditPortGroup) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	updatePortGroupFromParams(portGroupID, updateParams)
}

// UpdatePortGroupFromParams - Updates PortGroup given an EditPortGroup payload
func updatePortGroupFromParams(portGroupID string, updateParams *types.EditPortGroup) {
	updatePortGroup(portGroupID, updateParams.EditPortGroupActionParam) // #nosec G20
}

func DeletePortGroup(portGroupID string) (*types.PortGroup, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return deletePortGroup(portGroupID)
}

// DeletePortGroup - Remove PortGroup by ID 'portGroupID'
func deletePortGroup(portGroupID string) (*types.PortGroup, error) {
	pg, ok := Data.PortGroupIDToPortGroup[portGroupID]
	if !ok {
		return nil, fmt.Errorf("error! PortGroup %s does not exist", portGroupID)
	}
	delete(Data.PortGroupIDToPortGroup, portGroupID)
	return pg, nil
}

func AddPortGroupFromCreateParams(createParams *types.CreatePortGroupParams) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addPortGroupFromCreateParams(createParams)
}

// AddPortGroupFromCreateParams - Adds a storage group from create params
func addPortGroupFromCreateParams(createParams *types.CreatePortGroupParams) {
	portGroupID := createParams.PortGroupID
	portKeys := createParams.SymmetrixPortKey
	addPortGroup(portGroupID, "Fibre", portKeys) // #nosec G20
}

func AddPortGroupWithPortID(portGroupID string, portGroupType string, portIdentifiers []string) (*types.PortGroup, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addPortGroupWithPortID(portGroupID, portGroupType, portIdentifiers)
}

// AddPortGroupWithPortID - Adds a port group to the mock data cache
func addPortGroupWithPortID(portGroupID string, portGroupType string, portIdentifiers []string) (*types.PortGroup, error) {
	portKeys := make([]types.PortKey, 0)
	for _, dirPortKey := range portIdentifiers {
		dirPortDetails := strings.Split(dirPortKey, ":")
		if len(dirPortDetails) != 2 {
			errormsg := fmt.Errorf("invalid dir port specified: %s", dirPortKey)
			log.Error(errormsg)
			return nil, errormsg
		}
		portKey := types.PortKey{
			DirectorID: dirPortDetails[0],
			PortID:     dirPortKey,
		}
		portKeys = append(portKeys, portKey)
	}
	if _, ok := Data.PortGroupIDToPortGroup[portGroupID]; ok {
		return nil, errors.New("Error! Port Group already exists")
	}
	newPortGroup(portGroupID, portGroupType, portKeys)
	return Data.PortGroupIDToPortGroup[portGroupID], nil
}

func AddStorageGroupFromCreateParams(createParams *types.CreateStorageGroupParam) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addStorageGroupFromCreateParams(createParams)
}

// AddStorageGroupFromCreateParams - Adds a storage group from create params
func addStorageGroupFromCreateParams(createParams *types.CreateStorageGroupParam) {
	sgID := createParams.StorageGroupID
	srpID := createParams.SRPID
	serviceLevel := "None"
	if srpID != "None" {
		sloBasedParams := createParams.SLOBasedStorageGroupParam
		serviceLevel = sloBasedParams[0].SLOID
	} else {
		srpID = ""
	}
	addStorageGroup(sgID, srpID, serviceLevel) // #nosec G20
}

// keys - Return keys of the given map
func keys(m map[string]*types.StorageGroup) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ReturnStorageGroup - Returns storage group information from mock cache
func ReturnStorageGroup(w http.ResponseWriter, sgID string, remote bool) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnStorageGroup(w, sgID, remote)
}

func returnStorageGroup(w http.ResponseWriter, sgID string, remote bool) {
	if sgID != "" {
		if InducedErrors.GetSGOnRemote && remote {
			sg := Data.StorageGroupIDToStorageGroup["CSI-Test-Fake-Remote-SG"]
			fmt.Printf("Fake remote SG: %#v\n", sg)
			writeJSON(w, sg)
			return
		}
		if InducedErrors.GetSGWithVolOnRemote && remote {
			sg := Data.StorageGroupIDToStorageGroup["CSI-Test-Fake-Remote-SG"]
			sg.NumOfVolumes = 1
			fmt.Printf("Fake remote SG: %#v\n", sg)
			writeJSON(w, sg)
			return
		}
		if _, ok := Data.StorageGroupIDToRDFStorageGroup[sgID]; remote && !ok {
			writeError(w, "StorageGroup not found", http.StatusNotFound)
			return
		}
		if sg, ok := Data.StorageGroupIDToStorageGroup[sgID]; ok {
			fmt.Printf("SG: %#v\n", sg)
			writeJSON(w, sg)
			return
		}
		writeError(w, "StorageGroup not found", http.StatusNotFound)
	} else {
		storageGroupIDs := keys(Data.StorageGroupIDToStorageGroup)
		storageGroupIDList := &types.StorageGroupIDList{
			StorageGroupIDs: storageGroupIDs,
		}
		writeJSON(w, storageGroupIDList)
	}
}

func returnMaskingView(w http.ResponseWriter, mvID string) {
	if mvID != "" {
		if mv, ok := Data.MaskingViewIDToMaskingView[mvID]; ok {
			writeJSON(w, mv)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		maskingViewIDs := make([]string, 0)
		for k := range Data.MaskingViewIDToMaskingView {
			maskingViewIDs = append(maskingViewIDs, k)
		}
		maskingViewIDList := &types.MaskingViewList{
			MaskingViewIDs: maskingViewIDs,
		}
		writeJSON(w, maskingViewIDList)
	}
}

func writeJSON(w http.ResponseWriter, val interface{}) {
	if InducedErrors.InvalidResponse {
		fmt.Println("Inducing error")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		fmt.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Couldn't write to ResponseWriter")
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// AddOneVolumeToStorageGroup - Adds volume to a storage group in the mock cache
func AddOneVolumeToStorageGroup(volumeID, volumeIdentifier, sgID string, size int) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addOneVolumeToStorageGroup(volumeID, volumeIdentifier, sgID, size)
}

func addOneVolumeToStorageGroup(volumeID, volumeIdentifier, sgID string, size int) error {
	if _, ok := Data.StorageGroupIDToStorageGroup[sgID]; !ok {
		return errors.New("The requested storage group doesn't exist")
	}
	if _, ok := Data.VolumeIDToVolume[volumeID]; ok {
		// Found the volume in cache
		// We are adding it to another storage group
		if volumes, ok := Data.StorageGroupIDToVolumes[sgID]; ok {
			found := false
			for _, volume := range volumes {
				if strings.Contains(volume, volumeID) {
					found = true
					break
				}
			}
			if found {
				return errors.New("Volume is already a part of the SG")
			}
			// Update the volume cache
			currentStorageGroups := Data.VolumeIDToVolume[volumeID].StorageGroupIDList
			newStorageGroups := append(currentStorageGroups, sgID)
			Data.VolumeIDToVolume[volumeID].StorageGroupIDList = newStorageGroups
			// Update volume's replication details in case the storage-group is replicated
			if _, ok := Data.StorageGroupIDToRDFStorageGroup[sgID]; ok {
				Data.VolumeIDToVolume[volumeID].Type = "RDF1+TDEV"
				if strings.Contains(sgID, "ASYNC") {
					Data.VolumeIDToVolume[volumeID].RDFGroupIDList = []types.RDFGroupID{
						{RDFGroupNumber: Data.AsyncRDFGroup.RdfgNumber},
					}
				} else {
					Data.VolumeIDToVolume[volumeID].RDFGroupIDList = []types.RDFGroupID{
						{RDFGroupNumber: Data.MetroRDFGroup.RdfgNumber},
					}
				}
			}

			// Update the Storage Group caches
			Data.StorageGroupIDToStorageGroup[sgID].NumOfVolumes++
			currentVolumes := Data.StorageGroupIDToVolumes[sgID]
			newVolumes := append(currentVolumes, volumeID)
			Data.StorageGroupIDToVolumes[sgID] = newVolumes

			// Check if the volume was added to a masking view
			if Data.StorageGroupIDToStorageGroup[sgID].NumOfMaskingViews > 0 {
				Data.VolumeIDToVolume[volumeID].NumberOfFrontEndPaths = 1
			}
			Data.VolumeIDToVolume[volumeID].NumberOfStorageGroups++
		}
	} else {
		// We are adding a new volume
		addNewVolume(volumeID, volumeIdentifier, size, sgID) // #nosec G20
	}
	return nil
}

// AddVolumeToStorageGroupTest - Adds volume to storage group and updates mock cache
func AddVolumeToStorageGroupTest(w http.ResponseWriter, name, size, sgID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addVolumeToStorageGroupTest(w, name, size, sgID)
}

func addVolumeToStorageGroupTest(w http.ResponseWriter, name, size, sgID string) {
	if name == "" || size == "" {
		writeError(w, "null name or size", http.StatusBadRequest)
	}
	id := strconv.Itoa(time.Now().Nanosecond())
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		writeError(w, "unable to convert size string to integer", http.StatusBadRequest)
	}
	if InducedErrors.VolumeNotCreatedError == false {
		addOneVolumeToStorageGroup(id, name, sgID, sizeInt) // #nosec G20
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("sloprovisioning/system/%s/storagegroup/%s", DefaultSymmetrixID, sgID)
	if InducedErrors.JobFailedError {
		newMockJob(id, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		newMockJob(id, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
	}
	returnJobByID(w, id)
}

// AddSpecificVolumeToStorageGroup - Add volume based on volumeids to storage group mock cache
func AddSpecificVolumeToStorageGroup(w http.ResponseWriter, volumeIDs []string, sgID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addSpecificVolumeToStorageGroup(w, volumeIDs, sgID)
}

func addSpecificVolumeToStorageGroup(w http.ResponseWriter, volumeIDs []string, sgID string) {
	if len(volumeIDs) == 0 {
		writeError(w, "empty list", http.StatusBadRequest)
	}
	jobID := strconv.Itoa(time.Now().Nanosecond())
	if InducedErrors.VolumeNotAddedError {
		writeError(w, "Error adding volume to the SG", http.StatusRequestTimeout)
		return
	}
	for _, volumeID := range volumeIDs {
		addOneVolumeToStorageGroup(volumeID, "TestVol", sgID, 0) // #nosec G20
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("sloprovisioning/system/%s/storagegroup/%s", DefaultSymmetrixID, sgID)
	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
	}
	returnJobByID(w, jobID)
}

func removeOneVolumeFromStorageGroup(volumeID, storageGroupID string) error {
	if _, ok := Data.StorageGroupIDToStorageGroup[storageGroupID]; !ok {
		return errors.New("The requested storage group doesn't exist")
	}
	if _, ok := Data.StorageGroupIDToVolumes[storageGroupID]; !ok {
		return errors.New("Storage Group to volume mapping doesn't exist")
	}
	vol, ok := Data.VolumeIDToVolume[volumeID]
	if !ok {
		return errors.New("The requested volume doesn't exist")
	}
	// Remove SG from the volume's SG list
	currentSGList := vol.StorageGroupIDList
	newStorageGroupList := make([]string, 0)
	for _, sgID := range currentSGList {
		if sgID != storageGroupID {
			newStorageGroupList = append(newStorageGroupList, sgID)
		}
	}
	vol.StorageGroupIDList = newStorageGroupList
	vol.NumberOfStorageGroups--
	// Change Volume's replication properties if replicated
	removeReplicationProps := false
	_, removeReplicationProps = Data.StorageGroupIDToRDFStorageGroup[storageGroupID]
	for _, sgID := range vol.StorageGroupIDList {
		if _, ok := Data.StorageGroupIDToRDFStorageGroup[sgID]; ok {
			removeReplicationProps = false
			break
		}
	}
	if removeReplicationProps {
		vol.Type = "TDEV"
		vol.RDFGroupIDList = nil
	}
	// Remove volume from the SG's volume list
	currentVolumeIDs := Data.StorageGroupIDToVolumes[storageGroupID]
	newVolumeIDList := make([]string, 0)
	for _, volID := range currentVolumeIDs {
		if volID != volumeID {
			newVolumeIDList = append(newVolumeIDList, volID)
		}
	}
	Data.StorageGroupIDToVolumes[storageGroupID] = newVolumeIDList
	// Update the count of volumes in SG
	Data.StorageGroupIDToStorageGroup[storageGroupID].NumOfVolumes--
	// Check if we need to update the number of front end paths for this particular volume
	update := true
	for _, sgID := range vol.StorageGroupIDList {
		if sg, ok := Data.StorageGroupIDToStorageGroup[sgID]; ok {
			if sg.NumOfMaskingViews > 1 {
				update = false
				break
			}
		}
	}
	if update {
		vol.NumberOfFrontEndPaths = 0
	}
	return nil
}

// RemoveVolumeFromStorageGroup - Remove volumes from storage group mock cache
func RemoveVolumeFromStorageGroup(w http.ResponseWriter, volumeIDs []string, sgID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeVolumeFromStorageGroup(w, volumeIDs, sgID)
}

func removeVolumeFromStorageGroup(w http.ResponseWriter, volumeIDs []string, sgID string) {
	for _, volID := range volumeIDs {
		fmt.Println("Volume ID: " + volID)
		removeOneVolumeFromStorageGroup(volID, sgID) // #nosec G20
	}
	returnStorageGroup(w, sgID, false)
}

func HandlePortGroup(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handlePortGroup(w, r)
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/portgroup/{id}
// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/portgroup
func handlePortGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pgID := vars["id"]
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetPortGroupError {
			writeError(w, "Error retrieving Port Group(s): induced error", http.StatusRequestTimeout)
			return
		}
		returnPortGroup(w, pgID)

	case http.MethodPost:
		if InducedErrors.CreatePortGroupError {
			writeError(w, "Error creating Port Group: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createPortGroupParams := &types.CreatePortGroupParams{}
		err := decoder.Decode(createPortGroupParams)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		addPortGroupFromCreateParams(createPortGroupParams)
		returnPortGroup(w, createPortGroupParams.PortGroupID)
	case http.MethodPut:
		if InducedErrors.UpdatePortGroupError {
			writeError(w, "Error updating Port Group: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updatePortGroupParams := &types.EditPortGroup{}
		err := decoder.Decode(updatePortGroupParams)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		updatePortGroupFromParams(pgID, updatePortGroupParams)
		returnPortGroup(w, pgID)
	case http.MethodDelete:
		if InducedErrors.DeletePortGroupError {
			writeError(w, "Error deleting Port Group: induced error", http.StatusRequestTimeout)
			return
		}
		_, err := deletePortGroup(pgID)
		if err != nil {
			writeError(w, "Error deletePortGroup", http.StatusRequestTimeout)
			return
		}
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// /univmax/restapi/90/system/symmetrix/{symid}/director/{director}/port/{id}
// /univmax/restapi/90/system/symmetrix/{symid}/director/{director}/port
func HandlePort(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handlePort(w, r)
}

func handlePort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dID := vars["director"]
	pID := vars["id"]
	queryString := r.URL.Query()
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetPortError {
			writeError(w, "Error retrieving Port(s): induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.GetPortGigEError {
			queryType, ok := queryString["type"]
			if ok {
				if queryType[0] == "Gige" { // The first ?type=<value>
					writeError(w, "Error retrieving GigE ports: induced error", http.StatusRequestTimeout)
					return
				}
				if queryType[0] == "OSHostAndRDF" { // The first ?type=<value>
					writeError(w, "Error retrieving OSHostAndRDF ports: induced error", http.StatusRequestTimeout)
					return
				}
			}
		}
		if InducedErrors.GetPortISCSITargetError {
			queryType, ok := queryString["iscsi_target"]
			if ok {
				if queryType[0] == "true" { // The first ?iscsi_target=<value>
					writeError(w, "Error retrieving ISCSI targets: induced error", http.StatusRequestTimeout)
					return
				}
			}
		}
		if InducedErrors.GetPortNVMeTCPTargetError {
			queryType, ok := queryString["nvmetcp_endpoint"]
			if ok {
				if queryType[0] == "true" { // The first ?nvmetcp_endpoint=<value>
					writeError(w, "Error retrieving NVMeTCP targets: induced error", http.StatusRequestTimeout)
					return
				}
			}
		}
		// if we asked for a specific Port, return those details
		if pID != "" {
			if InducedErrors.GetSpecificPortError {
				writeError(w, "Error retrieving Specific Port: induced error", http.StatusRequestTimeout)
				return
			}
			// Specific ports can be modeleted
			portName := dID + ":" + pID
			if Data.PortIDToSymmetrixPortType[portName] != nil {
				port := Data.PortIDToSymmetrixPortType[portName]
				if port == nil || port.Type == "" {
					writeError(w, "port not found", http.StatusNotFound)
				} else {
					symPort := &types.Port{
						SymmetrixPort: *port,
					}
					encoder := json.NewEncoder(w)
					encoder.Encode(symPort) // #nosec G20
				}
				return
			}
			returnPort(w, dID, pID)
		}
		// return a list of Ports
		returnPortIDList(w, dID)
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

func AddPort(id, identifier, portType string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addPort(id, identifier, portType)
}

// AddPort adds a port entry. Port type can either be "FibreChannel" or "GigE", or "" for a non existent port.
func addPort(id, identifier, portType string) {
	port := &types.SymmetrixPortType{
		Type:       portType,
		Identifier: identifier,
	}
	Data.PortIDToSymmetrixPortType[id] = port
}

func returnPort(w http.ResponseWriter, dID, pID string) {
	replacements := make(map[string]string)
	replacements["__PORT_ID__"] = pID
	replacements["__DIRECTOR_ID__"] = dID
	returnJSONFile(Data.JSONDir, "port_template.json", w, replacements)
}

func returnPortIDList(w http.ResponseWriter, dID string) {
	replacements := make(map[string]string)
	replacements["__DIRECTOR_ID__"] = dID
	returnJSONFile(Data.JSONDir, "portIDList.json", w, replacements)
}

// /univmax/restapi/90/system/symmetrix/{symid}/director/{{id}
// /univmax/restapi/90/system/symmetrix/{symid}/director
func HandleDirector(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleDirector(w, r)
}

func handleDirector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dID := vars["id"]
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetDirectorError {
			writeError(w, "Error retrieving Director(s): induced error", http.StatusRequestTimeout)
			return
		}
		// if we asked for a specific Director, return those details
		if dID != "" {
			returnDirector(w, dID)
		}
		// return a list of Directors
		returnDirectorIDList(w)

	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

func returnDirector(w http.ResponseWriter, dID string) {
	replacements := make(map[string]string)
	replacements["__DIRECTOR_ID__"] = dID
	returnJSONFile(Data.JSONDir, "director_template.json", w, replacements)
}

func returnDirectorIDList(w http.ResponseWriter) {
	replacements := make(map[string]string)
	returnJSONFile(Data.JSONDir, "directorIDList.json", w, replacements)
}

func HandleInitiator(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleInitiator(w, r)
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/initiator/{id}
// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/initiator
func handleInitiator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	initID := vars["id"]
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetInitiatorError {
			writeError(w, "Error retrieving Initiator(s): induced error", http.StatusRequestTimeout)
			return
		}
		if initID != "" {
			if InducedErrors.GetInitiatorByIDError {
				writeError(w, "Error retrieving Initiator By ID: induced error", http.StatusRequestTimeout)
				return
			}
		}
		returnInitiator(w, initID)

	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

func HandleHost(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleHost(w, r)
}

// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/host/{id}
// /univmax/restapi/90/sloprovisioning/symmetrix/{symid}/host
func handleHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostID := vars["id"]
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetHostError {
			writeError(w, "Error retrieving Host(s): induced error", http.StatusRequestTimeout)
			return
		}
		returnHost(w, hostID)

	case http.MethodPost:
		if InducedErrors.CreateHostError {
			writeError(w, "Error creating Host: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createHostParam := &types.CreateHostParam{}
		err := decoder.Decode(createHostParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		// Scan the initiators to see if there are any non iqn ones; then assume
		// host type Fibre.
		isFibre := false
		for _, initiator := range createHostParam.InitiatorIDs {
			if !strings.HasPrefix(initiator, "iqn.") {
				isFibre = true
			}
		}
		if isFibre {
			// Might need to add the Port information here
			addHost(createHostParam.HostID, "Fibre", createHostParam.InitiatorIDs) // #nosec G20
		} else {
			// initNode := make([]string, 0)
			// initNode = append(initNode, "iqn.1993-08.org.centos:01:5ae577b352a7")
			addHost(createHostParam.HostID, "iSCSI", createHostParam.InitiatorIDs) // #nosec G20
		}
		returnHost(w, createHostParam.HostID)

	case http.MethodPut:
		if hasError(&InducedErrors.UpdateHostError) {
			// if InducedErrors.UpdateHostError {
			writeError(w, "Error updating Host: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updateHostParam := &types.UpdateHostParam{}
		err := decoder.Decode(updateHostParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		returnHost(w, hostID)

	case http.MethodDelete:
		if InducedErrors.DeleteHostError {
			writeError(w, "Error deleting Host: induced error", http.StatusRequestTimeout)
			return
		}
		err := removeHost(hostID)
		if err != nil {
			writeError(w, "error removeHost", http.StatusBadRequest)
			return
		}
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// ReturnHost - Returns a host from cache
func ReturnHost(w http.ResponseWriter, hostID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnHost(w, hostID)
}

func returnHost(w http.ResponseWriter, hostID string) {
	if hostID != "" {
		if host, ok := Data.HostIDToHost[hostID]; ok {
			writeJSON(w, host)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		hostIDs := make([]string, 0)
		for k := range Data.HostIDToHost {
			hostIDs = append(hostIDs, k)
		}
		hostIDList := &types.HostList{
			HostIDs: hostIDs,
		}
		writeJSON(w, hostIDList)
	}
}

// ReturnPortGroup - Returns port group information from cache
func ReturnPortGroup(w http.ResponseWriter, portGroupID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnPortGroup(w, portGroupID)
}

func returnPortGroup(w http.ResponseWriter, portGroupID string) {
	if portGroupID != "" {
		if pg, ok := Data.PortGroupIDToPortGroup[portGroupID]; ok {
			fmt.Printf("\n%v\n", pg)
			writeJSON(w, pg)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		portGroupIDs := make([]string, 0)
		for k := range Data.PortGroupIDToPortGroup {
			portGroupIDs = append(portGroupIDs, k)
		}
		portGroupList := &types.PortGroupList{
			PortGroupIDs: portGroupIDs,
		}
		writeJSON(w, portGroupList)
	}
}

// /univmax/restapi/performance/StorageGroup/metrics
func HandleStorageGroupMetrics(w http.ResponseWriter, _ *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleStorageGroupMetrics(w, nil)
}

func handleStorageGroupMetrics(w http.ResponseWriter, _ *http.Request) {
	if InducedErrors.GetStorageGroupMetricsError {
		writeError(w, "Error getting storage group metrics: induced error", http.StatusRequestTimeout)
		return
	}
	sgMetric := types.StorageGroupMetric{
		HostReads:         0.0,
		HostWrites:        0.0,
		HostMBReads:       0.0,
		HostMBWritten:     0.0,
		ReadResponseTime:  0.0,
		WriteResponseTime: 0.0,
		AllocatedCapacity: 0.0,
		AvgIOSize:         0.0,
		Timestamp:         1671091500000,
	}
	metricsIterator := &types.StorageGroupMetricsIterator{
		ResultList: types.StorageGroupMetricsResultList{
			Result: []types.StorageGroupMetric{sgMetric},
			From:   1,
			To:     1,
		},
		ID:             "query_id",
		Count:          1,
		ExpirationTime: 1671091597409,
		MaxPageSize:    1000,
	}
	writeJSON(w, metricsIterator)
}

// /univmax/restapi/performance/Volume/metrics
func HandleVolumeMetrics(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleVolumeMetrics(w, r)
}

func handleVolumeMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commaSeparatedStorageGroupList := vars["commaSeparatedStorageGroupList"]
	if InducedErrors.GetVolumesMetricsError {
		writeError(w, "Error getting volume metrics: induced error", http.StatusRequestTimeout)
		return
	}
	volumeMetric := types.VolumeMetric{
		MBRead:            0.0,
		MBWritten:         0.0,
		Reads:             0.0,
		Writes:            0.0,
		ReadResponseTime:  0.0,
		WriteResponseTime: 0.0,
		IoRate:            5.0,
		Timestamp:         1671091500000,
	}
	if InducedErrors.GetFreshMetrics {
		volumeMetric.Timestamp = time.Now().UnixMilli()
	}
	volumeResult := types.VolumeResult{
		VolumeResult:  []types.VolumeMetric{volumeMetric},
		VolumeID:      "002C8",
		StorageGroups: commaSeparatedStorageGroupList,
	}
	metricsIterator := &types.VolumeMetricsIterator{
		ResultList: types.VolumeMetricsResultList{
			Result: []types.VolumeResult{volumeResult},
			From:   1,
			To:     1,
		},
		ID:             "query_id",
		Count:          1,
		ExpirationTime: 1671091597409,
		MaxPageSize:    1000,
	}
	writeJSON(w, metricsIterator)
}

// /univmax/restapi/performance/file/filesystem/metrics
func HandleFileSysMetrics(w http.ResponseWriter, _ *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleFileSysMetrics(w, nil)
}

func handleFileSysMetrics(w http.ResponseWriter, _ *http.Request) {
	if InducedErrors.GetFileSysMetricsError {
		writeError(w, "Error getting volume metrics: induced error", http.StatusRequestTimeout)
		return
	}
	fileMetric := types.FileSystemResult{
		PercentBusy: 1,
		Timestamp:   1671091500000,
	}

	metricsIterator := &types.FileSystemMetricsIterator{
		ResultList: types.FileSystemMetricsResultList{
			Result: []types.FileSystemResult{fileMetric},
			From:   1,
			To:     1,
		},
		ID:             "query_id",
		Count:          1,
		ExpirationTime: 1671091597409,
		MaxPageSize:    1000,
	}
	writeJSON(w, metricsIterator)
}

// /univmax/restapi/performance/StorageGroup/keys
func HandleStorageGroupPerfKeys(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleStorageGroupPerfKeys(w, r)
}

func handleStorageGroupPerfKeys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storageGroupID := vars["storageGroupId"]
	if InducedErrors.GetStorageGroupPerfKeyError {
		writeError(w, "Error getting storage group perf key: induced error", http.StatusRequestTimeout)
		return
	}
	sgInfo := types.StorageGroupInfo{
		StorageGroupID:     storageGroupID,
		FirstAvailableDate: 0,
		LastAvailableDate:  1671091597409,
	}
	perfKeys := &types.StorageGroupKeysResult{
		StorageGroupInfos: []types.StorageGroupInfo{sgInfo},
	}
	writeJSON(w, perfKeys)
}

// /univmax/restapi/performance/Array/keys
func HandleArrayPerfKeys(w http.ResponseWriter, _ *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleArrayPerfKeys(w, nil)
}

func handleArrayPerfKeys(w http.ResponseWriter, _ *http.Request) {
	if InducedErrors.GetArrayPerfKeyError {
		writeError(w, "Error getting array perf key: induced error", http.StatusRequestTimeout)
		return
	}
	arrayInfo := types.ArrayInfo{
		SymmetrixID:        DefaultSymmetrixID,
		FirstAvailableDate: 0,
		LastAvailableDate:  1671091597409,
	}
	perfKeys := &types.ArrayKeysResult{
		ArrayInfos: []types.ArrayInfo{arrayInfo},
	}
	writeJSON(w, perfKeys)
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleNotFound(w, r)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, "URL not found: "+r.URL.String(), http.StatusNotFound)
}

// Write an error code to the response writer
func writeError(w http.ResponseWriter, message string, httpStatus int) {
	w.WriteHeader(httpStatus)
	resp := new(types.Error)
	resp.Message = message
	// The following aren't used by the hardware but could be used internally
	// resp.HTTPStatusCode = http.StatusNotFound
	// resp.ErrorCode = int(errorCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)
	if err != nil {
		log.Printf("error encoding json: %s", err.Error())
	}
}

// Return content from a JSON file. Arguments are:
//
//	 directory, filename  of the file
//	wrriter ResponseWriter where data is output
//
// An optional replacement map. If supplied every instance of a key in the JSON file will be replaced with the corresponding value.
func returnJSONFile(directory, filename string, w http.ResponseWriter, replacements map[string]string) (jsonBytes []byte) {
	jsonBytes, err := os.ReadFile(filepath.Join(directory, filename)) // #nosec G20
	if err != nil {
		log.Printf("Couldn't read %s/%s", directory, filename)
		if w != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		return make([]byte, 0)
	}
	if replacements != nil {
		jsonString := string(jsonBytes)
		for key, value := range replacements {
			jsonString = strings.Replace(jsonString, key, value, -1)
		}
		if Debug {
			log.Printf("Edited payload:%s", jsonString)
		}
		jsonBytes = []byte(jsonString)
	}
	if Debug {
		log.Printf("jsonBytes:%s", jsonBytes)
	}
	if w != nil {
		_, err = w.Write(jsonBytes)
		if err != nil {
			log.Printf("Couldn't write to ResponseWriter")
			w.WriteHeader(http.StatusInternalServerError)
			return make([]byte, 0)
		}
	}
	return jsonBytes
}

func AddTempSnapshots() {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addTempSnapshots()
}

// AddTempSnapshots adds marked for deletion snapshots into mock to help snapcleanup thread to be functional
func addTempSnapshots() {
	for i := 1; i <= 2; i++ {
		id := fmt.Sprintf("%05d", i)
		size := 7
		volumeIdentifier := "Vol" + id
		addNewVolume(id, volumeIdentifier, size, DefaultStorageGroup) // #nosec G20
		SnapID := fmt.Sprintf("%s-%s-%d", "DEL", "snapshot", i)
		addNewSnapshot(id, SnapID)
	}
}

// univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/snapshot/{SnapID}
func HandleSnapshot(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleSnapshot(w, r)
}

func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// volID := vars["volID"]
	SnapID := vars["SnapID"]
	switch r.Method {
	case http.MethodPost:
		if InducedErrors.CreateSnapshotError {
			writeError(w, "Failed to create snapshot: induced error", http.StatusBadRequest)
			return
		}
		if InducedErrors.MaxSnapSessionError {
			writeError(w, "The maximum number of sessions has been exceeded for the specified Source device", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createSnapParam := &types.CreateVolumesSnapshot{}
		err := decoder.Decode(createSnapParam)
		if err != nil {
			writeError(w, "problem decoding POST Snapshot payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		createSnapshot(w, r, vars["SnapID"], createSnapParam.ExecutionOption, createSnapParam.SourceVolumeList)
		return
	case http.MethodPut:
		if SnapID == "" {
			writeError(w, "Snapshot name must be supplied", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updateSnapParam := &types.ModifyVolumeSnapshot{}
		err := decoder.Decode(updateSnapParam)
		if err != nil {
			writeError(w, "problem decoding PUT Snapshot payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("PUT Snapshot payload: %#v\n", updateSnapParam)
		executionOption := updateSnapParam.ExecutionOption

		if updateSnapParam.Action == "Rename" {
			if InducedErrors.RenameSnapshotError {
				writeError(w, "error renaming the snapshot: induced error", http.StatusBadRequest)
				return
			}
			renameSnapshot(w, r, updateSnapParam.VolumeNameListSource, executionOption, SnapID, updateSnapParam.NewSnapshotName)
			return
		}
		if updateSnapParam.Action == "Link" {
			if InducedErrors.MaxSnapSessionError {
				writeError(w, "The maximum number of sessions has been exceeded for the specified Source device", http.StatusBadRequest)
				return
			}
			if InducedErrors.LinkSnapshotError {
				writeError(w, "error linking the snapshot: induced error", http.StatusBadRequest)
				return
			}
			linkSnapshot(w, r, updateSnapParam.VolumeNameListSource, updateSnapParam.VolumeNameListTarget, executionOption, SnapID)
			return
		}
		if updateSnapParam.Action == "Unlink" {
			if InducedErrors.LinkSnapshotError {
				writeError(w, "error unlinking the snapshot: induced error", http.StatusBadRequest)
				return
			}
			unlinkSnapshot(w, r, updateSnapParam.VolumeNameListSource, updateSnapParam.VolumeNameListTarget, executionOption, SnapID)
			return
		}
		if updateSnapParam.Action == "Restore" {
			// restoreSnapshot(w, r, updateSnapParam.VolumeNameListSource, updateSnapParam.VolumeNameListTarget, executionOption, SnapID)
			// return
			fmt.Printf("Not yet implemented")
		}
	case http.MethodDelete:
		decoder := json.NewDecoder(r.Body)
		deleteSnapParam := &types.DeleteVolumeSnapshot{}
		err := decoder.Decode(deleteSnapParam)
		if err != nil {
			writeError(w, "problem decoding Delete Snapshot payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		deleteSnapshot(w, r, vars["SnapID"], deleteSnapParam.ExecutionOption, deleteSnapParam.DeviceNameListSource, deleteSnapParam.Generation)
		return
	}
}

// CreateSnapshot - Creates a snapshot and updates mock cache
func CreateSnapshot(w http.ResponseWriter, r *http.Request, SnapID, executionOption string, sourceVolumeList []types.VolumeList) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	createSnapshot(w, r, SnapID, executionOption, sourceVolumeList)
}

func createSnapshot(w http.ResponseWriter, _ *http.Request, SnapID, executionOption string, sourceVolumeList []types.VolumeList) {
	if strings.Contains(SnapID, ":") {
		writeError(w, "error, invalid snapshot name", http.StatusBadRequest)
		return
	}
	if executionOption != types.ExecutionOptionSynchronous {
		writeError(w, "expected SYNCHRONOUS", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(sourceVolumeList) {
		writeError(w, "few devices not available", http.StatusBadRequest)
		return
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("/replication/symmetrix/%s/snapshot/%s", DefaultSymmetrixID, SnapID)
	jobID := fmt.Sprintf("SnapID-%d", time.Now().Nanosecond())
	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
		returnJobByID(w, jobID)
		return
	}
	for i := 0; i < len(sourceVolumeList); i++ {
		source := sourceVolumeList[i].Name
		if !duplicateSnapshotCreationRequest(source, SnapID) {
			// Snapshot with unique name
			addNewSnapshot(source, SnapID)
		}
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
	}
	returnJobByID(w, jobID)
}

// AddNewSnapshot adds a snapshot to the mock cache
func AddNewSnapshot(source, SnapID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addNewSnapshot(source, SnapID)
}

func addNewSnapshot(source, SnapID string) {
	time := time.Now().Nanosecond()
	snapshot := &types.Snapshot{
		Name:       SnapID,
		Generation: 0,
		State:      "Established",
		Timestamp:  strconv.Itoa(time),
	}
	snapIDtoSnap := Data.VolIDToSnapshots[source]
	if snapIDtoSnap == nil {
		snapIDtoSnap = map[string]*types.Snapshot{}
	}
	snapIDtoSnap[SnapID] = snapshot
	Data.VolIDToSnapshots[source] = snapIDtoSnap
	Data.VolumeIDToVolume[source].SnapSource = true
	fmt.Printf("*****added** %v***", Data.VolIDToSnapshots[source][SnapID])
	fmt.Printf("****Total Snaps on %s are: %d****", source, len(Data.VolIDToSnapshots[source]))
}

// DeleteSnapshot - Deletes a snapshot and updates mock cache
func DeleteSnapshot(w http.ResponseWriter, r *http.Request, SnapID string, executionOption string, deviceNameListSource []types.VolumeList, genID int64) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	deleteSnapshot(w, r, SnapID, executionOption, deviceNameListSource, genID)
}

func deleteSnapshot(w http.ResponseWriter, _ *http.Request, SnapID string, _ string, deviceNameListSource []types.VolumeList, _ int64) {
	if InducedErrors.DeleteSnapshotError {
		writeError(w, "error deleting the snapshot: induced error", http.StatusBadRequest)
		return
	}
	if deviceNameListSource[0].Name == "" {
		writeError(w, "no source volume names given to link the snapshot", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(deviceNameListSource) {
		writeError(w, "few devices not available", http.StatusBadRequest)
		return
	}
	resourceLink := fmt.Sprintf("/replication/symmetrix/%s/snapshot/%s", DefaultSymmetrixID, SnapID)
	jobID := fmt.Sprintf("SnapID-%d", time.Now().Nanosecond())
	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		for i := 0; i < len(deviceNameListSource); i++ {
			source := deviceNameListSource[i].Name

			// volume exists, check for availability of snapshot on it i.e, check if snapshot is found in snapIDtoSnap map "SnapID": Snapshot
			snapIDtoSnap := Data.VolIDToSnapshots[source]
			if _, ok := snapIDtoSnap[SnapID]; !ok {
				// snapshot is not found
				writeError(w, "no snapshot information", http.StatusBadRequest)
				return
			}

			// snapshot exists, check if it is linked to any target device/volumes
			snapIDtoLinkedVolKey := SnapID + ":" + source
			linkedVolume := Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey]
			if len(linkedVolume) > 0 {
				// snapshot is linked to some volumes, can not delete
				writeError(w, "delete cannot be attempted because the snapshot has a link", http.StatusBadRequest)
				return
			}

			// all checks done: volume exists, snapshot existing without links -> it can be deleted
			delete(snapIDtoSnap, SnapID)
			Data.VolumeIDToVolume[source].SnapSource = false
			newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
		}
	}
	returnJobByID(w, jobID)
}

// RenameSnapshot - Renames a snapshot and updates mock cache
func RenameSnapshot(w http.ResponseWriter, r *http.Request, sourceVolumeList []types.VolumeList, executionOption, oldSnapID, newSnapID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	renameSnapshot(w, r, sourceVolumeList, executionOption, oldSnapID, newSnapID)
}

func renameSnapshot(w http.ResponseWriter, _ *http.Request, sourceVolumeList []types.VolumeList, _, oldSnapID, newSnapID string) {
	if fewVolumeUnavalaible(sourceVolumeList) {
		writeError(w, "few devices not available", http.StatusBadRequest)
		return
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("/replication/symmetrix/%s/snapshot/%s", DefaultSymmetrixID, oldSnapID)
	jobID := fmt.Sprintf("SnapID-%d", time.Now().Nanosecond())
	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		for _, volID := range sourceVolumeList {
			if Data.VolIDToSnapshots[volID.Name][oldSnapID] == nil {
				writeError(w, "no snapshot information, Snapshot cannot be found", http.StatusBadRequest)
				return
			}
			for _, snap := range Data.VolIDToSnapshots[volID.Name] {
				if snap.Name == oldSnapID {
					snap.Name = newSnapID
					Data.VolIDToSnapshots[volID.Name] = map[string]*types.Snapshot{newSnapID: snap}
					newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
					Data.VolIDToSnapshots[volID.Name] = map[string]*types.Snapshot{newSnapID: snap}
					newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
				}
			}
		}
		returnJobByID(w, jobID)
	}
}

// LinkSnapshot - Links a snapshot and updates a mock cache
func LinkSnapshot(w http.ResponseWriter, r *http.Request, sourceVolumeList []types.VolumeList, targetVolumeList []types.VolumeList, executionOption, SnapID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	linkSnapshot(w, r, sourceVolumeList, targetVolumeList, executionOption, SnapID)
}

func linkSnapshot(w http.ResponseWriter, _ *http.Request, sourceVolumeList []types.VolumeList, targetVolumeList []types.VolumeList, _, SnapID string) {
	if sourceVolumeList[0].Name == "" {
		writeError(w, "no source volume names given to link the snapshot", http.StatusBadRequest)
		return
	}
	if targetVolumeList[0].Name == "" {
		writeError(w, "no link volume names given to link the snapshot to", http.StatusBadRequest)
		return
	}
	if len(sourceVolumeList) != len(targetVolumeList) {
		writeError(w, "cannot link snapshot, the number of source and devices should be same", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(sourceVolumeList) {
		writeError(w, "few source devices not available", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(targetVolumeList) {
		writeError(w, "few target devices not available", http.StatusBadRequest)
		return
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("/replication/symmetrix/%s/snapshot/%s", DefaultSymmetrixID, SnapID)
	jobID := fmt.Sprintf("SnapID-%d", time.Now().Nanosecond())

	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		for key, volID := range sourceVolumeList {
			snapIDtoSnap := Data.VolIDToSnapshots[volID.Name]
			targetVolID := targetVolumeList[key].Name
			if snapIDtoSnap[SnapID] == nil {
				writeError(w, "no snapshot information, snopshot cannot be found on this device", http.StatusBadRequest)
				return
			}
			// all devices exist, #source=#target, snapshot exist, check if target already linked
			snapIDtoLinkedVolKey := SnapID + ":" + volID.Name
			volIDToLinkedVols := Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey]
			if volIDToLinkedVols == nil {
				// No Linked Volume, first link request for this SnapID
				volIDToLinkedVols = map[string]*types.LinkedVolumes{}
			} else {
				// snapshot is linked to few devices, check if target is already linked
				if !(volIDToLinkedVols[targetVolID] == nil) {
					// duplicate link request
					writeError(w, "devices already in desired state", http.StatusBadRequest)
					return
				}
			}
			// all devices exist, #source=#target, snapshot exist, target is not linked -> ideal for Linking
			time := time.Now().Nanosecond()
			linkedVolume := &types.LinkedVolumes{
				TargetDevice: targetVolID,
				Timestamp:    strconv.Itoa(time),
				State:        "Linked",
				Copy:         false,
				Restored:     false,
				Linked:       true,
				Defined:      true,
			}
			if InducedErrors.TargetNotDefinedError {
				linkedVolume.Defined = false
			}

			volIDToLinkedVols[targetVolID] = linkedVolume
			Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey] = volIDToLinkedVols
			Data.VolumeIDToVolume[targetVolID].SnapTarget = true
			newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
		}
	}
	returnJobByID(w, jobID)
}

// UnlinkSnapshot - Unlinks a snapshot and updates mock cache
func UnlinkSnapshot(w http.ResponseWriter, r *http.Request, sourceVolumeList []types.VolumeList, targetVolumeList []types.VolumeList, executionOption, SnapID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	unlinkSnapshot(w, r, sourceVolumeList, targetVolumeList, executionOption, SnapID)
}

func unlinkSnapshot(w http.ResponseWriter, _ *http.Request, sourceVolumeList []types.VolumeList, targetVolumeList []types.VolumeList, _, SnapID string) {
	if sourceVolumeList[0].Name == "" {
		writeError(w, "no source volume names given to unlink the snapshot", http.StatusBadRequest)
		return
	}
	if targetVolumeList[0].Name == "" {
		writeError(w, "no target volume names given to unlink the snapshot to", http.StatusBadRequest)
		return
	}
	if len(sourceVolumeList) != len(targetVolumeList) {
		writeError(w, "cannot unlink snapshot, the number of source and devices should be same", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(sourceVolumeList) {
		writeError(w, "few source devices not available", http.StatusBadRequest)
		return
	}
	if fewVolumeUnavalaible(targetVolumeList) {
		writeError(w, "few target devices not available", http.StatusBadRequest)
		return
	}
	// Make a job to return
	resourceLink := fmt.Sprintf("/replication/symmetrix/%s/snapshot/%s", DefaultSymmetrixID, SnapID)
	jobID := fmt.Sprintf("SnapID-%d", time.Now().Nanosecond())

	if InducedErrors.JobFailedError {
		newMockJob(jobID, types.JobStatusRunning, types.JobStatusFailed, resourceLink)
	} else {
		for key, volID := range sourceVolumeList {
			snapIDtoSnap := Data.VolIDToSnapshots[volID.Name]
			targetVolID := targetVolumeList[key].Name
			if snapIDtoSnap[SnapID] == nil {
				writeError(w, "no snapshot information, snopshot cannot be found on this device", http.StatusBadRequest)
				return
			}
			// all devices exist, #source=#target, snapshot exist, check if source is linked to target
			snapIDtoLinkedVolKey := SnapID + ":" + volID.Name
			volIDToLinkedVolumes := Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey]
			if _, ok := volIDToLinkedVolumes[targetVolID]; ok {
				// source volume is linked to target, ideal for unlink
				delete(volIDToLinkedVolumes, targetVolID)
				volIDToLinkedVolumes = Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey]
				Data.VolumeIDToVolume[targetVolID].SnapTarget = false
				newMockJob(jobID, types.JobStatusRunning, types.JobStatusSucceeded, resourceLink)
			} else {
				// already unlinked
				writeError(w, "devices already in desired state", http.StatusBadRequest)
				return
			}
		}
	}
	returnJobByID(w, jobID)
}

// check if all the devices exist in the Mock VolumeIDToVolume or check if any unvailable devices
func fewVolumeUnavalaible(sourceVolumeList []types.VolumeList) bool {
	for _, volID := range sourceVolumeList {
		if Data.VolumeIDToVolume[volID.Name] == nil {
			return true
		}
	}
	return false
}

// returns true for Snapshot Creation if a snpshot with same name already there, false otherwise
func duplicateSnapshotCreationRequest(source, SnapID string) bool {
	_, ok := Data.VolIDToSnapshots[source][SnapID]
	return ok
}

// GET univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/volume
func HandleSymVolumes(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleSymVolumes(w, r)
}

func handleSymVolumes(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.GetSymVolumeError {
		writeError(w, "error fetching the list: induced error", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	symVolumeList := new(types.SymVolumeList)
	if details := queryParams.Get("includeDetails"); details == "true" {
		for key, snapshots := range Data.VolIDToSnapshots {
			symVolumeList.Name = append(symVolumeList.Name, key)
			var snapList []types.Snapshot
			for _, snap := range snapshots {
				snapshotName := fmt.Sprintf("%s-SRC-%s-%d", symVolumeList.Name[0], snap.Name, snap.Generation)
				if InducedErrors.InvalidSnapshotName {
					snapshotName = "InvalidSnapshot"
				}
				snapshot := types.Snapshot{
					Name:       snapshotName,
					Generation: snap.Generation,
					Timestamp:  snap.Timestamp,
					State:      snap.State,
				}
				snapList = append(snapList, snapshot)
			}
			symDevice := types.SymDevice{
				SymmetrixID: DefaultSymmetrixID,
				Name:        key,
				Snapshot:    snapList,
			}
			symVolumeList.SymDevice = append(symVolumeList.SymDevice, symDevice)
		}
	} else {
		for key := range Data.VolIDToSnapshots {
			symVolumeList.Name = append(symVolumeList.Name, key)
		}
	}
	writeJSON(w, symVolumeList)
}

// GET univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/volume/{volID}/snapshot/
// GET univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}
func HandleVolSnaps(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleVolSnaps(w, r)
}

func handleVolSnaps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volID := vars["volID"]
	SnapID := vars["SnapID"]
	if InducedErrors.GetVolSnapsError {
		writeError(w, "error fetching the Snapshot Info: induced error", http.StatusBadRequest)
		return
	}
	if Data.VolumeIDToVolume[volID] == nil {
		writeError(w, "Could not find volume: "+volID, http.StatusNotFound)
		return
	}

	volumeSnapshotSource, _ := returnSnapshotObjectList(volID)
	volumeSnapshotLink := returnVolumeSnapshotLink(volID)

	if SnapID == "" {
		// Both Volume Snapshots exist
		// for /{symid}/volume/{volID}/snapshot/
		snaphotVolumeGeneration := new(types.SnapshotVolumeGeneration)
		snaphotVolumeGeneration.DeviceName = volID
		snaphotVolumeGeneration.VolumeSnapshotSource = volumeSnapshotSource
		snaphotVolumeGeneration.VolumeSnapshotLink = volumeSnapshotLink
		writeJSON(w, snaphotVolumeGeneration)
	} else {
		// Both Volume Snapshots exist
		// for /{symid}/volume/{volID}/snapshot/{SnapID}
		volumeSnapshot := new(types.VolumeSnapshot)
		volumeSnapshot.DeviceName = volID
		volumeSnapshot.SnapshotName = SnapID
		for _, snapSrc := range volumeSnapshotSource {
			if snapSrc.SnapshotName == SnapID {
				volumeSnapshot.VolumeSnapshotSource = append(volumeSnapshot.VolumeSnapshotSource, types.VolumeSnapshotSource{
					SnapshotName: snapSrc.SnapshotName,
					Generation:   snapSrc.Generation,
					TimeStamp:    snapSrc.TimeStamp,
					State:        snapSrc.State,
				})
			}
		}
		volumeSnapshot.VolumeSnapshotLink = volumeSnapshotLink
		writeJSON(w, volumeSnapshot)
	}
}

// returns the List of VolumesSnapshot objects derived based on existing mock Snapshot object
func returnSnapshotObjectList(volID string) ([]types.VolumeSnapshotSource, []int64) {
	var volumeSnapshotSrc []types.VolumeSnapshotSource
	var generations []int64
	for _, snap := range Data.VolIDToSnapshots[volID] {
		snapshotSrc := types.VolumeSnapshotSource{
			SnapshotName:  snap.Name,
			Generation:    snap.Generation,
			TimeStamp:     snap.Timestamp,
			State:         snap.State,
			LinkedVolumes: returnLinkedVolumes(snap.Name + ":" + volID),
		}
		if InducedErrors.SnapshotExpired {
			snapshotSrc.Expired = true
		}
		volumeSnapshotSrc = append(volumeSnapshotSrc, snapshotSrc)
		generations = append(generations, snap.Generation)
	}

	return volumeSnapshotSrc, generations
}

// returns the List of Linked Volumes to Snapshots of a volume
func returnLinkedVolumes(snapIDtoLinkedVolKey string) []types.LinkedVolumes {
	var linkedVolumes []types.LinkedVolumes
	for _, volume := range Data.SnapIDToLinkedVol[snapIDtoLinkedVolKey] {
		linkedVolumes = append(linkedVolumes, *volume)
	}
	return linkedVolumes
}

// returns the List of volumeSnapshotLink to a Snapshot
func returnVolumeSnapshotLink(targetVolID string) []types.VolumeSnapshotLink {
	var snapshotLnk []types.VolumeSnapshotLink
	for _, volume := range Data.SnapIDToLinkedVol {
		if target, ok := volume[targetVolID]; ok {
			snapshotLnk = append(snapshotLnk, types.VolumeSnapshotLink{
				TargetDevice:     target.TargetDevice,
				Timestamp:        target.Timestamp,
				State:            target.State,
				TrackSize:        target.TrackSize,
				Tracks:           target.Tracks,
				PercentageCopied: target.PercentageCopied,
				Linked:           target.Linked,
				Restored:         target.Restored,
				Defined:          target.Defined,
				Copy:             target.Copy,
				Destage:          target.Destage,
				Modified:         target.Modified,
			})
		}
	}
	return snapshotLnk
}

// GET univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}/generation
// GET univmax/restapi/private/APIVersion/replication/symmetrix/{symid}/volume/{volID}/snapshot/{SnapID}/generation/{genID}
func HandleGenerations(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleGenerations(w, r)
}

func handleGenerations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volID := vars["volID"]
	SnapID := vars["SnapID"]
	genID := vars["genID"]
	if Data.VolumeIDToVolume[volID] == nil {
		writeError(w, "Could not find volume: "+volID, http.StatusNotFound)
		return
	}

	volumeSnapshotSource, generations := returnSnapshotObjectList(volID)
	volumeSnapshotLink := returnVolumeSnapshotLink(volID)

	if genID == "" {
		// Both Volume Snapshots exist
		// for /{symid}/volume/{volID}/snapshot/{SnapID}/generation/
		volumeSnapshotGenerations := new(types.VolumeSnapshotGenerations)
		volumeSnapshotGenerations.DeviceName = volID
		volumeSnapshotGenerations.Generation = generations
		volumeSnapshotGenerations.SnapshotName = SnapID
		volumeSnapshotGenerations.VolumeSnapshotSource = volumeSnapshotSource
		volumeSnapshotGenerations.VolumeSnapshotLink = volumeSnapshotLink
		writeJSON(w, volumeSnapshotGenerations)
		return
	}
	// Both Volume Snapshots exist
	// for /{symid}/volume/{volID}/snapshot/{SnapID}/generation/{genID}
	volumeSnapshotGeneration := new(types.VolumeSnapshotGeneration)
	volumeSnapshotGeneration.DeviceName = volID
	volumeSnapshotGeneration.SnapshotName = SnapID
	volumeSnapshotGeneration.VolumeSnapshotLink = volumeSnapshotLink
	// volumeSnapshotGeneration.VolumeSnapshotSource = returnSnapshotGenerationInfo(volID, SnapID, genID)
	gID, _ := strconv.ParseInt(genID, 10, 64)
	for _, snapSrc := range volumeSnapshotSource {
		if snapSrc.SnapshotName == SnapID && snapSrc.Generation == gID {
			volumeSnapshotGeneration.VolumeSnapshotSource = snapSrc
			volumeSnapshotGeneration.Generation = snapSrc.Generation
			break
		}
	}
	writeJSON(w, volumeSnapshotGeneration)
	return
}

func HandleCapabilities(w http.ResponseWriter, _ *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleCapabilities(w, nil)
}

func handleCapabilities(w http.ResponseWriter, _ *http.Request) {
	var jsonBytes []byte
	if InducedErrors.SnapshotNotLicensed {
		jsonBytes = []byte("{\"symmetrixCapability\":[{\"symmetrixId\":\"000197900046\",\"snapVxCapable\":false,\"rdfCapable\":true,\"virtualWitnessCapable\":false}]}")
	} else if InducedErrors.InvalidResponse {
		writeError(w, "something went wrong: induced error", http.StatusBadRequest)
		return
	} else if InducedErrors.UnisphereMismatchError {
		jsonBytes = []byte("{\"symmetrixCapability\":[{\"symmetrixId\":\"000000000000\",\"snapVxCapable\":true,\"rdfCapable\":true,\"virtualWitnessCapable\":false}]}")
	} else {
		jsonBytes = []byte("{\"symmetrixCapability\":[{\"symmetrixId\":\"000197900046\",\"snapVxCapable\":true,\"rdfCapable\":true,\"virtualWitnessCapable\":false}]}")
	}
	_, err := w.Write(jsonBytes)
	if err != nil {
		log.Printf("Couldn't write to ResponseWriter")
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func HandlePrivVolume(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handlePrivVolume(w, r)
}

func handlePrivVolume(w http.ResponseWriter, r *http.Request) {
	if InducedErrors.GetPrivVolumeByIDError {
		writeError(w, "error fetching the Volume structure: induced error", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	privateVolumeIterator := new(types.PrivVolumeIterator)
	if wwn := queryParams.Get("wwn"); wwn != "" {
		volID := wwn[27:]
		volume := Data.VolumeIDToVolume[volID]
		volumeHeader := parseVolumetoVolumeHeader(volume)
		timeFinderInfo := returnTimeFinderInfo(volID)
		var result []types.VolumeResultPrivate
		result = append(result, types.VolumeResultPrivate{
			VolumeHeader:   *volumeHeader,
			TimeFinderInfo: *timeFinderInfo,
		})
		privVolumeResultList := types.PrivVolumeResultList{
			PrivVolumeList: result,
			From:           1,
			To:             1,
		}
		privateVolumeIterator.ResultList = privVolumeResultList
		privateVolumeIterator.ID = "70e15d35-baaf-43d3-865a-bf3300684895_0"
		privateVolumeIterator.ExpirationTime = 1576137450163
		privateVolumeIterator.MaxPageSize = 1000
		privateVolumeIterator.Count = 1
	}
	writeJSON(w, privateVolumeIterator)
}

func parseVolumetoVolumeHeader(volume *types.Volume) *types.VolumeHeader {
	volumeHeader := &types.VolumeHeader{
		VolumeID:     volume.VolumeID,
		CapGB:        volume.CapacityGB,
		CapMB:        volume.FloatCapacityMB,
		Status:       volume.Status,
		SSID:         volume.SSID,
		EffectiveWWN: volume.WWN,
		Encapsulated: volume.Encapsulated,
	}

	return volumeHeader
}

func returnTimeFinderInfo(volID string) *types.TimeFinderInfo {
	timeFinder := new(types.TimeFinderInfo)
	if _, ok := Data.VolIDToSnapshots[volID]; ok {
		timeFinder.SnapVXSrc = ok
	}
	for _, volIDToLinkedVols := range Data.SnapIDToLinkedVol {
		if _, ok := volIDToLinkedVols[volID]; ok {
			timeFinder.SnapVXTgt = ok
			break
		}
	}
	if timeFinder.SnapVXSrc || timeFinder.SnapVXTgt {
		timeFinder.SnapVXSession = append(timeFinder.SnapVXSession, returnSnapVXSession(volID, timeFinder.SnapVXSrc, timeFinder.SnapVXTgt))
	}
	return timeFinder
}

func returnSnapVXSession(volID string, isSource, isTarget bool) types.SnapVXSession {
	var snapVXSession types.SnapVXSession
	if isSource {
		snapVXSession.SourceSnapshotGenInfo = returnSrcSnapshotGenInfo(volID)
	}

	if isTarget {
		for snapIDtoLinkedVolKey, volIDToLinkedVolumes := range Data.SnapIDToLinkedVol {
			sourceVolID := strings.Split(snapIDtoLinkedVolKey, ":")[1]
			SnapID := strings.Split(snapIDtoLinkedVolKey, ":")[0]
			if _, ok := volIDToLinkedVolumes[volID]; ok {
				snapVXSession.TargetSourceSnapshotGenInfo = &types.TargetSourceSnapshotGenInfo{
					TargetDevice: volID,
					SourceDevice: sourceVolID,
					SnapshotName: SnapID,
				}
			}
		}
	}
	return snapVXSession
}

func returnSrcSnapshotGenInfo(volID string) []types.SourceSnapshotGenInfo {
	var srcSnapGenInfo []types.SourceSnapshotGenInfo

	for _, snapIDtoSnap := range Data.VolIDToSnapshots[volID] {
		timestamp, _ := strconv.ParseInt(snapIDtoSnap.Timestamp, 10, 64)
		srcSnapGenInfo = append(srcSnapGenInfo, types.SourceSnapshotGenInfo{
			SnapshotHeader: types.SnapshotHeader{
				Device:       volID,
				SnapshotName: snapIDtoSnap.Name,
				Generation:   snapIDtoSnap.Generation,
				Timestamp:    timestamp,
			},
		})
	}

	return srcSnapGenInfo
}

func HandleHostGroup(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleHostGroup(w, r)
}

// /univmax/restapi/100/sloprovisioning/symmetrix/{symid}/hostgroup/{id}
// /univmax/restapi/100/sloprovisioning/symmetrix/{symid}/hostgroup
func handleHostGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostGroupID := vars["id"]
	switch r.Method {

	case http.MethodGet:
		if InducedErrors.GetHostGroupError {
			writeError(w, "Error retrieving HostGroup: induced error", http.StatusRequestTimeout)
			return
		} else if InducedErrors.GetHostGroupListError {
			writeError(w, "Error retrieving HostGroupList: induced error", http.StatusRequestTimeout)
			return
		}
		returnHostGroup(w, hostGroupID)

	case http.MethodPost:
		if InducedErrors.CreateHostGroupError {
			writeError(w, "Error creating HostGroup: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createHostGroupParam := &types.CreateHostGroupParam{}
		err := decoder.Decode(createHostGroupParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		addHostGroup(createHostGroupParam.HostGroupID, createHostGroupParam.HostIDs, createHostGroupParam.HostFlags) // #nosec G20
		returnHostGroup(w, createHostGroupParam.HostGroupID)

	case http.MethodPut:
		if hasError(&InducedErrors.UpdateHostGroupError) {
			writeError(w, "Error updating HostGroup: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updateHostGroupParam := &types.UpdateHostGroupParam{}
		err := decoder.Decode(updateHostGroupParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		updateHostGroupFromParams(hostGroupID, updateHostGroupParam)
		returnHostGroup(w, hostGroupID)

	case http.MethodDelete:
		if InducedErrors.DeleteHostGroupError {
			writeError(w, "Error deleting HostGroup: induced error", http.StatusRequestTimeout)
			return
		}
		err := removeHostGroup(hostGroupID)
		if err != nil {
			writeError(w, "error removeHostGroup", http.StatusBadRequest)
			return
		}
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// ReturnHostGroup - Returns a hostGroup from cache
func ReturnHostGroup(w http.ResponseWriter, hostGroupID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnHostGroup(w, hostGroupID)
}

func returnHostGroup(w http.ResponseWriter, hostGroupID string) {
	if hostGroupID != "" {
		if hostGroup, ok := Data.HostGroupIDToHostGroup[hostGroupID]; ok {
			writeJSON(w, hostGroup)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		hostGroupIDs := make([]string, 0)
		for k := range Data.HostGroupIDToHostGroup {
			hostGroupIDs = append(hostGroupIDs, k)
		}
		hostgroupIDList := &types.HostGroupList{
			HostGroupIDs: hostGroupIDs,
		}
		writeJSON(w, hostgroupIDList)
	}
}

func AddHostGroup(hostGroupID string, hostIDs []string, hostFlags *types.HostFlags) (*types.HostGroup, error) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return addHostGroup(hostGroupID, hostIDs, hostFlags)
}

// AddHostGroup - Adds a host group to the mock data cache
func addHostGroup(hostGroupID string, hostIDs []string, hostFlags *types.HostFlags) (*types.HostGroup, error) {
	if _, ok := Data.HostGroupIDToHostGroup[hostGroupID]; ok {
		return nil, errors.New("error! Host Group already exists")
	}
	newHostGroup(hostGroupID, hostIDs, hostFlags)
	return Data.HostGroupIDToHostGroup[hostGroupID], nil
}

func newHostGroup(hostGroupID string, hostIDs []string, hostFlags *types.HostFlags) {
	hostSummaries := []types.HostSummary{}

	for _, hostID := range hostIDs {
		Host := types.HostSummary{
			HostID: hostID,
		}
		hostSummaries = append(hostSummaries, Host)
	}

	hostGroup := &types.HostGroup{
		HostGroupID: hostGroupID,
		Hosts:       hostSummaries,
	}

	if hostFlags != nil {
		handleFlags(hostGroup, hostFlags)
	}

	Data.HostGroupIDToHostGroup[hostGroupID] = hostGroup
}

// RemoveHostGroup - Removes hostGroup from mock cache
func RemoveHostGroup(hostGroupID string) error {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	return removeHostGroup(hostGroupID)
}

// removeHostGroup - Remove a hostGroup from the mock data cache
func removeHostGroup(hostGroupID string) error {
	_, ok := Data.HostGroupIDToHostGroup[hostGroupID]
	if !ok {
		return errors.New("error! Host doesn't exist")
	}
	Data.HostGroupIDToHostGroup[hostGroupID] = nil
	return nil
}

func UpdateHostGroupFromParams(hostGroupID string, updateParams *types.UpdateHostGroupParam) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	updateHostGroupFromParams(hostGroupID, updateParams)
}

// UpdateHostGroupFromParams - Updates HostGroup given an UpdateHostGroupParam payload
func updateHostGroupFromParams(hostGroupID string, updateParams *types.UpdateHostGroupParam) {
	updateHostGroup(hostGroupID, updateParams.EditHostGroupAction) // #nosec G20
}

// updateHostGroup - Update HostGroup
func updateHostGroup(hostGroupID string, editPayload *types.EditHostGroupActionParams) (*types.HostGroup, error) {
	hostGroup, ok := Data.HostGroupIDToHostGroup[hostGroupID]
	if !ok {
		return nil, fmt.Errorf("error! HostGroup %s does not exist", hostGroupID)
	}

	if editPayload.RemoveHostParam != nil {
		hostSummaries := []types.HostSummary{}
		for _, host := range hostGroup.Hosts {
			if !stringInSlice(host.HostID, editPayload.RemoveHostParam.Host) {
				hostSummary := types.HostSummary{
					HostID: host.HostID,
				}
				hostSummaries = append(hostSummaries, hostSummary)
			}
		}
		hostGroup.Hosts = hostSummaries
	}

	if editPayload.AddHostParam != nil {
		for _, host := range editPayload.AddHostParam.Host {
			hostSummary := types.HostSummary{
				HostID: host,
			}
			hostGroup.Hosts = append(hostGroup.Hosts, hostSummary)
		}
	}

	if editPayload.SetHostGroupFlags != nil {
		handleFlags(hostGroup, editPayload.SetHostGroupFlags.HostFlags)
	}

	if editPayload.RenameHostGroupParam != nil {
		hostGroupID = editPayload.RenameHostGroupParam.NewHostGroupName
		hostGroup.HostGroupID = hostGroupID
	}

	// Update the HostGroup mapping with the update HostGroup
	Data.HostGroupIDToHostGroup[hostGroupID] = hostGroup
	return hostGroup, nil
}

func handleFlags(hostGroup *types.HostGroup, flagPayload *types.HostFlags) {
	var enabledFlags, disabledFlags []string
	if flagPayload.VolumeSetAddressing.Override {
		if flagPayload.VolumeSetAddressing.Enabled {
			enabledFlags = append(enabledFlags, "Volume_Set_Addressing(V)")
		}
		disabledFlags = append(disabledFlags, "Volume_Set_Addressing(V)")
	}

	if flagPayload.AvoidResetBroadcast.Override {
		if flagPayload.AvoidResetBroadcast.Enabled {
			enabledFlags = append(enabledFlags, "Avoid_Reset_Broadcast(ARB)")
		}
		disabledFlags = append(disabledFlags, "Avoid_Reset_Broadcast(ARB)")
	}

	if flagPayload.DisableQResetOnUA.Override {
		if flagPayload.DisableQResetOnUA.Enabled {
			enabledFlags = append(enabledFlags, "Disable_Q_Reset_on_UA(D)")
		}
		disabledFlags = append(disabledFlags, "Disable_Q_Reset_on_UA(D)")
	}

	if flagPayload.EnvironSet.Override {
		if flagPayload.EnvironSet.Enabled {
			enabledFlags = append(enabledFlags, "Environ_Set(E)")
		}
		disabledFlags = append(disabledFlags, "Environ_Set(E)")
	}

	if flagPayload.OpenVMS.Override {
		if flagPayload.OpenVMS.Enabled {
			enabledFlags = append(enabledFlags, "OpenVMS(OVMS)")
		}
		disabledFlags = append(disabledFlags, "OpenVMS(OVMS)")
	}

	if flagPayload.SCSI3.Override {
		if flagPayload.SCSI3.Enabled {
			enabledFlags = append(enabledFlags, "SCSI_3(SC3)")
		}
		disabledFlags = append(disabledFlags, "SCSI_3(SC3)")
	}

	if flagPayload.SCSISupport1.Override {
		if flagPayload.SCSISupport1.Enabled {
			enabledFlags = append(enabledFlags, "SCSI_Support1(OS2007)")
		}
		disabledFlags = append(disabledFlags, "SCSI_Support1(OS2007)")
	}

	if flagPayload.Spc2ProtocolVersion.Override {
		if flagPayload.Spc2ProtocolVersion.Enabled {
			enabledFlags = append(enabledFlags, "SPC2_Protocol_Version(SPC2)")
		}
		disabledFlags = append(disabledFlags, "SPC2_Protocol_Version(SPC2)")
	}

	enabledFlag := strings.Join(enabledFlags, ",")
	disabledFlag := strings.Join(disabledFlags, ",")

	hostGroup.EnabledFlags = enabledFlag
	hostGroup.DisabledFlags = disabledFlag

	if flagPayload.ConsistentLUN {
		hostGroup.ConsistentLun = true
	}
}

// /univmax/restapi/100/file/symmetrix/{symID}/nas_server/
// /univmax/restapi/100/file/symmetrix/{symID}/nas_server/{nasID}
func HandleNASServer(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleNASServer(w, r)
}

func handleNASServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nasID := vars["nasID"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetNASServerListError {
			writeError(w, "Error retrieving NAS server: induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.GetNASServerError {
			writeError(w, "Error retrieving NAS server: induced error", http.StatusRequestTimeout)
			return
		}
		returnNASServer(w, nasID)
	case http.MethodPut:
		if InducedErrors.UpdateNASServerError {
			writeError(w, "Error updating NAS server: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		modifyNASServerParam := &types.ModifyNASServer{}
		err := decoder.Decode(modifyNASServerParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		updateNASServer(nasID, modifyNASServerParam.Name)
		returnNASServer(w, nasID)
	case http.MethodDelete:
		if InducedErrors.DeleteNASServerError {
			writeError(w, "Error deleting NAS server: induced error", http.StatusRequestTimeout)
			return
		}
		removeNASServer(w, nasID)
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// UpdateNASServer updates NAS server
func UpdateNASServer(nasID string, payload types.ModifyNASServer) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	newName := payload.Name
	updateNASServer(nasID, newName)
}

func updateNASServer(nasID, newName string) {
	nas := Data.NASServerIDToNASServer[nasID]
	nas.Name = newName
	Data.NASServerIDToNASServer[nasID] = nas
}

// ReturnNASServer returns NAS server object
func ReturnNASServer(w http.ResponseWriter, nasID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnNASServer(w, nasID)
}

func returnNASServer(w http.ResponseWriter, nasID string) {
	if nasID == "" {
		// return NAS ServerList
		nasIter := &types.NASServerIterator{
			Entries: []types.NASServerList{
				{
					ID:   DefaultNASServerID,
					Name: "nas-1",
				},
				{
					ID:   "54xxx7a6-03b5-xxx-xxx-0zzzz8200209",
					Name: "nas-2",
				},
			},
		}
		writeJSON(w, nasIter)
		return
	}
	var nasServer *types.NASServer
	found := false
	for _, ns := range Data.NASServerIDToNASServer {
		if ns.ID == nasID {
			found = true
			nasServer = ns
		}
	}
	if !found {
		writeError(w, "NASServer cannot be found", http.StatusNotFound)
		return
	}
	writeJSON(w, nasServer)
}

// /univmax/restapi/100/file/symmetrix/{symID}/nfs_export/
// /univmax/restapi/100/file/symmetrix/{symID}/nfs_export/{nfsID}
func HandleNFSExport(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleNFSExport(w, r)
}

func handleNFSExport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nfsID := vars["nfsID"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetNFSExportListError {
			writeError(w, "Error retrieving NFS Export: induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.GetNFSExportError {
			writeError(w, "Error retrieving NFS Export: induced error", http.StatusNotFound)
			return
		}
		returnNFSExport(w, nfsID)
	case http.MethodPost:
		if InducedErrors.CreateNFSExportError {
			writeError(w, "Error creating NFS Export: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createNFSExportParam := &types.CreateNFSExport{}
		err := decoder.Decode(createNFSExportParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		addNewNFSExport("id-3", createNFSExportParam.Name)
		returnNFSExport(w, "id-3")
	case http.MethodPut:
		if InducedErrors.UpdateNFSExportError {
			writeError(w, "Error updating NFS Export: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		modifyNFSExportParam := &types.ModifyNFSExport{}
		err := decoder.Decode(modifyNFSExportParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		updateNFSExport(nfsID, modifyNFSExportParam.Name)
		returnNFSExport(w, nfsID)
	case http.MethodDelete:
		if InducedErrors.DeleteNFSExportError {
			writeError(w, "Error deleting NFS Export: induced error", http.StatusRequestTimeout)
			return
		}
		removeNFSExport(w, nfsID)
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// RemoveNFSExport removes NFS export object
func RemoveNFSExport(w http.ResponseWriter, nfsID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeNFSExport(w, nfsID)
}

func removeNFSExport(w http.ResponseWriter, nfsID string) {
	_, ok := Data.NFSExportIDToNFSExport[nfsID]
	if !ok {
		writeError(w, "error! fileSystem doesn't exist", http.StatusNotFound)
	}
	delete(Data.NFSExportIDToNFSExport, nfsID)
	return
}

// UpdateNFSExport updates NFS Export
func UpdateNFSExport(id string, payload types.ModifyNFSExport) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	newName := payload.Name
	updateNFSExport(id, newName)
}

func updateNFSExport(nfsID, newName string) {
	nfs := Data.NFSExportIDToNFSExport[nfsID]
	nfs.Name = newName
	Data.NFSExportIDToNFSExport[nfsID] = nfs
}

// ReturnNFSExport NFS export
func ReturnNFSExport(w http.ResponseWriter, nfsID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnNFSExport(w, nfsID)
}

func returnNFSExport(w http.ResponseWriter, nfsID string) {
	if nfsID == "" {
		// return NFS Export list
		nfsExportIter := &types.NFSExportIterator{
			ResultList: types.NFSExportList{
				NFSExportList: []types.NFSExportIDName{
					{
						ID:   "64xxx7a6-03b5-xxx-xxx-0zzzz8200208",
						Name: "nfs-ds-1",
					},
					{
						ID:   "64xxx7a6-03b5-xxx-xxx-0zzzz8200209",
						Name: "nfs-ds-2",
					},
				},
				From: 1,
				To:   2,
			},
			ID:             "52248851-fd6b-42c8-b7c7-2a9c0e40441a_0",
			Count:          2,
			ExpirationTime: 1688114398468,
			MaxPageSize:    1000,
		}
		writeJSON(w, nfsExportIter)
		return
	}
	var nfsExport *types.NFSExport
	found := false
	for _, nfs := range Data.NFSExportIDToNFSExport {
		if nfs.ID == nfsID {
			found = true
			nfsExport = nfs
		}
	}
	if !found {
		writeError(w, "NFSExport cannot be found", http.StatusNotFound)
		return
	}
	writeJSON(w, nfsExport)
}

// /univmax/restapi/100/file/symmetrix/{symID}/file_system/
// /univmax/restapi/100/file/symmetrix/{symID}/file_system/{fsID}
func HandleFileSystem(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleFileSystem(w, r)
}

func handleFileSystem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fsID := vars["fsID"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetFileSystemListError {
			writeError(w, "Error retrieving file systems: induced error", http.StatusRequestTimeout)
			return
		}
		if InducedErrors.GetFileSystemError {
			writeError(w, "Error retrieving file system, Could not find: induced error", http.StatusNotFound)
			return
		}
		if fsID == "" {
			// send in a list of file Syste
			// Here we want a volume iterator.
			queryParams := r.URL.Query()
			fileIdentifier := queryParams.Get("name")
			if fileIdentifier != "" {
				fsIDNameList := make([]types.FileSystemIDName, 0)
				for _, fs := range Data.FileSysIDToFileSystem {
					if fs.Name == fileIdentifier {
						fsIDName := types.FileSystemIDName{
							ID:   fs.ID,
							Name: fs.Name,
						}
						fsIDNameList = append(fsIDNameList, fsIDName)
					}
				}
				fileSysIter := &types.FileSystemIterator{
					ResultList: types.FileSystemList{
						FileSystemList: fsIDNameList,
						From:           1,
						To:             len(fsIDNameList),
					},
					ID:             "52248851-fd6b-42c8-b7c7-2a9c0e40441a_0",
					Count:          len(fsIDNameList),
					ExpirationTime: 1688114398468,
					MaxPageSize:    1000,
				}
				writeJSON(w, fileSysIter)
			} else {
				fileSysIter := &types.FileSystemIterator{
					ResultList: types.FileSystemList{
						FileSystemList: []types.FileSystemIDName{
							{
								ID:   DefaultFSID,
								Name: DefaultFSName,
							},
							{
								ID:   "64xxx7a6-03b5-xxx-xxx-0zzzz8200209",
								Name: "fs-ds-2",
							},
						},
						From: 1,
						To:   2,
					},
					ID:             "52248851-fd6b-42c8-b7c7-2a9c0e40441a_0",
					Count:          2,
					ExpirationTime: 1688114398468,
					MaxPageSize:    1000,
				}
				writeJSON(w, fileSysIter)
			}
		}
		returnFileSystem(w, fsID)
	case http.MethodPost:
		if InducedErrors.CreateFileSystemError {
			writeError(w, "Error creating file system: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		createFileSystemParam := &types.CreateFileSystem{}
		err := decoder.Decode(createFileSystemParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		id := strconv.Itoa(time.Now().Nanosecond())
		fsID := fmt.Sprintf("%s-%s-%d-%s", "649112ce-742b", "id", len(Data.FileSysIDToFileSystem), id)
		addNewFileSystem(fsID, createFileSystemParam.Name, createFileSystemParam.SizeTotal)
		returnFileSystem(w, fsID)
	case http.MethodPut:
		if InducedErrors.UpdateFileSystemError {
			writeError(w, "Error updating file system: induced error", http.StatusRequestTimeout)
			return
		}
		decoder := json.NewDecoder(r.Body)
		modifyFileSystemParam := &types.ModifyFileSystem{}
		err := decoder.Decode(modifyFileSystemParam)
		if err != nil {
			writeError(w, "InvalidJson", http.StatusBadRequest)
			return
		}
		updateFileSystem(fsID, modifyFileSystemParam.SizeTotal)
		returnFileSystem(w, fsID)
	case http.MethodDelete:
		if InducedErrors.DeleteFileSystemError {
			writeError(w, "Error deleting file system: induced error", http.StatusRequestTimeout)
			return
		}
		removeFileSystem(w, fsID)
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// ReturnFileSystem returns File System object
func ReturnFileSystem(w http.ResponseWriter, fsID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnFileSystem(w, fsID)
}

func returnFileSystem(w http.ResponseWriter, fsID string) {
	if fsID != "" {
		var fileSys *types.FileSystem
		found := false
		for _, fs := range Data.FileSysIDToFileSystem {
			if fs.ID == fsID {
				found = true
				fileSys = fs
			}
		}
		if !found {
			writeError(w, "FileSystem cannot be found", http.StatusNotFound)
			return
		}
		writeJSON(w, fileSys)
	}
}

// RemoveFileSystem removes File system from mock
func RemoveFileSystem(w http.ResponseWriter, fsID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeFileSystem(w, fsID)
}

func removeFileSystem(w http.ResponseWriter, fsID string) {
	_, ok := Data.FileSysIDToFileSystem[fsID]
	if !ok {
		writeError(w, "error! fileSystem doesn't exist", http.StatusNotFound)
	}
	delete(Data.FileSysIDToFileSystem, fsID)
	return
}

// RemoveNASServer removes NAS server
func RemoveNASServer(w http.ResponseWriter, nasID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	removeNASServer(w, nasID)
}

func removeNASServer(w http.ResponseWriter, nasID string) {
	_, ok := Data.NASServerIDToNASServer[nasID]
	if !ok {
		writeError(w, "error! fileSystem doesn't exist", http.StatusNotFound)
	}
	delete(Data.NASServerIDToNASServer, nasID)
	return
}

// AddNewNASServer adds new NAS server into mock
func AddNewNASServer(id, name string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addNewNASServer(id, name)
}

func addNewNASServer(nasID, nasName string) {
	Data.NASServerIDToNASServer[nasID] = newNASServer(nasID, nasName)
}

func newNASServer(nasID, nasName string) *types.NASServer {
	return &types.NASServer{
		ID:                          nasID,
		Health:                      types.Health{HealthStatus: "ok"},
		Name:                        nasName,
		StorageResourcePool:         "SRP_1",
		OperationalStatus:           "Started",
		PrimaryNode:                 "1",
		BackupNode:                  "2",
		Cluster:                     "64860dc2-571a-15d7",
		ProductionMode:              true,
		CurrentUnixDirectoryService: "Nonw",
		UsernameTranslation:         false,
		AutoUserMapping:             false,
		FileInterfaces:              nil,
		NFSServer:                   "6488552e-b863-6c30-xxx-1234xxxx",
		RootFSWWN:                   "wwn_0x540a001d01234567890",
		ConfigFSWWN:                 "wwn_0x540a001d01234567890",
	}
}

// AddNewNFSExport new NFS Export into mock
func AddNewNFSExport(id, name string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addNewNFSExport(id, name)
}

func addNewNFSExport(nfsID, nfsName string) {
	Data.NFSExportIDToNFSExport[nfsID] = newNFSExport(nfsID, nfsName)
}

func newNFSExport(nfsID, nfsName string) *types.NFSExport {
	return &types.NFSExport{
		ID:                 nfsID,
		Type:               "Nfs_VMWare",
		Filesystem:         "id1",
		NASServer:          "id1",
		Name:               nfsName,
		Path:               fmt.Sprintf("/%s", nfsName),
		Description:        "mock nfs export",
		DefaultAccess:      "NoAccess",
		MinSecurity:        "Sys",
		NFSOwnerUsername:   "root",
		NoAccessHosts:      nil,
		ReadOnlyHosts:      nil,
		ReadOnlyRootHosts:  nil,
		ReadWriteHosts:     nil,
		ReadWriteRootHosts: []string{"172.125.0.123"},
		AnonymousUID:       -2,
		AnonymousGID:       -2,
		NoSUID:             false,
	}
}

// AddNewFileSystem adds a new file system into mock
func AddNewFileSystem(id, name string, sizeInMiB int64) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addNewFileSystem(id, name, sizeInMiB)
}

func addNewFileSystem(fsID, fsName string, sizeInMiB int64) {
	Data.FileSysIDToFileSystem[fsID] = newFileSystem(fsID, fsName, sizeInMiB)
}

// UpdateFileSystem updates an existing FileSystem
func UpdateFileSystem(id string, payload types.ModifyFileSystem) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	sizeInMiB := payload.SizeTotal
	updateFileSystem(id, sizeInMiB)
}

func updateFileSystem(fsID string, sizeInMiB int64) {
	fs := Data.FileSysIDToFileSystem[fsID]
	fs.SizeTotal = sizeInMiB
	Data.FileSysIDToFileSystem[fsID] = fs
}

func newFileSystem(fsID, fsName string, sizeInMiB int64) *types.FileSystem {
	return &types.FileSystem{
		ID:          fsID,
		ParentOID:   "00000-0000",
		Name:        fsName,
		StorageWWN:  "0x540a001d000001740000976007271234",
		ExportFSID:  "648af4d7-92b5-d347-989c-026048200208",
		Description: "mock fs system",
		SizeTotal:   sizeInMiB,
		SizeUsed:    0,
		Health: struct {
			HealthStatus string `json:"health_status"`
		}{HealthStatus: "OK"},
		ReadOnly:                  false,
		FsType:                    "General",
		MountState:                "Mounted",
		AccessPolicy:              "Native",
		LockingPolicy:             "Advisory",
		FolderRenamePolicy:        "SMB_Rename_Forbidden",
		HostIOBlockSize:           8192,
		NasServer:                 "id1",
		SmbSyncWrites:             false,
		SmbOpLocks:                false,
		SmbNoNotify:               false,
		SmbNotifyOnAccess:         false,
		SmbNotifyOnWrite:          false,
		SmbNotifyOnChangeDirDepth: 0,
		AsyncMtime:                false,
		FlrMode:                   "None",
		FlrMinRet:                 "0D",
		FlrDefRet:                 "0D",
		FlrMaxRet:                 "0D",
		FlrAutoLock:               false,
		FlrAutoDelete:             false,
		FlrPolicyInterval:         0,
		FlrEnabled:                false,
		FlrClockTime:              "",
		FlrMaxRetentionDate:       "",
		FlrHasProtectedFiles:      false,
		QuotaConfig: &types.QuotaConfig{
			QuotaEnabled:     false,
			GracePeriod:      60480,
			DefaultHardLimit: 0,
			DefaultSoftLimit: 0,
		},
		EventNotifications: "",
		InfoThreshold:      0,
		HighThreshold:      75,
		WarningThreshold:   95,
		ServiceLevel:       "Optimized",
		DataReduction:      true,
	}
}

// AddNewFileInterface adds a new file interface into mock
func AddNewFileInterface(id, name string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	addNewFileInterface(id, name)
}

func addNewFileInterface(interfaceID, interfaceName string) {
	Data.FileIntIDtoFileInterface[interfaceID] = newFileInterface(interfaceID, interfaceName)
}

func newFileInterface(interfaceID, interfaceName string) *types.FileInterface {
	return &types.FileInterface{
		ID:         interfaceID,
		NasServer:  DefaultNASServerID,
		NetDevice:  "eth-1",
		MacAddress: "01:01:ab:01:01:zx",
		IPAddress:  "100.125.0.109",
		Netmask:    "255.255.255.0",
		Gateway:    "172.125.0.1",
		VlanID:     0,
		Name:       interfaceName,
		Role:       "Production",
		IsDisabled: false,
		Override:   false,
	}
}

// /univmax/restapi/100/file/symmetrix/{symID}//file_interface/
// /univmax/restapi/100/file/symmetrix/file_interface/{interfaceID}
func HandleFileInterface(w http.ResponseWriter, r *http.Request) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	handleFileInterface(w, r)
}

func handleFileInterface(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	interfaceID := vars["interfaceID"]
	switch r.Method {
	case http.MethodGet:
		if InducedErrors.GetFileInterfaceError {
			writeError(w, "Error retrieving FileSystemInterface: induced error", http.StatusNotFound)
			return
		}
		returnFileInterface(w, interfaceID)
	default:
		writeError(w, "Invalid Method", http.StatusBadRequest)
	}
}

// ReturnFileInterface returns File Interface object
func ReturnFileInterface(w http.ResponseWriter, interfaceID string) {
	mockCacheMutex.Lock()
	defer mockCacheMutex.Unlock()
	returnFileInterface(w, interfaceID)
}

func returnFileInterface(w http.ResponseWriter, interfaceID string) {
	if interfaceID != "" {
		if fi, ok := Data.FileIntIDtoFileInterface[interfaceID]; ok {
			writeJSON(w, fi)
		} else {
			writeError(w, "Could not find FileInterface", http.StatusNotFound)
			return
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
