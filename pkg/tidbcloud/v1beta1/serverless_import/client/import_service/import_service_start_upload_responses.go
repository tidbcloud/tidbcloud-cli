// Code generated by go-swagger; DO NOT EDIT.

package import_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"
)

// ImportServiceStartUploadReader is a Reader for the ImportServiceStartUpload structure.
type ImportServiceStartUploadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ImportServiceStartUploadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewImportServiceStartUploadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewImportServiceStartUploadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewImportServiceStartUploadOK creates a ImportServiceStartUploadOK with default headers values
func NewImportServiceStartUploadOK() *ImportServiceStartUploadOK {
	return &ImportServiceStartUploadOK{}
}

/*
ImportServiceStartUploadOK describes a response with status code 200, with default header values.

A successful response.
*/
type ImportServiceStartUploadOK struct {
	Payload *models.V1beta1StartUploadResponse
}

// IsSuccess returns true when this import service start upload o k response has a 2xx status code
func (o *ImportServiceStartUploadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this import service start upload o k response has a 3xx status code
func (o *ImportServiceStartUploadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this import service start upload o k response has a 4xx status code
func (o *ImportServiceStartUploadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this import service start upload o k response has a 5xx status code
func (o *ImportServiceStartUploadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this import service start upload o k response a status code equal to that given
func (o *ImportServiceStartUploadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the import service start upload o k response
func (o *ImportServiceStartUploadOK) Code() int {
	return 200
}

func (o *ImportServiceStartUploadOK) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports:start-upload][%d] importServiceStartUploadOK  %+v", 200, o.Payload)
}

func (o *ImportServiceStartUploadOK) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports:start-upload][%d] importServiceStartUploadOK  %+v", 200, o.Payload)
}

func (o *ImportServiceStartUploadOK) GetPayload() *models.V1beta1StartUploadResponse {
	return o.Payload
}

func (o *ImportServiceStartUploadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1StartUploadResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewImportServiceStartUploadDefault creates a ImportServiceStartUploadDefault with default headers values
func NewImportServiceStartUploadDefault(code int) *ImportServiceStartUploadDefault {
	return &ImportServiceStartUploadDefault{
		_statusCode: code,
	}
}

/*
ImportServiceStartUploadDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ImportServiceStartUploadDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this import service start upload default response has a 2xx status code
func (o *ImportServiceStartUploadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this import service start upload default response has a 3xx status code
func (o *ImportServiceStartUploadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this import service start upload default response has a 4xx status code
func (o *ImportServiceStartUploadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this import service start upload default response has a 5xx status code
func (o *ImportServiceStartUploadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this import service start upload default response a status code equal to that given
func (o *ImportServiceStartUploadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the import service start upload default response
func (o *ImportServiceStartUploadDefault) Code() int {
	return o._statusCode
}

func (o *ImportServiceStartUploadDefault) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports:start-upload][%d] ImportService_StartUpload default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceStartUploadDefault) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports:start-upload][%d] ImportService_StartUpload default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceStartUploadDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *ImportServiceStartUploadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
