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

import types "github.com/dell/gopowermax/types/v90"

// Debug is a boolean, when enabled, that enables logging of send payloads, and other debug information. Default to false.
// It is set true by unit testing.
var Debug = false

// ConfigConnect is an argument structure that can be passed to Authenticate.
// It contains the Endpoint, API Version (which should not be used), Username, and Password.
type ConfigConnect struct {
	Endpoint string
	Version  string
	Username string
	Password string
}

const (
	// DefaultAPIVersion is the default API version you will get if not specificed to NewClientWithArgs.
	// The other supported versions are listed here.
	DefaultAPIVersion = "90"
	APIVersion90      = "90"
	APIVersion91      = "91"
)

// Pmax interface has all the externally available functions provided by the pmax client library for the Powermax accessed through Unisphere.
type Pmax interface {
	// Authenticate causes authentication and tests the connection
	Authenticate(configConnect *ConfigConnect) error

	// SLO provisioning are the methods for SLO provisioning. All the methods requre a
	// symID to identify the Symmetrix.

	// GetVolumeIDsIterator generates a VolumeIterator containing the ids of either all or a selected set volumes.
	// The volumeIdentifierMatch string can be used to find a specific volume, or if the like bool is set, all the
	// volumes containing match as part of their VolumeIdentifier.
	GetVolumeIDsIterator(symID string, volumeIdentifierMatch string, like bool) (*types.VolumeIterator, error)

	// GetVolumesInStorageGroupIterator returns a list of volumes for a given StorageGroup
	GetVolumesInStorageGroupIterator(symID string, storageGroupId string) (*types.VolumeIterator, error)

	// GetVolumeIDsIteraotrPage gets a page of volume ids from a Volume iterator.
	GetVolumeIDsIteratorPage(iter *types.VolumeIterator, from, to int) ([]string, error)

	// DeleteVolumeIDsIterator deletes a Volume iterator.
	DeleteVolumeIDsIterator(iter *types.VolumeIterator) error

	// GetVolumeIDList provides a simpler interface that returns a []string of volume ids
	// of volumes matching the volumeIdentifierMatch (and like) criteria. It is
	// implemented in terms of GetVolumeIDsIterator, GetVolumeIDsIteratorPage, and DeleteVolumeIDsIterator
	// and handles all the details of the iteration for you.
	GetVolumeIDList(symID string, volumeIdentifierMatch string, like bool) ([]string, error)

	// GetVolumeIDListInStorageGroup returns a list of volume IDs that are associated with the StorageGroup
	GetVolumeIDListInStorageGroup(symID string, storageGroupId string) ([]string, error)

	// GetVolumeById returns a Volume given the volumeID.
	GetVolumeByID(symID string, volumeID string) (*types.Volume, error)

	// GetStorageGroupIDList returns a list of all the StorageGroup ids.
	GetStorageGroupIDList(symID string) (*types.StorageGroupIDList, error)

	// GetStorageGroup returns a storage group given the StorageGroup id.
	GetStorageGroup(symID string, storageGroupID string) (*types.StorageGroup, error)

	// GetStoragePool returns a storage pool given the GetStoragePoolID and SymID.
	GetStoragePool(symID string, storagePoolID string) (*types.StoragePool, error)

	// CreateStorageGroup creates a storage group given the Storage group id
	// and returns the storage group object. The storage group can be configured for thick volumes as an option.
	CreateStorageGroup(symID string, storageGroupID string, srpID string, serviceLevel string, thickVolumes bool) (*types.StorageGroup, error)
	// UpdateStorageGroup updates a storage group (i.e. a PUT operation) and should support all the defined
	// operations (but many have not been tested).
	UpdateStorageGroup(symID string, storageGroupID string, payload *types.UpdateStorageGroupPayload) (*types.Job, error)

	// CreateVolumeInStorageGroup takes simplified input arguments to create a volume of a give name and size in a particular storage group.
	// This method creates a job and waits on the job to complete.
	CreateVolumeInStorageGroup(symID string, storageGroupID string, volumeName string, sizeInCylinders int) (*types.Volume, error)

	// DeleteStorageGroup deletes a storage group given a storage group id
	DeleteStorageGroup(symID string, storageGroupID string) error

	// DeleteMaskingView deletes a masking view given a masking view id
	DeleteMaskingView(symID string, maskingViewID string) error

	// Get the list of Storage Pools
	GetStoragePoolList(symid string) (*types.StoragePoolList, error)

	// Rename a Volume given the volumeID
	RenameVolume(symID string, volumeID string, newName string) (*types.Volume, error)

	// Add volume(s) asynchronously to a StorageGroup
	AddVolumesToStorageGroup(symID string, storageGroupID string, volumeIDs ...string) error

	// Remove volume(s) synchronously from a StorageGroup
	RemoveVolumesFromStorageGroup(symID string, storageGroupID string, volumeIDs ...string) (*types.StorageGroup, error)

	// Initiate a job to remove storage space from the volume.
	InitiateDeallocationOfTracksFromVolume(symID string, volumeID string) (*types.Job, error)

	// Deletes a volume
	DeleteVolume(symID string, volumeID string) error

	// GetMaskingViewList  returns a list of the MaskingView names.
	GetMaskingViewList(symid string) (*types.MaskingViewList, error)

	// GetMaskingViewByID returns a masking view given it's identifier (which is the name)
	GetMaskingViewByID(symid string, maskingViewID string) (*types.MaskingView, error)

	// GetMaskingViewConnections returns the connections of a masking view (optionally for a specific volume id.)
	// Here volume id is the 5 digit volume ID.
	GetMaskingViewConnections(symid string, maskingViewID string, volumeID string) ([]*types.MaskingViewConnection, error)

	// CreateMaskingView creates a masking view given the Masking view id, Storage group id,
	// host id and the port id and returns the masking view object
	CreateMaskingView(symID string, maskingViewID string, storageGroupID string, hostOrhostGroupID string, isHost bool, portGroupID string) (*types.MaskingView, error)

	// CreatePortGroup creates a port group given the Port Group id and a list of dir/port ids
	CreatePortGroup(symID string, portGroupID string, dirPorts []types.PortKey) (*types.PortGroup, error)

	// System
	GetSymmetrixIDList() (*types.SymmetrixIDList, error)
	GetSymmetrixByID(id string) (*types.Symmetrix, error)

	// GetJobIDList retrieves the list of jobs on a given Symmetrix.
	// If optional parameter statusQuery is a types.JobStatusRunning or similar string, will search for jobs
	// with a particular status.
	GetJobIDList(symID string, statusQuery string) ([]string, error)
	GetJobByID(symID string, jobID string) (*types.Job, error)
	WaitOnJobCompletion(symID string, jobID string) (*types.Job, error)
	JobToString(job *types.Job) string

	// GetPortGroupList returns a list of all the Port Group ids.
	GetPortGroupList(symID string, portGroupType string) (*types.PortGroupList, error)
	// GetPortGroupByID returns a port group given the PortGroup id.
	GetPortGroupByID(symID string, portGroupID string) (*types.PortGroup, error)

	// GetInitiatorList returns a list of all the Initiator ids based on filters supplied
	GetInitiatorList(symID string, initiatorHBA string, isISCSI bool, inHost bool) (*types.InitiatorList, error)
	// GetInitiatorByID returns an Initiator given the Initiator id.
	GetInitiatorByID(symID string, initID string) (*types.Initiator, error)

	// GetHostList returns a list of all the Host ids.
	GetHostList(symID string) (*types.HostList, error)
	// GetHostByID returns a Host given the Host id.
	GetHostByID(symID string, hostID string) (*types.Host, error)
	// CreateHost creates a host from a list of InitiatorIDs (and optional HostFlags) return returns a types.Host.
	// Initiator IDs do not contain the storage port designations, just the IQN string or FC WWN.
	// Initiator IDs cannot be a member of more than one host.
	CreateHost(symID string, hostID string, initiatorIDs []string, hostFlags *types.HostFlags) (*types.Host, error)
	// DeleteHost deletes a host given the hostID.
	DeleteHost(symID string, hostID string) error
	// UpdateHostInitiators will update the inititators
	UpdateHostInitiators(symID string, host *types.Host, initiatorIDs []string) (*types.Host, error)
	// GetDirectorIDList returns a list of directors
	GetDirectorIDList(symID string) (*types.DirectorIDList, error)
	// GetPortList returns a list of all the ports on a specified director/array.
	GetPortList(symID string, directorID string, query string) (*types.PortList, error)
	// GetPort returns port details.
	GetPort(symID string, directorID string, portID string) (*types.Port, error)
	// GetListOfTargetAddresses returns an array of all IP addresses which expose iscsi targets.
	GetListOfTargetAddresses(symID string) ([]string, error)

	// SetAllowedArrays sets the list of arrays which can be manipulated
	// an empty list will allow all arrays to be accessed
	SetAllowedArrays(arrays []string) error
	// GetAllowedArrays returns a slice of arrays that can be manipulated
	GetAllowedArrays() []string
	// IsAllowedArray checks to see if we can manipulate the specified array
	IsAllowedArray(array string) (bool, error)

	// GetSnapVolumeList returns a list of all snapshot volumes on the array.
	GetSnapVolumeList(symID string, queryParams types.QueryParams) (*types.SymVolumeList, error)
	// GetVolumeSnapInfo returns snapVx information associated with a volume.
	GetVolumeSnapInfo(symID string, volume string) (*types.SnapshotVolumeGeneration, error)
	// GetSnapshotInfo returns snapVx information of the specified volume
	GetSnapshotInfo(symID, volume, SnapID string) (*types.VolumeSnapshot, error)
	// CreateSnapshot creates a snapVx snapshot of a volume using the input parameters
	CreateSnapshot(symID string, SnapID string, sourceVolumeList []types.VolumeList, ttl int64) error
	//ModifySnapshot executes actions on a snapshot
	ModifySnapshot(symID string, sourceVol []types.VolumeList,
		targetVol []types.VolumeList, SnapID string, action string,
		newSnapID string, generation int64) error
	// DeleteSnapshot deletes a snapshot from a volume
	DeleteSnapshot(symID, SnapID string, sourceVolumes []types.VolumeList, generation int64) error
	// GetSnapshotGenerations returns a list of all the snapshot generation on a specific snapshot
	GetSnapshotGenerations(symID, volume, SnapID string) (*types.VolumeSnapshotGenerations, error)
	// GetSnapshotGenerationInfo returns the specific generation info related to a snapshot
	GetSnapshotGenerationInfo(symID, volume, SnapID string, generation int64) (*types.VolumeSnapshotGeneration, error)
	// GetReplicationCapabilities returns details about SnapVX and SRDF execution capabilities on the Symmetrix array
	GetReplicationCapabilities() (*types.SymReplicationCapabilities, error)
	// GetPrivVolumeByID returns a Volume structure given the symmetrix and volume ID (volume ID is in WWN format)
	GetPrivVolumeByID(symID string, volumeID string) (*types.VolumeResultPrivate, error)

	// Delete PortGroup
	DeletePortGroup(symID string, portGroupID string) error
	// Update PortGroup
	UpdatePortGroup(symID string, portGroupId string, ports []types.PortKey) (*types.PortGroup, error)

	// Expand the size of an existing volume
	ExpandVolume(symID string, volumeID string, newSizeGB int) (*types.Volume, error)
}
