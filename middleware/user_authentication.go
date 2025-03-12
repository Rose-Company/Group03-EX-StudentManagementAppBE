package middleware

import (

	"Group03-EX-StudentManagementAppBE/common"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your_secret_key")

func UserAuthentication(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": "Authorization header is required"})
		return
	}

	arr := strings.Split(authorization, "Bearer ")
	if len(arr) < 2 {
		c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": "Token is required"})
		return
	}

	tokenString := arr[1]

	if tokenString == "" {
		c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": "Invalid token"})
		return
	}

	var claims common.UserJWTProfile
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": err.Error()})
		return
	}
	if claims, ok := token.Claims.(*common.UserJWTProfile); ok && token.Valid {
		//ex := claims.Exp
		if claims.Exp > time.Now().Unix() {
			c.Set(common.USER_JWT_KEY, claims)
			// convert string to uuid
			claims.Id = strings.ToLower(claims.Id)
			c.Set(common.UserId, claims.Id)
			c.Next()
		} else {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": "Token expired"})
		}
	} else {
		c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, gin.H{"error": "Token unauthorized"})
	}
}

func OptionalUserAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization != "" {
			arr := strings.Split(authorization, "Bearer ")
			if len(arr) >= 2 {
				tokenString := arr[1]
				if tokenString != "" {
					var claims common.UserJWTProfile
					token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return secretKey, nil
					})

					if err == nil {
						if claims, ok := token.Claims.(*common.UserJWTProfile); ok && token.Valid {
							if claims.Exp > time.Now().Unix() {
								c.Set(common.USER_JWT_KEY, claims)
							}
						}
					}

				}

			}
		}
		c.Next()

	}
}
