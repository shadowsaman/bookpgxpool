package storage

import (
	"app/models"
	"context"
)

type StorageI interface {
	CloseDB()
	Book() BookRepoI
	Category() CategoryRepoI
}

type BookRepoI interface {
	Insert(context.Context, *models.CreateBook) (string, error)
	GetByID(context.Context, *models.BookPrimeryKey) (*models.Book, error)
	GetList(context.Context, *models.GetListBookRequest) (*models.GetListBookResponse, error)
	Update(context.Context, *models.UpdateBook) error
	Delete(context.Context, *models.BookPrimeryKey) error
}

type CategoryRepoI interface {
	Insert(context.Context, *models.CreateCategory) (string, error)
	GetByID(context.Context, *models.CategoryPrimeryKey) (*models.Category, error)
	GetList(context.Context, *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error)
	Update(context.Context, *models.UpdateCategory) error
	Delete(context.Context, *models.CategoryPrimeryKey) error
}
