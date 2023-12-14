// Code generated by go-swagger; DO NOT EDIT.

package branch_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/models"
)

// BranchServiceGetBranchReader is a Reader for the BranchServiceGetBranch structure.
type BranchServiceGetBranchReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *BranchServiceGetBranchReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewBranchServiceGetBranchOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewBranchServiceGetBranchDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewBranchServiceGetBranchOK creates a BranchServiceGetBranchOK with default headers values
func NewBranchServiceGetBranchOK() *BranchServiceGetBranchOK {
	return &BranchServiceGetBranchOK{}
}

/*
BranchServiceGetBranchOK describes a response with status code 200, with default header values.

A successful response.
*/
type BranchServiceGetBranchOK struct {
	Payload *models.V1beta1Branch
}

// IsSuccess returns true when this branch service get branch o k response has a 2xx status code
func (o *BranchServiceGetBranchOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this branch service get branch o k response has a 3xx status code
func (o *BranchServiceGetBranchOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this branch service get branch o k response has a 4xx status code
func (o *BranchServiceGetBranchOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this branch service get branch o k response has a 5xx status code
func (o *BranchServiceGetBranchOK) IsServerError() bool {
	return false
}

// IsCode returns true when this branch service get branch o k response a status code equal to that given
func (o *BranchServiceGetBranchOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the branch service get branch o k response
func (o *BranchServiceGetBranchOK) Code() int {
	return 200
}

func (o *BranchServiceGetBranchOK) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] branchServiceGetBranchOK  %+v", 200, o.Payload)
}

func (o *BranchServiceGetBranchOK) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] branchServiceGetBranchOK  %+v", 200, o.Payload)
}

func (o *BranchServiceGetBranchOK) GetPayload() *models.V1beta1Branch {
	return o.Payload
}

func (o *BranchServiceGetBranchOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1Branch)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBranchServiceGetBranchDefault creates a BranchServiceGetBranchDefault with default headers values
func NewBranchServiceGetBranchDefault(code int) *BranchServiceGetBranchDefault {
	return &BranchServiceGetBranchDefault{
		_statusCode: code,
	}
}

/*
BranchServiceGetBranchDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type BranchServiceGetBranchDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this branch service get branch default response has a 2xx status code
func (o *BranchServiceGetBranchDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this branch service get branch default response has a 3xx status code
func (o *BranchServiceGetBranchDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this branch service get branch default response has a 4xx status code
func (o *BranchServiceGetBranchDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this branch service get branch default response has a 5xx status code
func (o *BranchServiceGetBranchDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this branch service get branch default response a status code equal to that given
func (o *BranchServiceGetBranchDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the branch service get branch default response
func (o *BranchServiceGetBranchDefault) Code() int {
	return o._statusCode
}

func (o *BranchServiceGetBranchDefault) Error() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] BranchService_GetBranch default  %+v", o._statusCode, o.Payload)
}

func (o *BranchServiceGetBranchDefault) String() string {
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] BranchService_GetBranch default  %+v", o._statusCode, o.Payload)
}

func (o *BranchServiceGetBranchDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *BranchServiceGetBranchDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}