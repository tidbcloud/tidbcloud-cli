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

// checks if the CSVFormat type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CSVFormat{}

// CSVFormat struct for CSVFormat
type CSVFormat struct {
	// Separator of each value in CSV files. Default is ','.
	Separator *string `json:"separator,omitempty"`
	// Delimiter of string type variables in CSV files. Default is '\"'.
	Delimiter NullableString `json:"delimiter,omitempty"`
	// Import CSV files of the tables with header. Default is true.
	Header NullableBool `json:"header,omitempty"`
	// Whether the columns in CSV files can be null. Default is false.
	NotNull NullableBool `json:"notNull,omitempty"`
	// Representation of null values in CSV files. Default is \"\\N\".
	Null NullableString `json:"null,omitempty"`
	// Whether to escape backslashes in CSV files. Default is true.
	BackslashEscape NullableBool `json:"backslashEscape,omitempty"`
	// Whether to trim the last separator in CSV files. Default is false.
	TrimLastSeparator NullableBool `json:"trimLastSeparator,omitempty"`
}

// NewCSVFormat instantiates a new CSVFormat object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCSVFormat() *CSVFormat {
	this := CSVFormat{}
	return &this
}

// NewCSVFormatWithDefaults instantiates a new CSVFormat object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCSVFormatWithDefaults() *CSVFormat {
	this := CSVFormat{}
	return &this
}

// GetSeparator returns the Separator field value if set, zero value otherwise.
func (o *CSVFormat) GetSeparator() string {
	if o == nil || IsNil(o.Separator) {
		var ret string
		return ret
	}
	return *o.Separator
}

// GetSeparatorOk returns a tuple with the Separator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CSVFormat) GetSeparatorOk() (*string, bool) {
	if o == nil || IsNil(o.Separator) {
		return nil, false
	}
	return o.Separator, true
}

// HasSeparator returns a boolean if a field has been set.
func (o *CSVFormat) HasSeparator() bool {
	if o != nil && !IsNil(o.Separator) {
		return true
	}

	return false
}

// SetSeparator gets a reference to the given string and assigns it to the Separator field.
func (o *CSVFormat) SetSeparator(v string) {
	o.Separator = &v
}

// GetDelimiter returns the Delimiter field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetDelimiter() string {
	if o == nil || IsNil(o.Delimiter.Get()) {
		var ret string
		return ret
	}
	return *o.Delimiter.Get()
}

// GetDelimiterOk returns a tuple with the Delimiter field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetDelimiterOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Delimiter.Get(), o.Delimiter.IsSet()
}

// HasDelimiter returns a boolean if a field has been set.
func (o *CSVFormat) HasDelimiter() bool {
	if o != nil && o.Delimiter.IsSet() {
		return true
	}

	return false
}

// SetDelimiter gets a reference to the given NullableString and assigns it to the Delimiter field.
func (o *CSVFormat) SetDelimiter(v string) {
	o.Delimiter.Set(&v)
}

// SetDelimiterNil sets the value for Delimiter to be an explicit nil
func (o *CSVFormat) SetDelimiterNil() {
	o.Delimiter.Set(nil)
}

// UnsetDelimiter ensures that no value is present for Delimiter, not even an explicit nil
func (o *CSVFormat) UnsetDelimiter() {
	o.Delimiter.Unset()
}

// GetHeader returns the Header field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetHeader() bool {
	if o == nil || IsNil(o.Header.Get()) {
		var ret bool
		return ret
	}
	return *o.Header.Get()
}

// GetHeaderOk returns a tuple with the Header field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetHeaderOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.Header.Get(), o.Header.IsSet()
}

// HasHeader returns a boolean if a field has been set.
func (o *CSVFormat) HasHeader() bool {
	if o != nil && o.Header.IsSet() {
		return true
	}

	return false
}

// SetHeader gets a reference to the given NullableBool and assigns it to the Header field.
func (o *CSVFormat) SetHeader(v bool) {
	o.Header.Set(&v)
}

// SetHeaderNil sets the value for Header to be an explicit nil
func (o *CSVFormat) SetHeaderNil() {
	o.Header.Set(nil)
}

// UnsetHeader ensures that no value is present for Header, not even an explicit nil
func (o *CSVFormat) UnsetHeader() {
	o.Header.Unset()
}

// GetNotNull returns the NotNull field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetNotNull() bool {
	if o == nil || IsNil(o.NotNull.Get()) {
		var ret bool
		return ret
	}
	return *o.NotNull.Get()
}

// GetNotNullOk returns a tuple with the NotNull field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetNotNullOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.NotNull.Get(), o.NotNull.IsSet()
}

// HasNotNull returns a boolean if a field has been set.
func (o *CSVFormat) HasNotNull() bool {
	if o != nil && o.NotNull.IsSet() {
		return true
	}

	return false
}

// SetNotNull gets a reference to the given NullableBool and assigns it to the NotNull field.
func (o *CSVFormat) SetNotNull(v bool) {
	o.NotNull.Set(&v)
}

// SetNotNullNil sets the value for NotNull to be an explicit nil
func (o *CSVFormat) SetNotNullNil() {
	o.NotNull.Set(nil)
}

