package app

import (
	"errors"
	"strings"

	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/Sohbetbackend/selfProject/internal/store"
)

func BooksList(f models.BookFilterRequest) ([]*models.BookResponse, int, error) {
	books, total, err := store.Store().BooksFindBy(f)
	if err != nil {
		return nil, 0, err
	}
	booksResponse := []*models.BookResponse{}
	for _, book := range books {
		b := models.BookResponse{}
		b.FromModel(&book)
		booksResponse = append(booksResponse, &b)
	}
	return booksResponse, total, err
}

func BookUpdate(data models.BookRequest) (*models.BookResponse, error) {
	model := &models.Book{}
	data.ToModel(model)

	var err error
	model, err = store.Store().BookUpdate(model)
	if err != nil {
		return nil, err
	}
	res := &models.BookResponse{}
	res.FromModel(model)
	return res, nil
}

func BookCreate(data models.BookRequest) (*models.BookResponse, error) {
	model := &models.Book{}
	data.ToModel(model)
	res := &models.BookResponse{}
	var err error
	model, err = store.Store().BookCreate(model)
	if err != nil {
		return nil, err
	}
	res.FromModel(model)
	return res, nil
}

func BookDelete(ids []string) ([]*models.BookResponse, error) {
	books, err := store.Store().BookFindByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(books) < 1 {
		return nil, errors.New("model not found: " + strings.Join(ids, ","))
	}
	books, err = store.Store().BookDelete(books)
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return make([]*models.BookResponse, 0), nil
	}
	var booksResponse []*models.BookResponse
	for _, book := range books {
		var b models.BookResponse
		b.FromModel(book)
		booksResponse = append(booksResponse, &b)
	}
	return booksResponse, nil
}
