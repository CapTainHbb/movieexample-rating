package mysql

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/captainhbb/movieexample-rating/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type Repository struct {
	db *gorm.DB
}

func New() (*Repository, error) {
	user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&model.Rating{}); err != nil {
		log.Fatal("Migration failed:", err)
		return nil, nil
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	var ratings []model.Rating
	result := r.db.Where(&ratings, "record_id = ? AND record_type = ?", recordID, recordType)
	if result.Error != nil {
		log.Fatalf("failed to fetch record %s %s\n", recordID, recordType)
		return nil, result.Error
	}

	return ratings, nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	ratingToSave := model.Rating{
		RecordType: string(recordType),
		RecordID: string(recordID),
		UserID: rating.UserID,
		Value: rating.Value,
	}
	result := r.db.Create(&ratingToSave)
	return result.Error
}