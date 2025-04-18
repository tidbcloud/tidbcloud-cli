/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
)

// checks if the V1beta1NodeInstance type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1NodeInstance{}

// V1beta1NodeInstance All fields are output only.
type V1beta1NodeInstance struct {
	Name          *string                        `json:"name,omitempty"`
	ClusterId     *string                        `json:"clusterId,omitempty"`
	InstanceId    *string                        `json:"instanceId,omitempty"`
	ComponentType *Dedicatedv1beta1ComponentType `json:"componentType,omitempty"`
	// the state of the instance, e.g. \"Available\".
	State *V1beta1NodeInstanceState `json:"state,omitempty"`
	// the cpu size of the instance, e.g. 2.
	VCpu *int32 `json:"vCpu,omitempty"`
	// the memory size of the instance, e.g. 8.
	MemorySizeGi *int32 `json:"memorySizeGi,omitempty"`
	// the availability zone of the instance, e.g. \"us-west1-a\".
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
	// the storage size of the instance, e.g. 100.
	StorageSizeGi            *int32         `json:"storageSizeGi,omitempty"`
	TidbNodeGroupId          NullableString `json:"tidbNodeGroupId,omitempty"`
	TidbNodeGroupDisplayName NullableString `json:"tidbNodeGroupDisplayName,omitempty"`
	IsDefaultTidbNodeGroup   NullableBool   `json:"isDefaultTidbNodeGroup,omitempty"`
	// Only available for instances which have storage. If raft_store_iops is not set, the default IOPS of raft store will be used.
	RaftStoreIops NullableInt32 `json:"raftStoreIops,omitempty"`
	// Only available for instances which have storage.
	StorageType          *StorageNodeSettingStorageType `json:"storageType,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1NodeInstance V1beta1NodeInstance

// NewV1beta1NodeInstance instantiates a new V1beta1NodeInstance object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1NodeInstance() *V1beta1NodeInstance {
	this := V1beta1NodeInstance{}
	return &this
}

// NewV1beta1NodeInstanceWithDefaults instantiates a new V1beta1NodeInstance object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1NodeInstanceWithDefaults() *V1beta1NodeInstance {
	this := V1beta1NodeInstance{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1NodeInstance) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *V1beta1NodeInstance) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetInstanceId returns the InstanceId field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetInstanceId() string {
	if o == nil || IsNil(o.InstanceId) {
		var ret string
		return ret
	}
	return *o.InstanceId
}

// GetInstanceIdOk returns a tuple with the InstanceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetInstanceIdOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceId) {
		return nil, false
	}
	return o.InstanceId, true
}

// HasInstanceId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasInstanceId() bool {
	if o != nil && !IsNil(o.InstanceId) {
		return true
	}

	return false
}

// SetInstanceId gets a reference to the given string and assigns it to the InstanceId field.
func (o *V1beta1NodeInstance) SetInstanceId(v string) {
	o.InstanceId = &v
}

// GetComponentType returns the ComponentType field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetComponentType() Dedicatedv1beta1ComponentType {
	if o == nil || IsNil(o.ComponentType) {
		var ret Dedicatedv1beta1ComponentType
		return ret
	}
	return *o.ComponentType
}

// GetComponentTypeOk returns a tuple with the ComponentType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetComponentTypeOk() (*Dedicatedv1beta1ComponentType, bool) {
	if o == nil || IsNil(o.ComponentType) {
		return nil, false
	}
	return o.ComponentType, true
}

// HasComponentType returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasComponentType() bool {
	if o != nil && !IsNil(o.ComponentType) {
		return true
	}

	return false
}

// SetComponentType gets a reference to the given Dedicatedv1beta1ComponentType and assigns it to the ComponentType field.
func (o *V1beta1NodeInstance) SetComponentType(v Dedicatedv1beta1ComponentType) {
	o.ComponentType = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetState() V1beta1NodeInstanceState {
	if o == nil || IsNil(o.State) {
		var ret V1beta1NodeInstanceState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStateOk() (*V1beta1NodeInstanceState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta1NodeInstanceState and assigns it to the State field.
func (o *V1beta1NodeInstance) SetState(v V1beta1NodeInstanceState) {
	o.State = &v
}

// GetVCpu returns the VCpu field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetVCpu() int32 {
	if o == nil || IsNil(o.VCpu) {
		var ret int32
		return ret
	}
	return *o.VCpu
}

// GetVCpuOk returns a tuple with the VCpu field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetVCpuOk() (*int32, bool) {
	if o == nil || IsNil(o.VCpu) {
		return nil, false
	}
	return o.VCpu, true
}

// HasVCpu returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasVCpu() bool {
	if o != nil && !IsNil(o.VCpu) {
		return true
	}

	return false
}

// SetVCpu gets a reference to the given int32 and assigns it to the VCpu field.
func (o *V1beta1NodeInstance) SetVCpu(v int32) {
	o.VCpu = &v
}

// GetMemorySizeGi returns the MemorySizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetMemorySizeGi() int32 {
	if o == nil || IsNil(o.MemorySizeGi) {
		var ret int32
		return ret
	}
	return *o.MemorySizeGi
}

// GetMemorySizeGiOk returns a tuple with the MemorySizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetMemorySizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.MemorySizeGi) {
		return nil, false
	}
	return o.MemorySizeGi, true
}

// HasMemorySizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasMemorySizeGi() bool {
	if o != nil && !IsNil(o.MemorySizeGi) {
		return true
	}

	return false
}

// SetMemorySizeGi gets a reference to the given int32 and assigns it to the MemorySizeGi field.
func (o *V1beta1NodeInstance) SetMemorySizeGi(v int32) {
	o.MemorySizeGi = &v
}

// GetAvailabilityZone returns the AvailabilityZone field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetAvailabilityZone() string {
	if o == nil || IsNil(o.AvailabilityZone) {
		var ret string
		return ret
	}
	return *o.AvailabilityZone
}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetAvailabilityZoneOk() (*string, bool) {
	if o == nil || IsNil(o.AvailabilityZone) {
		return nil, false
	}
	return o.AvailabilityZone, true
}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasAvailabilityZone() bool {
	if o != nil && !IsNil(o.AvailabilityZone) {
		return true
	}

	return false
}

// SetAvailabilityZone gets a reference to the given string and assigns it to the AvailabilityZone field.
func (o *V1beta1NodeInstance) SetAvailabilityZone(v string) {
	o.AvailabilityZone = &v
}

// GetStorageSizeGi returns the StorageSizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetStorageSizeGi() int32 {
	if o == nil || IsNil(o.StorageSizeGi) {
		var ret int32
		return ret
	}
	return *o.StorageSizeGi
}

// GetStorageSizeGiOk returns a tuple with the StorageSizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStorageSizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.StorageSizeGi) {
		return nil, false
	}
	return o.StorageSizeGi, true
}

// HasStorageSizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasStorageSizeGi() bool {
	if o != nil && !IsNil(o.StorageSizeGi) {
		return true
	}

	return false
}

// SetStorageSizeGi gets a reference to the given int32 and assigns it to the StorageSizeGi field.
func (o *V1beta1NodeInstance) SetStorageSizeGi(v int32) {
	o.StorageSizeGi = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId.Get()) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId.Get()
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.TidbNodeGroupId.Get(), o.TidbNodeGroupId.IsSet()
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasTidbNodeGroupId() bool {
	if o != nil && o.TidbNodeGroupId.IsSet() {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given NullableString and assigns it to the TidbNodeGroupId field.
func (o *V1beta1NodeInstance) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId.Set(&v)
}

// SetTidbNodeGroupIdNil sets the value for TidbNodeGroupId to be an explicit nil
func (o *V1beta1NodeInstance) SetTidbNodeGroupIdNil() {
	o.TidbNodeGroupId.Set(nil)
}

// UnsetTidbNodeGroupId ensures that no value is present for TidbNodeGroupId, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetTidbNodeGroupId() {
	o.TidbNodeGroupId.Unset()
}

// GetTidbNodeGroupDisplayName returns the TidbNodeGroupDisplayName field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetTidbNodeGroupDisplayName() string {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName.Get()) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupDisplayName.Get()
}

// GetTidbNodeGroupDisplayNameOk returns a tuple with the TidbNodeGroupDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetTidbNodeGroupDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.TidbNodeGroupDisplayName.Get(), o.TidbNodeGroupDisplayName.IsSet()
}

// HasTidbNodeGroupDisplayName returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasTidbNodeGroupDisplayName() bool {
	if o != nil && o.TidbNodeGroupDisplayName.IsSet() {
		return true
	}

	return false
}

// SetTidbNodeGroupDisplayName gets a reference to the given NullableString and assigns it to the TidbNodeGroupDisplayName field.
func (o *V1beta1NodeInstance) SetTidbNodeGroupDisplayName(v string) {
	o.TidbNodeGroupDisplayName.Set(&v)
}

// SetTidbNodeGroupDisplayNameNil sets the value for TidbNodeGroupDisplayName to be an explicit nil
func (o *V1beta1NodeInstance) SetTidbNodeGroupDisplayNameNil() {
	o.TidbNodeGroupDisplayName.Set(nil)
}

// UnsetTidbNodeGroupDisplayName ensures that no value is present for TidbNodeGroupDisplayName, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetTidbNodeGroupDisplayName() {
	o.TidbNodeGroupDisplayName.Unset()
}

// GetIsDefaultTidbNodeGroup returns the IsDefaultTidbNodeGroup field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetIsDefaultTidbNodeGroup() bool {
	if o == nil || IsNil(o.IsDefaultTidbNodeGroup.Get()) {
		var ret bool
		return ret
	}
	return *o.IsDefaultTidbNodeGroup.Get()
}

// GetIsDefaultTidbNodeGroupOk returns a tuple with the IsDefaultTidbNodeGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetIsDefaultTidbNodeGroupOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.IsDefaultTidbNodeGroup.Get(), o.IsDefaultTidbNodeGroup.IsSet()
}

// HasIsDefaultTidbNodeGroup returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasIsDefaultTidbNodeGroup() bool {
	if o != nil && o.IsDefaultTidbNodeGroup.IsSet() {
		return true
	}

	return false
}

// SetIsDefaultTidbNodeGroup gets a reference to the given NullableBool and assigns it to the IsDefaultTidbNodeGroup field.
func (o *V1beta1NodeInstance) SetIsDefaultTidbNodeGroup(v bool) {
	o.IsDefaultTidbNodeGroup.Set(&v)
}

// SetIsDefaultTidbNodeGroupNil sets the value for IsDefaultTidbNodeGroup to be an explicit nil
func (o *V1beta1NodeInstance) SetIsDefaultTidbNodeGroupNil() {
	o.IsDefaultTidbNodeGroup.Set(nil)
}

// UnsetIsDefaultTidbNodeGroup ensures that no value is present for IsDefaultTidbNodeGroup, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetIsDefaultTidbNodeGroup() {
	o.IsDefaultTidbNodeGroup.Unset()
}

// GetRaftStoreIops returns the RaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetRaftStoreIops() int32 {
	if o == nil || IsNil(o.RaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.RaftStoreIops.Get()
}

// GetRaftStoreIopsOk returns a tuple with the RaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.RaftStoreIops.Get(), o.RaftStoreIops.IsSet()
}

// HasRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasRaftStoreIops() bool {
	if o != nil && o.RaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the RaftStoreIops field.
func (o *V1beta1NodeInstance) SetRaftStoreIops(v int32) {
	o.RaftStoreIops.Set(&v)
}

// SetRaftStoreIopsNil sets the value for RaftStoreIops to be an explicit nil
func (o *V1beta1NodeInstance) SetRaftStoreIopsNil() {
	o.RaftStoreIops.Set(nil)
}

// UnsetRaftStoreIops ensures that no value is present for RaftStoreIops, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetRaftStoreIops() {
	o.RaftStoreIops.Unset()
}

// GetStorageType returns the StorageType field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetStorageType() StorageNodeSettingStorageType {
	if o == nil || IsNil(o.StorageType) {
		var ret StorageNodeSettingStorageType
		return ret
	}
	return *o.StorageType
}

// GetStorageTypeOk returns a tuple with the StorageType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStorageTypeOk() (*StorageNodeSettingStorageType, bool) {
	if o == nil || IsNil(o.StorageType) {
		return nil, false
	}
	return o.StorageType, true
}

// HasStorageType returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasStorageType() bool {
	if o != nil && !IsNil(o.StorageType) {
		return true
	}

	return false
}

// SetStorageType gets a reference to the given StorageNodeSettingStorageType and assigns it to the StorageType field.
func (o *V1beta1NodeInstance) SetStorageType(v StorageNodeSettingStorageType) {
	o.StorageType = &v
}

func (o V1beta1NodeInstance) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1NodeInstance) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	if !IsNil(o.InstanceId) {
		toSerialize["instanceId"] = o.InstanceId
	}
	if !IsNil(o.ComponentType) {
		toSerialize["componentType"] = o.ComponentType
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.VCpu) {
		toSerialize["vCpu"] = o.VCpu
	}
	if !IsNil(o.MemorySizeGi) {
		toSerialize["memorySizeGi"] = o.MemorySizeGi
	}
	if !IsNil(o.AvailabilityZone) {
		toSerialize["availabilityZone"] = o.AvailabilityZone
	}
	if !IsNil(o.StorageSizeGi) {
		toSerialize["storageSizeGi"] = o.StorageSizeGi
	}
	if o.TidbNodeGroupId.IsSet() {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId.Get()
	}
	if o.TidbNodeGroupDisplayName.IsSet() {
		toSerialize["tidbNodeGroupDisplayName"] = o.TidbNodeGroupDisplayName.Get()
	}
	if o.IsDefaultTidbNodeGroup.IsSet() {
		toSerialize["isDefaultTidbNodeGroup"] = o.IsDefaultTidbNodeGroup.Get()
	}
	if o.RaftStoreIops.IsSet() {
		toSerialize["raftStoreIops"] = o.RaftStoreIops.Get()
	}
	if !IsNil(o.StorageType) {
		toSerialize["storageType"] = o.StorageType
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1NodeInstance) UnmarshalJSON(data []byte) (err error) {
	varV1beta1NodeInstance := _V1beta1NodeInstance{}

	err = json.Unmarshal(data, &varV1beta1NodeInstance)

	if err != nil {
		return err
	}

	*o = V1beta1NodeInstance(varV1beta1NodeInstance)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "instanceId")
		delete(additionalProperties, "componentType")
		delete(additionalProperties, "state")
		delete(additionalProperties, "vCpu")
		delete(additionalProperties, "memorySizeGi")
		delete(additionalProperties, "availabilityZone")
		delete(additionalProperties, "storageSizeGi")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "tidbNodeGroupDisplayName")
		delete(additionalProperties, "isDefaultTidbNodeGroup")
		delete(additionalProperties, "raftStoreIops")
		delete(additionalProperties, "storageType")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1NodeInstance struct {
	value *V1beta1NodeInstance
	isSet bool
}

func (v NullableV1beta1NodeInstance) Get() *V1beta1NodeInstance {
	return v.value
}

func (v *NullableV1beta1NodeInstance) Set(val *V1beta1NodeInstance) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1NodeInstance) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1NodeInstance) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1NodeInstance(val *V1beta1NodeInstance) *NullableV1beta1NodeInstance {
	return &NullableV1beta1NodeInstance{value: val, isSet: true}
}

func (v NullableV1beta1NodeInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1NodeInstance) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
