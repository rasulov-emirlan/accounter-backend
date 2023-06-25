package auth

import "github.com/golang-jwt/jwt"

type (
	RegisterInput struct {
		FullName    string `json:"fullName" validate:"required"`
		Username    string `json:"username" validate:"required,min=6,max=500"`
		Password    string `json:"password" validate:"required,min=6,max=500"`
		PhoneNumber string `json:"phoneNumber" validate:"max=500"`
	}

	LoginInput struct {
		Username string `json:"username" validate:"required,min=6,max=500"`
		Password string `json:"password" validate:"required,min=6,max=500"`
	}

	RequestSellerLoginInput struct {
		Username string `json:"username" validate:"required,min=6,max=500"`
	}

	Session struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	AccessKey struct {
		UserID string `json:"userID"`
		Role   string `json:"role"` // owner/seller
		jwt.StandardClaims
	}

	RefreshKey struct {
		UserID string `json:"userID"`
		jwt.StandardClaims
	}
)
