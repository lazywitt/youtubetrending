package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpService struct {
	FetchService *Fetch
}

type VideoPageRequest struct {
	Token string `json:"token"`
}

func (s *HttpService) VideoPage(w http.ResponseWriter, r *http.Request) {
	v := VideoPageRequest{}

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getVideoRes, err := s.FetchService.GetVideo(context.Background(), v.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "videos: %v", getVideoRes)
}

type VideoSearchRequest struct {
	SearchKey string `json:"searchkey"`
}

func (s *HttpService) VideoSearch(w http.ResponseWriter, r *http.Request) {
	v := VideoSearchRequest{}

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getVideoRes, err := s.FetchService.SearchVideo(context.Background(), v.SearchKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "videos: %v", getVideoRes)

}
