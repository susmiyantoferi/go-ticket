package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"ticket/utils"
	"ticket/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "unauthorized", nil)
			ctx.Abort()
			return
		}

		tokenPart := strings.Split(authHeader, " ")
		if len(tokenPart) != 2 || tokenPart[0] != "Bearer" {
			web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenStr := tokenPart[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		claim := &utils.TokenClaim{}

		token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "invalid token", nil)
			ctx.Abort()
			return
		}

		ctx.Set("user", claim)
		ctx.Next()
	}

}

func RoleAccessMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userClaims, exist := ctx.Get("user")
		if !exist {
			web.ResponseJSON(ctx, http.StatusUnauthorized, "error", "user not found", nil)
			ctx.Abort()
			return
		}

		user := userClaims.(*utils.TokenClaim)
		role := user.Role

		for _, v := range allowedRoles {
			if role == v {
				ctx.Next()
				return
			}
		}

		web.ResponseJSON(ctx, http.StatusForbidden, "error", "no permission", nil)
		ctx.Abort()
	}
}
