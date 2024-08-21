# ExportOptionsFilter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Sql** | Pointer to **string** | Optional. Use SQL to filter the export. | [optional] 
**Table** | Pointer to [**ExportOptionsFilterTable**](ExportOptionsFilterTable.md) | Optional. Use table-filter to filter the export. | [optional] 

## Methods

### NewExportOptionsFilter

`func NewExportOptionsFilter() *ExportOptionsFilter`

NewExportOptionsFilter instantiates a new ExportOptionsFilter object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportOptionsFilterWithDefaults

`func NewExportOptionsFilterWithDefaults() *ExportOptionsFilter`

NewExportOptionsFilterWithDefaults instantiates a new ExportOptionsFilter object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSql

`func (o *ExportOptionsFilter) GetSql() string`

GetSql returns the Sql field if non-nil, zero value otherwise.

### GetSqlOk

`func (o *ExportOptionsFilter) GetSqlOk() (*string, bool)`

GetSqlOk returns a tuple with the Sql field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSql

`func (o *ExportOptionsFilter) SetSql(v string)`

SetSql sets Sql field to given value.

### HasSql

`func (o *ExportOptionsFilter) HasSql() bool`

HasSql returns a boolean if a field has been set.

### GetTable

`func (o *ExportOptionsFilter) GetTable() ExportOptionsFilterTable`

GetTable returns the Table field if non-nil, zero value otherwise.

### GetTableOk

`func (o *ExportOptionsFilter) GetTableOk() (*ExportOptionsFilterTable, bool)`

GetTableOk returns a tuple with the Table field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTable

`func (o *ExportOptionsFilter) SetTable(v ExportOptionsFilterTable)`

SetTable sets Table field to given value.

### HasTable

`func (o *ExportOptionsFilter) HasTable() bool`

HasTable returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


