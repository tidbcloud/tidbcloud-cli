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

// checks if the AzureBlobSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AzureBlobSource{}

// AzureBlobSource struct for AzureBlobSource
type AzureBlobSource struct {
	// The Azure Blob URI of the import source.
	AuthType ImportAzureBlobAuthTypeEnum `json:"authType"`
	// The sas token. This field is input-only.
	SasToken *string `json:"sasToken,omitempty"`
	// The Azure Blob URI of the import source. For example: azure://<account>.blob.core.windows.net/<container>/<path> or https://<account>.blob.core.windows.net/<container>/<path>.
	Uri                  string `json:"uri"`
	AdditionalProperties map[string]interface{}
}

type _AzureBlobSource AzureBlobSource

// NewAzureBlobSource instantiates a new AzureBlobSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAzureBlobSource(authType ImportAzureBlobAuthTypeEnum, uri string) *AzureBlobSource {
	this := AzureBlobSource{}
	this.AuthType = authType
	this.Uri = uri
	return &this
}

// NewAzureBlobSourceWithDefaults instantiates a new AzureBlobSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAzureBlobSourceWithDefaults() *AzureBlobSource {
	this := AzureBlobSource{}
	return &this
}

// GetAuthType returns the AuthType field value
func (o *AzureBlobSource) GetAuthType() ImportAzureBlobAuthTypeEnum {
	if o == nil {
		var ret ImportAzureBlobAuthTypeEnum
		return ret
	}

	return o.AuthType
}

// GetAuthTypeOk returns a tuple with the AuthType field value
// and a boolean to check if the value has been set.
func (o *AzureBlobSource) GetAuthTypeOk() (*ImportAzureBlobAuthTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthType, true
}

// SetAuthType sets field value
func (o *AzureBlobSource) SetAuthType(v ImportAzureBlobAuthTypeEnum) {
	o.AuthType = v
}

// GetSasToken returns the SasToken field value if set, zero value otherwise.
func (o *AzureBlobSource) GetSasToken() string {
	if o == nil || IsNil(o.SasToken) {
		var ret string
		return ret
	}
	return *o.SasToken
}

// GetSasTokenOk returns a tuple with the SasToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureBlobSource) GetSasTokenOk() (*string, bool) {
	if o == nil || IsNil(o.SasToken) {
		return nil, false
	}
	return o.SasToken, true
}

// HasSasToken returns a boolean if a field has been set.
func (o *AzureBlobSource) HasSasToken() bool {
	if o != nil && !IsNil(o.SasToken) {
		return true
	}

	return false
}

// SetSasToken gets a reference to the given string and assigns it to the SasToken field.
func (o *AzureBlobSource) SetSasToken(v string) {
	o.SasToken = &v
}

// GetUri returns the Uri field value
func (o *AzureBlobSource) GetUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uri
}

// GetUriOk returns a tuple with the Uri field value
// and a boolean to check if the value has been set.
func (o *AzureBlobSource) GetUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uri, true
}

// SetUri sets field value
func (o *AzureBlobSource) SetUri(v string) {
	o.Uri = v
}

func (o AzureBlobSource) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AzureBlobSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["authType"] = o.AuthType
	if !IsNil(o.SasToken) {
		toSerialize["sasToken"] = o.SasToken
	}
	toSerialize["uri"] = o.Uri

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *AzureBlobSource) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"authType",
		"uri",
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

	varAzureBlobSource := _AzureBlobSource{}

	err = json.Unmarshal(data, &varAzureBlobSource)

	if err != nil {
		return err
	}

	*o = AzureBlobSource(varAzureBlobSource)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "authType")
		delete(additionalProperties, "sasToken")
		delete(additionalProperties, "uri")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableAzureBlobSource struct {
	value *AzureBlobSource
	isSet bool
}

func (v NullableAzureBlobSource) Get() *AzureBlobSource {
	return v.value
}

func (v *NullableAzureBlobSource) Set(val *AzureBlobSource) {
	v.value = val
	v.isSet = true
}

func (v NullableAzureBlobSource) IsSet() bool {
	return v.isSet
}

func (v *NullableAzureBlobSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAzureBlobSource(val *AzureBlobSource) *NullableAzureBlobSource {
	return &NullableAzureBlobSource{value: val, isSet: true}
}

func (v NullableAzureBlobSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAzureBlobSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
