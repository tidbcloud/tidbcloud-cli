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

// checks if the Dedicatedv1beta1MaintenanceWindow type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1MaintenanceWindow{}

// Dedicatedv1beta1MaintenanceWindow MaintenanceWindow is a singleton resource that represents the maintenance window under a project.
type Dedicatedv1beta1MaintenanceWindow struct {
	Name                *string `json:"name,omitempty"`
	MaintenanceWindowId *string `json:"maintenanceWindowId,omitempty"`
	ProjectId           string  `json:"projectId"`
	// 0-6, 0 is Sunday.
	WeekDay int32 `json:"weekDay"`
	// 0-23 in UTC.
	DayHour int32 `json:"dayHour"`
	// 0-59 in UTC.
	HourMinute                int32                             `json:"hourMinute"`
	NextMaintenanceDate       *time.Time                        `json:"nextMaintenanceDate,omitempty"`
	UnchangedMaintenanceTasks []Dedicatedv1beta1MaintenanceTask `json:"unchangedMaintenanceTasks,omitempty"`
	AdditionalProperties      map[string]interface{}
}

type _Dedicatedv1beta1MaintenanceWindow Dedicatedv1beta1MaintenanceWindow

// NewDedicatedv1beta1MaintenanceWindow instantiates a new Dedicatedv1beta1MaintenanceWindow object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1MaintenanceWindow(projectId string, weekDay int32, dayHour int32, hourMinute int32) *Dedicatedv1beta1MaintenanceWindow {
	this := Dedicatedv1beta1MaintenanceWindow{}
	this.ProjectId = projectId
	this.WeekDay = weekDay
	this.DayHour = dayHour
	this.HourMinute = hourMinute
	return &this
}

// NewDedicatedv1beta1MaintenanceWindowWithDefaults instantiates a new Dedicatedv1beta1MaintenanceWindow object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1MaintenanceWindowWithDefaults() *Dedicatedv1beta1MaintenanceWindow {
	this := Dedicatedv1beta1MaintenanceWindow{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceWindow) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1MaintenanceWindow) SetName(v string) {
	o.Name = &v
}

// GetMaintenanceWindowId returns the MaintenanceWindowId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceWindow) GetMaintenanceWindowId() string {
	if o == nil || IsNil(o.MaintenanceWindowId) {
		var ret string
		return ret
	}
	return *o.MaintenanceWindowId
}

// GetMaintenanceWindowIdOk returns a tuple with the MaintenanceWindowId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetMaintenanceWindowIdOk() (*string, bool) {
	if o == nil || IsNil(o.MaintenanceWindowId) {
		return nil, false
	}
	return o.MaintenanceWindowId, true
}

// HasMaintenanceWindowId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) HasMaintenanceWindowId() bool {
	if o != nil && !IsNil(o.MaintenanceWindowId) {
		return true
	}

	return false
}

// SetMaintenanceWindowId gets a reference to the given string and assigns it to the MaintenanceWindowId field.
func (o *Dedicatedv1beta1MaintenanceWindow) SetMaintenanceWindowId(v string) {
	o.MaintenanceWindowId = &v
}

// GetProjectId returns the ProjectId field value
func (o *Dedicatedv1beta1MaintenanceWindow) GetProjectId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetProjectIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProjectId, true
}

// SetProjectId sets field value
func (o *Dedicatedv1beta1MaintenanceWindow) SetProjectId(v string) {
	o.ProjectId = v
}

// GetWeekDay returns the WeekDay field value
func (o *Dedicatedv1beta1MaintenanceWindow) GetWeekDay() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.WeekDay
}

// GetWeekDayOk returns a tuple with the WeekDay field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetWeekDayOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WeekDay, true
}

// SetWeekDay sets field value
func (o *Dedicatedv1beta1MaintenanceWindow) SetWeekDay(v int32) {
	o.WeekDay = v
}

// GetDayHour returns the DayHour field value
func (o *Dedicatedv1beta1MaintenanceWindow) GetDayHour() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.DayHour
}

// GetDayHourOk returns a tuple with the DayHour field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetDayHourOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DayHour, true
}

// SetDayHour sets field value
func (o *Dedicatedv1beta1MaintenanceWindow) SetDayHour(v int32) {
	o.DayHour = v
}

// GetHourMinute returns the HourMinute field value
func (o *Dedicatedv1beta1MaintenanceWindow) GetHourMinute() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.HourMinute
}

// GetHourMinuteOk returns a tuple with the HourMinute field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetHourMinuteOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.HourMinute, true
}

// SetHourMinute sets field value
func (o *Dedicatedv1beta1MaintenanceWindow) SetHourMinute(v int32) {
	o.HourMinute = v
}

