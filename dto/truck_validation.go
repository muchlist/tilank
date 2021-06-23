package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (t TruckRequest) Validate() error {
	if err := validation.ValidateStruct(&t,
		validation.Field(&t.NoIdentity, validation.Required),
		validation.Field(&t.NoPol, validation.Required),
		validation.Field(&t.Mark, validation.Required),
		validation.Field(&t.Owner, validation.Required),
		validation.Field(&t.Hp, validation.Required),
		validation.Field(&t.Email, validation.Required, is.Email),
	); err != nil {
		return err
	}
	return nil
}

func (t TruckEditRequest) Validate() error {
	if err := validation.ValidateStruct(&t,
		validation.Field(&t.NoIdentity, validation.Required),
		validation.Field(&t.NoPol, validation.Required),
		validation.Field(&t.Mark, validation.Required),
		validation.Field(&t.Owner, validation.Required),
		validation.Field(&t.Email, validation.Required, is.Email),
		validation.Field(&t.FilterTimestamp, validation.Required),
	); err != nil {
		return err
	}
	return nil
}
