# ExportServiceCreateExportBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExportOptions** | Pointer to [**ExportOptions**](ExportOptions.md) | Optional. The options of the export. | [optional] 
**Target** | Pointer to [**ExportTarget**](ExportTarget.md) | Optional. The target of the export. | [optional] 
**DisplayName** | Pointer to **string** | Optional. The display name of the export. Default: SNAPSHOT_{snapshot_time}. | [optional] 

## Methods

### NewExportServiceCreateExportBody

`func NewExportServiceCreateExportBody() *ExportServiceCreateExportBody`

NewExportServiceCreateExportBody instantiates a new ExportServiceCreateExportBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportServiceCreateExportBodyWithDefaults

`func NewExportServiceCreateExportBodyWithDefaults() *ExportServiceCreateExportBody`

NewExportServiceCreateExportBodyWithDefaults instantiates a new ExportServiceCreateExportBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExportOptions

`func (o *ExportServiceCreateExportBody) GetExportOptions() ExportOptions`

GetExportOptions returns the ExportOptions field if non-nil, zero value otherwise.

### GetExportOptionsOk

`func (o *ExportServiceCreateExportBody) GetExportOptionsOk() (*ExportOptions, bool)`

GetExportOptionsOk returns a tuple with the ExportOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportOptions

`func (o *ExportServiceCreateExportBody) SetExportOptions(v ExportOptions)`

SetExportOptions sets ExportOptions field to given value.

### HasExportOptions

`func (o *ExportServiceCreateExportBody) HasExportOptions() bool`

HasExportOptions returns a boolean if a field has been set.

### GetTarget

`func (o *ExportServiceCreateExportBody) GetTarget() ExportTarget`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *ExportServiceCreateExportBody) GetTargetOk() (*ExportTarget, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *ExportServiceCreateExportBody) SetTarget(v ExportTarget)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *ExportServiceCreateExportBody) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### GetDisplayName

`func (o *ExportServiceCreateExportBody) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *ExportServiceCreateExportBody) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *ExportServiceCreateExportBody) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *ExportServiceCreateExportBody) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


