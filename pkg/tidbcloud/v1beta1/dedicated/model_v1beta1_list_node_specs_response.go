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

// checks if the V1beta1ListNodeSpecsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ListNodeSpecsResponse{}

// V1beta1ListNodeSpecsResponse struct for V1beta1ListNodeSpecsResponse
type V1beta1ListNodeSpecsResponse struct {
	NodeSpecs []Dedicatedv1beta1NodeSpec `json:"nodeSpecs,omitempty"`
	// The total number of node specs that matched the query.
	TotalSize *int32 `json:"totalSize,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page. If this field is omitted, there are no subsequent pages.
	NextPageToken *string `json:"nextPageToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ListNodeSpecsResponse V1beta1ListNodeSpecsResponse

// NewV1beta1ListNodeSpecsResponse instantiates a new V1beta1ListNodeSpecsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ListNodeSpecsResponse() *V1beta1ListNodeSpecsResponse {
	this := V1beta1ListNodeSpecsResponse{}
	return &this
}

// NewV1beta1ListNodeSpecsResponseWithDefaults instantiates a new V1beta1ListNodeSpecsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ListNodeSpecsResponseWithDefaults() *V1beta1ListNodeSpecsResponse {
	this := V1beta1ListNodeSpecsResponse{}
	return &this
}

// GetNodeSpecs returns the NodeSpecs field value if set, zero value otherwise.
func (o *V1beta1ListNodeSpecsResponse) GetNodeSpecs() []Dedicatedv1beta1NodeSpec {
	if o == nil || IsNil(o.NodeSpecs) {
		var ret []Dedicatedv1beta1NodeSpec
		return ret
	}
	return o.NodeSpecs
}

// GetNodeSpecsOk returns a tuple with the NodeSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNodeSpecsResponse) GetNodeSpecsOk() ([]Dedicatedv1beta1NodeSpec, bool) {
	if o == nil || IsNil(o.NodeSpecs) {
		return nil, false
	}
	return o.NodeSpecs, true
}

// HasNodeSpecs returns a boolean if a field has been set.
func (o *V1beta1ListNodeSpecsResponse) HasNodeSpecs() bool {
	if o != nil && !IsNil(o.NodeSpecs) {
		return true
	}

	return false
}

// SetNodeSpecs gets a reference to the given []Dedicatedv1beta1NodeSpec and assigns it to the NodeSpecs field.
func (o *V1beta1ListNodeSpecsResponse) SetNodeSpecs(v []Dedicatedv1beta1NodeSpec) {
	o.NodeSpecs = v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *V1beta1ListNodeSpecsResponse) GetTotalSize() int32 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNodeSpecsResponse) GetTotalSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *V1beta1ListNodeSpecsResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *V1beta1ListNodeSpecsResponse) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *V1beta1ListNodeSpecsResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ListNodeSpecsResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *V1beta1ListNodeSpecsResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *V1beta1ListNodeSpecsResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

func (o V1beta1ListNodeSpecsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ListNodeSpecsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.NodeSpecs) {
		toSerialize["nodeSpecs"] = o.NodeSpecs
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

func (o *V1beta1ListNodeSpecsResponse) UnmarshalJSON(data []byte) (err error) {
	varV1beta1ListNodeSpecsResponse := _V1beta1ListNodeSpecsResponse{}

	err = json.Unmarshal(data, &varV1beta1ListNodeSpecsResponse)

	if err != nil {
		return err
	}

	*o = V1beta1ListNodeSpecsResponse(varV1beta1ListNodeSpecsResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "nodeSpecs")
		delete(additionalProperties, "totalSize")
		delete(additionalProperties, "nextPageToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ListNodeSpecsResponse struct {
	value *V1beta1ListNodeSpecsResponse
	isSet bool
}

func (v NullableV1beta1ListNodeSpecsResponse) Get() *V1beta1ListNodeSpecsResponse {
	return v.value
}

func (v *NullableV1beta1ListNodeSpecsResponse) Set(val *V1beta1ListNodeSpecsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ListNodeSpecsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ListNodeSpecsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ListNodeSpecsResponse(val *V1beta1ListNodeSpecsResponse) *NullableV1beta1ListNodeSpecsResponse {
	return &NullableV1beta1ListNodeSpecsResponse{value: val, isSet: true}
}

func (v NullableV1beta1ListNodeSpecsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ListNodeSpecsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

