package server

import "github.com/gin-gonic/gin"

type ApiError = string

func createError(e ApiError) map[string]any {
	return gin.H{"error": e}
}
