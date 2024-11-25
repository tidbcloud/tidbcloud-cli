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

// checks if the UpdateClusterRequestTidbNodeSettingTidbNodeGroup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateClusterRequestTidbNodeSettingTidbNodeGroup{}

// UpdateClusterRequestTidbNodeSettingTidbNodeGroup struct for UpdateClusterRequestTidbNodeSettingTidbNodeGroup
type UpdateClusterRequestTidbNodeSettingTidbNodeGroup struct {
	TidbNodeGroupId      *string       `json:"tidbNodeGroupId,omitempty"`
	NodeCount            NullableInt32 `json:"nodeCount,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UpdateClusterRequestTidbNodeSettingTidbNodeGroup UpdateClusterRequestTidbNodeSettingTidbNodeGroup

// NewUpdateClusterRequestTidbNodeSettingTidbNodeGroup instantiates a new UpdateClusterRequestTidbNodeSettingTidbNodeGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateClusterRequestTidbNodeSettingTidbNodeGroup() *UpdateClusterRequestTidbNodeSettingTidbNodeGroup {
	this := UpdateClusterRequestTidbNodeSettingTidbNodeGroup{}
	return &this
}

// NewUpdateClusterRequestTidbNodeSettingTidbNodeGroupWithDefaults instantiates a new UpdateClusterRequestTidbNodeSettingTidbNodeGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateClusterRequestTidbNodeSettingTidbNodeGroupWithDefaults() *UpdateClusterRequestTidbNodeSettingTidbNodeGroup {
	this := UpdateClusterRequestTidbNodeSettingTidbNodeGroup{}
	return &this
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		return nil, false
	}
	return o.TidbNodeGroupId, true
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) HasTidbNodeGroupId() bool {
	if o != nil && !IsNil(o.TidbNodeGroupId) {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given string and assigns it to the TidbNodeGroupId field.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId = &v
}

// GetNodeCount returns the NodeCount field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) GetNodeCount() int32 {
	if o == nil || IsNil(o.NodeCount.Get()) {
		var ret int32
		return ret
	}
	return *o.NodeCount.Get()
}

// GetNodeCountOk returns a tuple with the NodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.NodeCount.Get(), o.NodeCount.IsSet()
}

// HasNodeCount returns a boolean if a field has been set.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) HasNodeCount() bool {
	if o != nil && o.NodeCount.IsSet() {
		return true
	}

	return false
}

// SetNodeCount gets a reference to the given NullableInt32 and assigns it to the NodeCount field.
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) SetNodeCount(v int32) {
	o.NodeCount.Set(&v)
}

// SetNodeCountNil sets the value for NodeCount to be an explicit nil
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) SetNodeCountNil() {
	o.NodeCount.Set(nil)
}

// UnsetNodeCount ensures that no value is present for NodeCount, not even an explicit nil
func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) UnsetNodeCount() {
	o.NodeCount.Unset()
}

func (o UpdateClusterRequestTidbNodeSettingTidbNodeGroup) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateClusterRequestTidbNodeSettingTidbNodeGroup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.TidbNodeGroupId) {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId
	}
	if o.NodeCount.IsSet() {
		toSerialize["nodeCount"] = o.NodeCount.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) UnmarshalJSON(data []byte) (err error) {
	varUpdateClusterRequestTidbNodeSettingTidbNodeGroup := _UpdateClusterRequestTidbNodeSettingTidbNodeGroup{}

	err = json.Unmarshal(data, &varUpdateClusterRequestTidbNodeSettingTidbNodeGroup)

	if err != nil {
		return err
	}

	*o = UpdateClusterRequestTidbNodeSettingTidbNodeGroup(varUpdateClusterRequestTidbNodeSettingTidbNodeGroup)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "nodeCount")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup struct {
	value *UpdateClusterRequestTidbNodeSettingTidbNodeGroup
	isSet bool
}

func (v NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) Get() *UpdateClusterRequestTidbNodeSettingTidbNodeGroup {
	return v.value
}

func (v *NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) Set(val *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup(val *UpdateClusterRequestTidbNodeSettingTidbNodeGroup) *NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup {
	return &NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup{value: val, isSet: true}
}

func (v NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateClusterRequestTidbNodeSettingTidbNodeGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
