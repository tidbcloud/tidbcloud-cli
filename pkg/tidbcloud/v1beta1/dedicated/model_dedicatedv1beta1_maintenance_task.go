/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"time"
)

// checks if the Dedicatedv1beta1MaintenanceTask type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1MaintenanceTask{}

// Dedicatedv1beta1MaintenanceTask All fields are output only.
type Dedicatedv1beta1MaintenanceTask struct {
	Name              *string                               `json:"name,omitempty"`
	MaintenanceTaskId *string                               `json:"maintenanceTaskId,omitempty"`
	ProjectId         *string                               `json:"projectId,omitempty"`
	Description       *string                               `json:"description,omitempty"`
	State             *Dedicatedv1beta1MaintenanceTaskState `json:"state,omitempty"`
	// Timestamp when the task was created.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Timestamp when the task run.
	ScheduledApplyTime *time.Time `json:"scheduledApplyTime,omitempty"`
	// Timestamp when the task will be expired.
	Deadline             *time.Time `json:"deadline,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1MaintenanceTask Dedicatedv1beta1MaintenanceTask

// NewDedicatedv1beta1MaintenanceTask instantiates a new Dedicatedv1beta1MaintenanceTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1MaintenanceTask() *Dedicatedv1beta1MaintenanceTask {
	this := Dedicatedv1beta1MaintenanceTask{}
	return &this
}

// NewDedicatedv1beta1MaintenanceTaskWithDefaults instantiates a new Dedicatedv1beta1MaintenanceTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1MaintenanceTaskWithDefaults() *Dedicatedv1beta1MaintenanceTask {
	this := Dedicatedv1beta1MaintenanceTask{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1MaintenanceTask) SetName(v string) {
	o.Name = &v
}

// GetMaintenanceTaskId returns the MaintenanceTaskId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetMaintenanceTaskId() string {
	if o == nil || IsNil(o.MaintenanceTaskId) {
		var ret string
		return ret
	}
	return *o.MaintenanceTaskId
}

// GetMaintenanceTaskIdOk returns a tuple with the MaintenanceTaskId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetMaintenanceTaskIdOk() (*string, bool) {
	if o == nil || IsNil(o.MaintenanceTaskId) {
		return nil, false
	}
	return o.MaintenanceTaskId, true
}

// HasMaintenanceTaskId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasMaintenanceTaskId() bool {
	if o != nil && !IsNil(o.MaintenanceTaskId) {
		return true
	}

	return false
}

// SetMaintenanceTaskId gets a reference to the given string and assigns it to the MaintenanceTaskId field.
func (o *Dedicatedv1beta1MaintenanceTask) SetMaintenanceTaskId(v string) {
	o.MaintenanceTaskId = &v
}

// GetProjectId returns the ProjectId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetProjectId() string {
	if o == nil || IsNil(o.ProjectId) {
		var ret string
		return ret
	}
	return *o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProjectId) {
		return nil, false
	}
	return o.ProjectId, true
}

// HasProjectId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasProjectId() bool {
	if o != nil && !IsNil(o.ProjectId) {
		return true
	}

	return false
}

// SetProjectId gets a reference to the given string and assigns it to the ProjectId field.
func (o *Dedicatedv1beta1MaintenanceTask) SetProjectId(v string) {
	o.ProjectId = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Dedicatedv1beta1MaintenanceTask) SetDescription(v string) {
	o.Description = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetState() Dedicatedv1beta1MaintenanceTaskState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1MaintenanceTaskState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetStateOk() (*Dedicatedv1beta1MaintenanceTaskState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1MaintenanceTaskState and assigns it to the State field.
func (o *Dedicatedv1beta1MaintenanceTask) SetState(v Dedicatedv1beta1MaintenanceTaskState) {
	o.State = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *Dedicatedv1beta1MaintenanceTask) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetScheduledApplyTime returns the ScheduledApplyTime field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetScheduledApplyTime() time.Time {
	if o == nil || IsNil(o.ScheduledApplyTime) {
		var ret time.Time
		return ret
	}
	return *o.ScheduledApplyTime
}

// GetScheduledApplyTimeOk returns a tuple with the ScheduledApplyTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetScheduledApplyTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ScheduledApplyTime) {
		return nil, false
	}
	return o.ScheduledApplyTime, true
}

// HasScheduledApplyTime returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasScheduledApplyTime() bool {
	if o != nil && !IsNil(o.ScheduledApplyTime) {
		return true
	}

	return false
}

// SetScheduledApplyTime gets a reference to the given time.Time and assigns it to the ScheduledApplyTime field.
func (o *Dedicatedv1beta1MaintenanceTask) SetScheduledApplyTime(v time.Time) {
	o.ScheduledApplyTime = &v
}

// GetDeadline returns the Deadline field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetDeadline() time.Time {
	if o == nil || IsNil(o.Deadline) {
		var ret time.Time
		return ret
	}
	return *o.Deadline
}

// GetDeadlineOk returns a tuple with the Deadline field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetDeadlineOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Deadline) {
		return nil, false
	}
	return o.Deadline, true
}

// HasDeadline returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasDeadline() bool {
	if o != nil && !IsNil(o.Deadline) {
		return true
	}

	return false
}

// SetDeadline gets a reference to the given time.Time and assigns it to the Deadline field.
func (o *Dedicatedv1beta1MaintenanceTask) SetDeadline(v time.Time) {
	o.Deadline = &v
}

func (o Dedicatedv1beta1MaintenanceTask) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1MaintenanceTask) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.MaintenanceTaskId) {
		toSerialize["maintenanceTaskId"] = o.MaintenanceTaskId
	}
	if !IsNil(o.ProjectId) {
		toSerialize["projectId"] = o.ProjectId
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.ScheduledApplyTime) {
		toSerialize["scheduledApplyTime"] = o.ScheduledApplyTime
	}
	if !IsNil(o.Deadline) {
		toSerialize["deadline"] = o.Deadline
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1MaintenanceTask) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1MaintenanceTask := _Dedicatedv1beta1MaintenanceTask{}

	err = json.Unmarshal(data, &varDedicatedv1beta1MaintenanceTask)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1MaintenanceTask(varDedicatedv1beta1MaintenanceTask)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "maintenanceTaskId")
		delete(additionalProperties, "projectId")
		delete(additionalProperties, "description")
		delete(additionalProperties, "state")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "scheduledApplyTime")
		delete(additionalProperties, "deadline")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1MaintenanceTask struct {
	value *Dedicatedv1beta1MaintenanceTask
	isSet bool
}

func (v NullableDedicatedv1beta1MaintenanceTask) Get() *Dedicatedv1beta1MaintenanceTask {
	return v.value
}

func (v *NullableDedicatedv1beta1MaintenanceTask) Set(val *Dedicatedv1beta1MaintenanceTask) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1MaintenanceTask) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1MaintenanceTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1MaintenanceTask(val *Dedicatedv1beta1MaintenanceTask) *NullableDedicatedv1beta1MaintenanceTask {
	return &NullableDedicatedv1beta1MaintenanceTask{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1MaintenanceTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1MaintenanceTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
