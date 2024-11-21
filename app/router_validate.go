package app

import (
	"errors"
	"fmt"
	"github.com/gobwas/glob"
)

func ValidateBlock(block *RouteBlock) error {
	if block.ResourceId == "" {
		return errors.New("resource id is required for block")
	}

	if block.Hostname == "" {
		return errors.New("hostname is required")
	}

	if block.PathMatchModifier == "" {
		return errors.New("path match modifier is required")
	}

	if block.Path != "" && block.PathMatchModifier == "glob" {
		_, err := glob.Compile(block.Path)
		if err != nil {
			return fmt.Errorf("failed to compile glob: %s", err.Error())
		}
	}

	return nil
}
