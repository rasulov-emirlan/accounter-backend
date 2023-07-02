package domains

import (
	"fmt"
	"reflect"

	"github.com/SanaripEsep/esep-backend/internal/domains/auth"
	"github.com/SanaripEsep/esep-backend/internal/domains/categories"
	"github.com/SanaripEsep/esep-backend/internal/domains/stores"
	"github.com/SanaripEsep/esep-backend/pkg/logging"
	"github.com/SanaripEsep/esep-backend/pkg/validation"
)

type CommonDependencies struct {
	Log *logging.Logger
	Val *validation.Validator
}

func (d CommonDependencies) Validate() error {
	if d.Log == nil {
		return DependencyError{
			Dependency:       "CommonDependencies.Log",
			BrokenConstraint: "logger cannot be nil",
		}
	}
	if d.Val == nil {
		return DependencyError{
			Dependency:       "CommonDependencies.Val",
			BrokenConstraint: "validator cannot be nil",
		}
	}
	return nil
}

type AuthDependencies struct {
	OwnersRepo auth.OwnersRepository
	SecretKey  []byte
}

func (d AuthDependencies) Validate() error {
	if isNil(d.OwnersRepo) {
		return DependencyError{
			Dependency:       "AuthDependencies.OwnersRepo",
			BrokenConstraint: "owners repository cannot be nil",
		}
	}

	if len(d.SecretKey) < 4 {
		return DependencyError{
			Dependency:       "AuthDependencies.SecretKey",
			BrokenConstraint: "secret key must be at least 4 bytes long",
		}
	}

	return nil
}

type StoresDependencies struct {
	StoresRepo stores.StoresRepository
}

func (d StoresDependencies) Validate() error {
	if isNil(d.StoresRepo) {
		return DependencyError{
			Dependency:       "StoresDependencies.StoresRepo",
			BrokenConstraint: "stores repository cannot be nil",
		}
	}

	return nil
}

type CategoriesDependencies struct {
	CategoriesRepo categories.CategoriesRepository
}

func (d CategoriesDependencies) Validate() error {
	if isNil(d.CategoriesRepo) {
		return DependencyError{
			Dependency:       "CategoriesDependencies.CategoriesRepo",
			BrokenConstraint: "categories repository cannot be nil",
		}
	}

	return nil
}

type DependencyError struct {
	Dependency       string
	BrokenConstraint string
}

func (e DependencyError) Error() string {
	return fmt.Sprintf("dependency: %s, broke constraint: %s", e.Dependency, e.BrokenConstraint)
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
