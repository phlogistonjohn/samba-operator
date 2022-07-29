// SPDX-License-Identifier: Apache-2.0

package planner

import (
	"fmt"
)

type IncompatibleInstanceError struct {
	current  string
	existing string
}

func (iie IncompatibleInstanceError) Error() string {
	return fmt.Sprintf("Share resource %s is incompatible with %s",
		iie.current,
		iie.existing)
}

// CheckCompatible returns an error if the instance configurations are
// not compatible with each other. A compatible instance is one that
// can share the same smbd.
func CheckCompatible(current, existing InstanceConfiguration) error {
	// This returns an error rather than a boolean because in the future
	// we can choose to add details about what is incompatible to
	// the specific error type.
	return IncompatibleInstanceError{
		current.SmbShare.Name,
		existing.SmbShare.Name,
	}
}
