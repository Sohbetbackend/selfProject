package store

import "github.com/Sohbetbackend/selfProject/internal/models"

type IStore interface {
	BookFindById(ID string) (*models.Book, error)
	BookFindByIds(Ids []string) ([]*models.Book, error)
	BooksFindBy(f models.BookFilterRequest) (books []models.Book, total int, err error)
	BookCreate(model *models.Book) (*models.Book, error)
	BookUpdate(model *models.Book) (*models.Book, error)
	BookDelete(items []*models.Book) ([]*models.Book, error)
}