// UnsetNotNull ensures that no value is present for NotNull, not even an explicit nil
func (o *CSVFormat) UnsetNotNull() {
	o.NotNull.Unset()
}

// GetNull returns the Null field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetNull() string {
	if o == nil || IsNil(o.Null.Get()) {
		var ret string
		return ret
	}
	return *o.Null.Get()
}

// GetNullOk returns a tuple with the Null field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetNullOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Null.Get(), o.Null.IsSet()
}

// HasNull returns a boolean if a field has been set.
func (o *CSVFormat) HasNull() bool {
	if o != nil && o.Null.IsSet() {
		return true
	}

	return false
}

// SetNull gets a reference to the given NullableString and assigns it to the Null field.
func (o *CSVFormat) SetNull(v string) {
	o.Null.Set(&v)
}

// SetNullNil sets the value for Null to be an explicit nil
func (o *CSVFormat) SetNullNil() {
	o.Null.Set(nil)
}

// UnsetNull ensures that no value is present for Null, not even an explicit nil
func (o *CSVFormat) UnsetNull() {
	o.Null.Unset()
}

// GetBackslashEscape returns the BackslashEscape field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetBackslashEscape() bool {
	if o == nil || IsNil(o.BackslashEscape.Get()) {
		var ret bool
		return ret
	}
	return *o.BackslashEscape.Get()
}

// GetBackslashEscapeOk returns a tuple with the BackslashEscape field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetBackslashEscapeOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.BackslashEscape.Get(), o.BackslashEscape.IsSet()
}

// HasBackslashEscape returns a boolean if a field has been set.
func (o *CSVFormat) HasBackslashEscape() bool {
	if o != nil && o.BackslashEscape.IsSet() {
		return true
	}

	return false
}

// SetBackslashEscape gets a reference to the given NullableBool and assigns it to the BackslashEscape field.
func (o *CSVFormat) SetBackslashEscape(v bool) {
	o.BackslashEscape.Set(&v)
}

// SetBackslashEscapeNil sets the value for BackslashEscape to be an explicit nil
func (o *CSVFormat) SetBackslashEscapeNil() {
	o.BackslashEscape.Set(nil)
}

// UnsetBackslashEscape ensures that no value is present for BackslashEscape, not even an explicit nil
func (o *CSVFormat) UnsetBackslashEscape() {
	o.BackslashEscape.Unset()
}

// GetTrimLastSeparator returns the TrimLastSeparator field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CSVFormat) GetTrimLastSeparator() bool {
	if o == nil || IsNil(o.TrimLastSeparator.Get()) {
		var ret bool
		return ret
	}
	return *o.TrimLastSeparator.Get()
}

// GetTrimLastSeparatorOk returns a tuple with the TrimLastSeparator field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CSVFormat) GetTrimLastSeparatorOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.TrimLastSeparator.Get(), o.TrimLastSeparator.IsSet()
}

// HasTrimLastSeparator returns a boolean if a field has been set.
func (o *CSVFormat) HasTrimLastSeparator() bool {
	if o != nil && o.TrimLastSeparator.IsSet() {
		return true
	}

	return false
}

// SetTrimLastSeparator gets a reference to the given NullableBool and assigns it to the TrimLastSeparator field.
func (o *CSVFormat) SetTrimLastSeparator(v bool) {
	o.TrimLastSeparator.Set(&v)
}

// SetTrimLastSeparatorNil sets the value for TrimLastSeparator to be an explicit nil
func (o *CSVFormat) SetTrimLastSeparatorNil() {
	o.TrimLastSeparator.Set(nil)
}

// UnsetTrimLastSeparator ensures that no value is present for TrimLastSeparator, not even an explicit nil
func (o *CSVFormat) UnsetTrimLastSeparator() {
	o.TrimLastSeparator.Unset()
}

func (o CSVFormat) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CSVFormat) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Separator) {
		toSerialize["separator"] = o.Separator
	}
	if o.Delimiter.IsSet() {
		toSerialize["delimiter"] = o.Delimiter.Get()
	}
	if o.Header.IsSet() {
		toSerialize["header"] = o.Header.Get()
	}
	if o.NotNull.IsSet() {
		toSerialize["notNull"] = o.NotNull.Get()
	}
	if o.Null.IsSet() {
		toSerialize["null"] = o.Null.Get()
	}
	if o.BackslashEscape.IsSet() {
		toSerialize["backslashEscape"] = o.BackslashEscape.Get()
	}
	if o.TrimLastSeparator.IsSet() {
		toSerialize["trimLastSeparator"] = o.TrimLastSeparator.Get()
	}
	return toSerialize, nil
}

type NullableCSVFormat struct {
	value *CSVFormat
	isSet bool
}

func (v NullableCSVFormat) Get() *CSVFormat {
	return v.value
}

func (v *NullableCSVFormat) Set(val *CSVFormat) {
	v.value = val
	v.isSet = true
}

func (v NullableCSVFormat) IsSet() bool {
	return v.isSet
}

func (v *NullableCSVFormat) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCSVFormat(val *CSVFormat) *NullableCSVFormat {
	return &NullableCSVFormat{value: val, isSet: true}
}

func (v NullableCSVFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCSVFormat) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
