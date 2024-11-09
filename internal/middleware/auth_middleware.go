package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/provider"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RequireAuth(prov *provider.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		if tokenStr == "" {
			fmt.Println("RequireAuth: Missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Missing Authorization Header"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("RequireAuth: Unexpected signing method: %v\n", token.Header["alg"])
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Unexpected signing method"})
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.ENV.JWTKey), nil
		})

		if err != nil {
			fmt.Printf("RequireAuth: Error parsing token: %v\n", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Error parsing token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Printf("RequireAuth: Token expired\n")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Token expired"})
				return
			}

			var user models.User
			err := prov.DB.Where("id = ?", claims["sub"]).First(&user).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("RequireAuth: error getting user by id: %v\n", err)
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Error getting user by id"})
					return
				}
				fmt.Printf("RequireAuth: Error finding user: %v\n", err)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Error finding user"})
				return
			}

			ctx.Set("user_id", user.ID)
		} else {
			fmt.Printf("RequireAuth: Invalid token claims\n")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Invalid token claims"})
		}
		ctx.Next()
	}
}
