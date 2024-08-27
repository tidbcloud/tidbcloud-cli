# ApiListSqlUsersRsp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NextPageToken** | Pointer to **string** | &#x60;next_page_token&#x60; can be sent in a subsequent call to fetch more results | [optional] 
**SqlUsers** | Pointer to [**[]ApiSqlUser**](ApiSqlUser.md) | SqlUsers []*SqlUser &#x60;json:\&quot;sqlUsers\&quot;&#x60; | [optional] 

## Methods

### NewApiListSqlUsersRsp

`func NewApiListSqlUsersRsp() *ApiListSqlUsersRsp`

NewApiListSqlUsersRsp instantiates a new ApiListSqlUsersRsp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiListSqlUsersRspWithDefaults

`func NewApiListSqlUsersRspWithDefaults() *ApiListSqlUsersRsp`

NewApiListSqlUsersRspWithDefaults instantiates a new ApiListSqlUsersRsp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNextPageToken

`func (o *ApiListSqlUsersRsp) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *ApiListSqlUsersRsp) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *ApiListSqlUsersRsp) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *ApiListSqlUsersRsp) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.

### GetSqlUsers

`func (o *ApiListSqlUsersRsp) GetSqlUsers() []ApiSqlUser`

GetSqlUsers returns the SqlUsers field if non-nil, zero value otherwise.

### GetSqlUsersOk

`func (o *ApiListSqlUsersRsp) GetSqlUsersOk() (*[]ApiSqlUser, bool)`

GetSqlUsersOk returns a tuple with the SqlUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSqlUsers

`func (o *ApiListSqlUsersRsp) SetSqlUsers(v []ApiSqlUser)`

SetSqlUsers sets SqlUsers field to given value.

### HasSqlUsers

`func (o *ApiListSqlUsersRsp) HasSqlUsers() bool`

HasSqlUsers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


