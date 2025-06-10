package utils

import (
	"github.com/gin-gonic/gin"
)

func HandleUserID(c *gin.Context) *uint {
	userID, exists := c.Get("user_id")

	if !exists {
		return nil
	}

	userIDUint, ok := userID.(*uint)

	if !ok {
		return nil
	}

	return userIDUint
}

func IsOwner(resourceUserID *uint, currentUserID *uint) bool {
	return resourceUserID != nil && currentUserID != nil && *resourceUserID == *currentUserID
}
