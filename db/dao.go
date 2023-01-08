package db

import (
	"context"
	"fmt"
	"gorm.io/gorm/clause"
	"strings"

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
	// Create inserts the given videos into entity
	Create(ctx context.Context, videos []*models.Videos) error
	// GetById fetches entry by Id
	GetById(ctx context.Context, id string) (*models.Videos, error)
	// GetByToken fetches paginated response wrt given token
	GetByToken(ctx context.Context, token string) ([]*models.Videos, error)
	// GetBySearchKey returns videos which have title or description matching with given searchKey
	// matching results is performed by matching every token in searchKey with description and title
	// Example: "new hat" will match "hat in new york" and "old hat and new hat" both. match token are being created
	// with a combination of both title and description. Query size limited by 10
	GetBySearchKey(ctx context.Context, searchKey string) ([]models.Videos, error)
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

func (h *Handler) GetBySearchKey(ctx context.Context, searchKey string) ([]models.Videos, error) {
	var (
		modelVideos []models.Videos
		gormRes     *gorm.DB
	)

	gormRes = h.DB.Raw("select * from videos where to_tsvector( title || description) @@ to_tsquery( ? ) limit 10;",
		strings.ReplaceAll(searchKey, " ", "&")).Scan(&modelVideos)

	if gormRes != nil {
		return modelVideos, gormRes.Error
	}
	return modelVideos, nil
}
