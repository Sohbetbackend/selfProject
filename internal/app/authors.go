package app

import (
	"errors"
	"strings"

	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/Sohbetbackend/selfProject/internal/store"
)

// AuthorsFindBy(f models.AuthorsFilterRequest) (authors []*models.Author, total int, err error)

func AuthorsList(f models.AuthorsFilterRequest) ([]*models.AuthorResponse, int, error) {
	authors, total, err := store.Store().AuthorsFindBy(f)
	if err != nil {
		return nil, 0, err
	}
	authorsResponse := []*models.AuthorResponse{}
	for _, author := range authors {
		a := models.AuthorResponse{}
		a.FromModel(author)
		authorsResponse = append(authorsResponse, &a)
	}
	return authorsResponse, total, nil
}

func AuthorDetail(f models.AuthorsFilterRequest) (*models.AuthorResponse, error) {
	m, _, err := store.Store().AuthorsFindBy(f)
	if err != nil {
		return nil, err
	}
	if len(m) < 1 {
		return nil, ErrNotFound
	}
	res := &models.AuthorResponse{}
	res.FromModel(m[0])
	return res, nil
}

func AuthorsUpdate(data models.AuthorRequest) (*models.AuthorResponse, error) {
	model := &models.Author{
		LastName:  data.LastName,
		FirstName: data.FirstName,
	}
	data.ToModel(model)

	var err error
	model, err = store.Store().AuthorsUpdate(model)
	if err != nil {
		return nil, err
	}
	res := &models.AuthorResponse{}
	res.FromModel(model)
	return res, nil
}

func AuthorsCreate(data models.AuthorRequest) (*models.AuthorResponse, error) {
	model := &models.Author{}
	data.ToModel(model)
	res := &models.AuthorResponse{}
	var err error
	model, err = store.Store().AuthorsCreate(model)
	if err != nil {
		return nil, err
	}
	res.FromModel(model)
	return res, nil
}

func AuthorsDelete(ids []string) ([]*models.AuthorResponse, error) {
	authors, err := store.Store().AuthorFindByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(authors) < 1 {
		return nil, errors.New("model not found: " + strings.Join(ids, ","))
	}
	authors, err = store.Store().AuthorsDelete(authors)
	if err != nil {
		return nil, err
	}
	if len(authors) == 0 {
		return make([]*models.AuthorResponse, 0), nil
	}
	var authorsResponse []*models.AuthorResponse
	for _, author := range authors {
		var a models.AuthorResponse
		a.FromModel(author)
		authorsResponse = append(authorsResponse, &a)
	}
	return authorsResponse, nil
}
