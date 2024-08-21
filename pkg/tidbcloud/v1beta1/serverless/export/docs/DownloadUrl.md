# DownloadUrl

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the download file. | [optional] 
**Url** | Pointer to **string** | The download url. | [optional] 
**Size** | Pointer to **int64** | The size in bytes of the the download file. | [optional] 

## Methods

### NewDownloadUrl

`func NewDownloadUrl() *DownloadUrl`

NewDownloadUrl instantiates a new DownloadUrl object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDownloadUrlWithDefaults

`func NewDownloadUrlWithDefaults() *DownloadUrl`

NewDownloadUrlWithDefaults instantiates a new DownloadUrl object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DownloadUrl) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DownloadUrl) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DownloadUrl) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DownloadUrl) HasName() bool`

HasName returns a boolean if a field has been set.

### GetUrl

`func (o *DownloadUrl) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *DownloadUrl) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *DownloadUrl) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *DownloadUrl) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetSize

`func (o *DownloadUrl) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *DownloadUrl) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *DownloadUrl) SetSize(v int64)`

SetSize sets Size field to given value.

### HasSize

`func (o *DownloadUrl) HasSize() bool`

HasSize returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


