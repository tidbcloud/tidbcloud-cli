# ApiListProjectsRsp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NextPageToken** | Pointer to **string** | &#x60;next_page_token&#x60; can be sent in a subsequent call to fetch more results | [optional] 
**Projects** | Pointer to [**[]ApiProject**](ApiProject.md) |  | [optional] 

## Methods

### NewApiListProjectsRsp

`func NewApiListProjectsRsp() *ApiListProjectsRsp`

NewApiListProjectsRsp instantiates a new ApiListProjectsRsp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiListProjectsRspWithDefaults

`func NewApiListProjectsRspWithDefaults() *ApiListProjectsRsp`

NewApiListProjectsRspWithDefaults instantiates a new ApiListProjectsRsp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNextPageToken

`func (o *ApiListProjectsRsp) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *ApiListProjectsRsp) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *ApiListProjectsRsp) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *ApiListProjectsRsp) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.

### GetProjects

`func (o *ApiListProjectsRsp) GetProjects() []ApiProject`

GetProjects returns the Projects field if non-nil, zero value otherwise.

### GetProjectsOk

`func (o *ApiListProjectsRsp) GetProjectsOk() (*[]ApiProject, bool)`

GetProjectsOk returns a tuple with the Projects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjects

`func (o *ApiListProjectsRsp) SetProjects(v []ApiProject)`

SetProjects sets Projects field to given value.

### HasProjects

`func (o *ApiListProjectsRsp) HasProjects() bool`

HasProjects returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


