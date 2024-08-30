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

// checks if the ListImportsResp type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListImportsResp{}

// ListImportsResp struct for ListImportsResp
type ListImportsResp struct {
	// The imports.
	Imports []Import `json:"imports,omitempty"`
	// The total size of the imports.
	TotalSize *int64 `json:"totalSize,omitempty"`
	// The next page token.
	NextPageToken *string `json:"nextPageToken,omitempty"`
}

// NewListImportsResp instantiates a new ListImportsResp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListImportsResp() *ListImportsResp {
	this := ListImportsResp{}
	return &this
}

// NewListImportsRespWithDefaults instantiates a new ListImportsResp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListImportsRespWithDefaults() *ListImportsResp {
	this := ListImportsResp{}
	return &this
}

// GetImports returns the Imports field value if set, zero value otherwise.
func (o *ListImportsResp) GetImports() []Import {
	if o == nil || IsNil(o.Imports) {
		var ret []Import
		return ret
	}
	return o.Imports
}

// GetImportsOk returns a tuple with the Imports field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListImportsResp) GetImportsOk() ([]Import, bool) {
	if o == nil || IsNil(o.Imports) {
		return nil, false
	}
	return o.Imports, true
}

// HasImports returns a boolean if a field has been set.
func (o *ListImportsResp) HasImports() bool {
	if o != nil && !IsNil(o.Imports) {
		return true
	}

	return false
}

// SetImports gets a reference to the given []Import and assigns it to the Imports field.
func (o *ListImportsResp) SetImports(v []Import) {
	o.Imports = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *ListImportsResp) GetTotalSize() int64 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int64
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListImportsResp) GetTotalSizeOk() (*int64, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *ListImportsResp) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int64 and assigns it to the TotalSize field.
func (o *ListImportsResp) SetTotalSize(v int64) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *ListImportsResp) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListImportsResp) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *ListImportsResp) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *ListImportsResp) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o ListImportsResp) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListImportsResp) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Imports) {
		toSerialize["imports"] = o.Imports
	}
	if !IsNil(o.TotalSize) {
		toSerialize["totalSize"] = o.TotalSize
	}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}
	return toSerialize, nil
}

type NullableListImportsResp struct {
	value *ListImportsResp
	isSet bool
}

func (v NullableListImportsResp) Get() *ListImportsResp {
	return v.value
}

func (v *NullableListImportsResp) Set(val *ListImportsResp) {
	v.value = val
	v.isSet = true
}

func (v NullableListImportsResp) IsSet() bool {
	return v.isSet
}

func (v *NullableListImportsResp) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListImportsResp(val *ListImportsResp) *NullableListImportsResp {
	return &NullableListImportsResp{value: val, isSet: true}
}

func (v NullableListImportsResp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListImportsResp) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
