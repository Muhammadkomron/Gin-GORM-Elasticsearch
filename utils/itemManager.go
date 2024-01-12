package utils

import (
	"context"
	"encoding/json"
	"gin-gorm-tutorial/models"
	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
)

func GetItemList(db *gorm.DB, pagination *Query) (*[]models.Item, int, error) {
	var items []models.Item
	var total int64
	offset := (pagination.Page - 1) * pagination.Limit
	q := db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	r := q.Model(&models.Item{}).Find(&items).Count(&total)
	if r.Error != nil {
		return nil, 0, r.Error
	}
	return &items, int(total), nil
}

func GetOneItem(db *gorm.DB, id string) (*models.Item, error) {
	var item models.Item
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func GetESItemList(es *elasticsearch.TypedClient, pagination *Query) (*[]models.Item, int, error) {
	var items []models.Item
	offset := (pagination.Page - 1) * pagination.Limit
	r, err := es.Search().
		Index(models.ItemIndex).
		Sort(map[string]string{"CreatedAt": "asc"}).
		From(offset).
		Size(pagination.Limit).
		Do(context.Background())
	if err != nil {
		return nil, 0, err
	}
	for _, h := range r.Hits.Hits {
		var i models.Item
		json.Unmarshal(*&h.Source_, &i)
		items = append(items, i)
	}
	return &items, int(r.Hits.Total.Value), nil
}
