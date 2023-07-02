package domains

import (
	"github.com/SanaripEsep/esep-backend/internal/domains/auth"
	"github.com/SanaripEsep/esep-backend/internal/domains/categories"
	"github.com/SanaripEsep/esep-backend/internal/domains/stores"
)

type DomainCombiner struct {
	authService       auth.Service
	storesService     stores.Service
	categoriesService categories.Service
}

func NewDomainCombiner(
	cD CommonDependencies,
	aD AuthDependencies,
	sD StoresDependencies,
	categoryD CategoriesDependencies) (DomainCombiner, error) {
	if err := cD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	if err := aD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	if err := sD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	if err := categoryD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	return DomainCombiner{
		authService:       auth.NewService(aD.OwnersRepo, cD.Log, cD.Val, aD.SecretKey),
		storesService:     stores.NewService(sD.StoresRepo, cD.Log),
		categoriesService: categories.NewService(categoryD.CategoriesRepo, cD.Log),
	}, nil
}

func (d DomainCombiner) AuthService() auth.Service {
	return d.authService
}

func (d DomainCombiner) StoresService() stores.Service {
	return d.storesService
}

func (d DomainCombiner) CategoriesService() categories.Service {
	return d.categoriesService
}
