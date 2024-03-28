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

// ImportServiceListImportsReader is a Reader for the ImportServiceListImports structure.
type ImportServiceListImportsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ImportServiceListImportsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewImportServiceListImportsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewImportServiceListImportsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewImportServiceListImportsOK creates a ImportServiceListImportsOK with default headers values
func NewImportServiceListImportsOK() *ImportServiceListImportsOK {
	return &ImportServiceListImportsOK{}
}

/*
ImportServiceListImportsOK describes a response with status code 200, with default header values.

A successful response.
*/
type ImportServiceListImportsOK struct {
	Payload *models.V1beta1ListImportsResp
}

// IsSuccess returns true when this import service list imports o k response has a 2xx status code
func (o *ImportServiceListImportsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this import service list imports o k response has a 3xx status code
func (o *ImportServiceListImportsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this import service list imports o k response has a 4xx status code
func (o *ImportServiceListImportsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this import service list imports o k response has a 5xx status code
func (o *ImportServiceListImportsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this import service list imports o k response a status code equal to that given
func (o *ImportServiceListImportsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the import service list imports o k response
func (o *ImportServiceListImportsOK) Code() int {
	return 200
}

func (o *ImportServiceListImportsOK) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports][%d] importServiceListImportsOK  %+v", 200, o.Payload)
}

func (o *ImportServiceListImportsOK) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports][%d] importServiceListImportsOK  %+v", 200, o.Payload)
}

func (o *ImportServiceListImportsOK) GetPayload() *models.V1beta1ListImportsResp {
	return o.Payload
}

func (o *ImportServiceListImportsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1ListImportsResp)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewImportServiceListImportsDefault creates a ImportServiceListImportsDefault with default headers values
func NewImportServiceListImportsDefault(code int) *ImportServiceListImportsDefault {
	return &ImportServiceListImportsDefault{
		_statusCode: code,
	}
}

/*
ImportServiceListImportsDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ImportServiceListImportsDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this import service list imports default response has a 2xx status code
func (o *ImportServiceListImportsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this import service list imports default response has a 3xx status code
func (o *ImportServiceListImportsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this import service list imports default response has a 4xx status code
func (o *ImportServiceListImportsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this import service list imports default response has a 5xx status code
func (o *ImportServiceListImportsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this import service list imports default response a status code equal to that given
func (o *ImportServiceListImportsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the import service list imports default response
func (o *ImportServiceListImportsDefault) Code() int {
	return o._statusCode
}

func (o *ImportServiceListImportsDefault) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports][%d] ImportService_ListImports default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceListImportsDefault) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/imports][%d] ImportService_ListImports default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceListImportsDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *ImportServiceListImportsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
