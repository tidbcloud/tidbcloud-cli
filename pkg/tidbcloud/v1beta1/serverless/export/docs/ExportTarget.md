# ExportTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to [**ExportTargetTypeEnum**](ExportTargetTypeEnum.md) | Optional. The exported file type. Default is LOCAL. | [optional] 
**S3** | Pointer to [**S3Target**](S3Target.md) |  | [optional] 
**Gcs** | Pointer to [**GCSTarget**](GCSTarget.md) |  | [optional] 
**AzureBlob** | Pointer to [**AzureBlobTarget**](AzureBlobTarget.md) |  | [optional] 

## Methods

### NewExportTarget

`func NewExportTarget() *ExportTarget`

NewExportTarget instantiates a new ExportTarget object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportTargetWithDefaults

`func NewExportTargetWithDefaults() *ExportTarget`

NewExportTargetWithDefaults instantiates a new ExportTarget object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *ExportTarget) GetType() ExportTargetTypeEnum`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ExportTarget) GetTypeOk() (*ExportTargetTypeEnum, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ExportTarget) SetType(v ExportTargetTypeEnum)`

SetType sets Type field to given value.

### HasType

`func (o *ExportTarget) HasType() bool`

HasType returns a boolean if a field has been set.

### GetS3

`func (o *ExportTarget) GetS3() S3Target`

GetS3 returns the S3 field if non-nil, zero value otherwise.

### GetS3Ok

`func (o *ExportTarget) GetS3Ok() (*S3Target, bool)`

GetS3Ok returns a tuple with the S3 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetS3

`func (o *ExportTarget) SetS3(v S3Target)`

SetS3 sets S3 field to given value.

### HasS3

`func (o *ExportTarget) HasS3() bool`

HasS3 returns a boolean if a field has been set.

### GetGcs

`func (o *ExportTarget) GetGcs() GCSTarget`

GetGcs returns the Gcs field if non-nil, zero value otherwise.

### GetGcsOk

`func (o *ExportTarget) GetGcsOk() (*GCSTarget, bool)`

GetGcsOk returns a tuple with the Gcs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGcs

`func (o *ExportTarget) SetGcs(v GCSTarget)`

SetGcs sets Gcs field to given value.

### HasGcs

`func (o *ExportTarget) HasGcs() bool`

HasGcs returns a boolean if a field has been set.

### GetAzureBlob

`func (o *ExportTarget) GetAzureBlob() AzureBlobTarget`

GetAzureBlob returns the AzureBlob field if non-nil, zero value otherwise.

### GetAzureBlobOk

`func (o *ExportTarget) GetAzureBlobOk() (*AzureBlobTarget, bool)`

GetAzureBlobOk returns a tuple with the AzureBlob field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAzureBlob

`func (o *ExportTarget) SetAzureBlob(v AzureBlobTarget)`

SetAzureBlob sets AzureBlob field to given value.

### HasAzureBlob

`func (o *ExportTarget) HasAzureBlob() bool`

HasAzureBlob returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


