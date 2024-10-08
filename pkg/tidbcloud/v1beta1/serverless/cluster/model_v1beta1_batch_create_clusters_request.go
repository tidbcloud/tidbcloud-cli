/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
	"fmt"
)

// checks if the V1beta1BatchCreateClustersRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1BatchCreateClustersRequest{}

// V1beta1BatchCreateClustersRequest Message for requesting to create a batch of TiDB Cloud Serverless clusters.
type V1beta1BatchCreateClustersRequest struct {
	// The request message specifying the resources to create. A maximum of 1000 clusters can be created in a batch.
	Requests             []TidbCloudOpenApiserverlessv1beta1CreateClusterRequest `json:"requests"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1BatchCreateClustersRequest V1beta1BatchCreateClustersRequest

// NewV1beta1BatchCreateClustersRequest instantiates a new V1beta1BatchCreateClustersRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1BatchCreateClustersRequest(requests []TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) *V1beta1BatchCreateClustersRequest {
	this := V1beta1BatchCreateClustersRequest{}
	this.Requests = requests
	return &this
}

// NewV1beta1BatchCreateClustersRequestWithDefaults instantiates a new V1beta1BatchCreateClustersRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1BatchCreateClustersRequestWithDefaults() *V1beta1BatchCreateClustersRequest {
	this := V1beta1BatchCreateClustersRequest{}
	return &this
}

// GetRequests returns the Requests field value
func (o *V1beta1BatchCreateClustersRequest) GetRequests() []TidbCloudOpenApiserverlessv1beta1CreateClusterRequest {
	if o == nil {
		var ret []TidbCloudOpenApiserverlessv1beta1CreateClusterRequest
		return ret
	}

	return o.Requests
}

// GetRequestsOk returns a tuple with the Requests field value
// and a boolean to check if the value has been set.
func (o *V1beta1BatchCreateClustersRequest) GetRequestsOk() ([]TidbCloudOpenApiserverlessv1beta1CreateClusterRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.Requests, true
}

// SetRequests sets field value
func (o *V1beta1BatchCreateClustersRequest) SetRequests(v []TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) {
	o.Requests = v
}

func (o V1beta1BatchCreateClustersRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1BatchCreateClustersRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["requests"] = o.Requests

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1BatchCreateClustersRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"requests",
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

	varV1beta1BatchCreateClustersRequest := _V1beta1BatchCreateClustersRequest{}

	err = json.Unmarshal(data, &varV1beta1BatchCreateClustersRequest)

	if err != nil {
		return err
	}

	*o = V1beta1BatchCreateClustersRequest(varV1beta1BatchCreateClustersRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "requests")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1BatchCreateClustersRequest struct {
	value *V1beta1BatchCreateClustersRequest
	isSet bool
}

func (v NullableV1beta1BatchCreateClustersRequest) Get() *V1beta1BatchCreateClustersRequest {
	return v.value
}

func (v *NullableV1beta1BatchCreateClustersRequest) Set(val *V1beta1BatchCreateClustersRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1BatchCreateClustersRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1BatchCreateClustersRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1BatchCreateClustersRequest(val *V1beta1BatchCreateClustersRequest) *NullableV1beta1BatchCreateClustersRequest {
	return &NullableV1beta1BatchCreateClustersRequest{value: val, isSet: true}
}

func (v NullableV1beta1BatchCreateClustersRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1BatchCreateClustersRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
