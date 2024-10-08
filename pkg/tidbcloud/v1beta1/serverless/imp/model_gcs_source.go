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

// checks if the GCSSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GCSSource{}

// GCSSource struct for GCSSource
type GCSSource struct {
	// The GCS URI of the import source.
	Uri string `json:"uri"`
	// The auth method of the import source.
	AuthType             ImportGcsAuthTypeEnum `json:"authType"`
	ServiceAccountKey    *string               `json:"serviceAccountKey,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _GCSSource GCSSource

// NewGCSSource instantiates a new GCSSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGCSSource(uri string, authType ImportGcsAuthTypeEnum) *GCSSource {
	this := GCSSource{}
	this.Uri = uri
	this.AuthType = authType
	return &this
}

// NewGCSSourceWithDefaults instantiates a new GCSSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGCSSourceWithDefaults() *GCSSource {
	this := GCSSource{}
	return &this
}

// GetUri returns the Uri field value
func (o *GCSSource) GetUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uri
}

// GetUriOk returns a tuple with the Uri field value
// and a boolean to check if the value has been set.
func (o *GCSSource) GetUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uri, true
}

// SetUri sets field value
func (o *GCSSource) SetUri(v string) {
	o.Uri = v
}

// GetAuthType returns the AuthType field value
func (o *GCSSource) GetAuthType() ImportGcsAuthTypeEnum {
	if o == nil {
		var ret ImportGcsAuthTypeEnum
		return ret
	}

	return o.AuthType
}

// GetAuthTypeOk returns a tuple with the AuthType field value
// and a boolean to check if the value has been set.
func (o *GCSSource) GetAuthTypeOk() (*ImportGcsAuthTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthType, true
}

// SetAuthType sets field value
func (o *GCSSource) SetAuthType(v ImportGcsAuthTypeEnum) {
	o.AuthType = v
}

// GetServiceAccountKey returns the ServiceAccountKey field value if set, zero value otherwise.
func (o *GCSSource) GetServiceAccountKey() string {
	if o == nil || IsNil(o.ServiceAccountKey) {
		var ret string
		return ret
	}
	return *o.ServiceAccountKey
}

// GetServiceAccountKeyOk returns a tuple with the ServiceAccountKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GCSSource) GetServiceAccountKeyOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceAccountKey) {
		return nil, false
	}
	return o.ServiceAccountKey, true
}

// HasServiceAccountKey returns a boolean if a field has been set.
func (o *GCSSource) HasServiceAccountKey() bool {
	if o != nil && !IsNil(o.ServiceAccountKey) {
		return true
	}

	return false
}

// SetServiceAccountKey gets a reference to the given string and assigns it to the ServiceAccountKey field.
func (o *GCSSource) SetServiceAccountKey(v string) {
	o.ServiceAccountKey = &v
}

func (o GCSSource) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GCSSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uri"] = o.Uri
	toSerialize["authType"] = o.AuthType
	if !IsNil(o.ServiceAccountKey) {
		toSerialize["serviceAccountKey"] = o.ServiceAccountKey
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *GCSSource) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"uri",
		"authType",
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

	varGCSSource := _GCSSource{}

	err = json.Unmarshal(data, &varGCSSource)

	if err != nil {
		return err
	}

	*o = GCSSource(varGCSSource)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "uri")
		delete(additionalProperties, "authType")
		delete(additionalProperties, "serviceAccountKey")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableGCSSource struct {
	value *GCSSource
	isSet bool
}

func (v NullableGCSSource) Get() *GCSSource {
	return v.value
}

func (v *NullableGCSSource) Set(val *GCSSource) {
	v.value = val
	v.isSet = true
}

func (v NullableGCSSource) IsSet() bool {
	return v.isSet
}

func (v *NullableGCSSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGCSSource(val *GCSSource) *NullableGCSSource {
	return &NullableGCSSource{value: val, isSet: true}
}

func (v NullableGCSSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGCSSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
