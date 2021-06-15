package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c JptRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.OwnerName, validation.Required),
		validation.Field(&c.Hp, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Branch, validation.Required),
	); err != nil {
		return err
	}
	return nil
}

func (c JptEditRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.OwnerName, validation.Required),
		validation.Field(&c.Hp, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.FilterTimestamp, validation.Required),
	); err != nil {
		return err
	}
	return nil
}
