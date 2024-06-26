// Code generated by go-swagger; DO NOT EDIT.

package serverless_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"
)

// ServerlessServiceGetClusterReader is a Reader for the ServerlessServiceGetCluster structure.
type ServerlessServiceGetClusterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ServerlessServiceGetClusterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewServerlessServiceGetClusterOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewServerlessServiceGetClusterDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewServerlessServiceGetClusterOK creates a ServerlessServiceGetClusterOK with default headers values
func NewServerlessServiceGetClusterOK() *ServerlessServiceGetClusterOK {
	return &ServerlessServiceGetClusterOK{}
}

/*
ServerlessServiceGetClusterOK describes a response with status code 200, with default header values.

A successful response.
*/
type ServerlessServiceGetClusterOK struct {
	Payload *models.TidbCloudOpenApiserverlessv1beta1Cluster
}

// IsSuccess returns true when this serverless service get cluster o k response has a 2xx status code
func (o *ServerlessServiceGetClusterOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this serverless service get cluster o k response has a 3xx status code
func (o *ServerlessServiceGetClusterOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this serverless service get cluster o k response has a 4xx status code
func (o *ServerlessServiceGetClusterOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this serverless service get cluster o k response has a 5xx status code
func (o *ServerlessServiceGetClusterOK) IsServerError() bool {
	return false
}

// IsCode returns true when this serverless service get cluster o k response a status code equal to that given
func (o *ServerlessServiceGetClusterOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the serverless service get cluster o k response
func (o *ServerlessServiceGetClusterOK) Code() int {
	return 200
}

func (o *ServerlessServiceGetClusterOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /clusters/{clusterId}][%d] serverlessServiceGetClusterOK %s", 200, payload)
}

func (o *ServerlessServiceGetClusterOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /clusters/{clusterId}][%d] serverlessServiceGetClusterOK %s", 200, payload)
}

func (o *ServerlessServiceGetClusterOK) GetPayload() *models.TidbCloudOpenApiserverlessv1beta1Cluster {
	return o.Payload
}

func (o *ServerlessServiceGetClusterOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.TidbCloudOpenApiserverlessv1beta1Cluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewServerlessServiceGetClusterDefault creates a ServerlessServiceGetClusterDefault with default headers values
func NewServerlessServiceGetClusterDefault(code int) *ServerlessServiceGetClusterDefault {
	return &ServerlessServiceGetClusterDefault{
		_statusCode: code,
	}
}

/*
ServerlessServiceGetClusterDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ServerlessServiceGetClusterDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this serverless service get cluster default response has a 2xx status code
func (o *ServerlessServiceGetClusterDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this serverless service get cluster default response has a 3xx status code
func (o *ServerlessServiceGetClusterDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this serverless service get cluster default response has a 4xx status code
func (o *ServerlessServiceGetClusterDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this serverless service get cluster default response has a 5xx status code
func (o *ServerlessServiceGetClusterDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this serverless service get cluster default response a status code equal to that given
func (o *ServerlessServiceGetClusterDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the serverless service get cluster default response
func (o *ServerlessServiceGetClusterDefault) Code() int {
	return o._statusCode
}

func (o *ServerlessServiceGetClusterDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /clusters/{clusterId}][%d] ServerlessService_GetCluster default %s", o._statusCode, payload)
}

func (o *ServerlessServiceGetClusterDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /clusters/{clusterId}][%d] ServerlessService_GetCluster default %s", o._statusCode, payload)
}

func (o *ServerlessServiceGetClusterDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *ServerlessServiceGetClusterDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
