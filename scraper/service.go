package scraper

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazywitt/youtubetrending/db"
	"github.com/lazywitt/youtubetrending/db/models"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"
	"time"
)

type YoutubeSeed struct { //
	lastSearchTime   time.Time
	searchResource   []string
	searchKey        string
	orderBy          string
	maxResultPerCall int
	apiKey           []string
	dbClient         *db.Handler
}

type YoutubeApiService interface {
	// InitialiseScraper orchestrates periodic call to scrape YouTube data after a fixed interval
	InitialiseScraper(ctx context.Context) error
	// MakeApiCall executes the api call to youtubeV3 api
	MakeApiCall(ctx context.Context, service *youtube.SearchService) (*youtube.SearchListResponse, error)
	// PersistVideoData to save the video meta data in postgres
	PersistVideoData(ctx context.Context, searchRes *youtube.SearchListResponse) error
}

func GetYoutubeSeed(SearchResource []string, SearchKey string, OrderBy string, MaxResultPerCall int, ApiKey []string, LastSearchTime time.Time,
	dbClient *db.Handler) *YoutubeSeed {
	return &YoutubeSeed{
		searchResource:   SearchResource,
		searchKey:        SearchKey,
		orderBy:          OrderBy,
		maxResultPerCall: MaxResultPerCall,
		apiKey:           ApiKey,
		lastSearchTime:   LastSearchTime,
		dbClient:         dbClient,
	}
}

func (s *YoutubeSeed) MakeApiCall(ctx context.Context,
	youtubeSearchService *youtube.SearchService) (*youtube.SearchListResponse, error) {

	res := youtubeSearchService.List(s.searchResource)
	res.Q(s.searchKey)
	res.Order(s.orderBy)
	res.MaxResults(int64(s.maxResultPerCall))
	res.PublishedAfter(s.lastSearchTime.Format(time.RFC3339))
	res.Do()
	searchListRes, err := res.Do()
	if err != nil {
		log.Fatal("error fetching search results", err)
	}
	if searchListRes == nil {
		log.Fatal("nil pointer returned from search api")
	}
	return searchListRes, err
}

func (s *YoutubeSeed) InitialiseScraper(ctx context.Context) error {

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(strings.Join(s.apiKey[:], ",")))
	if err != nil {
		log.Fatal("error creating youtube service object", err)
	}
	youtubeSearchService := youtube.NewSearchService(youtubeService)

	for {
		searchRes, err := s.MakeApiCall(ctx, youtubeSearchService)
		if err != nil {
			fmt.Printf("error when executing api call %v", err)
			continue
		}

		err = s.PersistVideoData(ctx, searchRes)
		s.lastSearchTime = time.Now()

		time.Sleep(time.Second * 10) // wait for 10 seconds before triggering next api call
	}
	return nil
}

func (s *YoutubeSeed) PersistVideoData(ctx context.Context, searchRes *youtube.SearchListResponse) error {

	var scrapedVideos []*models.Videos
	for _, it := range searchRes.Items {

		scrapedVideos = append(scrapedVideos, &models.Videos{
			Id:          uuid.New(),
			YoutubeId:   it.Id.VideoId,
			Title:       it.Snippet.Title,
			Description: it.Snippet.Description,
		})
	}

	err := s.dbClient.Create(ctx, scrapedVideos)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error creating entry of scraper data: %v", err))
	}
	return nil
}
