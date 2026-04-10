package usecase

import (
	"project/internal/infrastructure"
	"project/internal/models"
	"project/internal/repository"

	"github.com/google/uuid"
)

type FileUsecase struct {
	repo repository.FileRepository
}

func NewFileUsecase(r repository.FileRepository) *FileUsecase {
	return &FileUsecase{repo: r}
}

func (u *FileUsecase) SaveFile(name string, data []byte) error {

	// save file locally
	url, err := infrastructure.SaveFileLocally(name, data)
	if err != nil {
		return err
	}

	// Generate unique ID for each file using UUID
	file := models.File{
		ID:   uuid.New().String(),
		Name: name,
		URL:  url,
	}

	return u.repo.Save(file)
}

func (u *FileUsecase) ListFiles() ([]models.File, error) {
	return u.repo.List()
}
