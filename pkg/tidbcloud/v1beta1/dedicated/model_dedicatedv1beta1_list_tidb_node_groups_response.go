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

// checks if the Dedicatedv1beta1ListTidbNodeGroupsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1ListTidbNodeGroupsResponse{}

// Dedicatedv1beta1ListTidbNodeGroupsResponse struct for Dedicatedv1beta1ListTidbNodeGroupsResponse
type Dedicatedv1beta1ListTidbNodeGroupsResponse struct {
	TidbNodeGroups []Dedicatedv1beta1TidbNodeGroup `json:"tidbNodeGroups,omitempty"`
	// The total number of TiDB groups that matched the query.
	TotalSize *int32 `json:"totalSize,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page. If this field is omitted, there are no subsequent pages.
	NextPageToken *string `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1ListTidbNodeGroupsResponse Dedicatedv1beta1ListTidbNodeGroupsResponse

// NewDedicatedv1beta1ListTidbNodeGroupsResponse instantiates a new Dedicatedv1beta1ListTidbNodeGroupsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1ListTidbNodeGroupsResponse() *Dedicatedv1beta1ListTidbNodeGroupsResponse {
	this := Dedicatedv1beta1ListTidbNodeGroupsResponse{}
	return &this
}

// NewDedicatedv1beta1ListTidbNodeGroupsResponseWithDefaults instantiates a new Dedicatedv1beta1ListTidbNodeGroupsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1ListTidbNodeGroupsResponseWithDefaults() *Dedicatedv1beta1ListTidbNodeGroupsResponse {
	this := Dedicatedv1beta1ListTidbNodeGroupsResponse{}
	return &this
}

// GetTidbNodeGroups returns the TidbNodeGroups field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetTidbNodeGroups() []Dedicatedv1beta1TidbNodeGroup {
	if o == nil || IsNil(o.TidbNodeGroups) {
		var ret []Dedicatedv1beta1TidbNodeGroup
		return ret
	}
	return o.TidbNodeGroups
}

// GetTidbNodeGroupsOk returns a tuple with the TidbNodeGroups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetTidbNodeGroupsOk() ([]Dedicatedv1beta1TidbNodeGroup, bool) {
	if o == nil || IsNil(o.TidbNodeGroups) {
		return nil, false
	}
	return o.TidbNodeGroups, true
}

// HasTidbNodeGroups returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) HasTidbNodeGroups() bool {
	if o != nil && !IsNil(o.TidbNodeGroups) {
		return true
	}

	return false
}

// SetTidbNodeGroups gets a reference to the given []Dedicatedv1beta1TidbNodeGroup and assigns it to the TidbNodeGroups field.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) SetTidbNodeGroups(v []Dedicatedv1beta1TidbNodeGroup) {
	o.TidbNodeGroups = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o Dedicatedv1beta1ListTidbNodeGroupsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1ListTidbNodeGroupsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.TidbNodeGroups) {
		toSerialize["tidbNodeGroups"] = o.TidbNodeGroups
	}
	if !IsNil(o.TotalSize) {
		toSerialize["totalSize"] = o.TotalSize
	}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1ListTidbNodeGroupsResponse) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1ListTidbNodeGroupsResponse := _Dedicatedv1beta1ListTidbNodeGroupsResponse{}

	err = json.Unmarshal(data, &varDedicatedv1beta1ListTidbNodeGroupsResponse)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1ListTidbNodeGroupsResponse(varDedicatedv1beta1ListTidbNodeGroupsResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "tidbNodeGroups")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1ListTidbNodeGroupsResponse struct {
	value *Dedicatedv1beta1ListTidbNodeGroupsResponse
	isSet bool
}

func (v NullableDedicatedv1beta1ListTidbNodeGroupsResponse) Get() *Dedicatedv1beta1ListTidbNodeGroupsResponse {
	return v.value
}

func (v *NullableDedicatedv1beta1ListTidbNodeGroupsResponse) Set(val *Dedicatedv1beta1ListTidbNodeGroupsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1ListTidbNodeGroupsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1ListTidbNodeGroupsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1ListTidbNodeGroupsResponse(val *Dedicatedv1beta1ListTidbNodeGroupsResponse) *NullableDedicatedv1beta1ListTidbNodeGroupsResponse {
	return &NullableDedicatedv1beta1ListTidbNodeGroupsResponse{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1ListTidbNodeGroupsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1ListTidbNodeGroupsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


