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

// ImportServiceCompleteMultipartUploadReader is a Reader for the ImportServiceCompleteMultipartUpload structure.
type ImportServiceCompleteMultipartUploadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ImportServiceCompleteMultipartUploadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewImportServiceCompleteMultipartUploadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewImportServiceCompleteMultipartUploadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewImportServiceCompleteMultipartUploadOK creates a ImportServiceCompleteMultipartUploadOK with default headers values
func NewImportServiceCompleteMultipartUploadOK() *ImportServiceCompleteMultipartUploadOK {
	return &ImportServiceCompleteMultipartUploadOK{}
}

/*
ImportServiceCompleteMultipartUploadOK describes a response with status code 200, with default header values.

A successful response.
*/
type ImportServiceCompleteMultipartUploadOK struct {
	Payload interface{}
}

// IsSuccess returns true when this import service complete multipart upload o k response has a 2xx status code
func (o *ImportServiceCompleteMultipartUploadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this import service complete multipart upload o k response has a 3xx status code
func (o *ImportServiceCompleteMultipartUploadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this import service complete multipart upload o k response has a 4xx status code
func (o *ImportServiceCompleteMultipartUploadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this import service complete multipart upload o k response has a 5xx status code
func (o *ImportServiceCompleteMultipartUploadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this import service complete multipart upload o k response a status code equal to that given
func (o *ImportServiceCompleteMultipartUploadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the import service complete multipart upload o k response
func (o *ImportServiceCompleteMultipartUploadOK) Code() int {
	return 200
}

func (o *ImportServiceCompleteMultipartUploadOK) Error() string {
	return fmt.Sprintf("[POST /v1beta1/clusters/{clusterId}/imports:completeUpload][%d] importServiceCompleteMultipartUploadOK  %+v", 200, o.Payload)
}

func (o *ImportServiceCompleteMultipartUploadOK) String() string {
	return fmt.Sprintf("[POST /v1beta1/clusters/{clusterId}/imports:completeUpload][%d] importServiceCompleteMultipartUploadOK  %+v", 200, o.Payload)
}

func (o *ImportServiceCompleteMultipartUploadOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ImportServiceCompleteMultipartUploadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewImportServiceCompleteMultipartUploadDefault creates a ImportServiceCompleteMultipartUploadDefault with default headers values
func NewImportServiceCompleteMultipartUploadDefault(code int) *ImportServiceCompleteMultipartUploadDefault {
	return &ImportServiceCompleteMultipartUploadDefault{
		_statusCode: code,
	}
}

/*
ImportServiceCompleteMultipartUploadDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ImportServiceCompleteMultipartUploadDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this import service complete multipart upload default response has a 2xx status code
func (o *ImportServiceCompleteMultipartUploadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this import service complete multipart upload default response has a 3xx status code
func (o *ImportServiceCompleteMultipartUploadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this import service complete multipart upload default response has a 4xx status code
func (o *ImportServiceCompleteMultipartUploadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this import service complete multipart upload default response has a 5xx status code
func (o *ImportServiceCompleteMultipartUploadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this import service complete multipart upload default response a status code equal to that given
func (o *ImportServiceCompleteMultipartUploadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the import service complete multipart upload default response
func (o *ImportServiceCompleteMultipartUploadDefault) Code() int {
	return o._statusCode
}

func (o *ImportServiceCompleteMultipartUploadDefault) Error() string {
	return fmt.Sprintf("[POST /v1beta1/clusters/{clusterId}/imports:completeUpload][%d] ImportService_CompleteMultipartUpload default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceCompleteMultipartUploadDefault) String() string {
	return fmt.Sprintf("[POST /v1beta1/clusters/{clusterId}/imports:completeUpload][%d] ImportService_CompleteMultipartUpload default  %+v", o._statusCode, o.Payload)
}

func (o *ImportServiceCompleteMultipartUploadDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *ImportServiceCompleteMultipartUploadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
