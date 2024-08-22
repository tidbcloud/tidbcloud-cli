# ExportOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FileType** | Pointer to [**ExportFileTypeEnum**](ExportFileTypeEnum.md) | Optional. The exported file type. Default: CSV with sql filter, SQL with other filter. | [optional] 
**Database** | Pointer to **string** |  | [optional] 
**Table** | Pointer to **string** |  | [optional] 
**Compression** | Pointer to [**ExportCompressionTypeEnum**](ExportCompressionTypeEnum.md) | Optional. The compression of the export. Default is GZIP. | [optional] 
**Filter** | Pointer to [**ExportOptionsFilter**](ExportOptionsFilter.md) | Optional. The filter of the export. Default is whole cluster. | [optional] 
**CsvFormat** | Pointer to [**ExportOptionsCSVFormat**](ExportOptionsCSVFormat.md) | Optional. The format of the csv. | [optional] 
**ParquetOptions** | Pointer to [**ExportOptionsParquetOptions**](ExportOptionsParquetOptions.md) | Optional. The options of the parquet. | [optional] 

## Methods

### NewExportOptions

`func NewExportOptions() *ExportOptions`

NewExportOptions instantiates a new ExportOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportOptionsWithDefaults

`func NewExportOptionsWithDefaults() *ExportOptions`

NewExportOptionsWithDefaults instantiates a new ExportOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFileType

`func (o *ExportOptions) GetFileType() ExportFileTypeEnum`

GetFileType returns the FileType field if non-nil, zero value otherwise.

### GetFileTypeOk

`func (o *ExportOptions) GetFileTypeOk() (*ExportFileTypeEnum, bool)`

GetFileTypeOk returns a tuple with the FileType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileType

`func (o *ExportOptions) SetFileType(v ExportFileTypeEnum)`

SetFileType sets FileType field to given value.

### HasFileType

`func (o *ExportOptions) HasFileType() bool`

HasFileType returns a boolean if a field has been set.

### GetDatabase

`func (o *ExportOptions) GetDatabase() string`

GetDatabase returns the Database field if non-nil, zero value otherwise.

### GetDatabaseOk

`func (o *ExportOptions) GetDatabaseOk() (*string, bool)`

GetDatabaseOk returns a tuple with the Database field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatabase

`func (o *ExportOptions) SetDatabase(v string)`

SetDatabase sets Database field to given value.

### HasDatabase

`func (o *ExportOptions) HasDatabase() bool`

HasDatabase returns a boolean if a field has been set.

### GetTable

`func (o *ExportOptions) GetTable() string`

GetTable returns the Table field if non-nil, zero value otherwise.

### GetTableOk

`func (o *ExportOptions) GetTableOk() (*string, bool)`

GetTableOk returns a tuple with the Table field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTable

`func (o *ExportOptions) SetTable(v string)`

SetTable sets Table field to given value.

### HasTable

`func (o *ExportOptions) HasTable() bool`

HasTable returns a boolean if a field has been set.

### GetCompression

`func (o *ExportOptions) GetCompression() ExportCompressionTypeEnum`

GetCompression returns the Compression field if non-nil, zero value otherwise.

### GetCompressionOk

`func (o *ExportOptions) GetCompressionOk() (*ExportCompressionTypeEnum, bool)`

GetCompressionOk returns a tuple with the Compression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompression

`func (o *ExportOptions) SetCompression(v ExportCompressionTypeEnum)`

SetCompression sets Compression field to given value.

### HasCompression

`func (o *ExportOptions) HasCompression() bool`

HasCompression returns a boolean if a field has been set.

### GetFilter

`func (o *ExportOptions) GetFilter() ExportOptionsFilter`

GetFilter returns the Filter field if non-nil, zero value otherwise.

### GetFilterOk

`func (o *ExportOptions) GetFilterOk() (*ExportOptionsFilter, bool)`

GetFilterOk returns a tuple with the Filter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilter

`func (o *ExportOptions) SetFilter(v ExportOptionsFilter)`

SetFilter sets Filter field to given value.

### HasFilter

`func (o *ExportOptions) HasFilter() bool`

HasFilter returns a boolean if a field has been set.

### GetCsvFormat

`func (o *ExportOptions) GetCsvFormat() ExportOptionsCSVFormat`

GetCsvFormat returns the CsvFormat field if non-nil, zero value otherwise.

### GetCsvFormatOk

`func (o *ExportOptions) GetCsvFormatOk() (*ExportOptionsCSVFormat, bool)`

GetCsvFormatOk returns a tuple with the CsvFormat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCsvFormat

`func (o *ExportOptions) SetCsvFormat(v ExportOptionsCSVFormat)`

SetCsvFormat sets CsvFormat field to given value.

### HasCsvFormat

`func (o *ExportOptions) HasCsvFormat() bool`

HasCsvFormat returns a boolean if a field has been set.

### GetParquetOptions

`func (o *ExportOptions) GetParquetOptions() ExportOptionsParquetOptions`

GetParquetOptions returns the ParquetOptions field if non-nil, zero value otherwise.

### GetParquetOptionsOk

`func (o *ExportOptions) GetParquetOptionsOk() (*ExportOptionsParquetOptions, bool)`

GetParquetOptionsOk returns a tuple with the ParquetOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParquetOptions

`func (o *ExportOptions) SetParquetOptions(v ExportOptionsParquetOptions)`

SetParquetOptions sets ParquetOptions field to given value.

### HasParquetOptions

`func (o *ExportOptions) HasParquetOptions() bool`

HasParquetOptions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


