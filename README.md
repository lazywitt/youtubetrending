# youtubetrending

# USAGE
- Host PGDB server in your local. Edit the credentials in the configs/pgdb-dev.yml for establishing a PGDB client. This is a dependency requirement for the project.
- Replace ApiKey field in configs/scraper-dev.yml with your youtubeV3 api key/s.
- Run server/server.go which is the main for this project. Build's working directory shoud be $GOPATH\src\github.com\lazywitt\youtubetrending
- Completing the above steps will trigger the youtube scraper in background and will expose an http server at localhost:4000

# Project Structure
GO project consists of 3 service packages -

* db 
  - dbService is dao level service interacting with PGDB to perform CRUD operations

* scraper 
  - this service interacts with the official youtube v3 api and stores the retreived data into the PGDB every 10 seconds

* fetch
  - fetchService is the top layer which exposes the two core api's which are Paginated response, search Video
  - this service also exposes an http server on top of service layer to expose 2 endpoints for providing REST interface.

- http://localhost:4000/videos/getpage

token field may also be left empty. The response json will provide with a new token. Expect empty token in case there are no pages left to serve.

REQUEST - 
{
  "token": "asf0faz"
}

- http://localhost:4000/videos/search

REQUEST - 
{
  "searchkey": "ronaldo shooting"
}

# sample via postman client

![image](https://user-images.githubusercontent.com/29565394/212305166-91a9f997-5ece-415f-b7ef-98570009ea7b.png)


# FEATURES

- Multiple api key support is implemented by combining multiple api keys together like this - "apiKey1, apiKey2"

- Text search is performed using to_tsvector queries with support for jumbled search, optimised via GIN index.
  - Example: searching for "new hat" will match "hat in new york" and "old hat and new hat" both. match token are being created with a combination of both title and description.

 

![youtubetrending (1)](https://user-images.githubusercontent.com/29565394/212297633-4c315b3a-cf9d-41c6-a191-3cd943b193a5.png)


 

TODO:
* dockerise
