# Export

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExportId** | Pointer to **string** | Output_only. The unique ID of the export. | [optional] [readonly] 
**Name** | Pointer to **string** | Output_only. The unique name of the export. | [optional] [readonly] 
**ClusterId** | **string** | Required. The cluster ID that export belong to. | 
**CreatedBy** | Pointer to **string** | Output_only. The creator of the export. | [optional] [readonly] 
**State** | Pointer to [**ExportStateEnum**](ExportStateEnum.md) | Output_only. The state of the export. | [optional] 
**ExportOptions** | Pointer to [**ExportOptions**](ExportOptions.md) | Optional. The options of the export. | [optional] 
**Target** | Pointer to [**ExportTarget**](ExportTarget.md) | Optional. The target of the export. | [optional] 
**Reason** | Pointer to **NullableString** | Optional. The failed reason of the export. | [optional] [readonly] 
**DisplayName** | Pointer to **string** | Optional. The display name of the export. Default: SNAPSHOT_{snapshot_time}. | [optional] 
**CreateTime** | Pointer to **time.Time** | Output_only. Timestamp when the export was created. | [optional] [readonly] 
**UpdateTime** | Pointer to **NullableTime** | Output_only. Timestamp when the export was updated. | [optional] [readonly] 
**CompleteTime** | Pointer to **NullableTime** | Output_only. Timestamp when the export was completed. | [optional] [readonly] 
**SnapshotTime** | Pointer to **NullableTime** | Output_only. Snapshot time of the export. | [optional] [readonly] 
**ExpireTime** | Pointer to **NullableTime** | Output_only. Expire time of the export. | [optional] [readonly] 

## Methods

### NewExport

`func NewExport(clusterId string, ) *Export`

NewExport instantiates a new Export object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportWithDefaults

`func NewExportWithDefaults() *Export`

NewExportWithDefaults instantiates a new Export object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExportId

`func (o *Export) GetExportId() string`

GetExportId returns the ExportId field if non-nil, zero value otherwise.

### GetExportIdOk

`func (o *Export) GetExportIdOk() (*string, bool)`

GetExportIdOk returns a tuple with the ExportId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportId

`func (o *Export) SetExportId(v string)`

SetExportId sets ExportId field to given value.

### HasExportId

`func (o *Export) HasExportId() bool`

HasExportId returns a boolean if a field has been set.

### GetName

`func (o *Export) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Export) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Export) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Export) HasName() bool`

HasName returns a boolean if a field has been set.

### GetClusterId

`func (o *Export) GetClusterId() string`

GetClusterId returns the ClusterId field if non-nil, zero value otherwise.

### GetClusterIdOk

`func (o *Export) GetClusterIdOk() (*string, bool)`

GetClusterIdOk returns a tuple with the ClusterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterId

`func (o *Export) SetClusterId(v string)`

SetClusterId sets ClusterId field to given value.


### GetCreatedBy

`func (o *Export) GetCreatedBy() string`

GetCreatedBy returns the CreatedBy field if non-nil, zero value otherwise.

### GetCreatedByOk

`func (o *Export) GetCreatedByOk() (*string, bool)`

GetCreatedByOk returns a tuple with the CreatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedBy

`func (o *Export) SetCreatedBy(v string)`

SetCreatedBy sets CreatedBy field to given value.

### HasCreatedBy

`func (o *Export) HasCreatedBy() bool`

HasCreatedBy returns a boolean if a field has been set.

### GetState

`func (o *Export) GetState() ExportStateEnum`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Export) GetStateOk() (*ExportStateEnum, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Export) SetState(v ExportStateEnum)`

SetState sets State field to given value.

### HasState

`func (o *Export) HasState() bool`

HasState returns a boolean if a field has been set.

### GetExportOptions

`func (o *Export) GetExportOptions() ExportOptions`

GetExportOptions returns the ExportOptions field if non-nil, zero value otherwise.

### GetExportOptionsOk

`func (o *Export) GetExportOptionsOk() (*ExportOptions, bool)`

GetExportOptionsOk returns a tuple with the ExportOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportOptions

`func (o *Export) SetExportOptions(v ExportOptions)`

SetExportOptions sets ExportOptions field to given value.

### HasExportOptions

`func (o *Export) HasExportOptions() bool`

HasExportOptions returns a boolean if a field has been set.

### GetTarget

`func (o *Export) GetTarget() ExportTarget`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *Export) GetTargetOk() (*ExportTarget, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *Export) SetTarget(v ExportTarget)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *Export) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### GetReason

`func (o *Export) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *Export) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *Export) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *Export) HasReason() bool`

HasReason returns a boolean if a field has been set.

### SetReasonNil

`func (o *Export) SetReasonNil(b bool)`

 SetReasonNil sets the value for Reason to be an explicit nil

### UnsetReason
`func (o *Export) UnsetReason()`

UnsetReason ensures that no value is present for Reason, not even an explicit nil
### GetDisplayName

`func (o *Export) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *Export) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *Export) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *Export) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetCreateTime

