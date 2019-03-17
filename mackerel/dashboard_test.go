package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestFindDashboard(t *testing.T) {
	tests := []struct {
		resp map[string]interface{} // the response of the mackerel api
		want *Dashboard
	}{
		{
			resp: map[string]interface{}{
				"id":      "foobar",
				"title":   "title",
				"urlPath": "url path",
				"widgets": []map[string]interface{}{
					{
						"type":  "graph",
						"title": "graph title",
						"graph": map[string]string{
							"type":   "host",
							"hostId": "host-foobar",
							"name":   "host-graph",
						},
					},
				},
				"createdAt": 1234567890,
				"updatedAt": 1234567890,
			},
			want: &Dashboard{
				ID:      "foobar",
				Title:   "title",
				URLPath: "url path",
				Widgets: []Widget{
					&WidgetGraph{
						Type:  WidgetTypeGraph,
						Title: "graph title",
						Graph: &GraphHost{
							Type:   GraphTypeHost,
							HostID: "host-foobar",
							Name:   "host-graph",
						},
					},
				},
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("FindDashboards-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode([]interface{}{tc.resp})
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			c := &Client{
				BaseURL:    u,
				HTTPClient: ts.Client(),
			}
			got, err := c.FindDashboards(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, []*Dashboard{tc.want}) {
				t.Errorf("want %#v, got %#v", []*Dashboard{tc.want}, got)
			}
		})

		t.Run(fmt.Sprintf("FindDashboard-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: want %s, got %s", http.MethodGet, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards/foobar" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards/foobar", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode(tc.resp)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			c := &Client{
				BaseURL:    u,
				HTTPClient: ts.Client(),
			}
			got, err := c.FindDashboard(context.Background(), "foobar")
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want %#v, got %#v", tc.want, got)
			}
		})

		t.Run(fmt.Sprintf("DeleteDashboard-%d", i), func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("unexpected method: want %s, got %s", http.MethodDelete, r.Method)
				}
				if r.URL.Path != "/api/v0/dashboards/foobar" {
					t.Errorf("unexpected path, want %s, got %s", "/api/v0/dashboards/foobar", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				enc.Encode(tc.resp)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			c := &Client{
				BaseURL:    u,
				HTTPClient: ts.Client(),
			}
			got, err := c.DeleteDashboard(context.Background(), "foobar")
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestCreateDashboard(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		BaseURL:    u,
		HTTPClient: ts.Client(),
	}

	ret, err := c.CreateDashboard(context.Background(), &Dashboard{
		Title:   "title",
		Memo:    "memo",
		URLPath: "url path",
		Widgets: []Widget{},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
