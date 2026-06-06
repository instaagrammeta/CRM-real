package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// JSONError shorthand.
func JSONError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg, "success": false})
}

// JSONOK shorthand.
func JSONOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// ParseUintParam — parse :id like params.
func ParseUintParam(c *gin.Context, name string) (uint, bool) {
	v := c.Param(name)
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid "+name)
		return 0, false
	}
	return uint(n), true
}

// QueryInt — parse query param как int с дефолтом.
func QueryInt(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
