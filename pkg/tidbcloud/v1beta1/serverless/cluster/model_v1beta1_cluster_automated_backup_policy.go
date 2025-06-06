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

// checks if the V1beta1ClusterAutomatedBackupPolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ClusterAutomatedBackupPolicy{}

// V1beta1ClusterAutomatedBackupPolicy Message for automated backup configuration for a cluster.
type V1beta1ClusterAutomatedBackupPolicy struct {
	// Optional. When automated backups should start, in HH:mm format, UTC.
	StartTime *string `json:"startTime,omitempty"`
	// OUTPUT_ONLY. Number of days to retain automated backups.
	RetentionDays        *int32 `json:"retentionDays,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ClusterAutomatedBackupPolicy V1beta1ClusterAutomatedBackupPolicy

// NewV1beta1ClusterAutomatedBackupPolicy instantiates a new V1beta1ClusterAutomatedBackupPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ClusterAutomatedBackupPolicy() *V1beta1ClusterAutomatedBackupPolicy {
	this := V1beta1ClusterAutomatedBackupPolicy{}
	return &this
}

// NewV1beta1ClusterAutomatedBackupPolicyWithDefaults instantiates a new V1beta1ClusterAutomatedBackupPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ClusterAutomatedBackupPolicyWithDefaults() *V1beta1ClusterAutomatedBackupPolicy {
	this := V1beta1ClusterAutomatedBackupPolicy{}
	return &this
}

// GetStartTime returns the StartTime field value if set, zero value otherwise.
func (o *V1beta1ClusterAutomatedBackupPolicy) GetStartTime() string {
	if o == nil || IsNil(o.StartTime) {
		var ret string
		return ret
	}
	return *o.StartTime
}

// GetStartTimeOk returns a tuple with the StartTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterAutomatedBackupPolicy) GetStartTimeOk() (*string, bool) {
	if o == nil || IsNil(o.StartTime) {
		return nil, false
	}
	return o.StartTime, true
}

// HasStartTime returns a boolean if a field has been set.
func (o *V1beta1ClusterAutomatedBackupPolicy) HasStartTime() bool {
	if o != nil && !IsNil(o.StartTime) {
		return true
	}

	return false
}

// SetStartTime gets a reference to the given string and assigns it to the StartTime field.
func (o *V1beta1ClusterAutomatedBackupPolicy) SetStartTime(v string) {
	o.StartTime = &v
}

// GetRetentionDays returns the RetentionDays field value if set, zero value otherwise.
func (o *V1beta1ClusterAutomatedBackupPolicy) GetRetentionDays() int32 {
	if o == nil || IsNil(o.RetentionDays) {
		var ret int32
		return ret
	}
	return *o.RetentionDays
}

// GetRetentionDaysOk returns a tuple with the RetentionDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterAutomatedBackupPolicy) GetRetentionDaysOk() (*int32, bool) {
	if o == nil || IsNil(o.RetentionDays) {
		return nil, false
	}
	return o.RetentionDays, true
}

// HasRetentionDays returns a boolean if a field has been set.
func (o *V1beta1ClusterAutomatedBackupPolicy) HasRetentionDays() bool {
	if o != nil && !IsNil(o.RetentionDays) {
		return true
	}

	return false
}

// SetRetentionDays gets a reference to the given int32 and assigns it to the RetentionDays field.
func (o *V1beta1ClusterAutomatedBackupPolicy) SetRetentionDays(v int32) {
	o.RetentionDays = &v
}

func (o V1beta1ClusterAutomatedBackupPolicy) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ClusterAutomatedBackupPolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.StartTime) {
		toSerialize["startTime"] = o.StartTime
	}
	if !IsNil(o.RetentionDays) {
		toSerialize["retentionDays"] = o.RetentionDays
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1ClusterAutomatedBackupPolicy) UnmarshalJSON(data []byte) (err error) {
	varV1beta1ClusterAutomatedBackupPolicy := _V1beta1ClusterAutomatedBackupPolicy{}

	err = json.Unmarshal(data, &varV1beta1ClusterAutomatedBackupPolicy)

	if err != nil {
		return err
	}

	*o = V1beta1ClusterAutomatedBackupPolicy(varV1beta1ClusterAutomatedBackupPolicy)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "startTime")
		delete(additionalProperties, "retentionDays")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ClusterAutomatedBackupPolicy struct {
	value *V1beta1ClusterAutomatedBackupPolicy
	isSet bool
}

func (v NullableV1beta1ClusterAutomatedBackupPolicy) Get() *V1beta1ClusterAutomatedBackupPolicy {
	return v.value
}

func (v *NullableV1beta1ClusterAutomatedBackupPolicy) Set(val *V1beta1ClusterAutomatedBackupPolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ClusterAutomatedBackupPolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ClusterAutomatedBackupPolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ClusterAutomatedBackupPolicy(val *V1beta1ClusterAutomatedBackupPolicy) *NullableV1beta1ClusterAutomatedBackupPolicy {
	return &NullableV1beta1ClusterAutomatedBackupPolicy{value: val, isSet: true}
}

func (v NullableV1beta1ClusterAutomatedBackupPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ClusterAutomatedBackupPolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
