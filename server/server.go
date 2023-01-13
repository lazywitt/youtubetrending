package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lazywitt/youtubetrending/db"
	"github.com/lazywitt/youtubetrending/db/models"
	"github.com/lazywitt/youtubetrending/fetch"
	"github.com/lazywitt/youtubetrending/scraper"
)

const (
	scraperConfigName = "scraper-dev"
	pgdbConfigName    = "pgdb-dev"
)

func main() {
	ctx := context.Background()

	pgdbConf := db.GetConf()
	scraperConf := scraper.GetConf()

	err := initDb(ctx, pgdbConf)
	if err != nil {
		fmt.Printf("error Initialising DB and entities: %v\n", err)
	}

	log.Println("start scraping")
	youtubeService := scraper.GetYoutubeSeed(scraperConf.SearchResource, scraperConf.SearchKey, scraperConf.OrderBy,
		scraperConf.MaxResultPerCall, scraperConf.ApiKey, time.Now().Add(-time.Hour), db.GetHandler(getDbClient(ctx, pgdbConf)))
	go youtubeService.InitialiseScraper(ctx)

	log.Println("start http server")
	initFetchHttpServer(ctx, pgdbConf)
}

func initFetchHttpServer(ctx context.Context, config *db.Config) {

	httpService := fetch.HttpService{FetchService: fetch.GetFetchService(ctx, db.GetHandler(getDbClient(ctx, config)))}

	mux := http.NewServeMux()
	mux.HandleFunc("/videos/search", httpService.VideoSearch)
	mux.HandleFunc("/videos/getpage", httpService.VideoPage)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// initDb Create entities in PGDB
func initDb(ctx context.Context, config *db.Config) error {
	dbClient := getDbClient(ctx, config)
	defer func() {
		sqlDB, _ := dbClient.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	err := dbClient.AutoMigrate(&models.Videos{})
	if err != nil {
		return err
	}
	// create indexes
	createIndexRes := dbClient.Exec(`CREATE INDEX IF NOT EXISTS  text_search_idx ON Videos USING GIN (to_tsvector('english', title || ' ' || description)); CREATE INDEX IF NOT EXISTS id_created_at_idx ON Videos (created_at, id);`)
	if createIndexRes != nil {
		return createIndexRes.Error
	}
	return nil
}

// getDbClient returns postgres client
func getDbClient(ctx context.Context, pgdbConf *db.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		pgdbConf.Host, pgdbConf.User, pgdbConf.Password, pgdbConf.Dbname, pgdbConf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to postgres", err)
	}
	fmt.Println("successfully connected to postgres")
	return db
}
