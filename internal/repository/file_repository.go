package repository

import "project/internal/models"

type FileRepository interface {
	Save(file models.File) error
	List() ([]models.File, error)
}
