/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package br

import (
	"encoding/json"
)

// checks if the RestoreRequestSnapshot type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RestoreRequestSnapshot{}

// RestoreRequestSnapshot struct for RestoreRequestSnapshot
type RestoreRequestSnapshot struct {
	BackupId *string `json:"backupId,omitempty"`
}

// NewRestoreRequestSnapshot instantiates a new RestoreRequestSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRestoreRequestSnapshot() *RestoreRequestSnapshot {
	this := RestoreRequestSnapshot{}
	return &this
}

// NewRestoreRequestSnapshotWithDefaults instantiates a new RestoreRequestSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRestoreRequestSnapshotWithDefaults() *RestoreRequestSnapshot {
	this := RestoreRequestSnapshot{}
	return &this
}

// GetBackupId returns the BackupId field value if set, zero value otherwise.
func (o *RestoreRequestSnapshot) GetBackupId() string {
	if o == nil || IsNil(o.BackupId) {
		var ret string
		return ret
	}
	return *o.BackupId
}

// GetBackupIdOk returns a tuple with the BackupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RestoreRequestSnapshot) GetBackupIdOk() (*string, bool) {
	if o == nil || IsNil(o.BackupId) {
		return nil, false
	}
	return o.BackupId, true
}

// HasBackupId returns a boolean if a field has been set.
func (o *RestoreRequestSnapshot) HasBackupId() bool {
	if o != nil && !IsNil(o.BackupId) {
		return true
	}

	return false
}

// SetBackupId gets a reference to the given string and assigns it to the BackupId field.
func (o *RestoreRequestSnapshot) SetBackupId(v string) {
	o.BackupId = &v
}

func (o RestoreRequestSnapshot) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RestoreRequestSnapshot) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.BackupId) {
		toSerialize["backupId"] = o.BackupId
	}
	return toSerialize, nil
}

type NullableRestoreRequestSnapshot struct {
	value *RestoreRequestSnapshot
	isSet bool
}

func (v NullableRestoreRequestSnapshot) Get() *RestoreRequestSnapshot {
	return v.value
}

func (v *NullableRestoreRequestSnapshot) Set(val *RestoreRequestSnapshot) {
	v.value = val
	v.isSet = true
}

func (v NullableRestoreRequestSnapshot) IsSet() bool {
	return v.isSet
}

func (v *NullableRestoreRequestSnapshot) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRestoreRequestSnapshot(val *RestoreRequestSnapshot) *NullableRestoreRequestSnapshot {
	return &NullableRestoreRequestSnapshot{value: val, isSet: true}
}

func (v NullableRestoreRequestSnapshot) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRestoreRequestSnapshot) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
