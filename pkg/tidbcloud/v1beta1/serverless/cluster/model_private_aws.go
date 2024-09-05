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

// checks if the PrivateAWS type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PrivateAWS{}

// PrivateAWS Message for AWS PrivateLink information.
type PrivateAWS struct {
	// Output_only. The AWS service name for private access.
	ServiceName *string `json:"serviceName,omitempty"`
	// Output_only. The availability zones that the service is available in.
	AvailabilityZone []string `json:"availabilityZone,omitempty"`
}

// NewPrivateAWS instantiates a new PrivateAWS object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateAWS() *PrivateAWS {
	this := PrivateAWS{}
	return &this
}

// NewPrivateAWSWithDefaults instantiates a new PrivateAWS object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateAWSWithDefaults() *PrivateAWS {
	this := PrivateAWS{}
	return &this
}

// GetServiceName returns the ServiceName field value if set, zero value otherwise.
func (o *PrivateAWS) GetServiceName() string {
	if o == nil || IsNil(o.ServiceName) {
		var ret string
		return ret
	}
	return *o.ServiceName
}

// GetServiceNameOk returns a tuple with the ServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateAWS) GetServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceName) {
		return nil, false
	}
	return o.ServiceName, true
}

// HasServiceName returns a boolean if a field has been set.
func (o *PrivateAWS) HasServiceName() bool {
	if o != nil && !IsNil(o.ServiceName) {
		return true
	}

	return false
}

// SetServiceName gets a reference to the given string and assigns it to the ServiceName field.
func (o *PrivateAWS) SetServiceName(v string) {
	o.ServiceName = &v
}

// GetAvailabilityZone returns the AvailabilityZone field value if set, zero value otherwise.
func (o *PrivateAWS) GetAvailabilityZone() []string {
	if o == nil || IsNil(o.AvailabilityZone) {
		var ret []string
		return ret
	}
	return o.AvailabilityZone
}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateAWS) GetAvailabilityZoneOk() ([]string, bool) {
	if o == nil || IsNil(o.AvailabilityZone) {
		return nil, false
	}
	return o.AvailabilityZone, true
}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *PrivateAWS) HasAvailabilityZone() bool {
	if o != nil && !IsNil(o.AvailabilityZone) {
		return true
	}

	return false
}

// SetAvailabilityZone gets a reference to the given []string and assigns it to the AvailabilityZone field.
func (o *PrivateAWS) SetAvailabilityZone(v []string) {
	o.AvailabilityZone = v
}

func (o PrivateAWS) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PrivateAWS) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ServiceName) {
		toSerialize["serviceName"] = o.ServiceName
	}
	if !IsNil(o.AvailabilityZone) {
		toSerialize["availabilityZone"] = o.AvailabilityZone
	}
	return toSerialize, nil
}

type NullablePrivateAWS struct {
	value *PrivateAWS
	isSet bool
}

func (v NullablePrivateAWS) Get() *PrivateAWS {
	return v.value
}

func (v *NullablePrivateAWS) Set(val *PrivateAWS) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivateAWS) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivateAWS) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivateAWS(val *PrivateAWS) *NullablePrivateAWS {
	return &NullablePrivateAWS{value: val, isSet: true}
}

func (v NullablePrivateAWS) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivateAWS) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
