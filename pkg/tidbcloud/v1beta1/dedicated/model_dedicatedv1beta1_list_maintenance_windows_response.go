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

// checks if the Dedicatedv1beta1ListMaintenanceWindowsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1ListMaintenanceWindowsResponse{}

// Dedicatedv1beta1ListMaintenanceWindowsResponse struct for Dedicatedv1beta1ListMaintenanceWindowsResponse
type Dedicatedv1beta1ListMaintenanceWindowsResponse struct {
	MaintenanceWindows []Dedicatedv1beta1MaintenanceWindow `json:"maintenanceWindows,omitempty"`
	// The total number of maintenance windows that matched the query.
	TotalSize *int32 `json:"totalSize,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page. If this field is omitted, there are no subsequent pages.
	NextPageToken        *string `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1ListMaintenanceWindowsResponse Dedicatedv1beta1ListMaintenanceWindowsResponse

// NewDedicatedv1beta1ListMaintenanceWindowsResponse instantiates a new Dedicatedv1beta1ListMaintenanceWindowsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1ListMaintenanceWindowsResponse() *Dedicatedv1beta1ListMaintenanceWindowsResponse {
	this := Dedicatedv1beta1ListMaintenanceWindowsResponse{}
	return &this
}

// NewDedicatedv1beta1ListMaintenanceWindowsResponseWithDefaults instantiates a new Dedicatedv1beta1ListMaintenanceWindowsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1ListMaintenanceWindowsResponseWithDefaults() *Dedicatedv1beta1ListMaintenanceWindowsResponse {
	this := Dedicatedv1beta1ListMaintenanceWindowsResponse{}
	return &this
}

// GetMaintenanceWindows returns the MaintenanceWindows field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetMaintenanceWindows() []Dedicatedv1beta1MaintenanceWindow {
	if o == nil || IsNil(o.MaintenanceWindows) {
		var ret []Dedicatedv1beta1MaintenanceWindow
		return ret
	}
	return o.MaintenanceWindows
}

// GetMaintenanceWindowsOk returns a tuple with the MaintenanceWindows field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetMaintenanceWindowsOk() ([]Dedicatedv1beta1MaintenanceWindow, bool) {
	if o == nil || IsNil(o.MaintenanceWindows) {
		return nil, false
	}
	return o.MaintenanceWindows, true
}

// HasMaintenanceWindows returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) HasMaintenanceWindows() bool {
	if o != nil && !IsNil(o.MaintenanceWindows) {
		return true
	}

	return false
}

// SetMaintenanceWindows gets a reference to the given []Dedicatedv1beta1MaintenanceWindow and assigns it to the MaintenanceWindows field.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) SetMaintenanceWindows(v []Dedicatedv1beta1MaintenanceWindow) {
	o.MaintenanceWindows = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o Dedicatedv1beta1ListMaintenanceWindowsResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1ListMaintenanceWindowsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.MaintenanceWindows) {
		toSerialize["maintenanceWindows"] = o.MaintenanceWindows
	}
	if !IsNil(o.TotalSize) {
		toSerialize["totalSize"] = o.TotalSize
	}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1ListMaintenanceWindowsResponse) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1ListMaintenanceWindowsResponse := _Dedicatedv1beta1ListMaintenanceWindowsResponse{}

	err = json.Unmarshal(data, &varDedicatedv1beta1ListMaintenanceWindowsResponse)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1ListMaintenanceWindowsResponse(varDedicatedv1beta1ListMaintenanceWindowsResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "maintenanceWindows")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1ListMaintenanceWindowsResponse struct {
	value *Dedicatedv1beta1ListMaintenanceWindowsResponse
	isSet bool
}

func (v NullableDedicatedv1beta1ListMaintenanceWindowsResponse) Get() *Dedicatedv1beta1ListMaintenanceWindowsResponse {
	return v.value
}

func (v *NullableDedicatedv1beta1ListMaintenanceWindowsResponse) Set(val *Dedicatedv1beta1ListMaintenanceWindowsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1ListMaintenanceWindowsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1ListMaintenanceWindowsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1ListMaintenanceWindowsResponse(val *Dedicatedv1beta1ListMaintenanceWindowsResponse) *NullableDedicatedv1beta1ListMaintenanceWindowsResponse {
	return &NullableDedicatedv1beta1ListMaintenanceWindowsResponse{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1ListMaintenanceWindowsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1ListMaintenanceWindowsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
