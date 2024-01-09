package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-gorm-tutorial/initializers"
	"gin-gorm-tutorial/models"
	"gin-gorm-tutorial/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var body struct {
	Title string `json:"title" binding:"required"`
	//Color       models.Color `json:"color" binding:"required" validate:"oneof=Red Blue Green Yellow Purple Orange Pink"`
	Price       int    `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func ItemCreate(c *gin.Context) {
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	tx := initializers.DB.Begin()
	item := models.Item{Title: body.Title, Price: body.Price, Description: body.Description}
	if err := initializers.DB.Create(&item).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create item in the DB",
		})
		return
	}
	_, err := initializers.ES.Index("item").
		Request(item).
		Do(context.Background())
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create item in the ES",
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, item)
}

func ItemFindAll(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)
	itemLists, total, err := utils.GetAllESItems(&pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"BookList": itemLists,
		"MaxPage":  total / pagination.Limit,
		"Total":    total,
	})
}

func ItemFindOne(c *gin.Context) {
	var item models.Item
	id := c.Param("id")
	if err := initializers.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	r, err := json.Marshal(item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't marshal to json: {}",
		})
	}
	c.JSON(http.StatusOK, r)
}

func ItemUpdate(c *gin.Context) {
	id := c.Param("id")
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var item models.Item
	if err := initializers.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	if err := initializers.DB.Model(&item).
		Updates(&models.Item{Title: body.Title, Price: body.Price, Description: body.Description}).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func ItemDelete(c *gin.Context) {
	id := c.Param("id")
	if (initializers.DB.First(&models.Item{}, id).RowsAffected == 0) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	initializers.DB.Delete(&models.Item{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
