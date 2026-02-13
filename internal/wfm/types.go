package wfm

type Item struct {
	Id             string               `json:"id"`
	Slug           string               `json:"slug"`
	GameRef        string               `json:"gameRef"`
	Tags           []string             `json:"tags,omitzero"`
	SetRoot        *bool                `json:"setRoot,omitempty"`
	SetParts       []string             `json:"setParts,omitempty"`
	QuantityInSet  int32                `json:"quantityInSet,omitempty"`
	Rarity         string               `json:"rarity,omitempty"`
	BulkTradable   bool                 `json:"bulkTradable,omitempty"`
	Subtypes       []string             `json:"subtypes,omitempty"`
	MaxRank        int32                `json:"maxRank,omitempty"`
	MaxCharges     int32                `json:"maxCharges,omitempty"`
	MaxAmberStars  int32                `json:"maxAmberStars,omitempty"`
	MaxCyanStars   int32                `json:"maxCyanStars,omitempty"`
	BaseEndo       int32                `json:"baseEndo,omitempty"`
	EndoMultiplier float32              `json:"endoMultiplier,omitempty"`
	Ducats         int32                `json:"ducats,omitempty"`
	Vosfor         int32                `json:"vosfor,omitempty"`
	ReqMasteryRank *int32               `json:"reqMasteryRank,omitempty"`
	Vaulted        *bool                `json:"vaulted,omitempty"`
	TradingTax     int32                `json:"tradingTax,omitempty"`
	Tradable       *bool                `json:"tradable,omitempty"`
	I18N           map[string]*ItemI18N `json:"i18n,omitempty"`
}

type ItemI18N struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	WikiLink    string `json:"wikiLink,omitempty"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
	SubIcon     string `json:"subIcon,omitempty"`
}

type Riven struct {
	Id             string                `json:"id"`
	Slug           string                `json:"slug"`
	GameRef        string                `json:"gameRef"`
	Group          string                `json:"group,omitempty"`
	RivenType      string                `json:"rivenType,omitempty"`
	Disposition    float64               `json:"disposition"`
	ReqMasteryRank int8                  `json:"reqMasteryRank"`
	I18N           map[string]*RivenI18N `json:"i18n,omitempty"`
}

type RivenI18N struct {
	Name     string `json:"itemName,omitempty"`
	WikiLink string `json:"wikiLink,omitempty"`
	Icon     string `json:"icon"`
	Thumb    string `json:"thumb"`
}

type RivenAttribute struct {
	Id                 string                         `json:"id"`
	Slug               string                         `json:"slug"`
	GameRef            string                         `json:"gameRef"`
	Group              string                         `json:"group,omitempty"`
	Prefix             string                         `json:"prefix"`
	Suffix             string                         `json:"suffix"`
	ExclusiveTo        []string                       `json:"exclusiveTo,omitempty"`
	PositiveIsNegative bool                           `json:"positiveIsNegative,omitempty"`
	Unit               string                         `json:"unit,omitempty"`
	PositiveOnly       bool                           `json:"positiveOnly,omitempty"`
	NegativeOnly       bool                           `json:"negativeOnly,omitempty"`
	I18N               map[string]*RivenAttributeI18N `json:"i18n,omitempty"`
}

type RivenAttributeI18N struct {
	Name  string `json:"effect"`
	Icon  string `json:"icon"`
	Thumb string `json:"thumb"`
}

type LichWeapon struct {
	Id             string                     `json:"id"`
	Slug           string                     `json:"slug"`
	GameRef        string                     `json:"gameRef"`
	ReqMasteryRank int8                       `json:"reqMasteryRank"`
	I18N           map[string]*LichWeaponI18N `json:"i18n,omitempty"`
}

type LichWeaponI18N struct {
	Name     string `json:"itemName"`
	WikiLink string `json:"wikiLink,omitempty"`
	Icon     string `json:"icon"`
	Thumb    string `json:"thumb"`
}

type LichEphemera struct {
	Id        string                       `json:"id"`
	Slug      string                       `json:"slug"`
	GameRef   string                       `json:"gameRef"`
	Animation string                       `json:"animation"`
	Element   string                       `json:"element"`
	I18N      map[string]*LichEphemeraI18N `json:"i18n,omitempty"`
}

type LichEphemeraI18N struct {
	Name  string `json:"itemName"`
	Icon  string `json:"icon"`
	Thumb string `json:"thumb"`
}

type LichQuirk struct {
	Id    string                    `json:"id"`
	Slug  string                    `json:"slug"`
	Group string                    `json:"group,omitempty"`
	I18N  map[string]*LichQuirkI18N `json:"i18n,omitempty"`
}

type LichQuirkI18N struct {
	Name        string `json:"itemName"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Thumb       string `json:"thumb,omitempty"`
}

