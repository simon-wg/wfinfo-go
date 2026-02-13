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

	u := c.baseURL.JoinPath("v2", "items")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Item]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	saveToCache(c, "items.json", resp.Data)
	saveToCache(c, "versions.json", currentVersions)

	return resp.Data, nil
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
	u := c.baseURL.JoinPath("v2", "versions")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*Versions]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}

	return resp.Data, nil
}

// FetchItem fetches full info about one particular item.
func (c *Client) FetchItem(slug string) (*Item, error) {
	u := c.baseURL.JoinPath("v2", "item", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*Item]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch item: %w", err)
	}

	return resp.Data, nil
}

// FetchItemSet retrieves information on item sets.
func (c *Client) FetchItemSet(slug string) (*ItemSet, error) {
	u := c.baseURL.JoinPath("v2", "item", slug, "set")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*ItemSet]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch item set: %w", err)
	}

	return resp.Data, nil
}

// FetchRivenWeapons fetches all tradable riven items.
func (c *Client) FetchRivenWeapons() ([]Riven, error) {
	u := c.baseURL.JoinPath("v2", "riven", "weapons")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Riven]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch riven weapons: %w", err)
	}

	return resp.Data, nil
}

// FetchRivenWeapon fetches full info about one particular riven item.
func (c *Client) FetchRivenWeapon(slug string) (*Riven, error) {
	u := c.baseURL.JoinPath("v2", "riven", "weapon", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*Riven]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch riven weapon: %w", err)
	}

	return resp.Data, nil
}

// FetchRivenAttributes fetches all attributes for riven weapons.
func (c *Client) FetchRivenAttributes() ([]RivenAttribute, error) {
	u := c.baseURL.JoinPath("v2", "riven", "attributes")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]RivenAttribute]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch riven attributes: %w", err)
	}

	return resp.Data, nil
}

// FetchLichWeapons fetches all tradable lich weapons.
func (c *Client) FetchLichWeapons() ([]LichWeapon, error) {
	u := c.baseURL.JoinPath("v2", "lich", "weapons")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]LichWeapon]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch lich weapons: %w", err)
	}

	return resp.Data, nil
}

// FetchLichWeapon fetches full info about one particular lich weapon.
func (c *Client) FetchLichWeapon(slug string) (*LichWeapon, error) {
	u := c.baseURL.JoinPath("v2", "lich", "weapon", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*LichWeapon]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch lich weapon: %w", err)
	}

	return resp.Data, nil
}

// FetchLichEphemeras fetches all tradable lich ephemeras.
func (c *Client) FetchLichEphemeras() ([]LichEphemera, error) {
	u := c.baseURL.JoinPath("v2", "lich", "ephemeras")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]LichEphemera]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch lich ephemeras: %w", err)
	}

	return resp.Data, nil
}

// FetchLichQuirks fetches all tradable lich quirks.
func (c *Client) FetchLichQuirks() ([]LichQuirk, error) {
	u := c.baseURL.JoinPath("v2", "lich", "quirks")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]LichQuirk]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch lich quirks: %w", err)
	}

	return resp.Data, nil
}

// FetchSisterWeapons fetches all tradable sister weapons.
func (c *Client) FetchSisterWeapons() ([]SisterWeapon, error) {
	u := c.baseURL.JoinPath("v2", "sister", "weapons")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]SisterWeapon]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch sister weapons: %w", err)
	}

	return resp.Data, nil
}

// FetchSisterWeapon fetches full info about one particular sister weapon.
func (c *Client) FetchSisterWeapon(slug string) (*SisterWeapon, error) {
	u := c.baseURL.JoinPath("v2", "sister", "weapon", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*SisterWeapon]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch sister weapon: %w", err)
	}

	return resp.Data, nil
}

// FetchSisterEphemeras fetches all tradable sister ephemeras.
func (c *Client) FetchSisterEphemeras() ([]SisterEphemera, error) {
	u := c.baseURL.JoinPath("v2", "sister", "ephemeras")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]SisterEphemera]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch sister ephemeras: %w", err)
	}

	return resp.Data, nil
}

// FetchSisterQuirks fetches all tradable sister quirks.
func (c *Client) FetchSisterQuirks() ([]SisterQuirk, error) {
	u := c.baseURL.JoinPath("v2", "sister", "quirks")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]SisterQuirk]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch sister quirks: %w", err)
	}

	return resp.Data, nil
}

