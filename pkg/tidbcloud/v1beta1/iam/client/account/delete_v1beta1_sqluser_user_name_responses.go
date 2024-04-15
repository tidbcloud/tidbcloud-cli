// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/iam/models"
)

// DeleteV1beta1SqluserUserNameReader is a Reader for the DeleteV1beta1SqluserUserName structure.
type DeleteV1beta1SqluserUserNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteV1beta1SqluserUserNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteV1beta1SqluserUserNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteV1beta1SqluserUserNameBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[DELETE /v1beta1/sqluser/{userName}] DeleteV1beta1SqluserUserName", response, response.Code())
	}
}

// NewDeleteV1beta1SqluserUserNameOK creates a DeleteV1beta1SqluserUserNameOK with default headers values
func NewDeleteV1beta1SqluserUserNameOK() *DeleteV1beta1SqluserUserNameOK {
	return &DeleteV1beta1SqluserUserNameOK{}
}

/*
DeleteV1beta1SqluserUserNameOK describes a response with status code 200, with default header values.

OK
*/
type DeleteV1beta1SqluserUserNameOK struct {
	Payload *models.APIBasicResp
}

// IsSuccess returns true when this delete v1beta1 sqluser user name o k response has a 2xx status code
func (o *DeleteV1beta1SqluserUserNameOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete v1beta1 sqluser user name o k response has a 3xx status code
func (o *DeleteV1beta1SqluserUserNameOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1beta1 sqluser user name o k response has a 4xx status code
func (o *DeleteV1beta1SqluserUserNameOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete v1beta1 sqluser user name o k response has a 5xx status code
func (o *DeleteV1beta1SqluserUserNameOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1beta1 sqluser user name o k response a status code equal to that given
func (o *DeleteV1beta1SqluserUserNameOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete v1beta1 sqluser user name o k response
func (o *DeleteV1beta1SqluserUserNameOK) Code() int {
	return 200
}

func (o *DeleteV1beta1SqluserUserNameOK) Error() string {
	return fmt.Sprintf("[DELETE /v1beta1/sqluser/{userName}][%d] deleteV1beta1SqluserUserNameOK  %+v", 200, o.Payload)
}

func (o *DeleteV1beta1SqluserUserNameOK) String() string {
	return fmt.Sprintf("[DELETE /v1beta1/sqluser/{userName}][%d] deleteV1beta1SqluserUserNameOK  %+v", 200, o.Payload)
}

func (o *DeleteV1beta1SqluserUserNameOK) GetPayload() *models.APIBasicResp {
	return o.Payload
}

func (o *DeleteV1beta1SqluserUserNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIBasicResp)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteV1beta1SqluserUserNameBadRequest creates a DeleteV1beta1SqluserUserNameBadRequest with default headers values
func NewDeleteV1beta1SqluserUserNameBadRequest() *DeleteV1beta1SqluserUserNameBadRequest {
	return &DeleteV1beta1SqluserUserNameBadRequest{}
}

/*
DeleteV1beta1SqluserUserNameBadRequest describes a response with status code 400, with default header values.

Bad Request
*/
type DeleteV1beta1SqluserUserNameBadRequest struct {
	Payload *models.APIOpenAPIError
}

// IsSuccess returns true when this delete v1beta1 sqluser user name bad request response has a 2xx status code
func (o *DeleteV1beta1SqluserUserNameBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete v1beta1 sqluser user name bad request response has a 3xx status code
func (o *DeleteV1beta1SqluserUserNameBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1beta1 sqluser user name bad request response has a 4xx status code
func (o *DeleteV1beta1SqluserUserNameBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete v1beta1 sqluser user name bad request response has a 5xx status code
func (o *DeleteV1beta1SqluserUserNameBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1beta1 sqluser user name bad request response a status code equal to that given
func (o *DeleteV1beta1SqluserUserNameBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the delete v1beta1 sqluser user name bad request response
func (o *DeleteV1beta1SqluserUserNameBadRequest) Code() int {
	return 400
}

func (o *DeleteV1beta1SqluserUserNameBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /v1beta1/sqluser/{userName}][%d] deleteV1beta1SqluserUserNameBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteV1beta1SqluserUserNameBadRequest) String() string {
	return fmt.Sprintf("[DELETE /v1beta1/sqluser/{userName}][%d] deleteV1beta1SqluserUserNameBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteV1beta1SqluserUserNameBadRequest) GetPayload() *models.APIOpenAPIError {
	return o.Payload
}

func (o *DeleteV1beta1SqluserUserNameBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIOpenAPIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}