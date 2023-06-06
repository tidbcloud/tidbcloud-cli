// Code generated by go-swagger; DO NOT EDIT.

package branch_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/branch/models"
)

// ListBranchesReader is a Reader for the ListBranches structure.
type ListBranchesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListBranchesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListBranchesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListBranchesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListBranchesOK creates a ListBranchesOK with default headers values
func NewListBranchesOK() *ListBranchesOK {
	return &ListBranchesOK{}
}

/*
ListBranchesOK describes a response with status code 200, with default header values.

A successful response.
*/
type ListBranchesOK struct {
	Payload *models.OpenapiListBranchesResp
}

// IsSuccess returns true when this list branches o k response has a 2xx status code
func (o *ListBranchesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list branches o k response has a 3xx status code
func (o *ListBranchesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list branches o k response has a 4xx status code
func (o *ListBranchesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list branches o k response has a 5xx status code
func (o *ListBranchesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list branches o k response a status code equal to that given
func (o *ListBranchesOK) IsCode(code int) bool {
	return code == 200
}

func (o *ListBranchesOK) Error() string {
	return fmt.Sprintf("[GET /api/v1beta/clusters/{cluster_id}/branches][%d] listBranchesOK  %+v", 200, o.Payload)
}

func (o *ListBranchesOK) String() string {
	return fmt.Sprintf("[GET /api/v1beta/clusters/{cluster_id}/branches][%d] listBranchesOK  %+v", 200, o.Payload)
}

func (o *ListBranchesOK) GetPayload() *models.OpenapiListBranchesResp {
	return o.Payload
}

func (o *ListBranchesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OpenapiListBranchesResp)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListBranchesDefault creates a ListBranchesDefault with default headers values
func NewListBranchesDefault(code int) *ListBranchesDefault {
	return &ListBranchesDefault{
		_statusCode: code,
	}
}

/*
ListBranchesDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ListBranchesDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// Code gets the status code for the list branches default response
func (o *ListBranchesDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this list branches default response has a 2xx status code
func (o *ListBranchesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list branches default response has a 3xx status code
func (o *ListBranchesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list branches default response has a 4xx status code
func (o *ListBranchesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list branches default response has a 5xx status code
func (o *ListBranchesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list branches default response a status code equal to that given
func (o *ListBranchesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ListBranchesDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1beta/clusters/{cluster_id}/branches][%d] ListBranches default  %+v", o._statusCode, o.Payload)
}

func (o *ListBranchesDefault) String() string {
	return fmt.Sprintf("[GET /api/v1beta/clusters/{cluster_id}/branches][%d] ListBranches default  %+v", o._statusCode, o.Payload)
}

func (o *ListBranchesDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *ListBranchesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
