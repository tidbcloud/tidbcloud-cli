// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Content content
//
// swagger:model Content
type Content struct {

	// comment
	Comment string `json:"comment,omitempty"`

	// connection example
	ConnectionExample string `json:"connection_example,omitempty"`

	// connection string
	ConnectionString string `json:"connection_string,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this content
func (m *Content) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this content based on context it is used
func (m *Content) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Content) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Content) UnmarshalBinary(b []byte) error {
	var res Content
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}