// GetNextMaintenanceDate returns the NextMaintenanceDate field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceWindow) GetNextMaintenanceDate() time.Time {
	if o == nil || IsNil(o.NextMaintenanceDate) {
		var ret time.Time
		return ret
	}
	return *o.NextMaintenanceDate
}

// GetNextMaintenanceDateOk returns a tuple with the NextMaintenanceDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetNextMaintenanceDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.NextMaintenanceDate) {
		return nil, false
	}
	return o.NextMaintenanceDate, true
}

// HasNextMaintenanceDate returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) HasNextMaintenanceDate() bool {
	if o != nil && !IsNil(o.NextMaintenanceDate) {
		return true
	}

	return false
}

// SetNextMaintenanceDate gets a reference to the given time.Time and assigns it to the NextMaintenanceDate field.
func (o *Dedicatedv1beta1MaintenanceWindow) SetNextMaintenanceDate(v time.Time) {
	o.NextMaintenanceDate = &v
}

// GetUnchangedMaintenanceTasks returns the UnchangedMaintenanceTasks field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceWindow) GetUnchangedMaintenanceTasks() []Dedicatedv1beta1MaintenanceTask {
	if o == nil || IsNil(o.UnchangedMaintenanceTasks) {
		var ret []Dedicatedv1beta1MaintenanceTask
		return ret
	}
	return o.UnchangedMaintenanceTasks
}

// GetUnchangedMaintenanceTasksOk returns a tuple with the UnchangedMaintenanceTasks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) GetUnchangedMaintenanceTasksOk() ([]Dedicatedv1beta1MaintenanceTask, bool) {
	if o == nil || IsNil(o.UnchangedMaintenanceTasks) {
		return nil, false
	}
	return o.UnchangedMaintenanceTasks, true
}

// HasUnchangedMaintenanceTasks returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceWindow) HasUnchangedMaintenanceTasks() bool {
	if o != nil && !IsNil(o.UnchangedMaintenanceTasks) {
		return true
	}

	return false
}

// SetUnchangedMaintenanceTasks gets a reference to the given []Dedicatedv1beta1MaintenanceTask and assigns it to the UnchangedMaintenanceTasks field.
func (o *Dedicatedv1beta1MaintenanceWindow) SetUnchangedMaintenanceTasks(v []Dedicatedv1beta1MaintenanceTask) {
	o.UnchangedMaintenanceTasks = v
}

func (o Dedicatedv1beta1MaintenanceWindow) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1MaintenanceWindow) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.MaintenanceWindowId) {
		toSerialize["maintenanceWindowId"] = o.MaintenanceWindowId
	}
	toSerialize["projectId"] = o.ProjectId
	toSerialize["weekDay"] = o.WeekDay
	toSerialize["dayHour"] = o.DayHour
	toSerialize["hourMinute"] = o.HourMinute
	if !IsNil(o.NextMaintenanceDate) {
		toSerialize["nextMaintenanceDate"] = o.NextMaintenanceDate
	}
	if !IsNil(o.UnchangedMaintenanceTasks) {
		toSerialize["unchangedMaintenanceTasks"] = o.UnchangedMaintenanceTasks
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1MaintenanceWindow) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"projectId",
		"weekDay",
		"dayHour",
		"hourMinute",
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

	varDedicatedv1beta1MaintenanceWindow := _Dedicatedv1beta1MaintenanceWindow{}

	err = json.Unmarshal(data, &varDedicatedv1beta1MaintenanceWindow)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1MaintenanceWindow(varDedicatedv1beta1MaintenanceWindow)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "maintenanceWindowId")
		delete(additionalProperties, "projectId")
		delete(additionalProperties, "weekDay")
		delete(additionalProperties, "dayHour")
		delete(additionalProperties, "hourMinute")
		delete(additionalProperties, "nextMaintenanceDate")
		delete(additionalProperties, "unchangedMaintenanceTasks")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1MaintenanceWindow struct {
	value *Dedicatedv1beta1MaintenanceWindow
	isSet bool
}

func (v NullableDedicatedv1beta1MaintenanceWindow) Get() *Dedicatedv1beta1MaintenanceWindow {
	return v.value
}

func (v *NullableDedicatedv1beta1MaintenanceWindow) Set(val *Dedicatedv1beta1MaintenanceWindow) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1MaintenanceWindow) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1MaintenanceWindow) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1MaintenanceWindow(val *Dedicatedv1beta1MaintenanceWindow) *NullableDedicatedv1beta1MaintenanceWindow {
	return &NullableDedicatedv1beta1MaintenanceWindow{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1MaintenanceWindow) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1MaintenanceWindow) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
