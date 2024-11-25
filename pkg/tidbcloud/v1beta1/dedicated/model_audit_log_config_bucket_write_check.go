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

// checks if the AuditLogConfigBucketWriteCheck type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuditLogConfigBucketWriteCheck{}

// AuditLogConfigBucketWriteCheck struct for AuditLogConfigBucketWriteCheck
type AuditLogConfigBucketWriteCheck struct {
	Writable *bool `json:"writable,omitempty"`
	// The reason why the bucket is not writable. Output only when `writable` is false.
	ErrorReason *string `json:"errorReason,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AuditLogConfigBucketWriteCheck AuditLogConfigBucketWriteCheck

// NewAuditLogConfigBucketWriteCheck instantiates a new AuditLogConfigBucketWriteCheck object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuditLogConfigBucketWriteCheck() *AuditLogConfigBucketWriteCheck {
	this := AuditLogConfigBucketWriteCheck{}
	return &this
}

// NewAuditLogConfigBucketWriteCheckWithDefaults instantiates a new AuditLogConfigBucketWriteCheck object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuditLogConfigBucketWriteCheckWithDefaults() *AuditLogConfigBucketWriteCheck {
	this := AuditLogConfigBucketWriteCheck{}
	return &this
}

// GetWritable returns the Writable field value if set, zero value otherwise.
func (o *AuditLogConfigBucketWriteCheck) GetWritable() bool {
	if o == nil || IsNil(o.Writable) {
		var ret bool
		return ret
	}
	return *o.Writable
}

// GetWritableOk returns a tuple with the Writable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLogConfigBucketWriteCheck) GetWritableOk() (*bool, bool) {
	if o == nil || IsNil(o.Writable) {
		return nil, false
	}
	return o.Writable, true
}

// HasWritable returns a boolean if a field has been set.
func (o *AuditLogConfigBucketWriteCheck) HasWritable() bool {
	if o != nil && !IsNil(o.Writable) {
		return true
	}

	return false
}

// SetWritable gets a reference to the given bool and assigns it to the Writable field.
func (o *AuditLogConfigBucketWriteCheck) SetWritable(v bool) {
	o.Writable = &v
}

// GetErrorReason returns the ErrorReason field value if set, zero value otherwise.
func (o *AuditLogConfigBucketWriteCheck) GetErrorReason() string {
	if o == nil || IsNil(o.ErrorReason) {
		var ret string
		return ret
	}
	return *o.ErrorReason
}

// GetErrorReasonOk returns a tuple with the ErrorReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLogConfigBucketWriteCheck) GetErrorReasonOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorReason) {
		return nil, false
	}
	return o.ErrorReason, true
}

// HasErrorReason returns a boolean if a field has been set.
func (o *AuditLogConfigBucketWriteCheck) HasErrorReason() bool {
	if o != nil && !IsNil(o.ErrorReason) {
		return true
	}

	return false
}

// SetErrorReason gets a reference to the given string and assigns it to the ErrorReason field.
func (o *AuditLogConfigBucketWriteCheck) SetErrorReason(v string) {
	o.ErrorReason = &v
}

func (o AuditLogConfigBucketWriteCheck) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuditLogConfigBucketWriteCheck) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Writable) {
		toSerialize["writable"] = o.Writable
	}
	if !IsNil(o.ErrorReason) {
		toSerialize["errorReason"] = o.ErrorReason
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *AuditLogConfigBucketWriteCheck) UnmarshalJSON(data []byte) (err error) {
	varAuditLogConfigBucketWriteCheck := _AuditLogConfigBucketWriteCheck{}

	err = json.Unmarshal(data, &varAuditLogConfigBucketWriteCheck)

	if err != nil {
		return err
	}

	*o = AuditLogConfigBucketWriteCheck(varAuditLogConfigBucketWriteCheck)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "writable")
		delete(additionalProperties, "errorReason")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableAuditLogConfigBucketWriteCheck struct {
	value *AuditLogConfigBucketWriteCheck
	isSet bool
}

func (v NullableAuditLogConfigBucketWriteCheck) Get() *AuditLogConfigBucketWriteCheck {
	return v.value
}

func (v *NullableAuditLogConfigBucketWriteCheck) Set(val *AuditLogConfigBucketWriteCheck) {
	v.value = val
	v.isSet = true
}

func (v NullableAuditLogConfigBucketWriteCheck) IsSet() bool {
	return v.isSet
}

func (v *NullableAuditLogConfigBucketWriteCheck) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuditLogConfigBucketWriteCheck(val *AuditLogConfigBucketWriteCheck) *NullableAuditLogConfigBucketWriteCheck {
	return &NullableAuditLogConfigBucketWriteCheck{value: val, isSet: true}
}

func (v NullableAuditLogConfigBucketWriteCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuditLogConfigBucketWriteCheck) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

