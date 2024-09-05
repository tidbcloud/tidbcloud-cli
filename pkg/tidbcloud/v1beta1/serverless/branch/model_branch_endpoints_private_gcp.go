/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package branch

import (
	"encoding/json"
)

// checks if the BranchEndpointsPrivateGCP type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BranchEndpointsPrivateGCP{}

// BranchEndpointsPrivateGCP Message for GCP Private Link Service.
type BranchEndpointsPrivateGCP struct {
	// Output Only. Target Service Account for Private Link Service.
	ServiceAttachmentName *string `json:"serviceAttachmentName,omitempty"`
	AdditionalProperties  map[string]interface{}
}

type _BranchEndpointsPrivateGCP BranchEndpointsPrivateGCP

// NewBranchEndpointsPrivateGCP instantiates a new BranchEndpointsPrivateGCP object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBranchEndpointsPrivateGCP() *BranchEndpointsPrivateGCP {
	this := BranchEndpointsPrivateGCP{}
	return &this
}

// NewBranchEndpointsPrivateGCPWithDefaults instantiates a new BranchEndpointsPrivateGCP object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBranchEndpointsPrivateGCPWithDefaults() *BranchEndpointsPrivateGCP {
	this := BranchEndpointsPrivateGCP{}
	return &this
}

// GetServiceAttachmentName returns the ServiceAttachmentName field value if set, zero value otherwise.
func (o *BranchEndpointsPrivateGCP) GetServiceAttachmentName() string {
	if o == nil || IsNil(o.ServiceAttachmentName) {
		var ret string
		return ret
	}
	return *o.ServiceAttachmentName
}

// GetServiceAttachmentNameOk returns a tuple with the ServiceAttachmentName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BranchEndpointsPrivateGCP) GetServiceAttachmentNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceAttachmentName) {
		return nil, false
	}
	return o.ServiceAttachmentName, true
}

// HasServiceAttachmentName returns a boolean if a field has been set.
func (o *BranchEndpointsPrivateGCP) HasServiceAttachmentName() bool {
	if o != nil && !IsNil(o.ServiceAttachmentName) {
		return true
	}

	return false
}

// SetServiceAttachmentName gets a reference to the given string and assigns it to the ServiceAttachmentName field.
func (o *BranchEndpointsPrivateGCP) SetServiceAttachmentName(v string) {
	o.ServiceAttachmentName = &v
}

func (o BranchEndpointsPrivateGCP) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BranchEndpointsPrivateGCP) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ServiceAttachmentName) {
		toSerialize["serviceAttachmentName"] = o.ServiceAttachmentName
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *BranchEndpointsPrivateGCP) UnmarshalJSON(data []byte) (err error) {
	varBranchEndpointsPrivateGCP := _BranchEndpointsPrivateGCP{}

	err = json.Unmarshal(data, &varBranchEndpointsPrivateGCP)

	if err != nil {
		return err
	}

	*o = BranchEndpointsPrivateGCP(varBranchEndpointsPrivateGCP)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "serviceAttachmentName")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableBranchEndpointsPrivateGCP struct {
	value *BranchEndpointsPrivateGCP
	isSet bool
}

func (v NullableBranchEndpointsPrivateGCP) Get() *BranchEndpointsPrivateGCP {
	return v.value
}

func (v *NullableBranchEndpointsPrivateGCP) Set(val *BranchEndpointsPrivateGCP) {
	v.value = val
	v.isSet = true
}

func (v NullableBranchEndpointsPrivateGCP) IsSet() bool {
	return v.isSet
}

func (v *NullableBranchEndpointsPrivateGCP) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBranchEndpointsPrivateGCP(val *BranchEndpointsPrivateGCP) *NullableBranchEndpointsPrivateGCP {
	return &NullableBranchEndpointsPrivateGCP{value: val, isSet: true}
}

func (v NullableBranchEndpointsPrivateGCP) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBranchEndpointsPrivateGCP) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
