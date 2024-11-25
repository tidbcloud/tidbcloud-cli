/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
	"time"
)

// checks if the Dedicatedv1beta1ClusterPausePlan type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1ClusterPausePlan{}

// Dedicatedv1beta1ClusterPausePlan struct for Dedicatedv1beta1ClusterPausePlan
type Dedicatedv1beta1ClusterPausePlan struct {
	PauseType            Dedicatedv1beta1ClusterPausePlanType `json:"pauseType"`
	ScheduledResumeTime  *time.Time                           `json:"scheduledResumeTime,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1ClusterPausePlan Dedicatedv1beta1ClusterPausePlan

// NewDedicatedv1beta1ClusterPausePlan instantiates a new Dedicatedv1beta1ClusterPausePlan object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1ClusterPausePlan(pauseType Dedicatedv1beta1ClusterPausePlanType) *Dedicatedv1beta1ClusterPausePlan {
	this := Dedicatedv1beta1ClusterPausePlan{}
	this.PauseType = pauseType
	return &this
}

// NewDedicatedv1beta1ClusterPausePlanWithDefaults instantiates a new Dedicatedv1beta1ClusterPausePlan object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1ClusterPausePlanWithDefaults() *Dedicatedv1beta1ClusterPausePlan {
	this := Dedicatedv1beta1ClusterPausePlan{}
	return &this
}

// GetPauseType returns the PauseType field value
func (o *Dedicatedv1beta1ClusterPausePlan) GetPauseType() Dedicatedv1beta1ClusterPausePlanType {
	if o == nil {
		var ret Dedicatedv1beta1ClusterPausePlanType
		return ret
	}

	return o.PauseType
}

// GetPauseTypeOk returns a tuple with the PauseType field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ClusterPausePlan) GetPauseTypeOk() (*Dedicatedv1beta1ClusterPausePlanType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PauseType, true
}

// SetPauseType sets field value
func (o *Dedicatedv1beta1ClusterPausePlan) SetPauseType(v Dedicatedv1beta1ClusterPausePlanType) {
	o.PauseType = v
}

// GetScheduledResumeTime returns the ScheduledResumeTime field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ClusterPausePlan) GetScheduledResumeTime() time.Time {
	if o == nil || IsNil(o.ScheduledResumeTime) {
		var ret time.Time
		return ret
	}
	return *o.ScheduledResumeTime
}

// GetScheduledResumeTimeOk returns a tuple with the ScheduledResumeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ClusterPausePlan) GetScheduledResumeTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ScheduledResumeTime) {
		return nil, false
	}
	return o.ScheduledResumeTime, true
}

// HasScheduledResumeTime returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ClusterPausePlan) HasScheduledResumeTime() bool {
	if o != nil && !IsNil(o.ScheduledResumeTime) {
		return true
	}

	return false
}

// SetScheduledResumeTime gets a reference to the given time.Time and assigns it to the ScheduledResumeTime field.
func (o *Dedicatedv1beta1ClusterPausePlan) SetScheduledResumeTime(v time.Time) {
	o.ScheduledResumeTime = &v
}

func (o Dedicatedv1beta1ClusterPausePlan) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1ClusterPausePlan) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["pauseType"] = o.PauseType
	if !IsNil(o.ScheduledResumeTime) {
		toSerialize["scheduledResumeTime"] = o.ScheduledResumeTime
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1ClusterPausePlan) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"pauseType",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varDedicatedv1beta1ClusterPausePlan := _Dedicatedv1beta1ClusterPausePlan{}

	err = json.Unmarshal(data, &varDedicatedv1beta1ClusterPausePlan)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1ClusterPausePlan(varDedicatedv1beta1ClusterPausePlan)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "pauseType")
		delete(additionalProperties, "scheduledResumeTime")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1ClusterPausePlan struct {
	value *Dedicatedv1beta1ClusterPausePlan
	isSet bool
}

func (v NullableDedicatedv1beta1ClusterPausePlan) Get() *Dedicatedv1beta1ClusterPausePlan {
	return v.value
}

func (v *NullableDedicatedv1beta1ClusterPausePlan) Set(val *Dedicatedv1beta1ClusterPausePlan) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1ClusterPausePlan) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1ClusterPausePlan) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1ClusterPausePlan(val *Dedicatedv1beta1ClusterPausePlan) *NullableDedicatedv1beta1ClusterPausePlan {
	return &NullableDedicatedv1beta1ClusterPausePlan{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1ClusterPausePlan) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1ClusterPausePlan) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
