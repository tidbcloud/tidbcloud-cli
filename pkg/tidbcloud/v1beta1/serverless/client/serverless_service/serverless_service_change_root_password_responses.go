// Code generated by go-swagger; DO NOT EDIT.

package serverless_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/models"
)

// ServerlessServiceChangeRootPasswordReader is a Reader for the ServerlessServiceChangeRootPassword structure.
type ServerlessServiceChangeRootPasswordReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ServerlessServiceChangeRootPasswordReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewServerlessServiceChangeRootPasswordOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewServerlessServiceChangeRootPasswordDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewServerlessServiceChangeRootPasswordOK creates a ServerlessServiceChangeRootPasswordOK with default headers values
func NewServerlessServiceChangeRootPasswordOK() *ServerlessServiceChangeRootPasswordOK {
	return &ServerlessServiceChangeRootPasswordOK{}
}

/*
ServerlessServiceChangeRootPasswordOK describes a response with status code 200, with default header values.

A successful response.
*/
type ServerlessServiceChangeRootPasswordOK struct {
	Payload models.V1beta1ChangeRootPasswordResponse
}

// IsSuccess returns true when this serverless service change root password o k response has a 2xx status code
func (o *ServerlessServiceChangeRootPasswordOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this serverless service change root password o k response has a 3xx status code
func (o *ServerlessServiceChangeRootPasswordOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this serverless service change root password o k response has a 4xx status code
func (o *ServerlessServiceChangeRootPasswordOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this serverless service change root password o k response has a 5xx status code
func (o *ServerlessServiceChangeRootPasswordOK) IsServerError() bool {
	return false
}

// IsCode returns true when this serverless service change root password o k response a status code equal to that given
func (o *ServerlessServiceChangeRootPasswordOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the serverless service change root password o k response
func (o *ServerlessServiceChangeRootPasswordOK) Code() int {
	return 200
}

func (o *ServerlessServiceChangeRootPasswordOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /clusters/{clusterId}/password][%d] serverlessServiceChangeRootPasswordOK %s", 200, payload)
}

func (o *ServerlessServiceChangeRootPasswordOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /clusters/{clusterId}/password][%d] serverlessServiceChangeRootPasswordOK %s", 200, payload)
}

func (o *ServerlessServiceChangeRootPasswordOK) GetPayload() models.V1beta1ChangeRootPasswordResponse {
	return o.Payload
}

func (o *ServerlessServiceChangeRootPasswordOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewServerlessServiceChangeRootPasswordDefault creates a ServerlessServiceChangeRootPasswordDefault with default headers values
func NewServerlessServiceChangeRootPasswordDefault(code int) *ServerlessServiceChangeRootPasswordDefault {
	return &ServerlessServiceChangeRootPasswordDefault{
		_statusCode: code,
	}
}

/*
ServerlessServiceChangeRootPasswordDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ServerlessServiceChangeRootPasswordDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this serverless service change root password default response has a 2xx status code
func (o *ServerlessServiceChangeRootPasswordDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this serverless service change root password default response has a 3xx status code
func (o *ServerlessServiceChangeRootPasswordDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this serverless service change root password default response has a 4xx status code
func (o *ServerlessServiceChangeRootPasswordDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this serverless service change root password default response has a 5xx status code
func (o *ServerlessServiceChangeRootPasswordDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this serverless service change root password default response a status code equal to that given
func (o *ServerlessServiceChangeRootPasswordDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the serverless service change root password default response
func (o *ServerlessServiceChangeRootPasswordDefault) Code() int {
	return o._statusCode
}

func (o *ServerlessServiceChangeRootPasswordDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /clusters/{clusterId}/password][%d] ServerlessService_ChangeRootPassword default %s", o._statusCode, payload)
}

func (o *ServerlessServiceChangeRootPasswordDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /clusters/{clusterId}/password][%d] ServerlessService_ChangeRootPassword default %s", o._statusCode, payload)
}

func (o *ServerlessServiceChangeRootPasswordDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *ServerlessServiceChangeRootPasswordDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
ServerlessServiceChangeRootPasswordBody Message for requesting to change the root password of a TiDB Serverless cluster.
swagger:model ServerlessServiceChangeRootPasswordBody
*/
type ServerlessServiceChangeRootPasswordBody struct {

	// Required. The new root password for the cluster.
	// Required: true
	Password *string `json:"password"`
}

// Validate validates this serverless service change root password body
func (o *ServerlessServiceChangeRootPasswordBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ServerlessServiceChangeRootPasswordBody) validatePassword(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"password", "body", o.Password); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this serverless service change root password body based on context it is used
func (o *ServerlessServiceChangeRootPasswordBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ServerlessServiceChangeRootPasswordBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ServerlessServiceChangeRootPasswordBody) UnmarshalBinary(b []byte) error {
	var res ServerlessServiceChangeRootPasswordBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
