// Code generated by go-swagger; DO NOT EDIT.

package serverless_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/serverless/models"
)

// ServerlessServiceDeleteClusterReader is a Reader for the ServerlessServiceDeleteCluster structure.
type ServerlessServiceDeleteClusterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ServerlessServiceDeleteClusterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewServerlessServiceDeleteClusterOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewServerlessServiceDeleteClusterDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewServerlessServiceDeleteClusterOK creates a ServerlessServiceDeleteClusterOK with default headers values
func NewServerlessServiceDeleteClusterOK() *ServerlessServiceDeleteClusterOK {
	return &ServerlessServiceDeleteClusterOK{}
}

/*
ServerlessServiceDeleteClusterOK describes a response with status code 200, with default header values.

A successful response.
*/
type ServerlessServiceDeleteClusterOK struct {
	Payload *models.V1Cluster
}

// IsSuccess returns true when this serverless service delete cluster o k response has a 2xx status code
func (o *ServerlessServiceDeleteClusterOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this serverless service delete cluster o k response has a 3xx status code
func (o *ServerlessServiceDeleteClusterOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this serverless service delete cluster o k response has a 4xx status code
func (o *ServerlessServiceDeleteClusterOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this serverless service delete cluster o k response has a 5xx status code
func (o *ServerlessServiceDeleteClusterOK) IsServerError() bool {
	return false
}

// IsCode returns true when this serverless service delete cluster o k response a status code equal to that given
func (o *ServerlessServiceDeleteClusterOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the serverless service delete cluster o k response
func (o *ServerlessServiceDeleteClusterOK) Code() int {
	return 200
}

func (o *ServerlessServiceDeleteClusterOK) Error() string {
	return fmt.Sprintf("[DELETE /v1/clusters/{clusterId}][%d] serverlessServiceDeleteClusterOK  %+v", 200, o.Payload)
}

func (o *ServerlessServiceDeleteClusterOK) String() string {
	return fmt.Sprintf("[DELETE /v1/clusters/{clusterId}][%d] serverlessServiceDeleteClusterOK  %+v", 200, o.Payload)
}

func (o *ServerlessServiceDeleteClusterOK) GetPayload() *models.V1Cluster {
	return o.Payload
}

func (o *ServerlessServiceDeleteClusterOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1Cluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewServerlessServiceDeleteClusterDefault creates a ServerlessServiceDeleteClusterDefault with default headers values
func NewServerlessServiceDeleteClusterDefault(code int) *ServerlessServiceDeleteClusterDefault {
	return &ServerlessServiceDeleteClusterDefault{
		_statusCode: code,
	}
}

/*
ServerlessServiceDeleteClusterDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ServerlessServiceDeleteClusterDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this serverless service delete cluster default response has a 2xx status code
func (o *ServerlessServiceDeleteClusterDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this serverless service delete cluster default response has a 3xx status code
func (o *ServerlessServiceDeleteClusterDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this serverless service delete cluster default response has a 4xx status code
func (o *ServerlessServiceDeleteClusterDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this serverless service delete cluster default response has a 5xx status code
func (o *ServerlessServiceDeleteClusterDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this serverless service delete cluster default response a status code equal to that given
func (o *ServerlessServiceDeleteClusterDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the serverless service delete cluster default response
func (o *ServerlessServiceDeleteClusterDefault) Code() int {
	return o._statusCode
}

func (o *ServerlessServiceDeleteClusterDefault) Error() string {
	return fmt.Sprintf("[DELETE /v1/clusters/{clusterId}][%d] ServerlessService_DeleteCluster default  %+v", o._statusCode, o.Payload)
}

func (o *ServerlessServiceDeleteClusterDefault) String() string {
	return fmt.Sprintf("[DELETE /v1/clusters/{clusterId}][%d] ServerlessService_DeleteCluster default  %+v", o._statusCode, o.Payload)
}

func (o *ServerlessServiceDeleteClusterDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *ServerlessServiceDeleteClusterDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
