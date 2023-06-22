package pkg

import (
	"ticket-api/internal/constants"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Paginate Gorm Scope
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := c.MustGet(constants.Limit).(int64)
		page := c.MustGet(constants.Page).(int64)

		offset := (page - 1) * limit

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
