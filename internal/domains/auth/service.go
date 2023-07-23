package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rasulov-emirlan/accounter-backend/internal/entities"
	"github.com/rasulov-emirlan/accounter-backend/pkg/logging"
	"github.com/rasulov-emirlan/accounter-backend/pkg/telemetry"
	"github.com/rasulov-emirlan/accounter-backend/pkg/validation"
)

type (
	OwnersRepository interface {
		Create(ctx context.Context, owner entities.Owner) (entities.Owner, error)
		Read(ctx context.Context, id string) (entities.Owner, error)
		ReadByUsername(ctx context.Context, username string) (entities.Owner, error)
		ReadAll(ctx context.Context) ([]entities.Owner, error) // TODO: add pagination
		Update(ctx context.Context, owner entities.Owner) (entities.Owner, error)
		Delete(ctx context.Context, id string) error
	}

	KeyValueRepository interface {
		Set(ctx context.Context, key, value string, ttl time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Delete(ctx context.Context, key string) error
	}

	Service interface {
		Register(ctx context.Context, input RegisterInput) (Session, error)
		Login(ctx context.Context, input LoginInput) (Session, error)
		Refresh(ctx context.Context, refreshToken string) (Session, error)

		ParseAccessKey(ctx context.Context, accessToken string) (AccessKey, error)
		ParseRefreshKey(ctx context.Context, refreshToken string) (RefreshKey, error)
		Me(ctx context.Context, accessKey AccessKey) (entities.Owner, error)
	}

	service struct {
		ownersRepo OwnersRepository
		log        *logging.Logger
		val        *validation.Validator

		secretKey []byte
	}
)

var _ Service = (*service)(nil)

func NewService(ownersRepo OwnersRepository, log *logging.Logger, val *validation.Validator, secretKey []byte) service {
	return service{
		ownersRepo: ownersRepo,
		log:        log,
		val:        val,
	}
}

func (s service) Register(ctx context.Context, input RegisterInput) (Session, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Register")).End()
	defer s.log.Sync()

	o, err := entities.NewOwner(input.PhoneNumber, input.FullName, input.Username, input.Password)
	if err != nil {
		s.log.Debug("auth:Register - failed to create owner", logging.String("stage", "validation"), logging.Error("err", err))
		return Session{}, err
	}
	if err := s.val.Validate(o); err != nil {
		s.log.Debug("auth:Register - failed to validate owner", logging.String("stage", "validation"), logging.Error("err", err))
		return Session{}, err
	}

	o, err = s.ownersRepo.Create(ctx, o)
	if err != nil {
		if err == ErrUsernameTaken {
			s.log.Debug("auth:Register - username is taken", logging.String("stage", "repository"), logging.Error("err", err))
			return Session{}, err
		}
		s.log.Error("auth:Register - failed to create owner", logging.String("stage", "repository"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	session, err := generateSession(o, s.secretKey)
	if err != nil {
		s.log.Error("auth:Register - failed to generate session", logging.String("stage", "jwt"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	s.log.Info("auth:Register - successfully registered", logging.String("stage", "success"), logging.String("username", o.Username))

	return session, nil
}

func (s service) Login(ctx context.Context, input LoginInput) (Session, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Login")).End()
	defer s.log.Sync()

	o, err := s.ownersRepo.ReadByUsername(ctx, input.Username)
	if err != nil {
		if err == ErrUsernameNotFound {
			s.log.Debug("auth:Login - owner not found", logging.String("stage", "repository"), logging.Error("err", err))
			return Session{}, err
		}
		s.log.Error("auth:Login - failed to read owner", logging.String("stage", "repository"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	if err := o.ComparePassword(input.Password); err != nil {
		s.log.Debug("auth:Login - wrong password", logging.String("stage", "validation"), logging.Error("err", err))
		return Session{}, ErrWrongPassword
	}

	session, err := generateSession(o, s.secretKey)
	if err != nil {
		s.log.Error("auth:Login - failed to generate session", logging.String("stage", "jwt"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	s.log.Info("auth:Login - successfully logged in", logging.String("stage", "success"), logging.String("username", o.Username))

	return session, nil
}

func (s service) Refresh(ctx context.Context, refreshToken string) (Session, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Refresh")).End()
	defer s.log.Sync()

	claims, err := s.ParseRefreshKey(ctx, refreshToken)
	if err != nil {
		s.log.Debug("auth:Refresh - failed to parse refresh key", logging.String("stage", "jwt"), logging.Error("err", err))
		return Session{}, ErrInvalidRefreshToken
	}

	o, err := s.ownersRepo.Read(ctx, claims.UserID)
	if err != nil {
		if err == ErrIdNotFound {
			s.log.Debug("auth:Refresh - owner not found", logging.String("stage", "repository"), logging.Error("err", err))
			return Session{}, err
		}
		s.log.Error("auth:Refresh - failed to read owner", logging.String("stage", "repository"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	session, err := generateSession(o, s.secretKey)
	if err != nil {
		s.log.Error("auth:Refresh - failed to generate session", logging.String("stage", "jwt"), logging.Error("err", err))
		return Session{}, ErrDefault
	}

	s.log.Info("auth:Refresh - successfully refreshed session", logging.String("stage", "success"), logging.String("username", o.Username))

	return session, nil
}

func generateSession(o entities.Owner, secretKey []byte) (Session, error) {
	claims := AccessKey{
		UserID: o.ID.String(),
		Role:   RoleOwner,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessKeyTTL).Unix(),
		},
	}

	accessToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secretKey)
	if err != nil {
		return Session{}, err
	}

	rClaims := RefreshKey{
		UserID: o.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshKeyTTL).Unix(),
		},
	}

	refreshToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, rClaims).
		SignedString(secretKey)
	if err != nil {
		return Session{}, err
	}

	return Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s service) ParseAccessKey(ctx context.Context, key string) (AccessKey, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.ParseAccessKey")).End()
	var claims AccessKey
	_, err := jwt.ParseWithClaims(key, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return AccessKey{}, err
	}

	return claims, nil
}

func (s service) ParseRefreshKey(ctx context.Context, key string) (RefreshKey, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.ParseRefreshKey")).End()
	var claims RefreshKey
	_, err := jwt.ParseWithClaims(key, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return RefreshKey{}, err
	}

	return claims, nil
}

func (s service) Me(ctx context.Context, accessKey AccessKey) (entities.Owner, error) {
	defer telemetry.NewSpan(ctx, telemetry.Name(PackageName+"service.Me")).End()
	defer s.log.Sync()

	o, err := s.ownersRepo.Read(ctx, accessKey.UserID)
	if err != nil {
		if err == ErrIdNotFound {
			s.log.Debug("auth:Me - owner not found", logging.String("stage", "repository"), logging.Error("err", err))
			return entities.Owner{}, err
		}
		s.log.Error("auth:Me - failed to read owner", logging.String("stage", "repository"), logging.Error("err", err))
		return entities.Owner{}, ErrDefault
	}

	s.log.Info("auth:Me - successfully read owner", logging.String("stage", "success"), logging.String("username", o.Username))
	return o, nil
}
