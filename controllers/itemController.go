package controllers

import (
	"context"
	"fmt"
	"gin-gorm-tutorial/models"
	"gin-gorm-tutorial/utils"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
)

type itemBody struct {
	//Color       models.Color `json:"color" binding:"required" validate:"oneof=Red Blue Green Yellow Purple Orange Pink"`
	Title       string `json:"title" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func ItemCreate(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	es, _ := c.Value("es").(*elasticsearch.TypedClient)
	var b itemBody
	if c.Bind(&b) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	var item models.Item
	db.First(&item, "title = ?", b.Title)
	if item.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item with this title already exists"})
		return
	}
	tx := db.Begin()
	item = models.Item{Title: b.Title, Price: b.Price, Description: b.Description}
	// Creating Item in DB
	if err := db.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create item in the DB"})
		return
	}
	// Creating Item in ES
	if _, err := es.Index("item").
		Id(strconv.Itoa(int(item.ID))).
		Request(item).
		Do(context.Background()); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create item in the ES"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, item)
}

func ItemFindAll(c *gin.Context) {
	//db, _ := c.Value("db").(*gorm.DB)
	es, _ := c.Value("es").(*elasticsearch.TypedClient)
	pagination, err := utils.GeneratePaginationFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read query"})
		return
	}
	// Fetching paginated Item list from DB
	//itemList, total, err := utils.GetItemList(db, &pagination)
	// Fetching paginated Item list from ES
	itemList, total, err := utils.GetESItemList(es, &pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	maxPage := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	c.JSON(http.StatusOK, gin.H{"BookList": itemList, "MaxPage": maxPage, "Total": total})
}

func ItemFindOne(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	id := c.Param("id")
	// Fetching Item from ES
	item, err := utils.GetOneItem(db, id)
	fmt.Println(item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, item)
}

func ItemUpdate(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	es, _ := c.Value("es").(*elasticsearch.TypedClient)
	id := c.Param("id")
	var b itemBody
	if c.Bind(&b) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	var item *models.Item
	db.Where("id != ?", id).First(&item, "title = ?", b.Title)
	if item.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item with this title already exists"})
		return
	}
	tx := db.Begin()
	// Fetching Item from DB
	item, err := utils.GetOneItem(db, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Updating Item in DB
	if err := db.Model(&item).
		Updates(&models.Item{Title: b.Title, Price: b.Price, Description: b.Description}).
		Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update item in the DB"})
		return
	}
	// Updating Item in ES
	if _, err := es.Update(models.ItemIndex, id).Doc(b).Do(context.Background()); err != nil {
		tx.Rollback()
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update item in the ES"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, item)
}

func ItemDelete(c *gin.Context) {
	db, _ := c.Value("db").(*gorm.DB)
	es, _ := c.Value("es").(*elasticsearch.TypedClient)
	id := c.Param("id")
	tx := db.Begin()
	if (db.First(&models.Item{}, id).RowsAffected == 0) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	// Deleting Item from DB
	if err := db.Delete(&models.Item{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Deleting Item from ES
	if _, err := es.Delete(models.ItemIndex, id).Do(context.Background()); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
