/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package imp

import (
	"encoding/json"
	"fmt"
)

// checks if the OSSSourceAccessKey type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OSSSourceAccessKey{}

// OSSSourceAccessKey struct for OSSSourceAccessKey
type OSSSourceAccessKey struct {
	// The access key id.
	Id string `json:"id"`
	// The secret access key. This field is input-only.
	Secret               string `json:"secret"`
	AdditionalProperties map[string]interface{}
}

type _OSSSourceAccessKey OSSSourceAccessKey

// NewOSSSourceAccessKey instantiates a new OSSSourceAccessKey object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOSSSourceAccessKey(id string, secret string) *OSSSourceAccessKey {
	this := OSSSourceAccessKey{}
	this.Id = id
	this.Secret = secret
	return &this
}

// NewOSSSourceAccessKeyWithDefaults instantiates a new OSSSourceAccessKey object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOSSSourceAccessKeyWithDefaults() *OSSSourceAccessKey {
	this := OSSSourceAccessKey{}
	return &this
}

// GetId returns the Id field value
func (o *OSSSourceAccessKey) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *OSSSourceAccessKey) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *OSSSourceAccessKey) SetId(v string) {
	o.Id = v
}

// GetSecret returns the Secret field value
func (o *OSSSourceAccessKey) GetSecret() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Secret
}

// GetSecretOk returns a tuple with the Secret field value
// and a boolean to check if the value has been set.
func (o *OSSSourceAccessKey) GetSecretOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Secret, true
}

// SetSecret sets field value
func (o *OSSSourceAccessKey) SetSecret(v string) {
	o.Secret = v
}

func (o OSSSourceAccessKey) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OSSSourceAccessKey) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["secret"] = o.Secret

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *OSSSourceAccessKey) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"secret",
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

	varOSSSourceAccessKey := _OSSSourceAccessKey{}

	err = json.Unmarshal(data, &varOSSSourceAccessKey)

	if err != nil {
		return err
	}

	*o = OSSSourceAccessKey(varOSSSourceAccessKey)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "secret")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableOSSSourceAccessKey struct {
	value *OSSSourceAccessKey
	isSet bool
}

func (v NullableOSSSourceAccessKey) Get() *OSSSourceAccessKey {
	return v.value
}

func (v *NullableOSSSourceAccessKey) Set(val *OSSSourceAccessKey) {
	v.value = val
	v.isSet = true
}

func (v NullableOSSSourceAccessKey) IsSet() bool {
	return v.isSet
}

func (v *NullableOSSSourceAccessKey) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOSSSourceAccessKey(val *OSSSourceAccessKey) *NullableOSSSourceAccessKey {
	return &NullableOSSSourceAccessKey{value: val, isSet: true}
}

func (v NullableOSSSourceAccessKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOSSSourceAccessKey) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
