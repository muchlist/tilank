package dto

import (
	"fmt"
	"tilank/config"
	"tilank/utils/sfunc"
)

func roleValidation(rolesIn []string) error {
	if len(rolesIn) > 0 {
		if !sfunc.ValueInSliceIsAvailable(rolesIn, config.GetRolesAvailable()) {
			return fmt.Errorf("role yang dimasukkan tidak tersedia. gunakan %s", config.GetRolesAvailable())
		}
	}
	return nil
}

func typeViolationValidation(typeViolation string) error {

	if !sfunc.InSlice(typeViolation, config.GetTypeAvailable()) {
		return fmt.Errorf("tile yang dimasukkan tidak tersedia. gunakan %s", config.GetTypeAvailable())
	}

	return nil
}