type SisterWeapon struct {
	Id             string                       `json:"id"`
	Slug           string                       `json:"slug"`
	GameRef        string                       `json:"gameRef"`
	ReqMasteryRank int8                         `json:"reqMasteryRank"`
	I18N           map[string]*SisterWeaponI18N `json:"i18n,omitempty"`
}

type SisterWeaponI18N struct {
	Name     string `json:"itemName"`
	WikiLink string `json:"wikiLink,omitempty"`
	Icon     string `json:"icon"`
	Thumb    string `json:"thumb"`
}

type SisterEphemera struct {
	Id        string                         `json:"id"`
	Slug      string                         `json:"slug"`
	GameRef   string                         `json:"gameRef"`
	Animation string                         `json:"animation"`
	Element   string                         `json:"element"`
	I18N      map[string]*SisterEphemeraI18N `json:"i18n,omitempty"`
}

type SisterEphemeraI18N struct {
	Name  string `json:"itemName"`
	Icon  string `json:"icon"`
	Thumb string `json:"thumb"`
}

type SisterQuirk struct {
	Id    string                      `json:"id"`
	Slug  string                      `json:"slug"`
	Group string                      `json:"group,omitempty"`
	I18N  map[string]*SisterQuirkI18N `json:"i18n,omitempty"`
}

type SisterQuirkI18N struct {
	Name        string `json:"itemName"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
}

type Npc struct {
	Id      string              `json:"id"`
	Slug    string              `json:"slug"`
	GameRef string              `json:"gameRef"`
	I18N    map[string]*NpcI18N `json:"i18n,omitempty"`
}
type NpcI18N struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Thumb string `json:"thumb"`
}

type Location struct {
	Id       string                   `json:"id"`
	Slug     string                   `json:"slug"`
	GameRef  string                   `json:"gameRef"`
	Faction  string                   `json:"faction,omitempty"`
	MinLevel int32                    `json:"minLevel,omitempty"`
	MaxLevel int32                    `json:"maxLevel,omitempty"`
	I18N     map[string]*LocationI18N `json:"i18n,omitempty"`
}

type LocationI18N struct {
	NodeName   string `json:"nodeName"`
	SystemName string `json:"systemName,omitempty"`
	Icon       string `json:"icon"`
	Thumb      string `json:"thumb"`
}

type Mission struct {
	Id      string                  `json:"id"`
	Slug    string                  `json:"slug"`
	GameRef string                  `json:"gameRef"`
	I18N    map[string]*MissionI18N `json:"i18n,omitempty"`
}

type MissionI18N struct {
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Thumb string `json:"thumb,omitempty"`
}

type Order struct {
	Id         string `json:"id"`                   // Is the unique identifier of the order.
	Type       string `json:"type"`                 // Specifies whether the order is a 'buy' or 'sell'.
	Platinum   int32  `json:"platinum"`             // Is the total platinum currency involved in the order.
	Quantity   int32  `json:"quantity"`             // Represents the number of items included in the order.
	PerTrade   int8   `json:"perTrade,omitempty"`   // (optional) indicates the items quantity per transaction.
	Rank       *int8  `json:"rank,omitempty"`       // (optional) specifies the rank or level of the item in the order.
	Charges    *int8  `json:"charges,omitempty"`    // (optional) specifies number of charges left (used in requiem mods).
	Subtype    string `json:"subtype,omitempty"`    // (optional) defines the specific subtype or category of the item.
	AmberStars *int8  `json:"amberStars,omitempty"` // (optional) denotes the count of amber stars in a sculpture order.
	CyanStars  *int8  `json:"cyanStars,omitempty"`  // (optional) denotes the count of cyan stars in a sculpture order.
	Visible    bool   `json:"visible"`              // (auth\mod) Indicates whether the order is publicly visible or not.
	CreatedAt  string `json:"createdAt"`            // Records the creation time of the order.
	UpdatedAt  string `json:"updatedAt"`            // Records the last modification time of the order.
	ItemId     string `json:"itemId"`               // Is the unique identifier of the item involved in the order.
	Group      string `json:"group"`                // User-defined group to which the order belongs
}

type OrderWithUser struct {
	Order
	User UserShort `json:"user"` // Represents the user who created the order, with basic profile information.
}

type TxItem struct {
	Id         string `json:"id,omitempty"`
	Rank       *int32 `json:"rank,omitempty"`
	Charges    *int32 `json:"charges,omitempty"`
	Subtype    string `json:"subtype,omitempty"`
	AmberStars *int32 `json:"amberStars,omitempty"`
	CyanStars  *int32 `json:"cyanStars,omitempty"`
}

type Transaction struct {
	Id        string     `json:"id"`
	Type      string     `json:"type"`
	OriginId  string     `json:"originId"`
	Platinum  int32      `json:"platinum"`
	Quantity  int32      `json:"quantity"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
	Item      *TxItem    `json:"item,omitempty"`
	User      *UserShort `json:"user,omitempty"`
}

