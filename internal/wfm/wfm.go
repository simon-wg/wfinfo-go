package wfm

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultBaseURL = "https://api.warframe.market"
	DefaultTimeout = 30 * time.Second
)

// Client is a Warframe Market API client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	ctx        context.Context
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the API.
func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client for the API.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new Warframe Market API client.
func NewClient(opts ...ClientOption) *Client {
	u, _ := url.Parse(DefaultBaseURL)
	c := &Client{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithContext returns a shallow copy of the client with the provided context.
func (c *Client) WithContext(ctx context.Context) *Client {
	if ctx == nil {
		return c
	}
	c2 := *c
	c2.ctx = ctx
	return &c2
}

func (c *Client) context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func fetchResource[T any](c *Client, query url.Values, path ...string) (T, error) {
	u := c.baseURL.JoinPath(path...)
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[T]
	if err := c.do(req, &resp); err != nil {
		var zero T
		return zero, err
	}

	return resp.Data, nil
}

func (c *Client) do(req *http.Request, v any) error {
	ctx := req.Context()
	var resp *http.Response
	var err error

	failureCount := 0
	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusTooManyRequests && failureCount < 5 {
			//nolint:errcheck
			resp.Body.Close()
			failureCount++
			sleepDuration := time.Duration(math.Pow(2, float64(failureCount))) * time.Second

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleepDuration):
				continue
			}
		}
		break
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// FetchItems fetches items from the Warframe Market API.
func (c *Client) FetchItems() ([]Item, error) {
	currentVersions, err := c.FetchVersions()
	if err != nil {
		if items, err := getFromCache[[]Item](c, "items.json"); err == nil {
			return *items, nil
		}
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}

	cachedVersions, _ := getFromCache[Versions](c, "versions.json")
	if cachedVersions != nil && *currentVersions == *cachedVersions {
		if items, err := getFromCache[[]Item](c, "items.json"); err == nil {
			return *items, nil
		}
	}

	items, err := fetchResource[[]Item](c, nil, "v2", "items")
	if err != nil {
		return nil, err
	}

	saveToCache(c, "items.json", items)
	saveToCache(c, "versions.json", currentVersions)

	return items, nil
}

func (c *Client) getCachePath(filename string) string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return filename
	}
	appCacheDir := filepath.Join(cacheDir, "wfm-go")
	_ = os.MkdirAll(appCacheDir, 0755)
	return filepath.Join(appCacheDir, filename)
}