// FetchLocations fetches all known locations.
func (c *Client) FetchLocations() ([]Location, error) {
	u := c.baseURL.JoinPath("v2", "locations")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Location]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}

	return resp.Data, nil
}

// FetchNpcs fetches all known NPCs.
func (c *Client) FetchNpcs() ([]Npc, error) {
	u := c.baseURL.JoinPath("v2", "npcs")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Npc]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch npcs: %w", err)
	}

	return resp.Data, nil
}

// FetchMissions fetches all known missions.
func (c *Client) FetchMissions() ([]Mission, error) {
	u := c.baseURL.JoinPath("v2", "missions")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Mission]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch missions: %w", err)
	}

	return resp.Data, nil
}

// FetchRecentOrders fetches the most recent orders.
func (c *Client) FetchRecentOrders() ([]OrderWithUser, error) {
	u := c.baseURL.JoinPath("v2", "orders", "recent")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]OrderWithUser]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch recent orders: %w", err)
	}

	return resp.Data, nil
}

// FetchItemOrders fetches all orders for an item from users online within the last 7 days.
func (c *Client) FetchItemOrders(slug string) ([]OrderWithUser, error) {
	u := c.baseURL.JoinPath("v2", "orders", "item", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]OrderWithUser]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch item orders: %w", err)
	}

	return resp.Data, nil
}

// FetchItemTopOrders fetches the top 5 buy and top 5 sell orders for a specific item.
func (c *Client) FetchItemTopOrders(slug string, params *TopOrdersParams) (*TopOrders, error) {
	u := c.baseURL.JoinPath("v2", "orders", "item", slug, "top")
	if params != nil {
		q := u.Query()
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
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*TopOrders]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch top orders: %w", err)
	}

	return resp.Data, nil
}

// FetchUserOrders fetches public orders from a specified user by slug.
func (c *Client) FetchUserOrders(slug string) ([]Order, error) {
	u := c.baseURL.JoinPath("v2", "orders", "user", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Order]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user orders: %w", err)
	}

	return resp.Data, nil
}

// FetchUserOrdersById fetches public orders from a specified user by ID.
func (c *Client) FetchUserOrdersById(userId string) ([]Order, error) {
	u := c.baseURL.JoinPath("v2", "orders", "userId", userId)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Order]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user orders by id: %w", err)
	}

	return resp.Data, nil
}

// FetchOrder fetches a single order by ID.
func (c *Client) FetchOrder(id string) (*OrderWithUser, error) {
	u := c.baseURL.JoinPath("v2", "order", id)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*OrderWithUser]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	return resp.Data, nil
}

// FetchUser fetches information about a particular user by slug.
func (c *Client) FetchUser(slug string) (*User, error) {
	u := c.baseURL.JoinPath("v2", "user", slug)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*User]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return resp.Data, nil
}

// FetchUserById fetches information about a particular user by ID.
func (c *Client) FetchUserById(userId string) (*User, error) {
	u := c.baseURL.JoinPath("v2", "userId", userId)
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*User]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}

	return resp.Data, nil
}

// FetchAchievements fetches all available achievements (except secret ones).
func (c *Client) FetchAchievements() ([]Achievement, error) {
	u := c.baseURL.JoinPath("v2", "achievements")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Achievement]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch achievements: %w", err)
	}

	return resp.Data, nil
}

// FetchUserAchievements fetches all user achievements by slug.
func (c *Client) FetchUserAchievements(slug string, featured bool) ([]Achievement, error) {
	u := c.baseURL.JoinPath("v2", "achievements", "user", slug)
	if featured {
		q := u.Query()
		q.Set("featured", "true")
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Achievement]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user achievements: %w", err)
	}

	return resp.Data, nil
}

// FetchUserAchievementsById fetches all user achievements by ID.
func (c *Client) FetchUserAchievementsById(userId string, featured bool) ([]Achievement, error) {
	u := c.baseURL.JoinPath("v2", "achievements", "userId", userId)
	if featured {
		q := u.Query()
		q.Set("featured", "true")
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[[]Achievement]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch user achievements by id: %w", err)
	}

	return resp.Data, nil
}

// FetchDashboardShowcase fetches featured items for the mobile app main screen.
func (c *Client) FetchDashboardShowcase() (*DashboardShowcase, error) {
	u := c.baseURL.JoinPath("v2", "dashboard", "showcase")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp genericResponse[*DashboardShowcase]
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch dashboard showcase: %w", err)
	}

	return resp.Data, nil
}
