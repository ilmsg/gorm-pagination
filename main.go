package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"gorm-pagination/internal/models"
	"gorm-pagination/pkg"
	"gorm-pagination/pkg/seed"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db := getConn()

	for _, s := range seed.All() {
		if err := s.Run(db); err != nil {
			log.Fatalf("Running seed '%s', failed with error\n", s.Name)
		}
	}

	// var categories []models.Category
	// if err := db.Model(&models.Category{}).Find(&categories).Error; err != nil {
	// 	log.Fatalf("find '%s', failed with error\n", err.Error())
	// }
	// fmt.Println("Categories:")
	// for _, category := range categories {
	// 	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// }

	limit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	page, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	cate := &CategoryGorm{db}
	pagging, err := cate.List(pkg.Pagination{Limit: limit, Page: page, Sort: "ID ASC"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Limt: %d\n", pagging.Limit)
	fmt.Printf("Page: %d\n", pagging.Page)
	fmt.Printf("Total Pages: %d\n", pagging.TotalPages)
	fmt.Printf("Total Rows: %d\n", pagging.TotalRows)

	fmt.Println("Categories:")
	for _, category := range pagging.Rows {
		fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	}
}

func getConn() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite db: %v", err)
	}

	db.Migrator().DropTable(&models.Category{})
	db.AutoMigrate(&models.Category{})

	return db
}

type CategoryGorm struct {
	db *gorm.DB
}

func (cg *CategoryGorm) List(pagination pkg.Pagination) (*pkg.Pagination, error) {
	var categories []*models.Category

	cg.db.Scopes(paginate(categories, &pagination, cg.db)).Find(&categories)
	pagination.Rows = categories

	return &pagination, nil
}

func paginate(value interface{}, pagination *pkg.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
