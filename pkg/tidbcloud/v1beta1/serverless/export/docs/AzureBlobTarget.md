# AzureBlobTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthType** | [**ExportAzureBlobAuthTypeEnum**](ExportAzureBlobAuthTypeEnum.md) | The Azure Blob URI of the export target. | 
**SasToken** | Pointer to **string** | The sas token. This field is input-only. | [optional] 
**Uri** | **string** | The Azure Blob URI of the export target. For example: https://accountname.blob.core.windows.net/container/folder. | 

## Methods

### NewAzureBlobTarget

`func NewAzureBlobTarget(authType ExportAzureBlobAuthTypeEnum, uri string, ) *AzureBlobTarget`

NewAzureBlobTarget instantiates a new AzureBlobTarget object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAzureBlobTargetWithDefaults

`func NewAzureBlobTargetWithDefaults() *AzureBlobTarget`

NewAzureBlobTargetWithDefaults instantiates a new AzureBlobTarget object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthType

`func (o *AzureBlobTarget) GetAuthType() ExportAzureBlobAuthTypeEnum`

GetAuthType returns the AuthType field if non-nil, zero value otherwise.

### GetAuthTypeOk

`func (o *AzureBlobTarget) GetAuthTypeOk() (*ExportAzureBlobAuthTypeEnum, bool)`

GetAuthTypeOk returns a tuple with the AuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthType

`func (o *AzureBlobTarget) SetAuthType(v ExportAzureBlobAuthTypeEnum)`

SetAuthType sets AuthType field to given value.


### GetSasToken

`func (o *AzureBlobTarget) GetSasToken() string`

GetSasToken returns the SasToken field if non-nil, zero value otherwise.

### GetSasTokenOk

`func (o *AzureBlobTarget) GetSasTokenOk() (*string, bool)`

GetSasTokenOk returns a tuple with the SasToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSasToken

`func (o *AzureBlobTarget) SetSasToken(v string)`

SetSasToken sets SasToken field to given value.

### HasSasToken

`func (o *AzureBlobTarget) HasSasToken() bool`

HasSasToken returns a boolean if a field has been set.

### GetUri

`func (o *AzureBlobTarget) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *AzureBlobTarget) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *AzureBlobTarget) SetUri(v string)`

SetUri sets Uri field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


