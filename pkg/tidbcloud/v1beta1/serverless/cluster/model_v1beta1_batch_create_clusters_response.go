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

// checks if the V1beta1BatchCreateClustersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1BatchCreateClustersResponse{}

// V1beta1BatchCreateClustersResponse Responses message to the request for creating a batch of TiDB Cloud Serverless clusters.
type V1beta1BatchCreateClustersResponse struct {
	// Clusters created.
	Clusters             []TidbCloudOpenApiserverlessv1beta1Cluster `json:"clusters,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1BatchCreateClustersResponse V1beta1BatchCreateClustersResponse

// NewV1beta1BatchCreateClustersResponse instantiates a new V1beta1BatchCreateClustersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1BatchCreateClustersResponse() *V1beta1BatchCreateClustersResponse {
	this := V1beta1BatchCreateClustersResponse{}
	return &this
}

// NewV1beta1BatchCreateClustersResponseWithDefaults instantiates a new V1beta1BatchCreateClustersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1BatchCreateClustersResponseWithDefaults() *V1beta1BatchCreateClustersResponse {
	this := V1beta1BatchCreateClustersResponse{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise.
func (o *V1beta1BatchCreateClustersResponse) GetClusters() []TidbCloudOpenApiserverlessv1beta1Cluster {
	if o == nil || IsNil(o.Clusters) {
		var ret []TidbCloudOpenApiserverlessv1beta1Cluster
		return ret
	}
	return o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1BatchCreateClustersResponse) GetClustersOk() ([]TidbCloudOpenApiserverlessv1beta1Cluster, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}
	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *V1beta1BatchCreateClustersResponse) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []TidbCloudOpenApiserverlessv1beta1Cluster and assigns it to the Clusters field.
func (o *V1beta1BatchCreateClustersResponse) SetClusters(v []TidbCloudOpenApiserverlessv1beta1Cluster) {
	o.Clusters = v
}

func (o V1beta1BatchCreateClustersResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1BatchCreateClustersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Clusters) {
		toSerialize["clusters"] = o.Clusters
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1BatchCreateClustersResponse) UnmarshalJSON(data []byte) (err error) {
	varV1beta1BatchCreateClustersResponse := _V1beta1BatchCreateClustersResponse{}

	err = json.Unmarshal(data, &varV1beta1BatchCreateClustersResponse)

	if err != nil {
		return err
	}

	*o = V1beta1BatchCreateClustersResponse(varV1beta1BatchCreateClustersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusters")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1BatchCreateClustersResponse struct {
	value *V1beta1BatchCreateClustersResponse
	isSet bool
}

func (v NullableV1beta1BatchCreateClustersResponse) Get() *V1beta1BatchCreateClustersResponse {
	return v.value
}

func (v *NullableV1beta1BatchCreateClustersResponse) Set(val *V1beta1BatchCreateClustersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1BatchCreateClustersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1BatchCreateClustersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1BatchCreateClustersResponse(val *V1beta1BatchCreateClustersResponse) *NullableV1beta1BatchCreateClustersResponse {
	return &NullableV1beta1BatchCreateClustersResponse{value: val, isSet: true}
}

func (v NullableV1beta1BatchCreateClustersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1BatchCreateClustersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}