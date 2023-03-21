package middleware

import (
	"net/http"
	h "nextlaundry_apis/helper"

	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, claim, err := h.ValidateToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		if _, ok := h.BlacklistedTokens[token]; ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			return
		}

		if claim.Role != "admin" && claim.Role == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden access"})
			return
		}
		ctx.Set("UID", claim.UID)
		ctx.Next()
	}
}

func Cashier() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, claim, err := h.ValidateToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		if _, ok := h.BlacklistedTokens[token]; ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			return
		}

		if claim.Role != "kasir" && claim.Role != "admin" && claim.Role == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden access"})
			return
		}
		ctx.Set("UID", claim.UID)
		ctx.Next()
	}
}

func Owner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, claim, err := h.ValidateToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		if _, ok := h.BlacklistedTokens[token]; ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			return
		}

		if claim.Role != "owner" && claim.Role != "kasir" && claim.Role != "admin" && claim.Role == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden access"})
			return
		}
		ctx.Set("UID", claim.UID)
		ctx.Next()
	}
}
