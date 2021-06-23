package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r RulesRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Score, validation.Required),
		validation.Field(&r.BlockTime, validation.Required, validation.Min(0)),
		validation.Field(&r.Description, validation.Required),
	); err != nil {
		return err
	}
	return nil
}

func (r RulesEditRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Score, validation.Required),
		validation.Field(&r.BlockTime, validation.Required, validation.Min(0)),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.FilterTimestamp, validation.Required),
	); err != nil {
		return err
	}
	return nil
}
