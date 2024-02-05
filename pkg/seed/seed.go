package seed

import (
	"gorm-pagination/internal/models"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func CreateCategory(db *gorm.DB, name string) error {
	return db.Create(&models.Category{Name: name}).Error
}

func All() []Seed {
	return []Seed{
		{
			Name: "CreateCategory",
			Run: func(db *gorm.DB) error {
				f := faker.New()
				for i := 0; i < 20; i++ {
					CreateCategory(db, f.Car().Category())
				}

				return nil
			},
		},
	}
}
