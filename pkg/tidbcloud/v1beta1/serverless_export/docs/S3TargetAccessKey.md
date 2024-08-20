# S3TargetAccessKey

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The access key id of the s3. | 
**Secret** | **string** | Input_Only. The secret access key of the s3. | 

## Methods

### NewS3TargetAccessKey

`func NewS3TargetAccessKey(id string, secret string, ) *S3TargetAccessKey`

NewS3TargetAccessKey instantiates a new S3TargetAccessKey object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewS3TargetAccessKeyWithDefaults

`func NewS3TargetAccessKeyWithDefaults() *S3TargetAccessKey`

NewS3TargetAccessKeyWithDefaults instantiates a new S3TargetAccessKey object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *S3TargetAccessKey) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *S3TargetAccessKey) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *S3TargetAccessKey) SetId(v string)`

SetId sets Id field to given value.


### GetSecret

`func (o *S3TargetAccessKey) GetSecret() string`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *S3TargetAccessKey) GetSecretOk() (*string, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *S3TargetAccessKey) SetSecret(v string)`

SetSecret sets Secret field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


