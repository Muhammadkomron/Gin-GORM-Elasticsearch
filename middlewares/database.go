package middlewares

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func ElasticsearchMiddleware(es *elasticsearch.TypedClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("es", es)
		c.Next()
	}
}
