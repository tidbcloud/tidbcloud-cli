// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PingchatLink pingchat link
//
// swagger:model pingchat.Link
type PingchatLink struct {

	// link
	Link string `json:"link,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this pingchat link
func (m *PingchatLink) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this pingchat link based on context it is used
func (m *PingchatLink) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PingchatLink) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PingchatLink) UnmarshalBinary(b []byte) error {
	var res PingchatLink
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
