package db

import (
	"context"
	"fmt"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"

	"github.com/lazywitt/youtubetrending/db/models"
)

type Handler struct {
	DB *gorm.DB
}

func GetHandler(dbClient *gorm.DB) *Handler {
	return &Handler{
		DB: dbClient,
	}
}

type dbService interface {
	// Create inserts the given video data into PGBD
	Create(ctx context.Context, videos []*models.Videos) error
	// GetById fetches entry by Id
	GetById(ctx context.Context, id string) (*models.Videos, error)
	// GetByToken fetches paginated response wrt given token
	GetByToken(ctx context.Context, token string) ([]*models.Videos, error)
}

func (h *Handler) Create(ctx context.Context, videos []*models.Videos) error {

	res := h.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(videos)
	if res.Error != nil {
		return fmt.Errorf("error creating entries in DB %v", res.Error)
	}
	return nil
}

func (h *Handler) GetById(ctx context.Context, id string) (*models.Videos, error) {
	var (
		getByTokenRes = &models.Videos{}
	)
	gormRes := h.DB.Where("Id = ?", id).First(getByTokenRes)
	return getByTokenRes, gormRes.Error
}

func (h *Handler) GetByToken(ctx context.Context, token string) ([]models.Videos, error) {
	var (
		modelVideos []models.Videos
		gormRes     *gorm.DB
	)

	if token != "" {
		tokenRes, tokenErr := h.GetById(ctx, token)
		if tokenErr != nil {
			return modelVideos, fmt.Errorf("error fetching details of token: %v", tokenErr)
		}
		fmt.Println(tokenRes)
		gormRes = h.DB.Where("((created_at = ? and Id >= ?) or (created_at < ?))", tokenRes.CreatedAt,
			tokenRes.Id, tokenRes.CreatedAt).Order("created_at DESC").Order("Id").Limit(11).Find(&modelVideos)
	} else {
		gormRes = h.DB.Order("created_at DESC").Order("Id").Limit(11).Find(&modelVideos)
	}

	if gormRes != nil {
		return modelVideos, gormRes.Error
	}
	return modelVideos, nil
}
