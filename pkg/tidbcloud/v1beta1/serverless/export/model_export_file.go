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

// checks if the ExportFile type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExportFile{}

// ExportFile struct for ExportFile
type ExportFile struct {
	// The name of the file.
	Name *string `json:"name,omitempty"`
	// The size in bytes of the file.
	Size *int64 `json:"size,omitempty"`
	// download url of the file.
	DownloadUrl *string `json:"downloadUrl,omitempty"`
}

// NewExportFile instantiates a new ExportFile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExportFile() *ExportFile {
	this := ExportFile{}
	return &this
}

// NewExportFileWithDefaults instantiates a new ExportFile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExportFileWithDefaults() *ExportFile {
	this := ExportFile{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ExportFile) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportFile) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ExportFile) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ExportFile) SetName(v string) {
	o.Name = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *ExportFile) GetSize() int64 {
	if o == nil || IsNil(o.Size) {
		var ret int64
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportFile) GetSizeOk() (*int64, bool) {
	if o == nil || IsNil(o.Size) {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *ExportFile) HasSize() bool {
	if o != nil && !IsNil(o.Size) {
		return true
	}

	return false
}

// SetSize gets a reference to the given int64 and assigns it to the Size field.
func (o *ExportFile) SetSize(v int64) {
	o.Size = &v
}

// GetDownloadUrl returns the DownloadUrl field value if set, zero value otherwise.
func (o *ExportFile) GetDownloadUrl() string {
	if o == nil || IsNil(o.DownloadUrl) {
		var ret string
		return ret
	}
	return *o.DownloadUrl
}

// GetDownloadUrlOk returns a tuple with the DownloadUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportFile) GetDownloadUrlOk() (*string, bool) {
	if o == nil || IsNil(o.DownloadUrl) {
		return nil, false
	}
	return o.DownloadUrl, true
}

// HasDownloadUrl returns a boolean if a field has been set.
func (o *ExportFile) HasDownloadUrl() bool {
	if o != nil && !IsNil(o.DownloadUrl) {
		return true
	}

	return false
}

// SetDownloadUrl gets a reference to the given string and assigns it to the DownloadUrl field.
func (o *ExportFile) SetDownloadUrl(v string) {
	o.DownloadUrl = &v
}

func (o ExportFile) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExportFile) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Size) {
		toSerialize["size"] = o.Size
	}
	if !IsNil(o.DownloadUrl) {
		toSerialize["downloadUrl"] = o.DownloadUrl
	}
	return toSerialize, nil
}

type NullableExportFile struct {
	value *ExportFile
	isSet bool
}

func (v NullableExportFile) Get() *ExportFile {
	return v.value
}

func (v *NullableExportFile) Set(val *ExportFile) {
	v.value = val
	v.isSet = true
}

func (v NullableExportFile) IsSet() bool {
	return v.isSet
}

func (v *NullableExportFile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportFile(val *ExportFile) *NullableExportFile {
	return &NullableExportFile{value: val, isSet: true}
}

func (v NullableExportFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportFile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
