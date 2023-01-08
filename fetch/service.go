package fetch

import (
	"context"
	"fmt"
	"github.com/lazywitt/youtubetrending/db"
	"github.com/lazywitt/youtubetrending/db/models"
)

type Fetch struct {
	dbClient *db.Handler
}

// fetchService is the service layer which will call the db layer to fetch video data
type fetchService interface {
	// GetVideo api to get paginated response of stored videos in reverse chronological order
	GetVideo(ctx context.Context, token string) (*GetVideoResponse, error)
	// SearchVideo api to search for videos wrt given keyword. default size limited by 10
	SearchVideo(ctx context.Context, searchKey string) ([]Video, error)
}

func GetFetchService(ctx context.Context, handler *db.Handler) *Fetch {
	return &Fetch{
		dbClient: handler,
	}
}

func (s *Fetch) GetVideo(ctx context.Context, token string) (*GetVideoResponse, error) {
	var getVideoRes = &GetVideoResponse{}

	getByTokenRes, getByTokenErr := s.dbClient.GetByToken(ctx, token)
	if getByTokenErr != nil {
		return nil, fmt.Errorf("error while calling dao layer GetByToken: %v", getByTokenErr)
	}

	if len(getByTokenRes) == 11 {
		getVideoRes.NextToken = getByTokenRes[10].Id.String()
		getByTokenRes = getByTokenRes[0:10]
	}

	getVideoRes.Videos = VideoModelToFetchVideoObject(ctx, getByTokenRes)
	return getVideoRes, nil
}

func (s *Fetch) SearchVideo(ctx context.Context, searchKey string) ([]Video, error) {

	getBySearchKeyRes, getBySearchErr := s.dbClient.GetBySearchKey(ctx, searchKey)
	if getBySearchErr != nil {
		return nil, fmt.Errorf("error querying by searchKey in db %v", getBySearchErr)
	}

	return VideoModelToFetchVideoObject(ctx, getBySearchKeyRes), nil
}

func VideoModelToFetchVideoObject(ctx context.Context, modelVideos []models.Videos) []Video {
	var videos []Video
	for _, modelVid := range modelVideos {
		videos = append(videos, Video{
			Id:          modelVid.Id.String(),
			YoutubeId:   modelVid.YoutubeId,
			Title:       modelVid.Title,
			Description: modelVid.Description,
		})
	}
	return videos
}

// Video is api level response structure
type Video struct {
	Id          string
	YoutubeId   string
	Title       string
	Description string
}

type GetVideoResponse struct {
	Videos    []Video
	NextToken string
}
