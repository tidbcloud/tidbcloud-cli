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

// checks if the TidbCloudOpenApiserverlessv1beta1ListClustersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}

// TidbCloudOpenApiserverlessv1beta1ListClustersResponse Responses message to the request for listing of TiDB Cloud Serverless clusters.
type TidbCloudOpenApiserverlessv1beta1ListClustersResponse struct {
	// A list of clusters.
	Clusters []TidbCloudOpenApiserverlessv1beta1Cluster `json:"clusters,omitempty"`
	// Token provided to retrieve the next page of results.
	NextPageToken *string `json:"nextPageToken,omitempty"`
	// Total number of available clusters.
	TotalSize            *int64 `json:"totalSize,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApiserverlessv1beta1ListClustersResponse TidbCloudOpenApiserverlessv1beta1ListClustersResponse

// NewTidbCloudOpenApiserverlessv1beta1ListClustersResponse instantiates a new TidbCloudOpenApiserverlessv1beta1ListClustersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApiserverlessv1beta1ListClustersResponse() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	this := TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}
	return &this
}

// NewTidbCloudOpenApiserverlessv1beta1ListClustersResponseWithDefaults instantiates a new TidbCloudOpenApiserverlessv1beta1ListClustersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApiserverlessv1beta1ListClustersResponseWithDefaults() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	this := TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetClusters() []TidbCloudOpenApiserverlessv1beta1Cluster {
	if o == nil || IsNil(o.Clusters) {
		var ret []TidbCloudOpenApiserverlessv1beta1Cluster
		return ret
	}
	return o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetClustersOk() ([]TidbCloudOpenApiserverlessv1beta1Cluster, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}
	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []TidbCloudOpenApiserverlessv1beta1Cluster and assigns it to the Clusters field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetClusters(v []TidbCloudOpenApiserverlessv1beta1Cluster) {
	o.Clusters = v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetTotalSize() int64 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int64
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetTotalSizeOk() (*int64, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int64 and assigns it to the TotalSize field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetTotalSize(v int64) {
	o.TotalSize = &v
}

func (o TidbCloudOpenApiserverlessv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApiserverlessv1beta1ListClustersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Clusters) {
		toSerialize["clusters"] = o.Clusters
	}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}
	if !IsNil(o.TotalSize) {
		toSerialize["totalSize"] = o.TotalSize
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) UnmarshalJSON(data []byte) (err error) {
	varTidbCloudOpenApiserverlessv1beta1ListClustersResponse := _TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}

	err = json.Unmarshal(data, &varTidbCloudOpenApiserverlessv1beta1ListClustersResponse)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApiserverlessv1beta1ListClustersResponse(varTidbCloudOpenApiserverlessv1beta1ListClustersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusters")
		delete(additionalProperties, "nextPageToken")
		delete(additionalProperties, "totalSize")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse struct {
	value *TidbCloudOpenApiserverlessv1beta1ListClustersResponse
	isSet bool
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Get() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	return v.value
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Set(val *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse(val *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	return &NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
