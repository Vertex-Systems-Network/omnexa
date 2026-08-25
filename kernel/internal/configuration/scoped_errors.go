package configuration

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeScopedServiceInvalid  failure.Code = "configuration.scoped.service.invalid"
	codeScopedPolicyInvalid   failure.Code = "configuration.scoped.policy.invalid"
	codeScopedContextInvalid  failure.Code = "configuration.scoped.context.invalid"
	codeScopedMutationInvalid failure.Code = "configuration.scoped.mutation.invalid"
	codeScopedValueInvalid    failure.Code = "configuration.scoped.value.invalid"
	codeScopedPersistence     failure.Code = "configuration.scoped.persistence.failed"
)

func scopedServiceInvalidFailure() error {
	return classifiedFailure(codeScopedServiceInvalid, failure.CategoryInvariant, "scoped setting service is invalid")
}

func scopedPolicyInvalidFailure() error {
	return classifiedFailure(codeScopedPolicyInvalid, failure.CategoryValidation, "scoped setting policy is invalid")
}

func scopedContextInvalidFailure() error {
	return classifiedFailure(codeScopedContextInvalid, failure.CategoryAuthorization, "scoped setting context is not trusted")
}

func scopedMutationInvalidFailure() error {
	return classifiedFailure(codeScopedMutationInvalid, failure.CategoryValidation, "scoped setting mutation metadata is invalid")
}

func scopedValueInvalidFailure() error {
	return classifiedFailure(codeScopedValueInvalid, failure.CategoryValidation, "scoped setting value is invalid")
}

func scopedPersistenceFailure(cause error) error {
	value, err := failure.Wrap(cause, codeScopedPersistence, failure.CategoryDependency, "scoped setting persistence failed")
	if err != nil {
		return errors.New("scoped setting persistence failure could not be classified safely")
	}
	return value
}
