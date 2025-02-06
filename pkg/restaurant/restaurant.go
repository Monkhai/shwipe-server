package restaurant

import "fmt"

type Restaurant struct {
	Name       string   `json:"name"`
	Rating     float64  `json:"rating"`
	PriceLevel int      `json:"price_level"`
	Photos     []string `json:"photos"`
}

func GetRestaurantPhotoURL(photoReference string) string {
	const API_KEY = "AIzaSyALFutkrFeGGS8jR_HVgO1xUqrlJ-_ZZm4"
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", photoReference, API_KEY)
}

type GetRestaurantResponse struct {
	Results          []RestaurantResponse `json:"results"`
	HtmlAttributions []string             `json:"html_attributions"`
	NextPageToken    string               `json:"next_page_token"`
}

type RestaurantResponse struct {
	BusinessStatus string `json:"business_status"`
	Geometry       struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		Viewport struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"viewport"`
	} `json:"geometry"`
	Icon                string `json:"icon"`
	IconBackgroundColor string `json:"icon_background_color"`
	IconMaskBaseURI     string `json:"icon_mask_base_uri"`
	Name                string `json:"name"`
	OpeningHours        struct {
		OpenNow bool `json:"open_now"`
	} `json:"opening_hours"`
	Photos []struct {
		Height           int      `json:"height"`
		HTMLAttributions []string `json:"html_attributions"`
		PhotoReference   string   `json:"photo_reference"`
		Width            int      `json:"width"`
	} `json:"photos"`
	PlaceID  string `json:"place_id"`
	PlusCode struct {
		CompoundCode string `json:"compound_code"`
		GlobalCode   string `json:"global_code"`
	} `json:"plus_code"`
	PriceLevel       int      `json:"price_level"`
	Rating           float64  `json:"rating"`
	Reference        string   `json:"reference"`
	Scope            string   `json:"scope"`
	Types            []string `json:"types"`
	UserRatingsTotal int      `json:"user_ratings_total"`
	Vicinity         string   `json:"vicinity"`
}
