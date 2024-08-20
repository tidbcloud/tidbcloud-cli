# ListExportsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Exports** | Pointer to [**[]Export**](Export.md) | A list of exports. | [optional] 
**NextPageToken** | Pointer to **string** | Token provided to retrieve the next page of results. | [optional] 
**TotalSize** | Pointer to **int64** | Total number of exports. | [optional] 

## Methods

### NewListExportsResponse

`func NewListExportsResponse() *ListExportsResponse`

NewListExportsResponse instantiates a new ListExportsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListExportsResponseWithDefaults

`func NewListExportsResponseWithDefaults() *ListExportsResponse`

NewListExportsResponseWithDefaults instantiates a new ListExportsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExports

`func (o *ListExportsResponse) GetExports() []Export`

GetExports returns the Exports field if non-nil, zero value otherwise.

### GetExportsOk

`func (o *ListExportsResponse) GetExportsOk() (*[]Export, bool)`

GetExportsOk returns a tuple with the Exports field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExports

`func (o *ListExportsResponse) SetExports(v []Export)`

SetExports sets Exports field to given value.

### HasExports

`func (o *ListExportsResponse) HasExports() bool`

HasExports returns a boolean if a field has been set.

### GetNextPageToken

`func (o *ListExportsResponse) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *ListExportsResponse) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *ListExportsResponse) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *ListExportsResponse) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.

### GetTotalSize

`func (o *ListExportsResponse) GetTotalSize() int64`

GetTotalSize returns the TotalSize field if non-nil, zero value otherwise.

### GetTotalSizeOk

`func (o *ListExportsResponse) GetTotalSizeOk() (*int64, bool)`

GetTotalSizeOk returns a tuple with the TotalSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalSize

`func (o *ListExportsResponse) SetTotalSize(v int64)`

SetTotalSize sets TotalSize field to given value.

### HasTotalSize

`func (o *ListExportsResponse) HasTotalSize() bool`

HasTotalSize returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


