# DownloadExportsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Downloads** | Pointer to [**[]DownloadUrl**](DownloadUrl.md) | The download urls of the export. | [optional] 

## Methods

### NewDownloadExportsResponse

`func NewDownloadExportsResponse() *DownloadExportsResponse`

NewDownloadExportsResponse instantiates a new DownloadExportsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDownloadExportsResponseWithDefaults

`func NewDownloadExportsResponseWithDefaults() *DownloadExportsResponse`

NewDownloadExportsResponseWithDefaults instantiates a new DownloadExportsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDownloads

`func (o *DownloadExportsResponse) GetDownloads() []DownloadUrl`

GetDownloads returns the Downloads field if non-nil, zero value otherwise.

### GetDownloadsOk

`func (o *DownloadExportsResponse) GetDownloadsOk() (*[]DownloadUrl, bool)`

GetDownloadsOk returns a tuple with the Downloads field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDownloads

`func (o *DownloadExportsResponse) SetDownloads(v []DownloadUrl)`

SetDownloads sets Downloads field to given value.

### HasDownloads

`func (o *DownloadExportsResponse) HasDownloads() bool`

HasDownloads returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


