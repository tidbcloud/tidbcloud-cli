# V1beta1Branch

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Output Only. The name of the resource. | [optional] [readonly] 
**BranchId** | Pointer to **string** | Output only. The system-generated ID of the resource. | [optional] [readonly] 
**DisplayName** | **string** | Required. User-settable and human-readable display name for the branch. | 
**ClusterId** | Pointer to **string** | Output only. The cluster ID of this branch. | [optional] [readonly] 
**ParentId** | Pointer to **string** | Optional. The parent ID of this branch. | [optional] 
**CreatedBy** | Pointer to **string** | Output only. The creator of the branch. | [optional] [readonly] 
**State** | Pointer to [**V1beta1BranchState**](V1beta1BranchState.md) | Output only. The state of this branch. | [optional] 
**Endpoints** | Pointer to [**BranchEndpoints**](BranchEndpoints.md) | Optional. The endpoints of this branch. | [optional] 
**UserPrefix** | Pointer to **NullableString** | Output only. User name prefix of this branch. For each TiDB Serverless branch, TiDB Cloud generates a unique prefix to distinguish it from other branches. Whenever you use or set a database user name, you must include the prefix in the user name. | [optional] [readonly] 
**Usage** | Pointer to [**BranchUsage**](BranchUsage.md) | Output only. Usage metrics of this branch. Only display in FULL view. | [optional] 
**CreateTime** | Pointer to **time.Time** |  | [optional] [readonly] 
**UpdateTime** | Pointer to **time.Time** |  | [optional] [readonly] 
**Annotations** | Pointer to **map[string]string** | Optional. The annotations of this branch.. | [optional] 
**ParentDisplayName** | Pointer to **string** | Output only. The parent display name of this branch. | [optional] [readonly] 

## Methods

### NewV1beta1Branch

`func NewV1beta1Branch(displayName string, ) *V1beta1Branch`

NewV1beta1Branch instantiates a new V1beta1Branch object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewV1beta1BranchWithDefaults

`func NewV1beta1BranchWithDefaults() *V1beta1Branch`

NewV1beta1BranchWithDefaults instantiates a new V1beta1Branch object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *V1beta1Branch) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *V1beta1Branch) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *V1beta1Branch) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *V1beta1Branch) HasName() bool`

HasName returns a boolean if a field has been set.

### GetBranchId

`func (o *V1beta1Branch) GetBranchId() string`

GetBranchId returns the BranchId field if non-nil, zero value otherwise.

### GetBranchIdOk

`func (o *V1beta1Branch) GetBranchIdOk() (*string, bool)`

GetBranchIdOk returns a tuple with the BranchId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBranchId

`func (o *V1beta1Branch) SetBranchId(v string)`

SetBranchId sets BranchId field to given value.

### HasBranchId

`func (o *V1beta1Branch) HasBranchId() bool`

HasBranchId returns a boolean if a field has been set.

### GetDisplayName

`func (o *V1beta1Branch) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *V1beta1Branch) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *V1beta1Branch) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetClusterId

`func (o *V1beta1Branch) GetClusterId() string`

GetClusterId returns the ClusterId field if non-nil, zero value otherwise.

### GetClusterIdOk

`func (o *V1beta1Branch) GetClusterIdOk() (*string, bool)`

GetClusterIdOk returns a tuple with the ClusterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterId

`func (o *V1beta1Branch) SetClusterId(v string)`

SetClusterId sets ClusterId field to given value.

### HasClusterId

`func (o *V1beta1Branch) HasClusterId() bool`

HasClusterId returns a boolean if a field has been set.

### GetParentId

`func (o *V1beta1Branch) GetParentId() string`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *V1beta1Branch) GetParentIdOk() (*string, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *V1beta1Branch) SetParentId(v string)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *V1beta1Branch) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetCreatedBy

`func (o *V1beta1Branch) GetCreatedBy() string`

GetCreatedBy returns the CreatedBy field if non-nil, zero value otherwise.

### GetCreatedByOk

`func (o *V1beta1Branch) GetCreatedByOk() (*string, bool)`

GetCreatedByOk returns a tuple with the CreatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedBy

`func (o *V1beta1Branch) SetCreatedBy(v string)`

SetCreatedBy sets CreatedBy field to given value.

### HasCreatedBy

`func (o *V1beta1Branch) HasCreatedBy() bool`

HasCreatedBy returns a boolean if a field has been set.

### GetState

`func (o *V1beta1Branch) GetState() V1beta1BranchState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *V1beta1Branch) GetStateOk() (*V1beta1BranchState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *V1beta1Branch) SetState(v V1beta1BranchState)`

SetState sets State field to given value.

### HasState

`func (o *V1beta1Branch) HasState() bool`

HasState returns a boolean if a field has been set.

### GetEndpoints

