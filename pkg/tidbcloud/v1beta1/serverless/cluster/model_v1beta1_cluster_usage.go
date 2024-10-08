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

// checks if the V1beta1ClusterUsage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ClusterUsage{}

// V1beta1ClusterUsage Message for usage statistics of a cluster.
type V1beta1ClusterUsage struct {
	// Output_only. The Request Units used in this month.
	RequestUnit *string `json:"requestUnit,omitempty"`
	// Output_only. The storage used on row-based storage in bytes.
	RowBasedStorage *float64 `json:"rowBasedStorage,omitempty"`
	// Output_only. The storage used on column-based storage in bytes.
	ColumnarStorage      *float64 `json:"columnarStorage,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ClusterUsage V1beta1ClusterUsage

// NewV1beta1ClusterUsage instantiates a new V1beta1ClusterUsage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ClusterUsage() *V1beta1ClusterUsage {
	this := V1beta1ClusterUsage{}
	return &this
}

// NewV1beta1ClusterUsageWithDefaults instantiates a new V1beta1ClusterUsage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ClusterUsageWithDefaults() *V1beta1ClusterUsage {
	this := V1beta1ClusterUsage{}
	return &this
}

// GetRequestUnit returns the RequestUnit field value if set, zero value otherwise.
func (o *V1beta1ClusterUsage) GetRequestUnit() string {
	if o == nil || IsNil(o.RequestUnit) {
		var ret string
		return ret
	}
	return *o.RequestUnit
}

// GetRequestUnitOk returns a tuple with the RequestUnit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterUsage) GetRequestUnitOk() (*string, bool) {
	if o == nil || IsNil(o.RequestUnit) {
		return nil, false
	}
	return o.RequestUnit, true
}

// HasRequestUnit returns a boolean if a field has been set.
func (o *V1beta1ClusterUsage) HasRequestUnit() bool {
	if o != nil && !IsNil(o.RequestUnit) {
		return true
	}

	return false
}

// SetRequestUnit gets a reference to the given string and assigns it to the RequestUnit field.
func (o *V1beta1ClusterUsage) SetRequestUnit(v string) {
	o.RequestUnit = &v
}

// GetRowBasedStorage returns the RowBasedStorage field value if set, zero value otherwise.
func (o *V1beta1ClusterUsage) GetRowBasedStorage() float64 {
	if o == nil || IsNil(o.RowBasedStorage) {
		var ret float64
		return ret
	}
	return *o.RowBasedStorage
}

// GetRowBasedStorageOk returns a tuple with the RowBasedStorage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterUsage) GetRowBasedStorageOk() (*float64, bool) {
	if o == nil || IsNil(o.RowBasedStorage) {
		return nil, false
	}
	return o.RowBasedStorage, true
}

// HasRowBasedStorage returns a boolean if a field has been set.
func (o *V1beta1ClusterUsage) HasRowBasedStorage() bool {
	if o != nil && !IsNil(o.RowBasedStorage) {
		return true
	}

	return false
}

// SetRowBasedStorage gets a reference to the given float64 and assigns it to the RowBasedStorage field.
func (o *V1beta1ClusterUsage) SetRowBasedStorage(v float64) {
	o.RowBasedStorage = &v
}

// GetColumnarStorage returns the ColumnarStorage field value if set, zero value otherwise.
func (o *V1beta1ClusterUsage) GetColumnarStorage() float64 {
	if o == nil || IsNil(o.ColumnarStorage) {
		var ret float64
		return ret
	}
	return *o.ColumnarStorage
}

// GetColumnarStorageOk returns a tuple with the ColumnarStorage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterUsage) GetColumnarStorageOk() (*float64, bool) {
	if o == nil || IsNil(o.ColumnarStorage) {
		return nil, false
	}
	return o.ColumnarStorage, true
}

// HasColumnarStorage returns a boolean if a field has been set.
func (o *V1beta1ClusterUsage) HasColumnarStorage() bool {
	if o != nil && !IsNil(o.ColumnarStorage) {
		return true
	}

	return false
}

// SetColumnarStorage gets a reference to the given float64 and assigns it to the ColumnarStorage field.
func (o *V1beta1ClusterUsage) SetColumnarStorage(v float64) {
	o.ColumnarStorage = &v
}

func (o V1beta1ClusterUsage) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ClusterUsage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.RequestUnit) {
		toSerialize["requestUnit"] = o.RequestUnit
	}
	if !IsNil(o.RowBasedStorage) {
		toSerialize["rowBasedStorage"] = o.RowBasedStorage
	}
	if !IsNil(o.ColumnarStorage) {
		toSerialize["columnarStorage"] = o.ColumnarStorage
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1ClusterUsage) UnmarshalJSON(data []byte) (err error) {
	varV1beta1ClusterUsage := _V1beta1ClusterUsage{}

	err = json.Unmarshal(data, &varV1beta1ClusterUsage)

	if err != nil {
		return err
	}

	*o = V1beta1ClusterUsage(varV1beta1ClusterUsage)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "requestUnit")
		delete(additionalProperties, "rowBasedStorage")
		delete(additionalProperties, "columnarStorage")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ClusterUsage struct {
	value *V1beta1ClusterUsage
	isSet bool
}

func (v NullableV1beta1ClusterUsage) Get() *V1beta1ClusterUsage {
	return v.value
}

func (v *NullableV1beta1ClusterUsage) Set(val *V1beta1ClusterUsage) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ClusterUsage) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ClusterUsage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ClusterUsage(val *V1beta1ClusterUsage) *NullableV1beta1ClusterUsage {
	return &NullableV1beta1ClusterUsage{value: val, isSet: true}
}

func (v NullableV1beta1ClusterUsage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ClusterUsage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
