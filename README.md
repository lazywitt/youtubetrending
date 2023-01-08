# youtubetrending

GO project consists of 3 service packages -

* fetch - fetchService is the top layer which exposes the two core api's which are Paginated response, search Video
* db - dbService is dao level service interacting with PGDB to perform CRUD operations
* scraper - this service interacts with the official youtube v3 api and stores the retreived data into the PGDB every 10 seconds


Multiple api key support is implemented by combining multiple api keys together like this - "apiKey1, apiKey2"

Text search is performed using to_tsvector queries, optimised via GIN index


 


 
![youtubetrending](https://user-images.githubusercontent.com/29565394/211211358-554e197a-12c6-4540-bc10-9487893cc8da.png)