func getFromCache[T any](c *Client, filename string) (*T, error) {
	file, err := os.ReadFile(c.getCachePath(filename))
	if err != nil {
		return nil, err
	}
	var v T
	if err := json.Unmarshal(file, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func saveToCache[T any](c *Client, filename string, v T) {
	data, _ := json.MarshalIndent(v, "", "  ")
	_ = os.WriteFile(c.getCachePath(filename), data, 0644)
}

// FetchVersions fetches the current version number of the server's resources.
func (c *Client) FetchVersions() (*Versions, error) {
	return fetchResource[*Versions](c, nil, "v2", "versions")
}

// FetchItem fetches full info about one particular item.
func (c *Client) FetchItem(slug string) (*Item, error) {
	return fetchResource[*Item](c, nil, "v2", "item", slug)
}

// FetchItemSet retrieves information on item sets.
func (c *Client) FetchItemSet(slug string) (*ItemSet, error) {
	return fetchResource[*ItemSet](c, nil, "v2", "item", slug, "set")
}

// FetchRivenWeapons fetches all tradable riven items.
func (c *Client) FetchRivenWeapons() ([]Riven, error) {
	return fetchResource[[]Riven](c, nil, "v2", "riven", "weapons")
}

// FetchRivenWeapon fetches full info about one particular riven item.
func (c *Client) FetchRivenWeapon(slug string) (*Riven, error) {
	return fetchResource[*Riven](c, nil, "v2", "riven", "weapon", slug)
}

// FetchRivenAttributes fetches all attributes for riven weapons.
func (c *Client) FetchRivenAttributes() ([]RivenAttribute, error) {
	return fetchResource[[]RivenAttribute](c, nil, "v2", "riven", "attributes")
}

// FetchLichWeapons fetches all tradable lich weapons.
func (c *Client) FetchLichWeapons() ([]LichWeapon, error) {
	return fetchResource[[]LichWeapon](c, nil, "v2", "lich", "weapons")
}

// FetchLichWeapon fetches full info about one particular lich weapon.
func (c *Client) FetchLichWeapon(slug string) (*LichWeapon, error) {
	return fetchResource[*LichWeapon](c, nil, "v2", "lich", "weapon", slug)
}

// FetchLichEphemeras fetches all tradable lich ephemeras.
func (c *Client) FetchLichEphemeras() ([]LichEphemera, error) {
	return fetchResource[[]LichEphemera](c, nil, "v2", "lich", "ephemeras")
}

// FetchLichQuirks fetches all tradable lich quirks.
func (c *Client) FetchLichQuirks() ([]LichQuirk, error) {
	return fetchResource[[]LichQuirk](c, nil, "v2", "lich", "quirks")
}

// FetchSisterWeapons fetches all tradable sister weapons.
func (c *Client) FetchSisterWeapons() ([]SisterWeapon, error) {
	return fetchResource[[]SisterWeapon](c, nil, "v2", "sister", "weapons")
}

// FetchSisterWeapon fetches full info about one particular sister weapon.
func (c *Client) FetchSisterWeapon(slug string) (*SisterWeapon, error) {
	return fetchResource[*SisterWeapon](c, nil, "v2", "sister", "weapon", slug)
}

// FetchSisterEphemeras fetches all tradable sister ephemeras.
func (c *Client) FetchSisterEphemeras() ([]SisterEphemera, error) {
	return fetchResource[[]SisterEphemera](c, nil, "v2", "sister", "ephemeras")
}

// FetchSisterQuirks fetches all tradable sister quirks.
func (c *Client) FetchSisterQuirks() ([]SisterQuirk, error) {
	return fetchResource[[]SisterQuirk](c, nil, "v2", "sister", "quirks")
}

// FetchLocations fetches all known locations.
func (c *Client) FetchLocations() ([]Location, error) {
	return fetchResource[[]Location](c, nil, "v2", "locations")
}

// FetchNpcs fetches all known NPCs.
func (c *Client) FetchNpcs() ([]Npc, error) {
	return fetchResource[[]Npc](c, nil, "v2", "npcs")
}

// FetchMissions fetches all known missions.
func (c *Client) FetchMissions() ([]Mission, error) {
	return fetchResource[[]Mission](c, nil, "v2", "missions")
}

// FetchRecentOrders fetches the most recent orders.
func (c *Client) FetchRecentOrders() ([]OrderWithUser, error) {
	return fetchResource[[]OrderWithUser](c, nil, "v2", "orders", "recent")
}

// FetchItemOrders fetches all orders for an item from users online within the last 7 days.
func (c *Client) FetchItemOrders(slug string) ([]OrderWithUser, error) {
	return fetchResource[[]OrderWithUser](c, nil, "v2", "orders", "item", slug)
}

// FetchItemTopOrders fetches the top 5 buy and top 5 sell orders for a specific item.
func (c *Client) FetchItemTopOrders(slug string, params *TopOrdersParams) (*TopOrders, error) {
	var q url.Values
	if params != nil {
		q = make(url.Values)
		if params.Rank != nil {
			q.Set("rank", fmt.Sprintf("%d", *params.Rank))
		}
		if params.RankLt != nil {
			q.Set("rankLt", fmt.Sprintf("%d", *params.RankLt))
		}
		if params.Charges != nil {
			q.Set("charges", fmt.Sprintf("%d", *params.Charges))
		}
		if params.ChargesLt != nil {
			q.Set("chargesLt", fmt.Sprintf("%d", *params.ChargesLt))
		}
		if params.AmberStars != nil {
			q.Set("amberStars", fmt.Sprintf("%d", *params.AmberStars))
		}
		if params.AmberStarsLt != nil {
			q.Set("amberStarsLt", fmt.Sprintf("%d", *params.AmberStarsLt))
		}
		if params.CyanStars != nil {
			q.Set("cyanStars", fmt.Sprintf("%d", *params.CyanStars))
		}
		if params.CyanStarsLt != nil {
			q.Set("cyanStarsLt", fmt.Sprintf("%d", *params.CyanStarsLt))
		}
		if params.Subtype != "" {
			q.Set("subtype", params.Subtype)
		}
	}

	return fetchResource[*TopOrders](c, q, "v2", "orders", "item", slug, "top")
}

// FetchUserOrders fetches public orders from a specified user by slug.
func (c *Client) FetchUserOrders(slug string) ([]Order, error) {
	return fetchResource[[]Order](c, nil, "v2", "orders", "user", slug)
}

// FetchUserOrdersById fetches public orders from a specified user by ID.
func (c *Client) FetchUserOrdersById(userId string) ([]Order, error) {
	return fetchResource[[]Order](c, nil, "v2", "orders", "userId", userId)
}

// FetchOrder fetches a single order by ID.
func (c *Client) FetchOrder(id string) (*OrderWithUser, error) {
	return fetchResource[*OrderWithUser](c, nil, "v2", "order", id)
}

// FetchUser fetches information about a particular user by slug.
func (c *Client) FetchUser(slug string) (*User, error) {
	return fetchResource[*User](c, nil, "v2", "user", slug)
}

// FetchUserById fetches information about a particular user by ID.
func (c *Client) FetchUserById(userId string) (*User, error) {
	return fetchResource[*User](c, nil, "v2", "userId", userId)
}

// FetchAchievements fetches all available achievements (except secret ones).
func (c *Client) FetchAchievements() ([]Achievement, error) {
	return fetchResource[[]Achievement](c, nil, "v2", "achievements")
}

// FetchUserAchievements fetches all user achievements by slug.
func (c *Client) FetchUserAchievements(slug string, featured bool) ([]Achievement, error) {
	var q url.Values
	if featured {
		q = make(url.Values)
		q.Set("featured", "true")
	}
	return fetchResource[[]Achievement](c, q, "v2", "achievements", "user", slug)
}

// FetchUserAchievementsById fetches all user achievements by ID.
func (c *Client) FetchUserAchievementsById(userId string, featured bool) ([]Achievement, error) {
	var q url.Values
	if featured {
		q = make(url.Values)
		q.Set("featured", "true")
	}
	return fetchResource[[]Achievement](c, q, "v2", "achievements", "userId", userId)
}

// FetchDashboardShowcase fetches featured items for the mobile app main screen.
func (c *Client) FetchDashboardShowcase() (*DashboardShowcase, error) {
	return fetchResource[*DashboardShowcase](c, nil, "v2", "dashboard", "showcase")
}
