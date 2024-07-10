# V1beta1ListBranchesResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Branches** | Pointer to [**[]V1beta1Branch**](V1beta1Branch.md) |  | [optional] 
**NextPageToken** | Pointer to **string** | A token identifying a page of results the server should return. | [optional] 
**TotalSize** | Pointer to **int64** |  | [optional] 

## Methods

### NewV1beta1ListBranchesResponse

`func NewV1beta1ListBranchesResponse() *V1beta1ListBranchesResponse`

NewV1beta1ListBranchesResponse instantiates a new V1beta1ListBranchesResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1beta1ListBranchesResponseWithDefaults

`func NewV1beta1ListBranchesResponseWithDefaults() *V1beta1ListBranchesResponse`

NewV1beta1ListBranchesResponseWithDefaults instantiates a new V1beta1ListBranchesResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBranches

`func (o *V1beta1ListBranchesResponse) GetBranches() []V1beta1Branch`

GetBranches returns the Branches field if non-nil, zero value otherwise.

### GetBranchesOk

`func (o *V1beta1ListBranchesResponse) GetBranchesOk() (*[]V1beta1Branch, bool)`

GetBranchesOk returns a tuple with the Branches field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBranches

`func (o *V1beta1ListBranchesResponse) SetBranches(v []V1beta1Branch)`

SetBranches sets Branches field to given value.

### HasBranches

`func (o *V1beta1ListBranchesResponse) HasBranches() bool`

HasBranches returns a boolean if a field has been set.

### GetNextPageToken

`func (o *V1beta1ListBranchesResponse) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *V1beta1ListBranchesResponse) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *V1beta1ListBranchesResponse) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *V1beta1ListBranchesResponse) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.

### GetTotalSize

`func (o *V1beta1ListBranchesResponse) GetTotalSize() int64`

GetTotalSize returns the TotalSize field if non-nil, zero value otherwise.

### GetTotalSizeOk

`func (o *V1beta1ListBranchesResponse) GetTotalSizeOk() (*int64, bool)`

GetTotalSizeOk returns a tuple with the TotalSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalSize

`func (o *V1beta1ListBranchesResponse) SetTotalSize(v int64)`

SetTotalSize sets TotalSize field to given value.

### HasTotalSize

`func (o *V1beta1ListBranchesResponse) HasTotalSize() bool`

HasTotalSize returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


