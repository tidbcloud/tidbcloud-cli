# ExportOptionsParquetOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Compression** | Pointer to [**ExportParquetCompressionTypeEnum**](ExportParquetCompressionTypeEnum.md) | Optional. The compression of the parquet. Default is ZSTD. | [optional] 

## Methods

### NewExportOptionsParquetOptions

`func NewExportOptionsParquetOptions() *ExportOptionsParquetOptions`

NewExportOptionsParquetOptions instantiates a new ExportOptionsParquetOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportOptionsParquetOptionsWithDefaults

`func NewExportOptionsParquetOptionsWithDefaults() *ExportOptionsParquetOptions`

NewExportOptionsParquetOptionsWithDefaults instantiates a new ExportOptionsParquetOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompression

`func (o *ExportOptionsParquetOptions) GetCompression() ExportParquetCompressionTypeEnum`

GetCompression returns the Compression field if non-nil, zero value otherwise.

### GetCompressionOk

`func (o *ExportOptionsParquetOptions) GetCompressionOk() (*ExportParquetCompressionTypeEnum, bool)`

GetCompressionOk returns a tuple with the Compression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompression

`func (o *ExportOptionsParquetOptions) SetCompression(v ExportParquetCompressionTypeEnum)`

SetCompression sets Compression field to given value.

### HasCompression

`func (o *ExportOptionsParquetOptions) HasCompression() bool`

HasCompression returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


