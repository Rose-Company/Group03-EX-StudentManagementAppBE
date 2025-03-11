package common

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserJWTProfile struct {
	jwt.RegisteredClaims
	Id          string `json:"id"`
	Role        string `json:"role"`
	AppAccess   bool   `json:"app_access"`
	AdminAccess bool   `json:"admin_access"`
	Iat         int    `json:"iat"`
	Exp         int64  `json:"exp"`
	Iss         string `json:"iss"`
}

type JWTCustomClaims struct {
	UID string `json:"id"`
	jwt.RegisteredClaims
}

type JWT struct {
	claims *JWTCustomClaims
}

func NewJWT(c *gin.Context, tokenKey string) *JWT {
	rawToken, ok := c.Get(tokenKey)
	if !ok {
		return nil
	}
	jwtToken, ok := rawToken.(*jwt.Token)
	if !ok {
		return nil
	}
	claims, ok := jwtToken.Claims.(*JWTCustomClaims)
	if !ok {
		return nil
	}
	return &JWT{claims: claims}
}

func (j *JWT) GetUID() string {
	return j.claims.UID
}

func ProfileFromJwt(c *gin.Context) (bool, *UserJWTProfile) {
	value, ok := c.Get(USER_JWT_KEY)
	if !ok {
		return false, nil
	}

	userJWTProfile, ok := value.(*UserJWTProfile)
	if !ok {

		return false, nil
	}
	return ok, userJWTProfile

}
