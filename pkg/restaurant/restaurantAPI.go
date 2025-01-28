package restaurant

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Monkhai/shwipe-server.git/secrets"
)

type RestaurantAPI struct {
	baseURL string
	apiKey  string
}

func NewRestaurantAPI(baseURL string, apiKey string) *RestaurantAPI {
	return &RestaurantAPI{baseURL: baseURL, apiKey: apiKey}
}

/*
GetResaturants retrieves a list of restaurants near the specified coordinates.

Parameters:

	lat: latitude of the location
	lng: longitude of the location
	nextPageTokenPtr: optional token for pagination (nil for first page)

Returns:

	[]Restaurant: list of restaurants near the specified location
	string: token for the next page of results, empty if no more results
*/
func (r *RestaurantAPI) GetResaturants(lat, lng string, nextPageTokenPtr *string) ([]Restaurant, *string, error) {
	paramsObj := map[string]string{
		"location": fmt.Sprintf("%s,%s", lat, lng),
		"rankby":   "distance",
		"type":     "restaurant",
		"key":      r.apiKey,
		"keyword":  "All",
	}

	if nextPageTokenPtr != nil {
		log.Println("Using next page token", *nextPageTokenPtr)
		paramsObj["pagetoken"] = *nextPageTokenPtr
	}

	query := url.Values{}
	for key, value := range paramsObj {
		query.Add(key, value)
	}

	queryString := query.Encode()

	requestUrl := fmt.Sprintf("%s?%s", secrets.BASE_URL, queryString)

	response, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("API request failed with status code: %d", response.StatusCode)
		return nil, nil, errors.New("API request failed")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObj GetRestaurantResponse
	err = json.Unmarshal(body, &responseObj)
	if err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return nil, nil, errors.New("JSON unmarshal error")
	}

	var restaurants []Restaurant
	for _, rest := range responseObj.Results {
		photos := []string{}

		for _, photo := range rest.Photos {
			photos = append(photos, photo.PhotoReference)
		}

		restaurants = append(restaurants, Restaurant{
			Name:       rest.Name,
			Rating:     rest.Rating,
			PriceLevel: rest.PriceLevel,
			Photos:     photos,
		})
	}

	return restaurants, &responseObj.NextPageToken, nil
}
