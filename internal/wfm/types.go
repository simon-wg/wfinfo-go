package wfm

type ItemResponse struct {
	ApiVersion string        `json:"api_version"`
	Data       []interface{} `json:"data"`
	Error      any           `json:"error,omitempty"`
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
