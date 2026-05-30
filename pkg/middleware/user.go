package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func ExtractUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user header"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id header"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, int32(userID))

		c.Next()
	}
}
