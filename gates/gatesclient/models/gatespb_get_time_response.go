// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GatespbGetTimeResponse gatespb get time response
//
// swagger:model gatespbGetTimeResponse
type GatespbGetTimeResponse struct {

	// time
	Time string `json:"time,omitempty"`
}

// Validate validates this gatespb get time response
func (m *GatespbGetTimeResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this gatespb get time response based on context it is used
func (m *GatespbGetTimeResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GatespbGetTimeResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GatespbGetTimeResponse) UnmarshalBinary(b []byte) error {
	var res GatespbGetTimeResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
