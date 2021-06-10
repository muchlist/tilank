package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c ViolationRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.NoIdentity, validation.Required),
		validation.Field(&c.OwnerID, validation.Required),
		validation.Field(&c.TypeViolation, validation.Required),
		validation.Field(&c.DetailViolation, validation.Required),
		validation.Field(&c.Location, validation.Required),
	); err != nil {
		return err
	}
	// validate type
	if err := typeViolationValidation(c.TypeViolation); err != nil {
		return err
	}
	return nil
}

func (c ViolationEditRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.NoIdentity, validation.Required),
		validation.Field(&c.OwnerID, validation.Required),
		validation.Field(&c.TypeViolation, validation.Required),
		validation.Field(&c.DetailViolation, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.FilterTimestamp, validation.Required),
	); err != nil {
		return err
	}
	// validate type
	if err := typeViolationValidation(c.TypeViolation); err != nil {
		return err
	}
	return nil
}