/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// checks if the TidbNodeGroupServiceCreateTidbNodeGroupRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbNodeGroupServiceCreateTidbNodeGroupRequest{}

// TidbNodeGroupServiceCreateTidbNodeGroupRequest struct for TidbNodeGroupServiceCreateTidbNodeGroupRequest
type TidbNodeGroupServiceCreateTidbNodeGroupRequest struct {
	Name *string `json:"name,omitempty"`
	// The unique ID of the TiDB group.
	TidbNodeGroupId *string `json:"tidbNodeGroupId,omitempty"`
	// The display name of the TiDB group.
	DisplayName *string `json:"displayName,omitempty"`
	// The number of TiDB nodes in the TiDB group.
	NodeCount            int32                                   `json:"nodeCount"`
	Endpoints            []Dedicatedv1beta1TidbNodeGroupEndpoint `json:"endpoints,omitempty"`
	NodeSpecKey          *string                                 `json:"nodeSpecKey,omitempty"`
	NodeSpecDisplayName  *string                                 `json:"nodeSpecDisplayName,omitempty"`
	IsDefaultGroup       *bool                                   `json:"isDefaultGroup,omitempty"`
	State                *Dedicatedv1beta1TidbNodeGroupState     `json:"state,omitempty"`
	NodeChangingProgress *ClusterNodeChangingProgress            `json:"nodeChangingProgress,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbNodeGroupServiceCreateTidbNodeGroupRequest TidbNodeGroupServiceCreateTidbNodeGroupRequest

// NewTidbNodeGroupServiceCreateTidbNodeGroupRequest instantiates a new TidbNodeGroupServiceCreateTidbNodeGroupRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbNodeGroupServiceCreateTidbNodeGroupRequest(nodeCount int32) *TidbNodeGroupServiceCreateTidbNodeGroupRequest {
	this := TidbNodeGroupServiceCreateTidbNodeGroupRequest{}
	this.NodeCount = nodeCount
	return &this
}

// NewTidbNodeGroupServiceCreateTidbNodeGroupRequestWithDefaults instantiates a new TidbNodeGroupServiceCreateTidbNodeGroupRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbNodeGroupServiceCreateTidbNodeGroupRequestWithDefaults() *TidbNodeGroupServiceCreateTidbNodeGroupRequest {
	this := TidbNodeGroupServiceCreateTidbNodeGroupRequest{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetName(v string) {
	o.Name = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		return nil, false
	}
	return o.TidbNodeGroupId, true
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasTidbNodeGroupId() bool {
	if o != nil && !IsNil(o.TidbNodeGroupId) {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given string and assigns it to the TidbNodeGroupId field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetNodeCount returns the NodeCount field value
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeCount, true
}

// SetNodeCount sets field value
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetNodeCount(v int32) {
	o.NodeCount = v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetEndpoints() []Dedicatedv1beta1TidbNodeGroupEndpoint {
	if o == nil || IsNil(o.Endpoints) {
		var ret []Dedicatedv1beta1TidbNodeGroupEndpoint
		return ret
	}
	return o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetEndpointsOk() ([]Dedicatedv1beta1TidbNodeGroupEndpoint, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []Dedicatedv1beta1TidbNodeGroupEndpoint and assigns it to the Endpoints field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetEndpoints(v []Dedicatedv1beta1TidbNodeGroupEndpoint) {
	o.Endpoints = v
}

// GetNodeSpecKey returns the NodeSpecKey field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeSpecKey() string {
	if o == nil || IsNil(o.NodeSpecKey) {
		var ret string
		return ret
	}
	return *o.NodeSpecKey
}

// GetNodeSpecKeyOk returns a tuple with the NodeSpecKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeSpecKeyOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecKey) {
		return nil, false
	}
	return o.NodeSpecKey, true
}

// HasNodeSpecKey returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasNodeSpecKey() bool {
	if o != nil && !IsNil(o.NodeSpecKey) {
		return true
	}

	return false
}

// SetNodeSpecKey gets a reference to the given string and assigns it to the NodeSpecKey field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetNodeSpecKey(v string) {
	o.NodeSpecKey = &v
}

// GetNodeSpecDisplayName returns the NodeSpecDisplayName field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeSpecDisplayName() string {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		var ret string
		return ret
	}
	return *o.NodeSpecDisplayName
}

// GetNodeSpecDisplayNameOk returns a tuple with the NodeSpecDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeSpecDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		return nil, false
	}
	return o.NodeSpecDisplayName, true
}

// HasNodeSpecDisplayName returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasNodeSpecDisplayName() bool {
	if o != nil && !IsNil(o.NodeSpecDisplayName) {
		return true
	}

	return false
}

// SetNodeSpecDisplayName gets a reference to the given string and assigns it to the NodeSpecDisplayName field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetNodeSpecDisplayName(v string) {
	o.NodeSpecDisplayName = &v
}

// GetIsDefaultGroup returns the IsDefaultGroup field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetIsDefaultGroup() bool {
	if o == nil || IsNil(o.IsDefaultGroup) {
		var ret bool
		return ret
	}
	return *o.IsDefaultGroup
}

// GetIsDefaultGroupOk returns a tuple with the IsDefaultGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetIsDefaultGroupOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDefaultGroup) {
		return nil, false
	}
	return o.IsDefaultGroup, true
}

// HasIsDefaultGroup returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasIsDefaultGroup() bool {
	if o != nil && !IsNil(o.IsDefaultGroup) {
		return true
	}

	return false
}

// SetIsDefaultGroup gets a reference to the given bool and assigns it to the IsDefaultGroup field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetIsDefaultGroup(v bool) {
	o.IsDefaultGroup = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetState() Dedicatedv1beta1TidbNodeGroupState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1TidbNodeGroupState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetStateOk() (*Dedicatedv1beta1TidbNodeGroupState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1TidbNodeGroupState and assigns it to the State field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetState(v Dedicatedv1beta1TidbNodeGroupState) {
	o.State = &v
}

// GetNodeChangingProgress returns the NodeChangingProgress field value if set, zero value otherwise.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeChangingProgress() ClusterNodeChangingProgress {
	if o == nil || IsNil(o.NodeChangingProgress) {
		var ret ClusterNodeChangingProgress
		return ret
	}
	return *o.NodeChangingProgress
}

// GetNodeChangingProgressOk returns a tuple with the NodeChangingProgress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) GetNodeChangingProgressOk() (*ClusterNodeChangingProgress, bool) {
	if o == nil || IsNil(o.NodeChangingProgress) {
		return nil, false
	}
	return o.NodeChangingProgress, true
}

// HasNodeChangingProgress returns a boolean if a field has been set.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) HasNodeChangingProgress() bool {
	if o != nil && !IsNil(o.NodeChangingProgress) {
		return true
	}

	return false
}

// SetNodeChangingProgress gets a reference to the given ClusterNodeChangingProgress and assigns it to the NodeChangingProgress field.
func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) SetNodeChangingProgress(v ClusterNodeChangingProgress) {
	o.NodeChangingProgress = &v
}

func (o TidbNodeGroupServiceCreateTidbNodeGroupRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbNodeGroupServiceCreateTidbNodeGroupRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbNodeGroupId) {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	toSerialize["nodeCount"] = o.NodeCount
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.NodeSpecKey) {
		toSerialize["nodeSpecKey"] = o.NodeSpecKey
	}
	if !IsNil(o.NodeSpecDisplayName) {
		toSerialize["nodeSpecDisplayName"] = o.NodeSpecDisplayName
	}
	if !IsNil(o.IsDefaultGroup) {
		toSerialize["isDefaultGroup"] = o.IsDefaultGroup
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.NodeChangingProgress) {
		toSerialize["nodeChangingProgress"] = o.NodeChangingProgress
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbNodeGroupServiceCreateTidbNodeGroupRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"nodeCount",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varTidbNodeGroupServiceCreateTidbNodeGroupRequest := _TidbNodeGroupServiceCreateTidbNodeGroupRequest{}

	err = json.Unmarshal(data, &varTidbNodeGroupServiceCreateTidbNodeGroupRequest)

	if err != nil {
		return err
	}

	*o = TidbNodeGroupServiceCreateTidbNodeGroupRequest(varTidbNodeGroupServiceCreateTidbNodeGroupRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "nodeCount")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "nodeSpecKey")
		delete(additionalProperties, "nodeSpecDisplayName")
		delete(additionalProperties, "isDefaultGroup")
		delete(additionalProperties, "state")
		delete(additionalProperties, "nodeChangingProgress")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest struct {
	value *TidbNodeGroupServiceCreateTidbNodeGroupRequest
	isSet bool
}

func (v NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) Get() *TidbNodeGroupServiceCreateTidbNodeGroupRequest {
	return v.value
}

func (v *NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) Set(val *TidbNodeGroupServiceCreateTidbNodeGroupRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbNodeGroupServiceCreateTidbNodeGroupRequest(val *TidbNodeGroupServiceCreateTidbNodeGroupRequest) *NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest {
	return &NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest{value: val, isSet: true}
}

func (v NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbNodeGroupServiceCreateTidbNodeGroupRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
