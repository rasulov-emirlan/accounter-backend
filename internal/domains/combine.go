package domains

import "github.com/rasulov-emirlan/esep-backend/internal/domains/auth"

type DomainCombiner struct {
	authService auth.Service
}

func NewDomainCombiner(cD CommonDependencies, aD AuthDependencies) (DomainCombiner, error) {
	if err := cD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	if err := aD.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	authService := auth.NewService(aD.OwnersRepo, cD.Log, cD.Val, aD.SecretKey)

	return DomainCombiner{
		authService: authService,
	}, nil
}

func (d DomainCombiner) AuthService() auth.Service {
	return d.authService
}
