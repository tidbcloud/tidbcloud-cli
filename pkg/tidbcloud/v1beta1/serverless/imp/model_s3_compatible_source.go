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

// checks if the S3CompatibleSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &S3CompatibleSource{}

// S3CompatibleSource struct for S3CompatibleSource
type S3CompatibleSource struct {
	// The S3 compatible URI of the import source.
	Uri string `json:"uri"`
	// The auth method of the import source.
	AuthType ImportS3CompatibleAuthTypeEnum `json:"authType"`
	// The access key.
	AccessKey *S3CompatibleSourceAccessKey `json:"accessKey,omitempty"`
	// The custom S3 endpoint (HTTPS only). Used for connecting to non-AWS S3-compatible storage, such as Cloudflare or other cloud providers. Ensure the endpoint is a valid HTTPS URL (e.g., \"https://custom-s3.example.com\").
	Endpoint             NullableString `json:"endpoint,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _S3CompatibleSource S3CompatibleSource

// NewS3CompatibleSource instantiates a new S3CompatibleSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewS3CompatibleSource(uri string, authType ImportS3CompatibleAuthTypeEnum) *S3CompatibleSource {
	this := S3CompatibleSource{}
	this.Uri = uri
	this.AuthType = authType
	return &this
}

// NewS3CompatibleSourceWithDefaults instantiates a new S3CompatibleSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewS3CompatibleSourceWithDefaults() *S3CompatibleSource {
	this := S3CompatibleSource{}
	return &this
}

// GetUri returns the Uri field value
func (o *S3CompatibleSource) GetUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uri
}

// GetUriOk returns a tuple with the Uri field value
// and a boolean to check if the value has been set.
func (o *S3CompatibleSource) GetUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uri, true
}

// SetUri sets field value
func (o *S3CompatibleSource) SetUri(v string) {
	o.Uri = v
}

// GetAuthType returns the AuthType field value
func (o *S3CompatibleSource) GetAuthType() ImportS3CompatibleAuthTypeEnum {
	if o == nil {
		var ret ImportS3CompatibleAuthTypeEnum
		return ret
	}

	return o.AuthType
}

// GetAuthTypeOk returns a tuple with the AuthType field value
// and a boolean to check if the value has been set.
func (o *S3CompatibleSource) GetAuthTypeOk() (*ImportS3CompatibleAuthTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthType, true
}

// SetAuthType sets field value
func (o *S3CompatibleSource) SetAuthType(v ImportS3CompatibleAuthTypeEnum) {
	o.AuthType = v
}

// GetAccessKey returns the AccessKey field value if set, zero value otherwise.
func (o *S3CompatibleSource) GetAccessKey() S3CompatibleSourceAccessKey {
	if o == nil || IsNil(o.AccessKey) {
		var ret S3CompatibleSourceAccessKey
		return ret
	}
	return *o.AccessKey
}

// GetAccessKeyOk returns a tuple with the AccessKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *S3CompatibleSource) GetAccessKeyOk() (*S3CompatibleSourceAccessKey, bool) {
	if o == nil || IsNil(o.AccessKey) {
		return nil, false
	}
	return o.AccessKey, true
}

// HasAccessKey returns a boolean if a field has been set.
func (o *S3CompatibleSource) HasAccessKey() bool {
	if o != nil && !IsNil(o.AccessKey) {
		return true
	}

	return false
}

// SetAccessKey gets a reference to the given S3CompatibleSourceAccessKey and assigns it to the AccessKey field.
func (o *S3CompatibleSource) SetAccessKey(v S3CompatibleSourceAccessKey) {
	o.AccessKey = &v
}

// GetEndpoint returns the Endpoint field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *S3CompatibleSource) GetEndpoint() string {
	if o == nil || IsNil(o.Endpoint.Get()) {
		var ret string
		return ret
	}
	return *o.Endpoint.Get()
}

// GetEndpointOk returns a tuple with the Endpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *S3CompatibleSource) GetEndpointOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Endpoint.Get(), o.Endpoint.IsSet()
}

// HasEndpoint returns a boolean if a field has been set.
func (o *S3CompatibleSource) HasEndpoint() bool {
	if o != nil && o.Endpoint.IsSet() {
		return true
	}

	return false
}

// SetEndpoint gets a reference to the given NullableString and assigns it to the Endpoint field.
func (o *S3CompatibleSource) SetEndpoint(v string) {
	o.Endpoint.Set(&v)
}

// SetEndpointNil sets the value for Endpoint to be an explicit nil
func (o *S3CompatibleSource) SetEndpointNil() {
	o.Endpoint.Set(nil)
}

// UnsetEndpoint ensures that no value is present for Endpoint, not even an explicit nil
func (o *S3CompatibleSource) UnsetEndpoint() {
	o.Endpoint.Unset()
}

func (o S3CompatibleSource) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o S3CompatibleSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["uri"] = o.Uri
	toSerialize["authType"] = o.AuthType
	if !IsNil(o.AccessKey) {
		toSerialize["accessKey"] = o.AccessKey
	}
	if o.Endpoint.IsSet() {
		toSerialize["endpoint"] = o.Endpoint.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *S3CompatibleSource) UnmarshalJSON(data []byte) (err error) {
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

	varS3CompatibleSource := _S3CompatibleSource{}

	err = json.Unmarshal(data, &varS3CompatibleSource)

	if err != nil {
		return err
	}

	*o = S3CompatibleSource(varS3CompatibleSource)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "uri")
		delete(additionalProperties, "authType")
		delete(additionalProperties, "accessKey")
		delete(additionalProperties, "endpoint")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableS3CompatibleSource struct {
	value *S3CompatibleSource
	isSet bool
}

func (v NullableS3CompatibleSource) Get() *S3CompatibleSource {
	return v.value
}

func (v *NullableS3CompatibleSource) Set(val *S3CompatibleSource) {
	v.value = val
	v.isSet = true
}

func (v NullableS3CompatibleSource) IsSet() bool {
	return v.isSet
}

func (v *NullableS3CompatibleSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableS3CompatibleSource(val *S3CompatibleSource) *NullableS3CompatibleSource {
	return &NullableS3CompatibleSource{value: val, isSet: true}
}

func (v NullableS3CompatibleSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableS3CompatibleSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
