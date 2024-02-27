// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1beta1GenerateUploadURLResponse v1beta1 generate upload Url response
//
// swagger:model v1beta1GenerateUploadUrlResponse
type V1beta1GenerateUploadURLResponse struct {

	// The ID of the upload
	UploadID string `json:"uploadId,omitempty"`

	// The URL to upload the file to
	UploadURL []string `json:"uploadUrl"`
}

// Validate validates this v1beta1 generate upload Url response
func (m *V1beta1GenerateUploadURLResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v1beta1 generate upload Url response based on context it is used
func (m *V1beta1GenerateUploadURLResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1beta1GenerateUploadURLResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1GenerateUploadURLResponse) UnmarshalBinary(b []byte) error {
	var res V1beta1GenerateUploadURLResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
