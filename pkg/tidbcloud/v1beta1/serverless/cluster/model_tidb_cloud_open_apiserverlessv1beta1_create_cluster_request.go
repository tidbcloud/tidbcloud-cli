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

// checks if the TidbCloudOpenApiserverlessv1beta1CreateClusterRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApiserverlessv1beta1CreateClusterRequest{}

// TidbCloudOpenApiserverlessv1beta1CreateClusterRequest Message for requesting to create a TiDB Cloud Serverless cluster.
type TidbCloudOpenApiserverlessv1beta1CreateClusterRequest struct {
	// Required. The cluster to be created.
	Cluster              TidbCloudOpenApiserverlessv1beta1Cluster `json:"cluster"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApiserverlessv1beta1CreateClusterRequest TidbCloudOpenApiserverlessv1beta1CreateClusterRequest

// NewTidbCloudOpenApiserverlessv1beta1CreateClusterRequest instantiates a new TidbCloudOpenApiserverlessv1beta1CreateClusterRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApiserverlessv1beta1CreateClusterRequest(cluster TidbCloudOpenApiserverlessv1beta1Cluster) *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest {
	this := TidbCloudOpenApiserverlessv1beta1CreateClusterRequest{}
	this.Cluster = cluster
	return &this
}

// NewTidbCloudOpenApiserverlessv1beta1CreateClusterRequestWithDefaults instantiates a new TidbCloudOpenApiserverlessv1beta1CreateClusterRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApiserverlessv1beta1CreateClusterRequestWithDefaults() *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest {
	this := TidbCloudOpenApiserverlessv1beta1CreateClusterRequest{}
	return &this
}

// GetCluster returns the Cluster field value
func (o *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) GetCluster() TidbCloudOpenApiserverlessv1beta1Cluster {
	if o == nil {
		var ret TidbCloudOpenApiserverlessv1beta1Cluster
		return ret
	}

	return o.Cluster
}

// GetClusterOk returns a tuple with the Cluster field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) GetClusterOk() (*TidbCloudOpenApiserverlessv1beta1Cluster, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Cluster, true
}

// SetCluster sets field value
func (o *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) SetCluster(v TidbCloudOpenApiserverlessv1beta1Cluster) {
	o.Cluster = v
}

func (o TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["cluster"] = o.Cluster

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"cluster",
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

	varTidbCloudOpenApiserverlessv1beta1CreateClusterRequest := _TidbCloudOpenApiserverlessv1beta1CreateClusterRequest{}

	err = json.Unmarshal(data, &varTidbCloudOpenApiserverlessv1beta1CreateClusterRequest)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApiserverlessv1beta1CreateClusterRequest(varTidbCloudOpenApiserverlessv1beta1CreateClusterRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "cluster")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest struct {
	value *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest
	isSet bool
}

func (v NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) Get() *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest {
	return v.value
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) Set(val *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest(val *TidbCloudOpenApiserverlessv1beta1CreateClusterRequest) *NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest {
	return &NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1CreateClusterRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