type UserShort struct {
	Id         string `json:"id"`
	IngameName string `json:"ingameName"`       // In-game name of the user.
	Avatar     string `json:"avatar,omitempty"` // Optional avatar image.
	Reputation int16  `json:"reputation"`       // Reputation score.
	Locale     string `json:"locale"`           // Preferred communication language (e.g., 'en', 'ko', 'es').
	Platform   string `json:"platform"`         // Gaming platform used by the user.
	Crossplay  bool   `json:"crossplay"`

	Status   string   `json:"status"`   // Current status of the user.
	Activity Activity `json:"activity"` // Addition to the status, current activity of the user.
	LastSeen string   `json:"lastSeen"` // Timestamp of the user's last online presence.
}

type User struct {
	Id           string `json:"id"`                     // Unique identifier of the user.
	IngameName   string `json:"ingameName"`             // User's in-game name.
	Avatar       string `json:"avatar,omitempty"`       // Optional link to the user's avatar image.
	Background   string `json:"background,omitempty"`   // Optional link to the user's profile background image.
	About        string `json:"about,omitempty"`        // Optional HTML-formatted text about the user.
	Reputation   int16  `json:"reputation"`             // User's reputation score.
	MasteryLevel int8   `json:"masteryLevel,omitempty"` // Optional in-game mastery level.

	Platform  string `json:"platform"`  // Platform the user plays on.
	Crossplay bool   `json:"crossplay"` // Indicates if the user is open to cross-platform trading.
	Locale    string `json:"locale"`    // User's locale or preferred language.

	AchievementShowcase []Achievement `json:"achievementShowcase"` // List of achievements the user chose to showcase.

	Status   string   `json:"status"`   // Current status of the user.
	Activity Activity `json:"activity"` // Current activity the user is engaged in.
	LastSeen string   `json:"lastSeen"` // Timestamp of the user's last online presence.

	Banned   bool   `json:"banned,omitempty"`   // Indicates whether the user is currently banned.
	BanUntil string `json:"banUntil,omitempty"` // End date of the current ban, if applicable.

	// Fields below are accessible only to moderators and admins.
	Warned      bool   `json:"warned,omitempty"`      // Indicates whether the user has been warned.
	WarnMessage string `json:"warnMessage,omitempty"` // Warning message, if any.
	BanMessage  string `json:"banMessage,omitempty"`  // Ban message or reason for the ban, if any.
}

