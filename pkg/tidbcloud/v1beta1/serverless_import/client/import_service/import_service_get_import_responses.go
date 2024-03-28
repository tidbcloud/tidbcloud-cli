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

// ImportServiceGetImportReader is a Reader for the ImportServiceGetImport structure.
type ImportServiceGetImportReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ImportServiceGetImportReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewImportServiceGetImportOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewImportServiceGetImportDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewImportServiceGetImportOK creates a ImportServiceGetImportOK with default headers values
func NewImportServiceGetImportOK() *ImportServiceGetImportOK {
	return &ImportServiceGetImportOK{}
}

/*
ImportServiceGetImportOK describes a response with status code 200, with default header values.

A successful response.
*/
type ImportServiceGetImportOK struct {
	Payload *models.V1beta1Import
}

// IsSuccess returns true when this import service get import o k response has a 2xx status code
func (o *ImportServiceGetImportOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this import service get import o k response has a 3xx status code
func (o *ImportServiceGetImportOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this import service get import o k response has a 4xx status code
func (o *ImportServiceGetImportOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this import service get import o k response has a 5xx status code
func (o *ImportServiceGetImportOK) IsServerError() bool {
	return false
}

// IsCode returns true when this import service get import o k response a status code equal to that given
func (o *ImportServiceGetImportOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the import service get import o k response
func (o *ImportServiceGetImportOK) Code() int {
	return 200
}

func (o *ImportServiceGetImportOK) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports/{id}][%d] importServiceGetImportOK  %+v", 200, o.Payload)
}

func (o *ImportServiceGetImportOK) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports/{id}][%d] importServiceGetImportOK  %+v", 200, o.Payload)
}

func (o *ImportServiceGetImportOK) GetPayload() *models.V1beta1Import {
	return o.Payload
}

func (o *ImportServiceGetImportOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1Import)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewImportServiceGetImportDefault creates a ImportServiceGetImportDefault with default headers values
func NewImportServiceGetImportDefault(code int) *ImportServiceGetImportDefault {
	return &ImportServiceGetImportDefault{
		_statusCode: code,
	}
}

/*
ImportServiceGetImportDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ImportServiceGetImportDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this import service get import default response has a 2xx status code
func (o *ImportServiceGetImportDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this import service get import default response has a 3xx status code
func (o *ImportServiceGetImportDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this import service get import default response has a 4xx status code
func (o *ImportServiceGetImportDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this import service get import default response has a 5xx status code
func (o *ImportServiceGetImportDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this import service get import default response a status code equal to that given
func (o *ImportServiceGetImportDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the import service get import default response
func (o *ImportServiceGetImportDefault) Code() int {
	return o._statusCode
}

func (o *ImportServiceGetImportDefault) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports/{id}][%d] ImportService_GetImport default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceGetImportDefault) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports/{id}][%d] ImportService_GetImport default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceGetImportDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *ImportServiceGetImportDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