`func (o *Export) GetCreateTime() time.Time`

GetCreateTime returns the CreateTime field if non-nil, zero value otherwise.

### GetCreateTimeOk

`func (o *Export) GetCreateTimeOk() (*time.Time, bool)`

GetCreateTimeOk returns a tuple with the CreateTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreateTime

`func (o *Export) SetCreateTime(v time.Time)`

SetCreateTime sets CreateTime field to given value.

### HasCreateTime

`func (o *Export) HasCreateTime() bool`

HasCreateTime returns a boolean if a field has been set.

### GetUpdateTime

`func (o *Export) GetUpdateTime() time.Time`

GetUpdateTime returns the UpdateTime field if non-nil, zero value otherwise.

### GetUpdateTimeOk

`func (o *Export) GetUpdateTimeOk() (*time.Time, bool)`

GetUpdateTimeOk returns a tuple with the UpdateTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdateTime

`func (o *Export) SetUpdateTime(v time.Time)`

SetUpdateTime sets UpdateTime field to given value.

### HasUpdateTime

`func (o *Export) HasUpdateTime() bool`

HasUpdateTime returns a boolean if a field has been set.

### SetUpdateTimeNil

`func (o *Export) SetUpdateTimeNil(b bool)`

 SetUpdateTimeNil sets the value for UpdateTime to be an explicit nil

### UnsetUpdateTime
`func (o *Export) UnsetUpdateTime()`

UnsetUpdateTime ensures that no value is present for UpdateTime, not even an explicit nil
### GetCompleteTime

`func (o *Export) GetCompleteTime() time.Time`

GetCompleteTime returns the CompleteTime field if non-nil, zero value otherwise.

### GetCompleteTimeOk

`func (o *Export) GetCompleteTimeOk() (*time.Time, bool)`

GetCompleteTimeOk returns a tuple with the CompleteTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleteTime

`func (o *Export) SetCompleteTime(v time.Time)`

SetCompleteTime sets CompleteTime field to given value.

### HasCompleteTime

`func (o *Export) HasCompleteTime() bool`

HasCompleteTime returns a boolean if a field has been set.

### SetCompleteTimeNil

`func (o *Export) SetCompleteTimeNil(b bool)`

 SetCompleteTimeNil sets the value for CompleteTime to be an explicit nil

### UnsetCompleteTime
`func (o *Export) UnsetCompleteTime()`

UnsetCompleteTime ensures that no value is present for CompleteTime, not even an explicit nil
### GetSnapshotTime

`func (o *Export) GetSnapshotTime() time.Time`

GetSnapshotTime returns the SnapshotTime field if non-nil, zero value otherwise.

### GetSnapshotTimeOk

`func (o *Export) GetSnapshotTimeOk() (*time.Time, bool)`

GetSnapshotTimeOk returns a tuple with the SnapshotTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnapshotTime

`func (o *Export) SetSnapshotTime(v time.Time)`

SetSnapshotTime sets SnapshotTime field to given value.

### HasSnapshotTime

`func (o *Export) HasSnapshotTime() bool`

HasSnapshotTime returns a boolean if a field has been set.

### SetSnapshotTimeNil

`func (o *Export) SetSnapshotTimeNil(b bool)`

 SetSnapshotTimeNil sets the value for SnapshotTime to be an explicit nil

### UnsetSnapshotTime
`func (o *Export) UnsetSnapshotTime()`

UnsetSnapshotTime ensures that no value is present for SnapshotTime, not even an explicit nil
### GetExpireTime

`func (o *Export) GetExpireTime() time.Time`

GetExpireTime returns the ExpireTime field if non-nil, zero value otherwise.

### GetExpireTimeOk

`func (o *Export) GetExpireTimeOk() (*time.Time, bool)`

GetExpireTimeOk returns a tuple with the ExpireTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpireTime

`func (o *Export) SetExpireTime(v time.Time)`

SetExpireTime sets ExpireTime field to given value.

### HasExpireTime

`func (o *Export) HasExpireTime() bool`

HasExpireTime returns a boolean if a field has been set.

### SetExpireTimeNil

`func (o *Export) SetExpireTimeNil(b bool)`

 SetExpireTimeNil sets the value for ExpireTime to be an explicit nil

### UnsetExpireTime
`func (o *Export) UnsetExpireTime()`

UnsetExpireTime ensures that no value is present for ExpireTime, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


