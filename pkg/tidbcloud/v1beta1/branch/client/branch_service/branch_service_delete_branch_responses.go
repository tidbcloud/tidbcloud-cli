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

// BranchServiceDeleteBranchReader is a Reader for the BranchServiceDeleteBranch structure.
type BranchServiceDeleteBranchReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *BranchServiceDeleteBranchReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewBranchServiceDeleteBranchOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewBranchServiceDeleteBranchDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewBranchServiceDeleteBranchOK creates a BranchServiceDeleteBranchOK with default headers values
func NewBranchServiceDeleteBranchOK() *BranchServiceDeleteBranchOK {
	return &BranchServiceDeleteBranchOK{}
}

/*
BranchServiceDeleteBranchOK describes a response with status code 200, with default header values.

A successful response.
*/
type BranchServiceDeleteBranchOK struct {
	Payload *models.V1beta1Branch
}

// IsSuccess returns true when this branch service delete branch o k response has a 2xx status code
func (o *BranchServiceDeleteBranchOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this branch service delete branch o k response has a 3xx status code
func (o *BranchServiceDeleteBranchOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this branch service delete branch o k response has a 4xx status code
func (o *BranchServiceDeleteBranchOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this branch service delete branch o k response has a 5xx status code
func (o *BranchServiceDeleteBranchOK) IsServerError() bool {
	return false
}

// IsCode returns true when this branch service delete branch o k response a status code equal to that given
func (o *BranchServiceDeleteBranchOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the branch service delete branch o k response
func (o *BranchServiceDeleteBranchOK) Code() int {
	return 200
}

func (o *BranchServiceDeleteBranchOK) Error() string {
	return fmt.Sprintf("[DELETE /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] branchServiceDeleteBranchOK  %+v", 200, o.Payload)
}

func (o *BranchServiceDeleteBranchOK) String() string {
	return fmt.Sprintf("[DELETE /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] branchServiceDeleteBranchOK  %+v", 200, o.Payload)
}

func (o *BranchServiceDeleteBranchOK) GetPayload() *models.V1beta1Branch {
	return o.Payload
}

func (o *BranchServiceDeleteBranchOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1Branch)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBranchServiceDeleteBranchDefault creates a BranchServiceDeleteBranchDefault with default headers values
func NewBranchServiceDeleteBranchDefault(code int) *BranchServiceDeleteBranchDefault {
	return &BranchServiceDeleteBranchDefault{
		_statusCode: code,
	}
}

/*
BranchServiceDeleteBranchDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type BranchServiceDeleteBranchDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this branch service delete branch default response has a 2xx status code
func (o *BranchServiceDeleteBranchDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this branch service delete branch default response has a 3xx status code
func (o *BranchServiceDeleteBranchDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this branch service delete branch default response has a 4xx status code
func (o *BranchServiceDeleteBranchDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this branch service delete branch default response has a 5xx status code
func (o *BranchServiceDeleteBranchDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this branch service delete branch default response a status code equal to that given
func (o *BranchServiceDeleteBranchDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the branch service delete branch default response
func (o *BranchServiceDeleteBranchDefault) Code() int {
	return o._statusCode
}

func (o *BranchServiceDeleteBranchDefault) Error() string {
	return fmt.Sprintf("[DELETE /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] BranchService_DeleteBranch default  %+v", o._statusCode, o.Payload)
}

func (o *BranchServiceDeleteBranchDefault) String() string {
	return fmt.Sprintf("[DELETE /v1beta1/clusters/{clusterId}/branches/{branchId}][%d] BranchService_DeleteBranch default  %+v", o._statusCode, o.Payload)
}

func (o *BranchServiceDeleteBranchDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *BranchServiceDeleteBranchDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
