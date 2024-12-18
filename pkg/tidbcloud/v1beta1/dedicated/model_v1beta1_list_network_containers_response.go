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

// checks if the V1beta1ListNetworkContainersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ListNetworkContainersResponse{}

// V1beta1ListNetworkContainersResponse struct for V1beta1ListNetworkContainersResponse
type V1beta1ListNetworkContainersResponse struct {
	NetworkContainers []V1beta1NetworkContainer `json:"networkContainers,omitempty"`
	// The total number of network containers that matched the query.
	TotalSize *int32 `json:"totalSize,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page. If this field is omitted, there are no subsequent pages.
	NextPageToken        *string `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ListNetworkContainersResponse V1beta1ListNetworkContainersResponse

// NewV1beta1ListNetworkContainersResponse instantiates a new V1beta1ListNetworkContainersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ListNetworkContainersResponse() *V1beta1ListNetworkContainersResponse {
	this := V1beta1ListNetworkContainersResponse{}
	return &this
}

// NewV1beta1ListNetworkContainersResponseWithDefaults instantiates a new V1beta1ListNetworkContainersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ListNetworkContainersResponseWithDefaults() *V1beta1ListNetworkContainersResponse {
	this := V1beta1ListNetworkContainersResponse{}
	return &this
}

// GetNetworkContainers returns the NetworkContainers field value if set, zero value otherwise.
func (o *V1beta1ListNetworkContainersResponse) GetNetworkContainers() []V1beta1NetworkContainer {
	if o == nil || IsNil(o.NetworkContainers) {
		var ret []V1beta1NetworkContainer
		return ret
	}
	return o.NetworkContainers
}

// GetNetworkContainersOk returns a tuple with the NetworkContainers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNetworkContainersResponse) GetNetworkContainersOk() ([]V1beta1NetworkContainer, bool) {
	if o == nil || IsNil(o.NetworkContainers) {
		return nil, false
	}
	return o.NetworkContainers, true
}

// HasNetworkContainers returns a boolean if a field has been set.
func (o *V1beta1ListNetworkContainersResponse) HasNetworkContainers() bool {
	if o != nil && !IsNil(o.NetworkContainers) {
		return true
	}

	return false
}

// SetNetworkContainers gets a reference to the given []V1beta1NetworkContainer and assigns it to the NetworkContainers field.
func (o *V1beta1ListNetworkContainersResponse) SetNetworkContainers(v []V1beta1NetworkContainer) {
	o.NetworkContainers = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *V1beta1ListNetworkContainersResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNetworkContainersResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *V1beta1ListNetworkContainersResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *V1beta1ListNetworkContainersResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *V1beta1ListNetworkContainersResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNetworkContainersResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *V1beta1ListNetworkContainersResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *V1beta1ListNetworkContainersResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o V1beta1ListNetworkContainersResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ListNetworkContainersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.NetworkContainers) {
		toSerialize["networkContainers"] = o.NetworkContainers
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

func (o *V1beta1ListNetworkContainersResponse) UnmarshalJSON(data []byte) (err error) {
	varV1beta1ListNetworkContainersResponse := _V1beta1ListNetworkContainersResponse{}

	err = json.Unmarshal(data, &varV1beta1ListNetworkContainersResponse)

	if err != nil {
		return err
	}

	*o = V1beta1ListNetworkContainersResponse(varV1beta1ListNetworkContainersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "networkContainers")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ListNetworkContainersResponse struct {
	value *V1beta1ListNetworkContainersResponse
	isSet bool
}

func (v NullableV1beta1ListNetworkContainersResponse) Get() *V1beta1ListNetworkContainersResponse {
	return v.value
}

func (v *NullableV1beta1ListNetworkContainersResponse) Set(val *V1beta1ListNetworkContainersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ListNetworkContainersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ListNetworkContainersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ListNetworkContainersResponse(val *V1beta1ListNetworkContainersResponse) *NullableV1beta1ListNetworkContainersResponse {
	return &NullableV1beta1ListNetworkContainersResponse{value: val, isSet: true}
}

func (v NullableV1beta1ListNetworkContainersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ListNetworkContainersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
