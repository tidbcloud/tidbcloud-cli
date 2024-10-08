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

// checks if the ExportServiceCreateExportBody type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExportServiceCreateExportBody{}

// ExportServiceCreateExportBody struct for ExportServiceCreateExportBody
type ExportServiceCreateExportBody struct {
	// Optional. The options of the export.
	ExportOptions *ExportOptions `json:"exportOptions,omitempty"`
	// Optional. The target of the export.
	Target *ExportTarget `json:"target,omitempty"`
	// Optional. The display name of the export. Default: SNAPSHOT_{snapshot_time}.
	DisplayName          *string `json:"displayName,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ExportServiceCreateExportBody ExportServiceCreateExportBody

// NewExportServiceCreateExportBody instantiates a new ExportServiceCreateExportBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExportServiceCreateExportBody() *ExportServiceCreateExportBody {
	this := ExportServiceCreateExportBody{}
	return &this
}

// NewExportServiceCreateExportBodyWithDefaults instantiates a new ExportServiceCreateExportBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExportServiceCreateExportBodyWithDefaults() *ExportServiceCreateExportBody {
	this := ExportServiceCreateExportBody{}
	return &this
}

// GetExportOptions returns the ExportOptions field value if set, zero value otherwise.
func (o *ExportServiceCreateExportBody) GetExportOptions() ExportOptions {
	if o == nil || IsNil(o.ExportOptions) {
		var ret ExportOptions
		return ret
	}
	return *o.ExportOptions
}

// GetExportOptionsOk returns a tuple with the ExportOptions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportServiceCreateExportBody) GetExportOptionsOk() (*ExportOptions, bool) {
	if o == nil || IsNil(o.ExportOptions) {
		return nil, false
	}
	return o.ExportOptions, true
}

// HasExportOptions returns a boolean if a field has been set.
func (o *ExportServiceCreateExportBody) HasExportOptions() bool {
	if o != nil && !IsNil(o.ExportOptions) {
		return true
	}

	return false
}

// SetExportOptions gets a reference to the given ExportOptions and assigns it to the ExportOptions field.
func (o *ExportServiceCreateExportBody) SetExportOptions(v ExportOptions) {
	o.ExportOptions = &v
}

// GetTarget returns the Target field value if set, zero value otherwise.
func (o *ExportServiceCreateExportBody) GetTarget() ExportTarget {
	if o == nil || IsNil(o.Target) {
		var ret ExportTarget
		return ret
	}
	return *o.Target
}

// GetTargetOk returns a tuple with the Target field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportServiceCreateExportBody) GetTargetOk() (*ExportTarget, bool) {
	if o == nil || IsNil(o.Target) {
		return nil, false
	}
	return o.Target, true
}

// HasTarget returns a boolean if a field has been set.
func (o *ExportServiceCreateExportBody) HasTarget() bool {
	if o != nil && !IsNil(o.Target) {
		return true
	}

	return false
}

// SetTarget gets a reference to the given ExportTarget and assigns it to the Target field.
func (o *ExportServiceCreateExportBody) SetTarget(v ExportTarget) {
	o.Target = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *ExportServiceCreateExportBody) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportServiceCreateExportBody) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *ExportServiceCreateExportBody) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *ExportServiceCreateExportBody) SetDisplayName(v string) {
	o.DisplayName = &v
}

func (o ExportServiceCreateExportBody) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExportServiceCreateExportBody) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ExportOptions) {
		toSerialize["exportOptions"] = o.ExportOptions
	}
	if !IsNil(o.Target) {
		toSerialize["target"] = o.Target
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ExportServiceCreateExportBody) UnmarshalJSON(data []byte) (err error) {
	varExportServiceCreateExportBody := _ExportServiceCreateExportBody{}

	err = json.Unmarshal(data, &varExportServiceCreateExportBody)

	if err != nil {
		return err
	}

	*o = ExportServiceCreateExportBody(varExportServiceCreateExportBody)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "exportOptions")
		delete(additionalProperties, "target")
		delete(additionalProperties, "displayName")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableExportServiceCreateExportBody struct {
	value *ExportServiceCreateExportBody
	isSet bool
}

func (v NullableExportServiceCreateExportBody) Get() *ExportServiceCreateExportBody {
	return v.value
}

func (v *NullableExportServiceCreateExportBody) Set(val *ExportServiceCreateExportBody) {
	v.value = val
	v.isSet = true
}

func (v NullableExportServiceCreateExportBody) IsSet() bool {
	return v.isSet
}

func (v *NullableExportServiceCreateExportBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportServiceCreateExportBody(val *ExportServiceCreateExportBody) *NullableExportServiceCreateExportBody {
	return &NullableExportServiceCreateExportBody{value: val, isSet: true}
}

func (v NullableExportServiceCreateExportBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportServiceCreateExportBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
