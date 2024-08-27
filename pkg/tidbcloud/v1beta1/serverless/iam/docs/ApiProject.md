# ApiProject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AwsCmekEnabled** | Pointer to **bool** | The AwsCmekEnabled of the project. | [optional] 
**ClusterCount** | Pointer to **int32** | Number of cluster_ in the project. | [optional] 
**CreateTimestamp** | Pointer to **string** | The create time key of the project. | [optional] 
**Id** | Pointer to **string** | The id of the project. | [optional] 
**Name** | Pointer to **string** | The name of the API key. | [optional] 
**OrgId** | Pointer to **string** | The org id  of the project. | [optional] 
**UserCount** | Pointer to **int32** | Number of users in the project. | [optional] 

## Methods

### NewApiProject

`func NewApiProject() *ApiProject`

NewApiProject instantiates a new ApiProject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiProjectWithDefaults

`func NewApiProjectWithDefaults() *ApiProject`

NewApiProjectWithDefaults instantiates a new ApiProject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAwsCmekEnabled

`func (o *ApiProject) GetAwsCmekEnabled() bool`

GetAwsCmekEnabled returns the AwsCmekEnabled field if non-nil, zero value otherwise.

### GetAwsCmekEnabledOk

`func (o *ApiProject) GetAwsCmekEnabledOk() (*bool, bool)`

GetAwsCmekEnabledOk returns a tuple with the AwsCmekEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAwsCmekEnabled

`func (o *ApiProject) SetAwsCmekEnabled(v bool)`

SetAwsCmekEnabled sets AwsCmekEnabled field to given value.

### HasAwsCmekEnabled

`func (o *ApiProject) HasAwsCmekEnabled() bool`

HasAwsCmekEnabled returns a boolean if a field has been set.

### GetClusterCount

`func (o *ApiProject) GetClusterCount() int32`

GetClusterCount returns the ClusterCount field if non-nil, zero value otherwise.

### GetClusterCountOk

`func (o *ApiProject) GetClusterCountOk() (*int32, bool)`

GetClusterCountOk returns a tuple with the ClusterCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterCount

`func (o *ApiProject) SetClusterCount(v int32)`

SetClusterCount sets ClusterCount field to given value.

### HasClusterCount

`func (o *ApiProject) HasClusterCount() bool`

HasClusterCount returns a boolean if a field has been set.

### GetCreateTimestamp

`func (o *ApiProject) GetCreateTimestamp() string`

GetCreateTimestamp returns the CreateTimestamp field if non-nil, zero value otherwise.

### GetCreateTimestampOk

`func (o *ApiProject) GetCreateTimestampOk() (*string, bool)`

GetCreateTimestampOk returns a tuple with the CreateTimestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreateTimestamp

`func (o *ApiProject) SetCreateTimestamp(v string)`

SetCreateTimestamp sets CreateTimestamp field to given value.

### HasCreateTimestamp

`func (o *ApiProject) HasCreateTimestamp() bool`

HasCreateTimestamp returns a boolean if a field has been set.

### GetId

`func (o *ApiProject) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ApiProject) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ApiProject) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ApiProject) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *ApiProject) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ApiProject) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ApiProject) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ApiProject) HasName() bool`

HasName returns a boolean if a field has been set.

### GetOrgId

`func (o *ApiProject) GetOrgId() string`

GetOrgId returns the OrgId field if non-nil, zero value otherwise.

### GetOrgIdOk

`func (o *ApiProject) GetOrgIdOk() (*string, bool)`

GetOrgIdOk returns a tuple with the OrgId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgId

`func (o *ApiProject) SetOrgId(v string)`

SetOrgId sets OrgId field to given value.

### HasOrgId

`func (o *ApiProject) HasOrgId() bool`

HasOrgId returns a boolean if a field has been set.

### GetUserCount

`func (o *ApiProject) GetUserCount() int32`

GetUserCount returns the UserCount field if non-nil, zero value otherwise.

### GetUserCountOk

`func (o *ApiProject) GetUserCountOk() (*int32, bool)`

GetUserCountOk returns a tuple with the UserCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserCount

`func (o *ApiProject) SetUserCount(v int32)`

SetUserCount sets UserCount field to given value.

### HasUserCount

`func (o *ApiProject) HasUserCount() bool`

HasUserCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


