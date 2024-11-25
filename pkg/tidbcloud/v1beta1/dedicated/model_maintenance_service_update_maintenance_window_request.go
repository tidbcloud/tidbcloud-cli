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

// checks if the MaintenanceServiceUpdateMaintenanceWindowRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MaintenanceServiceUpdateMaintenanceWindowRequest{}

// MaintenanceServiceUpdateMaintenanceWindowRequest struct for MaintenanceServiceUpdateMaintenanceWindowRequest
type MaintenanceServiceUpdateMaintenanceWindowRequest struct {
	WeekDay *int32 `json:"weekDay,omitempty"`
	DayHour *int32 `json:"dayHour,omitempty"`
	HourMinute *int32 `json:"hourMinute,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _MaintenanceServiceUpdateMaintenanceWindowRequest MaintenanceServiceUpdateMaintenanceWindowRequest

// NewMaintenanceServiceUpdateMaintenanceWindowRequest instantiates a new MaintenanceServiceUpdateMaintenanceWindowRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMaintenanceServiceUpdateMaintenanceWindowRequest() *MaintenanceServiceUpdateMaintenanceWindowRequest {
	this := MaintenanceServiceUpdateMaintenanceWindowRequest{}
	return &this
}

// NewMaintenanceServiceUpdateMaintenanceWindowRequestWithDefaults instantiates a new MaintenanceServiceUpdateMaintenanceWindowRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMaintenanceServiceUpdateMaintenanceWindowRequestWithDefaults() *MaintenanceServiceUpdateMaintenanceWindowRequest {
	this := MaintenanceServiceUpdateMaintenanceWindowRequest{}
	return &this
}

// GetWeekDay returns the WeekDay field value if set, zero value otherwise.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetWeekDay() int32 {
	if o == nil || IsNil(o.WeekDay) {
		var ret int32
		return ret
	}
	return *o.WeekDay
}

// GetWeekDayOk returns a tuple with the WeekDay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetWeekDayOk() (*int32, bool) {
	if o == nil || IsNil(o.WeekDay) {
		return nil, false
	}
	return o.WeekDay, true
}

// HasWeekDay returns a boolean if a field has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) HasWeekDay() bool {
	if o != nil && !IsNil(o.WeekDay) {
		return true
	}

	return false
}

// SetWeekDay gets a reference to the given int32 and assigns it to the WeekDay field.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) SetWeekDay(v int32) {
	o.WeekDay = &v
}

// GetDayHour returns the DayHour field value if set, zero value otherwise.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetDayHour() int32 {
	if o == nil || IsNil(o.DayHour) {
		var ret int32
		return ret
	}
	return *o.DayHour
}

// GetDayHourOk returns a tuple with the DayHour field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetDayHourOk() (*int32, bool) {
	if o == nil || IsNil(o.DayHour) {
		return nil, false
	}
	return o.DayHour, true
}

// HasDayHour returns a boolean if a field has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) HasDayHour() bool {
	if o != nil && !IsNil(o.DayHour) {
		return true
	}

	return false
}

// SetDayHour gets a reference to the given int32 and assigns it to the DayHour field.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) SetDayHour(v int32) {
	o.DayHour = &v
}

// GetHourMinute returns the HourMinute field value if set, zero value otherwise.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetHourMinute() int32 {
	if o == nil || IsNil(o.HourMinute) {
		var ret int32
		return ret
	}
	return *o.HourMinute
}

// GetHourMinuteOk returns a tuple with the HourMinute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) GetHourMinuteOk() (*int32, bool) {
	if o == nil || IsNil(o.HourMinute) {
		return nil, false
	}
	return o.HourMinute, true
}

// HasHourMinute returns a boolean if a field has been set.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) HasHourMinute() bool {
	if o != nil && !IsNil(o.HourMinute) {
		return true
	}

	return false
}

// SetHourMinute gets a reference to the given int32 and assigns it to the HourMinute field.
func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) SetHourMinute(v int32) {
	o.HourMinute = &v
}

func (o MaintenanceServiceUpdateMaintenanceWindowRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MaintenanceServiceUpdateMaintenanceWindowRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.WeekDay) {
		toSerialize["weekDay"] = o.WeekDay
	}
	if !IsNil(o.DayHour) {
		toSerialize["dayHour"] = o.DayHour
	}
	if !IsNil(o.HourMinute) {
		toSerialize["hourMinute"] = o.HourMinute
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *MaintenanceServiceUpdateMaintenanceWindowRequest) UnmarshalJSON(data []byte) (err error) {
	varMaintenanceServiceUpdateMaintenanceWindowRequest := _MaintenanceServiceUpdateMaintenanceWindowRequest{}

	err = json.Unmarshal(data, &varMaintenanceServiceUpdateMaintenanceWindowRequest)

	if err != nil {
		return err
	}

	*o = MaintenanceServiceUpdateMaintenanceWindowRequest(varMaintenanceServiceUpdateMaintenanceWindowRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "weekDay")
		delete(additionalProperties, "dayHour")
		delete(additionalProperties, "hourMinute")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableMaintenanceServiceUpdateMaintenanceWindowRequest struct {
	value *MaintenanceServiceUpdateMaintenanceWindowRequest
	isSet bool
}

func (v NullableMaintenanceServiceUpdateMaintenanceWindowRequest) Get() *MaintenanceServiceUpdateMaintenanceWindowRequest {
	return v.value
}

func (v *NullableMaintenanceServiceUpdateMaintenanceWindowRequest) Set(val *MaintenanceServiceUpdateMaintenanceWindowRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableMaintenanceServiceUpdateMaintenanceWindowRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableMaintenanceServiceUpdateMaintenanceWindowRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMaintenanceServiceUpdateMaintenanceWindowRequest(val *MaintenanceServiceUpdateMaintenanceWindowRequest) *NullableMaintenanceServiceUpdateMaintenanceWindowRequest {
	return &NullableMaintenanceServiceUpdateMaintenanceWindowRequest{value: val, isSet: true}
}

func (v NullableMaintenanceServiceUpdateMaintenanceWindowRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMaintenanceServiceUpdateMaintenanceWindowRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

