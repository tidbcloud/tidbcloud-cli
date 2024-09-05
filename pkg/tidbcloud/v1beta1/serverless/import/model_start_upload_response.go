/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package imp

import (
	"encoding/json"
)

// checks if the StartUploadResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StartUploadResponse{}

// StartUploadResponse struct for StartUploadResponse
type StartUploadResponse struct {
	UploadUrl            []string `json:"uploadUrl,omitempty"`
	UploadId             *string  `json:"uploadId,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _StartUploadResponse StartUploadResponse

// NewStartUploadResponse instantiates a new StartUploadResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStartUploadResponse() *StartUploadResponse {
	this := StartUploadResponse{}
	return &this
}

// NewStartUploadResponseWithDefaults instantiates a new StartUploadResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStartUploadResponseWithDefaults() *StartUploadResponse {
	this := StartUploadResponse{}
	return &this
}

// GetUploadUrl returns the UploadUrl field value if set, zero value otherwise.
func (o *StartUploadResponse) GetUploadUrl() []string {
	if o == nil || IsNil(o.UploadUrl) {
		var ret []string
		return ret
	}
	return o.UploadUrl
}

// GetUploadUrlOk returns a tuple with the UploadUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StartUploadResponse) GetUploadUrlOk() ([]string, bool) {
	if o == nil || IsNil(o.UploadUrl) {
		return nil, false
	}
	return o.UploadUrl, true
}

// HasUploadUrl returns a boolean if a field has been set.
func (o *StartUploadResponse) HasUploadUrl() bool {
	if o != nil && !IsNil(o.UploadUrl) {
		return true
	}

	return false
}

// SetUploadUrl gets a reference to the given []string and assigns it to the UploadUrl field.
func (o *StartUploadResponse) SetUploadUrl(v []string) {
	o.UploadUrl = v
}

// GetUploadId returns the UploadId field value if set, zero value otherwise.
func (o *StartUploadResponse) GetUploadId() string {
	if o == nil || IsNil(o.UploadId) {
		var ret string
		return ret
	}
	return *o.UploadId
}

// GetUploadIdOk returns a tuple with the UploadId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StartUploadResponse) GetUploadIdOk() (*string, bool) {
	if o == nil || IsNil(o.UploadId) {
		return nil, false
	}
	return o.UploadId, true
}

// HasUploadId returns a boolean if a field has been set.
func (o *StartUploadResponse) HasUploadId() bool {
	if o != nil && !IsNil(o.UploadId) {
		return true
	}

	return false
}

// SetUploadId gets a reference to the given string and assigns it to the UploadId field.
func (o *StartUploadResponse) SetUploadId(v string) {
	o.UploadId = &v
}

func (o StartUploadResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StartUploadResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.UploadUrl) {
		toSerialize["uploadUrl"] = o.UploadUrl
	}
	if !IsNil(o.UploadId) {
		toSerialize["uploadId"] = o.UploadId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *StartUploadResponse) UnmarshalJSON(data []byte) (err error) {
	varStartUploadResponse := _StartUploadResponse{}

	err = json.Unmarshal(data, &varStartUploadResponse)

	if err != nil {
		return err
	}

	*o = StartUploadResponse(varStartUploadResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "uploadUrl")
		delete(additionalProperties, "uploadId")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableStartUploadResponse struct {
	value *StartUploadResponse
	isSet bool
}

func (v NullableStartUploadResponse) Get() *StartUploadResponse {
	return v.value
}

func (v *NullableStartUploadResponse) Set(val *StartUploadResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableStartUploadResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableStartUploadResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStartUploadResponse(val *StartUploadResponse) *NullableStartUploadResponse {
	return &NullableStartUploadResponse{value: val, isSet: true}
}

func (v NullableStartUploadResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStartUploadResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
