package utils

import (
	"github.com/gin-gonic/gin"
)

type Query struct {
	Limit int    `form:"limit,default=2"`
	Page  int    `form:"page,default=1"`
	Sort  string `form:"sort,default=created_at asc"`
}

func GeneratePaginationFromRequest(c *gin.Context) (Query, error) {
	var q Query
	if err := c.Bind(&q); err != nil {
		return q, err
	}
	return q, nil
}
