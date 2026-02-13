package wfm

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

var MarketBaseUrl = "https://api.warframe.market"
var ItemsUrl = MarketBaseUrl + "/v2/items"

type ItemResponse struct {
	ApiVersion string     `json:"api_version"`
	Data       []ItemJson `json:"data"`
	Error      any        `json:"error,omitempty"`
}

type ItemJson struct {
	Id             string                   `json:"id"`
	Slug           string                   `json:"slug"`
	GameRef        string                   `json:"gameRef"`
	Tags           []string                 `json:"tags,omitzero"`
	SetRoot        *bool                    `json:"setRoot,omitempty"`
	SetParts       []string                 `json:"setParts,omitempty"`
	QuantityInSet  int32                    `json:"quantityInSet,omitempty"`
	Rarity         string                   `json:"rarity,omitempty"`
	BulkTradable   bool                     `json:"bulkTradable,omitempty"`
	Subtypes       []string                 `json:"subtypes,omitempty"`
	MaxRank        int32                    `json:"maxRank,omitempty"`
	MaxCharges     int32                    `json:"maxCharges,omitempty"`
	MaxAmberStars  int32                    `json:"maxAmberStars,omitempty"`
	MaxCyanStars   int32                    `json:"maxCyanStars,omitempty"`
	BaseEndo       int32                    `json:"baseEndo,omitempty"`
	EndoMultiplier float32                  `json:"endoMultiplier,omitempty"`
	Ducats         int32                    `json:"ducats,omitempty"`
	Vosfor         int32                    `json:"vosfor,omitempty"`
	ReqMasteryRank *int32                   `json:"reqMasteryRank,omitempty"`
	Vaulted        *bool                    `json:"vaulted,omitempty"`
	TradingTax     int32                    `json:"tradingTax,omitempty"`
	Tradable       *bool                    `json:"tradable,omitempty"`
	I18N           map[string]*ItemI18NJson `json:"i18n,omitempty"`
}

type ItemI18NJson struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	WikiLink    string `json:"wikiLink,omitempty"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
	SubIcon     string `json:"subIcon,omitempty"`
}

// FetchItems fetches items from the Warframe Market API.
func FetchItems() ([]ItemJson, error) {
	c := http.Client{}

	// Fetch items from the Warframe Market API
	resp, err := c.Get(ItemsUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch items: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limiting by retrying after a delay
	// If it fails 5 times, or has a different error, return an error
	failureCount := 0
	for resp.StatusCode == http.StatusTooManyRequests && failureCount <= 5 {
		failureCount++
		sleepDuration := math.Pow(2, float64(failureCount))
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		resp, err = c.Get(ItemsUrl)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch items on retry %d: %w", failureCount, err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch items, status code: %d", resp.StatusCode)
	}

	// Decode the JSON response from the API
	var itemResponse ItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&itemResponse); err != nil {
		return nil, fmt.Errorf("Failed to decode API response: %w", err)
	}

	return itemResponse.Data, nil
}
