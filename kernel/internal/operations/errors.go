package operations

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeDependencyInvalid   failure.Code = "operations.dependency.invalid"
	codeDependencyDuplicate failure.Code = "operations.dependency.duplicate"
	codeRegistryFrozen      failure.Code = "operations.registry.frozen"
	codeRegistryInvalid     failure.Code = "operations.registry.invalid"
)

func safeFailure(code failure.Code, title string) error {
	category := failure.CategoryValidation
	if code == codeDependencyDuplicate || code == codeRegistryFrozen {
		category = failure.CategoryConflict
	}
	value, err := failure.New(code, category, title)
	if err != nil {
		return errors.New("operations failure could not be classified safely")
	}
	return value
}
