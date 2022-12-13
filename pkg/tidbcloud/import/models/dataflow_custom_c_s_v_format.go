// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataflowCustomCSVFormat dataflow custom c s v format
//
// swagger:model dataflowCustomCSVFormat
type DataflowCustomCSVFormat struct {

	// backslash escape
	BackslashEscape bool `json:"backslash_escape,omitempty"`

	// delimiter
	Delimiter string `json:"delimiter,omitempty"`

	// header
	Header bool `json:"header,omitempty"`

	// not null
	NotNull bool `json:"not_null,omitempty"`

	// null
	Null string `json:"null,omitempty"`

	// separator
	Separator string `json:"separator,omitempty"`

	// trim last separator
	TrimLastSeparator bool `json:"trim_last_separator,omitempty"`
}

// Validate validates this dataflow custom c s v format
func (m *DataflowCustomCSVFormat) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dataflow custom c s v format based on context it is used
func (m *DataflowCustomCSVFormat) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DataflowCustomCSVFormat) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataflowCustomCSVFormat) UnmarshalBinary(b []byte) error {
	var res DataflowCustomCSVFormat
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
