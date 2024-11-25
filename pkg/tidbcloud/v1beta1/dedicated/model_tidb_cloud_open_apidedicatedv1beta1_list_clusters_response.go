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

// checks if the TidbCloudOpenApidedicatedv1beta1ListClustersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApidedicatedv1beta1ListClustersResponse{}

// TidbCloudOpenApidedicatedv1beta1ListClustersResponse struct for TidbCloudOpenApidedicatedv1beta1ListClustersResponse
type TidbCloudOpenApidedicatedv1beta1ListClustersResponse struct {
	Clusters []TidbCloudOpenApidedicatedv1beta1Cluster `json:"clusters,omitempty"`
	// The total number of clusters that matched the query.
	TotalSize *int32 `json:"totalSize,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page. If this field is omitted, there are no subsequent pages.
	NextPageToken *string `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApidedicatedv1beta1ListClustersResponse TidbCloudOpenApidedicatedv1beta1ListClustersResponse

// NewTidbCloudOpenApidedicatedv1beta1ListClustersResponse instantiates a new TidbCloudOpenApidedicatedv1beta1ListClustersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApidedicatedv1beta1ListClustersResponse() *TidbCloudOpenApidedicatedv1beta1ListClustersResponse {
	this := TidbCloudOpenApidedicatedv1beta1ListClustersResponse{}
	return &this
}

// NewTidbCloudOpenApidedicatedv1beta1ListClustersResponseWithDefaults instantiates a new TidbCloudOpenApidedicatedv1beta1ListClustersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApidedicatedv1beta1ListClustersResponseWithDefaults() *TidbCloudOpenApidedicatedv1beta1ListClustersResponse {
	this := TidbCloudOpenApidedicatedv1beta1ListClustersResponse{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetClusters() []TidbCloudOpenApidedicatedv1beta1Cluster {
	if o == nil || IsNil(o.Clusters) {
		var ret []TidbCloudOpenApidedicatedv1beta1Cluster
		return ret
	}
	return o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetClustersOk() ([]TidbCloudOpenApidedicatedv1beta1Cluster, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}
	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []TidbCloudOpenApidedicatedv1beta1Cluster and assigns it to the Clusters field.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) SetClusters(v []TidbCloudOpenApidedicatedv1beta1Cluster) {
	o.Clusters = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o TidbCloudOpenApidedicatedv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApidedicatedv1beta1ListClustersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Clusters) {
		toSerialize["clusters"] = o.Clusters
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

func (o *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) UnmarshalJSON(data []byte) (err error) {
	varTidbCloudOpenApidedicatedv1beta1ListClustersResponse := _TidbCloudOpenApidedicatedv1beta1ListClustersResponse{}

	err = json.Unmarshal(data, &varTidbCloudOpenApidedicatedv1beta1ListClustersResponse)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApidedicatedv1beta1ListClustersResponse(varTidbCloudOpenApidedicatedv1beta1ListClustersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusters")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse struct {
	value *TidbCloudOpenApidedicatedv1beta1ListClustersResponse
	isSet bool
}

func (v NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) Get() *TidbCloudOpenApidedicatedv1beta1ListClustersResponse {
	return v.value
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) Set(val *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse(val *TidbCloudOpenApidedicatedv1beta1ListClustersResponse) *NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse {
	return &NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1ListClustersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


