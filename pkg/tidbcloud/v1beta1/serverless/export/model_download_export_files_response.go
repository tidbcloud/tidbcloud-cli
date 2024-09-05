/*
TiDB Cloud Serverless Export Open API

TiDB Cloud Serverless Export Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package export

import (
	"encoding/json"
)

// checks if the DownloadExportFilesResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DownloadExportFilesResponse{}

// DownloadExportFilesResponse struct for DownloadExportFilesResponse
type DownloadExportFilesResponse struct {
	// The files with download url of the export.
	Files []ExportFile `json:"files,omitempty"`
}

// NewDownloadExportFilesResponse instantiates a new DownloadExportFilesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDownloadExportFilesResponse() *DownloadExportFilesResponse {
	this := DownloadExportFilesResponse{}
	return &this
}

// NewDownloadExportFilesResponseWithDefaults instantiates a new DownloadExportFilesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDownloadExportFilesResponseWithDefaults() *DownloadExportFilesResponse {
	this := DownloadExportFilesResponse{}
	return &this
}

// GetFiles returns the Files field value if set, zero value otherwise.
func (o *DownloadExportFilesResponse) GetFiles() []ExportFile {
	if o == nil || IsNil(o.Files) {
		var ret []ExportFile
		return ret
	}
	return o.Files
}

// GetFilesOk returns a tuple with the Files field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DownloadExportFilesResponse) GetFilesOk() ([]ExportFile, bool) {
	if o == nil || IsNil(o.Files) {
		return nil, false
	}
	return o.Files, true
}

// HasFiles returns a boolean if a field has been set.
func (o *DownloadExportFilesResponse) HasFiles() bool {
	if o != nil && !IsNil(o.Files) {
		return true
	}

	return false
}

// SetFiles gets a reference to the given []ExportFile and assigns it to the Files field.
func (o *DownloadExportFilesResponse) SetFiles(v []ExportFile) {
	o.Files = v
}

func (o DownloadExportFilesResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DownloadExportFilesResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Files) {
		toSerialize["files"] = o.Files
	}
	return toSerialize, nil
}

type NullableDownloadExportFilesResponse struct {
	value *DownloadExportFilesResponse
	isSet bool
}

func (v NullableDownloadExportFilesResponse) Get() *DownloadExportFilesResponse {
	return v.value
}

func (v *NullableDownloadExportFilesResponse) Set(val *DownloadExportFilesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableDownloadExportFilesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableDownloadExportFilesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDownloadExportFilesResponse(val *DownloadExportFilesResponse) *NullableDownloadExportFilesResponse {
	return &NullableDownloadExportFilesResponse{value: val, isSet: true}
}

func (v NullableDownloadExportFilesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDownloadExportFilesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
