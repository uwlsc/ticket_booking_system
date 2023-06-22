package middlewares

import (
	"strconv"

	"ticket-api/internal/constants"
	"ticket-api/pkg"

	"github.com/gin-gonic/gin"
)

type PaginationMiddleware struct {
	logger pkg.Logger
}

func NewPaginationMiddleware(logger pkg.Logger) PaginationMiddleware {
	return PaginationMiddleware{logger: logger}
}

func (p PaginationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		perPage, err := strconv.ParseInt(c.Query("per_page"), 10, 0)
		if err != nil {
			perPage = 10
		}

		page, err := strconv.ParseInt(c.Query("page"), 10, 0)
		if err != nil {
			page = 0
		}

		c.Set(constants.Limit, perPage)
		c.Set(constants.Page, page)

		c.Next()
	}
}
