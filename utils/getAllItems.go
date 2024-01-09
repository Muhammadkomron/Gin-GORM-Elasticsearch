package utils

import (
	"context"
	"encoding/json"
	"gin-gorm-tutorial/initializers"
	"gin-gorm-tutorial/models"
)

func GetAllItems(pagination *Pagination) (*[]models.Item, error) {
	var items []models.Item
	offset := (pagination.Page - 1) * pagination.Limit
	q := initializers.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	r := q.Model(&models.Item{}).Find(&items)
	if r.Error != nil {
		return nil, r.Error
	}
	return &items, nil
}

func GetAllESItems(pagination *Pagination) (*[]models.Item, int, error) {
	//res, err := initializers.ES.Search().
	//	Index("index_name").
	//	Request(&search.Request{
	//		Query: &types.Query{
	//			Match: map[string]types.MatchQuery{
	//				"name": {Query: "Foo"},
	//			},
	//		},
	//	}).Do(context.Background())
	var items []models.Item
	offset := (pagination.Page - 1) * pagination.Limit
	r, err := initializers.ES.Search().
		Index("item").
		Sort(map[string]string{"CreatedAt": "asc"}).
		From(offset).
		Size(pagination.Limit).
		Do(context.Background())
	if err != nil {
		return nil, 0, err
	}
	for _, h := range r.Hits.Hits {
		s := (*json.RawMessage)(&h.Source_)
		var i models.Item
		err := json.Unmarshal(*s, &i)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}
	return &items, int(r.Hits.Total.Value), nil
}
