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

// checks if the ExportOptionsFilter type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExportOptionsFilter{}

// ExportOptionsFilter struct for ExportOptionsFilter
type ExportOptionsFilter struct {
	// Optional. Use SQL to filter the export.
	Sql *string `json:"sql,omitempty"`
	// Optional. Use table-filter to filter the export.
	Table                *ExportOptionsFilterTable `json:"table,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ExportOptionsFilter ExportOptionsFilter

// NewExportOptionsFilter instantiates a new ExportOptionsFilter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExportOptionsFilter() *ExportOptionsFilter {
	this := ExportOptionsFilter{}
	return &this
}

// NewExportOptionsFilterWithDefaults instantiates a new ExportOptionsFilter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExportOptionsFilterWithDefaults() *ExportOptionsFilter {
	this := ExportOptionsFilter{}
	return &this
}

// GetSql returns the Sql field value if set, zero value otherwise.
func (o *ExportOptionsFilter) GetSql() string {
	if o == nil || IsNil(o.Sql) {
		var ret string
		return ret
	}
	return *o.Sql
}

// GetSqlOk returns a tuple with the Sql field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportOptionsFilter) GetSqlOk() (*string, bool) {
	if o == nil || IsNil(o.Sql) {
		return nil, false
	}
	return o.Sql, true
}

// HasSql returns a boolean if a field has been set.
func (o *ExportOptionsFilter) HasSql() bool {
	if o != nil && !IsNil(o.Sql) {
		return true
	}

	return false
}

// SetSql gets a reference to the given string and assigns it to the Sql field.
func (o *ExportOptionsFilter) SetSql(v string) {
	o.Sql = &v
}

// GetTable returns the Table field value if set, zero value otherwise.
func (o *ExportOptionsFilter) GetTable() ExportOptionsFilterTable {
	if o == nil || IsNil(o.Table) {
		var ret ExportOptionsFilterTable
		return ret
	}
	return *o.Table
}

// GetTableOk returns a tuple with the Table field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportOptionsFilter) GetTableOk() (*ExportOptionsFilterTable, bool) {
	if o == nil || IsNil(o.Table) {
		return nil, false
	}
	return o.Table, true
}

// HasTable returns a boolean if a field has been set.
func (o *ExportOptionsFilter) HasTable() bool {
	if o != nil && !IsNil(o.Table) {
		return true
	}

	return false
}

// SetTable gets a reference to the given ExportOptionsFilterTable and assigns it to the Table field.
func (o *ExportOptionsFilter) SetTable(v ExportOptionsFilterTable) {
	o.Table = &v
}

func (o ExportOptionsFilter) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExportOptionsFilter) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Sql) {
		toSerialize["sql"] = o.Sql
	}
	if !IsNil(o.Table) {
		toSerialize["table"] = o.Table
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ExportOptionsFilter) UnmarshalJSON(data []byte) (err error) {
	varExportOptionsFilter := _ExportOptionsFilter{}

	err = json.Unmarshal(data, &varExportOptionsFilter)

	if err != nil {
		return err
	}

	*o = ExportOptionsFilter(varExportOptionsFilter)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "sql")
		delete(additionalProperties, "table")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableExportOptionsFilter struct {
	value *ExportOptionsFilter
	isSet bool
}

func (v NullableExportOptionsFilter) Get() *ExportOptionsFilter {
	return v.value
}

func (v *NullableExportOptionsFilter) Set(val *ExportOptionsFilter) {
	v.value = val
	v.isSet = true
}

func (v NullableExportOptionsFilter) IsSet() bool {
	return v.isSet
}

func (v *NullableExportOptionsFilter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportOptionsFilter(val *ExportOptionsFilter) *NullableExportOptionsFilter {
	return &NullableExportOptionsFilter{value: val, isSet: true}
}

func (v NullableExportOptionsFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportOptionsFilter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