type UserPrivate struct {
	Id          string `json:"id"`
	Role        any    `json:"role"`                 // Role assigned to the user (e.g., moderator, user).
	IngameName  string `json:"ingameName"`           // In-game name.
	Avatar      string `json:"avatar,omitempty"`     // Optional avatar image.
	Background  string `json:"background,omitempty"` // Optional background image.
	About       string `json:"about,omitempty"`      // Optional about text in html
	AboutRaw    string `json:"aboutRaw,omitempty"`   // Optional about text in raw markdown.
	Reputation  int16  `json:"reputation"`           // Reputation score.
	MasteryRank int8   `json:"masteryRank"`          // In-game mastery level.
	Credits     int32  `json:"credits"`              // In-game currency balance.

	Platform  string `json:"platform"`  // Gaming platform.
	Crossplay bool   `json:"crossplay"` // Crossplay enabled \ disabled
	Locale    string `json:"locale"`    // Preferred communication language.
	Theme     string `json:"theme"`     // Preferred color scheme for UI.

	AchievementShowcase []Achievement `json:"achievementShowcase"` // List of achievements the user chose to showcase.

	Verification bool   `json:"verification"` // Verification status.
	CheckCode    string `json:"checkCode"`    // Unique check code for the user.

	Tier         any  `json:"tier"`         // Subscription tier.
	Subscription bool `json:"subscription"` // Subscription status.

	Warned      bool   `json:"warned,omitempty"`
	WarnMessage string `json:"warnMessage,omitempty"`
	Banned      bool   `json:"banned,omitempty"`     // Ban status.
	BanUntil    string `json:"banUntil,omitempty"`   // End date of the ban.
	BanMessage  string `json:"banMessage,omitempty"` // Reason for the ban.

	ReviewsLeft    int16    `json:"reviewsLeft"`    // How much reviews the user can still write today. (reset at midnight UTC)
	UnreadMessages int16    `json:"unreadMessages"` // Count of unread messages.
	IgnoreList     []string `json:"ignoreList"`     // List of ignored users.

	DeleteInProgress bool   `json:"deleteInProgress,omitempty"` // Flag for pending deletion of the account.
	DeleteAt         string `json:"deleteAt,omitempty"`         // Scheduled deletion date.

	LinkedAccounts any  `json:"linkedAccounts"` // Accounts linked with the user's profile.
	HasEmail       bool `json:"hasEmail"`       // If the user has an email address.

	LastSeen  string `json:"lastSeen"`  // Timestamp of the last online presence.
	CreatedAt string `json:"createdAt"` // Account creation date.
}

type Activity struct {
	Type      ActivityType `json:"type" `               // Name of the activity (e.g., 'on mission', 'dojo').
	Details   string       `json:"details,omitempty"`   // Optional specifics about the activity (e.g., mission name, solo/squad status).
	StartedAt string       `json:"startedAt,omitempty"` // Timestamp of the activity start.
}

type AchievementState struct {
	Featured    bool   `json:"featured,omitempty"`    // If true, the achievement is featured
	Hidden      bool   `json:"hidden,omitempty"`      // If true, the achievement is hidden from the public
	Progress    *int32 `json:"progress,omitempty"`    // Current progress towards the achievement goal
	CompletedAt string `json:"completedAt,omitempty"` // Timestamp when the achievement was achieved
}

type Achievement struct {
	Id              string                      `json:"id"`
	Slug            string                      `json:"slug"`
	Type            string                      `json:"type"`                      // Type of the achievement (e.g., "task", "event", etc.)
	Secret          bool                        `json:"secret,omitempty"`          // If true, the achievement is secret and not shown to public
	ReputationBonus int32                       `json:"reputationBonus,omitempty"` // Reputation bonus for the achievement
	Goal            int32                       `json:"goal,omitempty"`            // Goal to achieve
	I18N            map[string]*AchievementI18N `json:"i18n"`                      // Localized text for the achievement
	State           *AchievementState           `json:"state,omitempty"`           // Current state of the achievement
}