`func (o *V1beta1Branch) GetEndpoints() BranchEndpoints`

GetEndpoints returns the Endpoints field if non-nil, zero value otherwise.

### GetEndpointsOk

`func (o *V1beta1Branch) GetEndpointsOk() (*BranchEndpoints, bool)`

GetEndpointsOk returns a tuple with the Endpoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoints

`func (o *V1beta1Branch) SetEndpoints(v BranchEndpoints)`

SetEndpoints sets Endpoints field to given value.

### HasEndpoints

`func (o *V1beta1Branch) HasEndpoints() bool`

HasEndpoints returns a boolean if a field has been set.

### GetUserPrefix

`func (o *V1beta1Branch) GetUserPrefix() string`

GetUserPrefix returns the UserPrefix field if non-nil, zero value otherwise.

### GetUserPrefixOk

`func (o *V1beta1Branch) GetUserPrefixOk() (*string, bool)`

GetUserPrefixOk returns a tuple with the UserPrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserPrefix

`func (o *V1beta1Branch) SetUserPrefix(v string)`

SetUserPrefix sets UserPrefix field to given value.

### HasUserPrefix

`func (o *V1beta1Branch) HasUserPrefix() bool`

HasUserPrefix returns a boolean if a field has been set.

### SetUserPrefixNil

`func (o *V1beta1Branch) SetUserPrefixNil(b bool)`

 SetUserPrefixNil sets the value for UserPrefix to be an explicit nil

### UnsetUserPrefix
`func (o *V1beta1Branch) UnsetUserPrefix()`

UnsetUserPrefix ensures that no value is present for UserPrefix, not even an explicit nil
### GetUsage

`func (o *V1beta1Branch) GetUsage() BranchUsage`

GetUsage returns the Usage field if non-nil, zero value otherwise.

### GetUsageOk

`func (o *V1beta1Branch) GetUsageOk() (*BranchUsage, bool)`

GetUsageOk returns a tuple with the Usage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsage

`func (o *V1beta1Branch) SetUsage(v BranchUsage)`

SetUsage sets Usage field to given value.

### HasUsage

`func (o *V1beta1Branch) HasUsage() bool`

HasUsage returns a boolean if a field has been set.

### GetCreateTime

`func (o *V1beta1Branch) GetCreateTime() time.Time`

GetCreateTime returns the CreateTime field if non-nil, zero value otherwise.

### GetCreateTimeOk

`func (o *V1beta1Branch) GetCreateTimeOk() (*time.Time, bool)`

GetCreateTimeOk returns a tuple with the CreateTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreateTime

`func (o *V1beta1Branch) SetCreateTime(v time.Time)`

SetCreateTime sets CreateTime field to given value.

### HasCreateTime

`func (o *V1beta1Branch) HasCreateTime() bool`

HasCreateTime returns a boolean if a field has been set.

### GetUpdateTime

`func (o *V1beta1Branch) GetUpdateTime() time.Time`

GetUpdateTime returns the UpdateTime field if non-nil, zero value otherwise.

### GetUpdateTimeOk

`func (o *V1beta1Branch) GetUpdateTimeOk() (*time.Time, bool)`

GetUpdateTimeOk returns a tuple with the UpdateTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdateTime

`func (o *V1beta1Branch) SetUpdateTime(v time.Time)`

SetUpdateTime sets UpdateTime field to given value.

### HasUpdateTime

`func (o *V1beta1Branch) HasUpdateTime() bool`

HasUpdateTime returns a boolean if a field has been set.

### GetAnnotations

`func (o *V1beta1Branch) GetAnnotations() map[string]string`

GetAnnotations returns the Annotations field if non-nil, zero value otherwise.

### GetAnnotationsOk

`func (o *V1beta1Branch) GetAnnotationsOk() (*map[string]string, bool)`

GetAnnotationsOk returns a tuple with the Annotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnnotations

`func (o *V1beta1Branch) SetAnnotations(v map[string]string)`

SetAnnotations sets Annotations field to given value.

### HasAnnotations

`func (o *V1beta1Branch) HasAnnotations() bool`

HasAnnotations returns a boolean if a field has been set.

### GetParentDisplayName

`func (o *V1beta1Branch) GetParentDisplayName() string`

GetParentDisplayName returns the ParentDisplayName field if non-nil, zero value otherwise.

### GetParentDisplayNameOk

`func (o *V1beta1Branch) GetParentDisplayNameOk() (*string, bool)`

GetParentDisplayNameOk returns a tuple with the ParentDisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentDisplayName

`func (o *V1beta1Branch) SetParentDisplayName(v string)`

SetParentDisplayName sets ParentDisplayName field to given value.

### HasParentDisplayName

`func (o *V1beta1Branch) HasParentDisplayName() bool`

HasParentDisplayName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


