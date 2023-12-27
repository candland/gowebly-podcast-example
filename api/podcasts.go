package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Posdcast is a struct to represent a podcast.
// id	"40325177-02e5-4ecf-aec5-fc7345d3b83c"
// title	"The Secrets Hotline"
// images
// default	"https://img4.luminarypodcasts.com/v1/podcast/40325177-02e5-4ecf-aec5-fc7345d3b83c/default/h_500,w_500/image/s--YOMzvZ-B--/aHR0cHM6Ly9pbWFnZXMubWVnYXBob25lLmZtL0d3bnMwYW1RbzFpNmhlMTFDbkVqUnZUbTUzTm9ndlBjWmhzVTZ4czdLbmcvcGxhaW4vczM6Ly9tZWdhcGhvbmUtcHJvZC9wb2RjYXN0cy9jOTUyNjRjOC02MjRmLTExZWEtOGE2ZC1iMzk1ZTU2YTI0YWEvaW1hZ2UvdXBsb2Fkc18yRjE1ODY1NTA1Mjc3NzQtaWNic3pzY2pyeC1hZWQ4NjgxMWI5MTMyZmY0NTdiYjcyOTM0NDRkZjM2ZF8yRlNlY3JldHNfSG90bGluZV9UaHVtYm5haWwuanBn.jpg"
// featured	"https://img4.luminarypodcasts.com/v1/podcast/40325177-02e5-4ecf-aec5-fc7345d3b83c/featured/h_430,w_720/image/s--TEXiyJIj--/aHR0cHM6Ly9pbWFnZXMubHVtaW5hcnkuYXVkaW8vaW1hZ2VzLXVzLWVhc3QtMi5wcm9kLmx1bWluYXJ5LmF1ZGlvLzE1ODc3NjE4MTcxMDQ4ODgzODEtU2VjcmV0c19Ib3RsaW5lX0ZlYXR1cmVkLnBuZw==.png"
// thumbnail	"https://img4.luminarypodcasts.com/v1/podcast/40325177-02e5-4ecf-aec5-fc7345d3b83c/thumbnail/h_360,w_360/image/s--yD3hPQA9--/aHR0cHM6Ly9pbWFnZXMubHVtaW5hcnkuYXVkaW8vaW1hZ2VzLXVzLWVhc3QtMi5wcm9kLmx1bWluYXJ5LmF1ZGlvLzE1ODc3NjE3OTgwMzk2MTE4NDYtU2VjcmV0c19Ib3RsaW5lX1RodW1ibmFpbC5wbmc=.png"
// wide	"https://img4.luminarypodcasts.com/v1/podcast/40325177-02e5-4ecf-aec5-fc7345d3b83c/wide/h_140,w_225/image/s--kmKs_M15--/aHR0cHM6Ly9pbWFnZXMubHVtaW5hcnkuYXVkaW8vaW1hZ2VzLXVzLWVhc3QtMi5wcm9kLmx1bWluYXJ5LmF1ZGlvLzE1ODc3NjE4MDc3NjQ3MTAzMjAtU2VjcmV0c19Ib3RsaW5lX1dpZGUucG5n.png"
// isExclusive	true
// publisherName	"Love and Radio"
// publisherId	"7c8226ef-de61-439a-9c51-9865d30ac2d9"
// mediaType	"podcast"
// description	"What’s your secret? Listen as anonymous callers from across the globe share what is hidden in their lives. It’s both a voyeuristic and affirming experience. \nLeave your own at 1-929-SECRETS (929-732-7387) or secretshotline.org. New episodes coming this November."
// categoryId	"4fb6e6ed-d6ef-4d04-878b-40942df8d241"
// categoryName	"Performing Arts"
// hasFreeEpisodes	true
// playSequence	"episodic"
type Podcast struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Images struct {
		Thumbnail string `json:"thumbnail"`
	} `json:"images"`
	PublisherName string `json:"publisherName"`
	Description   string `json:"description"`
	CategoryName  string `json:"categoryName"`
	PlaySequence  string `json:"playSequence"`
}

type Podcasts struct {
	Podcasts []Podcast `json:"podcasts"`
	Page     string    `json:"page"`
	Search   string    `json:"search"`
	Status   int       `json:"status"`
}

func (p Podcasts) GetPrevPage() string {
	if pageInt, err := strconv.Atoi(p.Page); err == nil {
		if pageInt > 1 {
			return strconv.Itoa(pageInt - 1)
		}
	}
	return ""
}

func (p Podcasts) GetNextPage() string {
	if len(p.Podcasts) == 10 {
		if pageInt, err := strconv.Atoi(p.Page); err == nil {
			return strconv.Itoa(pageInt + 1)
		}
	}
	return ""
}

// GET /api/podcasts
// get a list of podcasts from https://601f1754b5a0e9001706a292.mockapi.io/podcasts and unmarshals the JSON response into a slice of Podcast structs
func GetPodcasts(search string, page string) (Podcasts, error) {
	result := Podcasts{}
	// Create a new HTTP client.
	client := &http.Client{}
	// Create a new HTTP request.
	req, err := http.NewRequest("GET", "https://601f1754b5a0e9001706a292.mockapi.io/podcasts", nil)
	if err != nil {
		return result, err
	}
	q := req.URL.Query()
	if search != "" {
		q.Add("search", search)
	}
	if page == "" {
		page = "1"
	}
	q.Add("limit", "10")
	q.Add("page", page)
	req.URL.RawQuery = q.Encode()

	result.Page = page
	result.Search = search

	fmt.Println(req.URL.String())

	// Set the HTTP request header to accept JSON.
	req.Header.Set("Accept", "application/json")
	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	// Close the HTTP response body.
	defer resp.Body.Close()

	result.Status = resp.StatusCode
	if resp.StatusCode == 200 {
		// Decode the JSON response into a slice of Podcast structs.
		var podcasts []Podcast

		if err := json.NewDecoder(resp.Body).Decode(&podcasts); err != nil {
			fmt.Println(resp.Body)
			return result, err
		}
		// Return the slice of Podcast structs.
		result.Podcasts = podcasts
	} else {
		result.Podcasts = []Podcast{}
	}
	return result, nil
}
