/*
TiDB Cloud Serverless Export Open API

TiDB Cloud Serverless Export Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package export

import (
	"encoding/json"
)

// checks if the ExportTarget type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExportTarget{}

// ExportTarget struct for ExportTarget
type ExportTarget struct {
	// Optional. The exported file type. Default is LOCAL.
	Type                 *ExportTargetTypeEnum `json:"type,omitempty"`
	S3                   *S3Target             `json:"s3,omitempty"`
	Gcs                  *GCSTarget            `json:"gcs,omitempty"`
	AzureBlob            *AzureBlobTarget      `json:"azureBlob,omitempty"`
	Oss                  *OSSTarget            `json:"oss,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ExportTarget ExportTarget

// NewExportTarget instantiates a new ExportTarget object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExportTarget() *ExportTarget {
	this := ExportTarget{}
	return &this
}

// NewExportTargetWithDefaults instantiates a new ExportTarget object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExportTargetWithDefaults() *ExportTarget {
	this := ExportTarget{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ExportTarget) GetType() ExportTargetTypeEnum {
	if o == nil || IsNil(o.Type) {
		var ret ExportTargetTypeEnum
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportTarget) GetTypeOk() (*ExportTargetTypeEnum, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ExportTarget) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given ExportTargetTypeEnum and assigns it to the Type field.
func (o *ExportTarget) SetType(v ExportTargetTypeEnum) {
	o.Type = &v
}

// GetS3 returns the S3 field value if set, zero value otherwise.
func (o *ExportTarget) GetS3() S3Target {
	if o == nil || IsNil(o.S3) {
		var ret S3Target
		return ret
	}
	return *o.S3
}

// GetS3Ok returns a tuple with the S3 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportTarget) GetS3Ok() (*S3Target, bool) {
	if o == nil || IsNil(o.S3) {
		return nil, false
	}
	return o.S3, true
}

// HasS3 returns a boolean if a field has been set.
func (o *ExportTarget) HasS3() bool {
	if o != nil && !IsNil(o.S3) {
		return true
	}

	return false
}

// SetS3 gets a reference to the given S3Target and assigns it to the S3 field.
func (o *ExportTarget) SetS3(v S3Target) {
	o.S3 = &v
}

// GetGcs returns the Gcs field value if set, zero value otherwise.
func (o *ExportTarget) GetGcs() GCSTarget {
	if o == nil || IsNil(o.Gcs) {
		var ret GCSTarget
		return ret
	}
	return *o.Gcs
}

// GetGcsOk returns a tuple with the Gcs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportTarget) GetGcsOk() (*GCSTarget, bool) {
	if o == nil || IsNil(o.Gcs) {
		return nil, false
	}
	return o.Gcs, true
}

// HasGcs returns a boolean if a field has been set.
func (o *ExportTarget) HasGcs() bool {
	if o != nil && !IsNil(o.Gcs) {
		return true
	}

	return false
}

// SetGcs gets a reference to the given GCSTarget and assigns it to the Gcs field.
func (o *ExportTarget) SetGcs(v GCSTarget) {
	o.Gcs = &v
}

// GetAzureBlob returns the AzureBlob field value if set, zero value otherwise.
func (o *ExportTarget) GetAzureBlob() AzureBlobTarget {
	if o == nil || IsNil(o.AzureBlob) {
		var ret AzureBlobTarget
		return ret
	}
	return *o.AzureBlob
}

// GetAzureBlobOk returns a tuple with the AzureBlob field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportTarget) GetAzureBlobOk() (*AzureBlobTarget, bool) {
	if o == nil || IsNil(o.AzureBlob) {
		return nil, false
	}
	return o.AzureBlob, true
}

// HasAzureBlob returns a boolean if a field has been set.
func (o *ExportTarget) HasAzureBlob() bool {
	if o != nil && !IsNil(o.AzureBlob) {
		return true
	}

	return false
}

// SetAzureBlob gets a reference to the given AzureBlobTarget and assigns it to the AzureBlob field.
func (o *ExportTarget) SetAzureBlob(v AzureBlobTarget) {
	o.AzureBlob = &v
}

// GetOss returns the Oss field value if set, zero value otherwise.
func (o *ExportTarget) GetOss() OSSTarget {
	if o == nil || IsNil(o.Oss) {
		var ret OSSTarget
		return ret
	}
	return *o.Oss
}

// GetOssOk returns a tuple with the Oss field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportTarget) GetOssOk() (*OSSTarget, bool) {
	if o == nil || IsNil(o.Oss) {
		return nil, false
	}
	return o.Oss, true
}

// HasOss returns a boolean if a field has been set.
func (o *ExportTarget) HasOss() bool {
	if o != nil && !IsNil(o.Oss) {
		return true
	}

	return false
}

// SetOss gets a reference to the given OSSTarget and assigns it to the Oss field.
func (o *ExportTarget) SetOss(v OSSTarget) {
	o.Oss = &v
}

func (o ExportTarget) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExportTarget) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.S3) {
		toSerialize["s3"] = o.S3
	}
	if !IsNil(o.Gcs) {
		toSerialize["gcs"] = o.Gcs
	}
	if !IsNil(o.AzureBlob) {
		toSerialize["azureBlob"] = o.AzureBlob
	}
	if !IsNil(o.Oss) {
		toSerialize["oss"] = o.Oss
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ExportTarget) UnmarshalJSON(data []byte) (err error) {
	varExportTarget := _ExportTarget{}

	err = json.Unmarshal(data, &varExportTarget)

	if err != nil {
		return err
	}

	*o = ExportTarget(varExportTarget)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "type")
		delete(additionalProperties, "s3")
		delete(additionalProperties, "gcs")
		delete(additionalProperties, "azureBlob")
		delete(additionalProperties, "oss")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableExportTarget struct {
	value *ExportTarget
	isSet bool
}

func (v NullableExportTarget) Get() *ExportTarget {
	return v.value
}

func (v *NullableExportTarget) Set(val *ExportTarget) {
	v.value = val
	v.isSet = true
}

func (v NullableExportTarget) IsSet() bool {
	return v.isSet
}

func (v *NullableExportTarget) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportTarget(val *ExportTarget) *NullableExportTarget {
	return &NullableExportTarget{value: val, isSet: true}
}

func (v NullableExportTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportTarget) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
