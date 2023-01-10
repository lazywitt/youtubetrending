package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lazywitt/youtubetrending/db"
	"github.com/lazywitt/youtubetrending/db/models"
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

	youtubeService := scraper.GetYoutubeSeed(scraperConf.SearchResource, scraperConf.SearchKey, scraperConf.OrderBy,
		scraperConf.MaxResultPerCall, scraperConf.ApiKey, time.Now().Add(-time.Hour), db.GetHandler(getDbClient(ctx, pgdbConf)))

	youtubeService.InitialiseScraper(ctx)
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
	createIndexRes := dbClient.Exec(`CREATE INDEX IF NOT EXISTS  text_search_idx ON Videos USING GIN (to_tsvector('english', title || ' ' || description)); CREATE INDEX IF NOT EXISTS id_created_at_idx ON Videos (id, created_at);
	CREATE INDEX IF NOT EXISTS id_created_at_idx ON Videos (created_at, id);`)
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
