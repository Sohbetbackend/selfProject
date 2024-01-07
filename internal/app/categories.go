package app

import (
	"errors"
	"strings"

	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/Sohbetbackend/selfProject/internal/store"
)

func CategoryList(f models.CategoryFilterRequest) ([]*models.CategoryResponse, int, error) {
	authors, total, err := store.Store().CategoryFindBy(f)
	if err != nil {
		return nil, 0, err
	}
	categoriesResponse := []*models.CategoryResponse{}
	for _, category := range authors {
		c := models.CategoryResponse{}
		c.FromModel(&category)
		categoriesResponse = append(categoriesResponse, &c)
	}
	return categoriesResponse, total, nil
}

func CategoryUpdate(data models.CategoryRequest) (*models.CategoryResponse, error) {
	model := &models.Category{}
	data.ToModel(model)

	var err error
	model, err = store.Store().CategoryUpdate(model)
	if err != nil {
		return nil, err
	}
	res := &models.CategoryResponse{}
	res.FromModel(model)
	return res, nil
}

func CategoryCreate(data models.CategoryRequest) (*models.CategoryResponse, error) {
	model := &models.Category{}
	data.ToModel(model)
	res := &models.CategoryResponse{}
	var err error
	model, err = store.Store().CategoryCreate(model)
	if err != nil {
		return nil, err
	}
	res.FromModel(model)
	return res, nil
}

func CategoryDelete(ids []string) ([]*models.CategoryResponse, error) {
	categories, err := store.Store().CategoryFindByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(categories) < 1 {
		return nil, errors.New("model not found: " + strings.Join(ids, ","))
	}
	categories, err = store.Store().CategoryDelete(categories)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return make([]*models.CategoryResponse, 0), nil
	}
	var categoriesResponse []*models.CategoryResponse
	for _, author := range categories {
		var ct models.CategoryResponse
		ct.FromModel(author)
		categoriesResponse = append(categoriesResponse, &ct)
	}
	return categoriesResponse, nil
}