type AchievementI18N struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Thumb       string `json:"thumb,omitempty"`
}

type DashboardShowcase struct {
	I18N  map[string]*DashboardShowcaseI18N `json:"i18n,omitempty"`
	Items []*DashboardShowcaseItem          `json:"items"`
}

type DashboardShowcaseI18N struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type DashboardShowcaseItem struct {
	Item       string `json:"item"`
	Background string `json:"background"`
	BigCard    bool   `json:"bigCard"`
}

type Status string

const (
	StatusOnline    Status = "online"
	StatusOffline   Status = "offline"
	StatusInvisible Status = "invisible"
	StatusInGame    Status = "ingame"
)

type ActivityType string

const (
	ActivityUnknown   ActivityType = "UNKNOWN"
	ActivityIdle      ActivityType = "IDLE"
	ActivityOnMission ActivityType = "ON_MISSION"
	ActivityInDojo    ActivityType = "IN_DOJO"
	ActivityInOrbiter ActivityType = "IN_ORBITER"
	ActivityInRelay   ActivityType = "IN_RELAY"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

type Tier string

const (
	TierNone    Tier = "none"
	TierBronze  Tier = "bronze"
	TierSilver  Tier = "silver"
	TierGold    Tier = "gold"
	TierDiamond Tier = "diamond"
)

type Language string

const (
	LangKO     Language = "ko"
	LangRU     Language = "ru"
	LangDE     Language = "de"
	LangFR     Language = "fr"
	LangPT     Language = "pt"
	LangZHSimp Language = "zh-hans"
	LangZHTrad Language = "zh-hant"
	LangES     Language = "es"
	LangIT     Language = "it"
	LangPL     Language = "pl"
	LangUK     Language = "uk"
	LangEN     Language = "en"
)

type Platform string

const (
	PlatformPC     Platform = "pc"
	PlatformPS4    Platform = "ps4"
	PlatformXbox   Platform = "xbox"
	PlatformSwitch Platform = "switch"
	PlatformMobile Platform = "mobile"
)

type Scope string

const (
	ScopeMe        Scope = "me"
	ScopeProfile   Scope = "profile"
	ScopeSettings  Scope = "settings"
	ScopeContracts Scope = "contracts"
	ScopeLedger    Scope = "ledger"
	ScopeReviews   Scope = "reviews"
)

type Versions struct {
	Apps        VersionsApps        `json:"apps"`
	Collections VersionsCollections `json:"collections"`
	UpdatedAt   string              `json:"updatedAt"`
}

type VersionsApps struct {
	Ios         string `json:"ios"`
	Android     string `json:"android"`
	MinIos      string `json:"minIos"`
	MinAndroid  string `json:"minAndroid"`
}

type VersionsCollections struct {
	Items     string `json:"items"`
	Rivens    string `json:"rivens"`
	Liches    string `json:"liches"`
	Sisters   string `json:"sisters"`
	Missions  string `json:"missions"`
	Npcs      string `json:"npcs"`
	Locations string `json:"locations"`
}

type ItemSet struct {
	Id    string `json:"id"`
	Items []Item `json:"items"`
}

type TopOrders struct {
	Buy  []OrderWithUser `json:"buy"`
	Sell []OrderWithUser `json:"sell"`
}

type TopOrdersParams struct {
	Rank         *int
	RankLt       *int
	Charges      *int
	ChargesLt    *int
	AmberStars   *int
	AmberStarsLt *int
	CyanStars    *int
	CyanStarsLt  *int
	Subtype      string
}

type genericResponse[T any] struct {
	ApiVersion string `json:"api_version"`
	Data       T      `json:"data"`
	Error      any    `json:"error,omitempty"`
}
