/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
)

// checks if the PrivateGCP type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PrivateGCP{}

// PrivateGCP Message for GCP Private Service information.
type PrivateGCP struct {
	// Output_only. The target GCP service attachment name for private access.
	ServiceAttachmentName *string `json:"serviceAttachmentName,omitempty"`
	AdditionalProperties  map[string]interface{}
}

type _PrivateGCP PrivateGCP

// NewPrivateGCP instantiates a new PrivateGCP object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateGCP() *PrivateGCP {
	this := PrivateGCP{}
	return &this
}

// NewPrivateGCPWithDefaults instantiates a new PrivateGCP object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateGCPWithDefaults() *PrivateGCP {
	this := PrivateGCP{}
	return &this
}

// GetServiceAttachmentName returns the ServiceAttachmentName field value if set, zero value otherwise.
func (o *PrivateGCP) GetServiceAttachmentName() string {
	if o == nil || IsNil(o.ServiceAttachmentName) {
		var ret string
		return ret
	}
	return *o.ServiceAttachmentName
}

// GetServiceAttachmentNameOk returns a tuple with the ServiceAttachmentName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateGCP) GetServiceAttachmentNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceAttachmentName) {
		return nil, false
	}
	return o.ServiceAttachmentName, true
}

// HasServiceAttachmentName returns a boolean if a field has been set.
func (o *PrivateGCP) HasServiceAttachmentName() bool {
	if o != nil && !IsNil(o.ServiceAttachmentName) {
		return true
	}

	return false
}

// SetServiceAttachmentName gets a reference to the given string and assigns it to the ServiceAttachmentName field.
func (o *PrivateGCP) SetServiceAttachmentName(v string) {
	o.ServiceAttachmentName = &v
}

func (o PrivateGCP) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PrivateGCP) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ServiceAttachmentName) {
		toSerialize["serviceAttachmentName"] = o.ServiceAttachmentName
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PrivateGCP) UnmarshalJSON(data []byte) (err error) {
	varPrivateGCP := _PrivateGCP{}

	err = json.Unmarshal(data, &varPrivateGCP)

	if err != nil {
		return err
	}

	*o = PrivateGCP(varPrivateGCP)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "serviceAttachmentName")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivateGCP struct {
	value *PrivateGCP
	isSet bool
}

func (v NullablePrivateGCP) Get() *PrivateGCP {
	return v.value
}

func (v *NullablePrivateGCP) Set(val *PrivateGCP) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivateGCP) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivateGCP) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivateGCP(val *PrivateGCP) *NullablePrivateGCP {
	return &NullablePrivateGCP{value: val, isSet: true}
}

func (v NullablePrivateGCP) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivateGCP) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
