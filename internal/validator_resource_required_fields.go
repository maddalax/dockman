package internal

import (
	"errors"
)

type RequiredFieldsValidator struct {
	Resource *Resource
}

func (v RequiredFieldsValidator) Validate() error {
	if v.Resource.Name == "" {
		return errors.New("name is required")
	}

	if v.Resource.Environment == "" {
		return errors.New("environment is required")
	}

	if v.Resource.RunType == RunTypeUnknown {
		return errors.New("run type is required")
	}

	return nil
}
