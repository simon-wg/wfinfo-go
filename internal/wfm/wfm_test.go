package wfm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_FetchItems(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/versions" {
			_ = json.NewEncoder(w).Encode(genericResponse[*Versions]{
				Data: &Versions{UpdatedAt: "2024-01-01T00:00:00Z"},
			})
			return
		}
		if r.URL.Path == "/v2/items" {
			_ = json.NewEncoder(w).Encode(genericResponse[[]Item]{
				Data: []Item{
					{Id: "1", Slug: "item-1"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := NewClient(WithBaseURL(u))

	// We need to bypass the cache for this test or use a temporary cache directory
	// For now, let's just test that it can fetch items
	items, err := client.FetchItems()
	if err != nil {
		t.Fatalf("FetchItems failed: %v", err)
	}

	if len(items) != 1 || items[0].Slug != "item-1" {
		t.Errorf("Unexpected items: %v", items)
	}
}

func TestClient_FetchItemTopOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/orders/item/ash-prime/top" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("rank") != "5" {
			t.Errorf("Expected rank=5, got %s", q.Get("rank"))
		}

		_ = json.NewEncoder(w).Encode(genericResponse[*TopOrders]{
			Data: &TopOrders{
				Buy: []OrderWithUser{{Order: Order{Platinum: 10}}},
			},
		})
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := NewClient(WithBaseURL(u))

	rank := 5
	orders, err := client.FetchItemTopOrders("ash-prime", &TopOrdersParams{Rank: &rank})
	if err != nil {
		t.Fatalf("FetchItemTopOrders failed: %v", err)
	}

	if len(orders.Buy) != 1 || orders.Buy[0].Platinum != 10 {
		t.Errorf("Unexpected orders: %v", orders)
	}
}

func TestWithBaseURL(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	client := NewClient(WithBaseURL(u))
	if client.baseURL.String() != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", client.baseURL.String())
	}
}

func TestGenericResponseDecoding(t *testing.T) {
	data := `{"api_version":"1.0","data":{"id":"123"}}`
	var resp genericResponse[map[string]string]
	err := json.Unmarshal([]byte(data), &resp)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	expected := map[string]string{"id": "123"}
	if !reflect.DeepEqual(resp.Data, expected) {
		t.Errorf("Expected %v, got %v", expected, resp.Data)
	}
}
