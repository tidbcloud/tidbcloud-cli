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

// checks if the V1beta1ShowCloudProvidersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ShowCloudProvidersResponse{}

// V1beta1ShowCloudProvidersResponse struct for V1beta1ShowCloudProvidersResponse
type V1beta1ShowCloudProvidersResponse struct {
	CloudProviders []V1beta1RegionCloudProvider `json:"cloudProviders,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ShowCloudProvidersResponse V1beta1ShowCloudProvidersResponse

// NewV1beta1ShowCloudProvidersResponse instantiates a new V1beta1ShowCloudProvidersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ShowCloudProvidersResponse() *V1beta1ShowCloudProvidersResponse {
	this := V1beta1ShowCloudProvidersResponse{}
	return &this
}

// NewV1beta1ShowCloudProvidersResponseWithDefaults instantiates a new V1beta1ShowCloudProvidersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ShowCloudProvidersResponseWithDefaults() *V1beta1ShowCloudProvidersResponse {
	this := V1beta1ShowCloudProvidersResponse{}
	return &this
}

// GetCloudProviders returns the CloudProviders field value if set, zero value otherwise.
func (o *V1beta1ShowCloudProvidersResponse) GetCloudProviders() []V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProviders) {
		var ret []V1beta1RegionCloudProvider
		return ret
	}
	return o.CloudProviders
}

// GetCloudProvidersOk returns a tuple with the CloudProviders field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ShowCloudProvidersResponse) GetCloudProvidersOk() ([]V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProviders) {
		return nil, false
	}
	return o.CloudProviders, true
}

// HasCloudProviders returns a boolean if a field has been set.
func (o *V1beta1ShowCloudProvidersResponse) HasCloudProviders() bool {
	if o != nil && !IsNil(o.CloudProviders) {
		return true
	}

	return false
}

// SetCloudProviders gets a reference to the given []V1beta1RegionCloudProvider and assigns it to the CloudProviders field.
func (o *V1beta1ShowCloudProvidersResponse) SetCloudProviders(v []V1beta1RegionCloudProvider) {
	o.CloudProviders = v
}

func (o V1beta1ShowCloudProvidersResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ShowCloudProvidersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CloudProviders) {
		toSerialize["cloudProviders"] = o.CloudProviders
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1ShowCloudProvidersResponse) UnmarshalJSON(data []byte) (err error) {
	varV1beta1ShowCloudProvidersResponse := _V1beta1ShowCloudProvidersResponse{}

	err = json.Unmarshal(data, &varV1beta1ShowCloudProvidersResponse)

	if err != nil {
		return err
	}

	*o = V1beta1ShowCloudProvidersResponse(varV1beta1ShowCloudProvidersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "cloudProviders")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ShowCloudProvidersResponse struct {
	value *V1beta1ShowCloudProvidersResponse
	isSet bool
}

func (v NullableV1beta1ShowCloudProvidersResponse) Get() *V1beta1ShowCloudProvidersResponse {
	return v.value
}

func (v *NullableV1beta1ShowCloudProvidersResponse) Set(val *V1beta1ShowCloudProvidersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ShowCloudProvidersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ShowCloudProvidersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ShowCloudProvidersResponse(val *V1beta1ShowCloudProvidersResponse) *NullableV1beta1ShowCloudProvidersResponse {
	return &NullableV1beta1ShowCloudProvidersResponse{value: val, isSet: true}
}

func (v NullableV1beta1ShowCloudProvidersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ShowCloudProvidersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

