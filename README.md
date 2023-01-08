# youtubetrending

GO project consists of 3 services/package -

* fetch - fetchService is the top layer which exposes the two core api's which are Paginated response, search Video
* db - dbService is dao level service interacting with PGDB to perform CRUD operations
* scraper - this service interacts with the official youtube v3 api and stores the retreived data into the PGDB every 10 seconds
 
