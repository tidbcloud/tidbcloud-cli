// Code generated by go-swagger; DO NOT EDIT.

package export_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_export/models"
)

// ExportServiceListExportsReader is a Reader for the ExportServiceListExports structure.
type ExportServiceListExportsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ExportServiceListExportsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewExportServiceListExportsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewExportServiceListExportsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewExportServiceListExportsOK creates a ExportServiceListExportsOK with default headers values
func NewExportServiceListExportsOK() *ExportServiceListExportsOK {
	return &ExportServiceListExportsOK{}
}

/*
ExportServiceListExportsOK describes a response with status code 200, with default header values.

A successful response.
*/
type ExportServiceListExportsOK struct {
	Payload *models.V1beta1ListExportsResponse
}

// IsSuccess returns true when this export service list exports o k response has a 2xx status code
func (o *ExportServiceListExportsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this export service list exports o k response has a 3xx status code
func (o *ExportServiceListExportsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this export service list exports o k response has a 4xx status code
func (o *ExportServiceListExportsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this export service list exports o k response has a 5xx status code
func (o *ExportServiceListExportsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this export service list exports o k response a status code equal to that given
func (o *ExportServiceListExportsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the export service list exports o k response
func (o *ExportServiceListExportsOK) Code() int {
	return 200
}

func (o *ExportServiceListExportsOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/exports][%d] exportServiceListExportsOK %s", 200, payload)
}

func (o *ExportServiceListExportsOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/exports][%d] exportServiceListExportsOK %s", 200, payload)
}

func (o *ExportServiceListExportsOK) GetPayload() *models.V1beta1ListExportsResponse {
	return o.Payload
}

func (o *ExportServiceListExportsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1ListExportsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewExportServiceListExportsDefault creates a ExportServiceListExportsDefault with default headers values
func NewExportServiceListExportsDefault(code int) *ExportServiceListExportsDefault {
	return &ExportServiceListExportsDefault{
		_statusCode: code,
	}
}

/*
ExportServiceListExportsDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type ExportServiceListExportsDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this export service list exports default response has a 2xx status code
func (o *ExportServiceListExportsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this export service list exports default response has a 3xx status code
func (o *ExportServiceListExportsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this export service list exports default response has a 4xx status code
func (o *ExportServiceListExportsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this export service list exports default response has a 5xx status code
func (o *ExportServiceListExportsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this export service list exports default response a status code equal to that given
func (o *ExportServiceListExportsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the export service list exports default response
func (o *ExportServiceListExportsDefault) Code() int {
	return o._statusCode
}

func (o *ExportServiceListExportsDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/exports][%d] ExportService_ListExports default %s", o._statusCode, payload)
}

func (o *ExportServiceListExportsDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1beta1/clusters/{clusterId}/exports][%d] ExportService_ListExports default %s", o._statusCode, payload)
}

func (o *ExportServiceListExportsDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *ExportServiceListExportsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
