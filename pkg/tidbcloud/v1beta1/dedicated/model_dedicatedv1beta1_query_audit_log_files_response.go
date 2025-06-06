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

// checks if the Dedicatedv1beta1QueryAuditLogFilesResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1QueryAuditLogFilesResponse{}

// Dedicatedv1beta1QueryAuditLogFilesResponse struct for Dedicatedv1beta1QueryAuditLogFilesResponse
type Dedicatedv1beta1QueryAuditLogFilesResponse struct {
	AuditLogFiles        []Dedicatedv1beta1AuditLogFile `json:"auditLogFiles,omitempty"`
	TotalSize            *int32                         `json:"totalSize,omitempty"`
	NextPageToken        *string                        `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1QueryAuditLogFilesResponse Dedicatedv1beta1QueryAuditLogFilesResponse

// NewDedicatedv1beta1QueryAuditLogFilesResponse instantiates a new Dedicatedv1beta1QueryAuditLogFilesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1QueryAuditLogFilesResponse() *Dedicatedv1beta1QueryAuditLogFilesResponse {
	this := Dedicatedv1beta1QueryAuditLogFilesResponse{}
	return &this
}

// NewDedicatedv1beta1QueryAuditLogFilesResponseWithDefaults instantiates a new Dedicatedv1beta1QueryAuditLogFilesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1QueryAuditLogFilesResponseWithDefaults() *Dedicatedv1beta1QueryAuditLogFilesResponse {
	this := Dedicatedv1beta1QueryAuditLogFilesResponse{}
	return &this
}

// GetAuditLogFiles returns the AuditLogFiles field value if set, zero value otherwise.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetAuditLogFiles() []Dedicatedv1beta1AuditLogFile {
	if o == nil || IsNil(o.AuditLogFiles) {
		var ret []Dedicatedv1beta1AuditLogFile
		return ret
	}
	return o.AuditLogFiles
}

// GetAuditLogFilesOk returns a tuple with the AuditLogFiles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetAuditLogFilesOk() ([]Dedicatedv1beta1AuditLogFile, bool) {
	if o == nil || IsNil(o.AuditLogFiles) {
		return nil, false
	}
	return o.AuditLogFiles, true
}

// HasAuditLogFiles returns a boolean if a field has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) HasAuditLogFiles() bool {
	if o != nil && !IsNil(o.AuditLogFiles) {
		return true
	}

	return false
}

// SetAuditLogFiles gets a reference to the given []Dedicatedv1beta1AuditLogFile and assigns it to the AuditLogFiles field.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) SetAuditLogFiles(v []Dedicatedv1beta1AuditLogFile) {
	o.AuditLogFiles = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o Dedicatedv1beta1QueryAuditLogFilesResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1QueryAuditLogFilesResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AuditLogFiles) {
		toSerialize["auditLogFiles"] = o.AuditLogFiles
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

func (o *Dedicatedv1beta1QueryAuditLogFilesResponse) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1QueryAuditLogFilesResponse := _Dedicatedv1beta1QueryAuditLogFilesResponse{}

	err = json.Unmarshal(data, &varDedicatedv1beta1QueryAuditLogFilesResponse)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1QueryAuditLogFilesResponse(varDedicatedv1beta1QueryAuditLogFilesResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "auditLogFiles")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1QueryAuditLogFilesResponse struct {
	value *Dedicatedv1beta1QueryAuditLogFilesResponse
	isSet bool
}

func (v NullableDedicatedv1beta1QueryAuditLogFilesResponse) Get() *Dedicatedv1beta1QueryAuditLogFilesResponse {
	return v.value
}

func (v *NullableDedicatedv1beta1QueryAuditLogFilesResponse) Set(val *Dedicatedv1beta1QueryAuditLogFilesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1QueryAuditLogFilesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1QueryAuditLogFilesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1QueryAuditLogFilesResponse(val *Dedicatedv1beta1QueryAuditLogFilesResponse) *NullableDedicatedv1beta1QueryAuditLogFilesResponse {
	return &NullableDedicatedv1beta1QueryAuditLogFilesResponse{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1QueryAuditLogFilesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1QueryAuditLogFilesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
