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
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetHostMetaData(t *testing.T) {
	const (
		hostID    = "9rxGOHfVF8F"
		namespace = "testing"
	)
	lastModified := time.Date(2018, 3, 6, 3, 0, 0, 0, time.UTC)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"type":12345,"region":"jp","env":"staging","instance_type":"c4.xlarge"}`)
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

	type metadata struct {
		Type         int    `json:"type"`
		Region       string `json:"region"`
		Env          string `json:"env"`
		InstanceType string `json:"instance_type"`
	}
	var got metadata
	metametadata, err := c.GetHostMetaData(context.Background(), hostID, namespace, &got)
	if err != nil {
		t.Error(err)
	}
	if !metametadata.LastModified.Equal(lastModified) {
		t.Errorf("want %s, got %s", lastModified, metametadata.LastModified)
	}

	want := metadata{
		Type:         12345,
		Region:       "jp",
		Env:          "staging",
		InstanceType: "c4.xlarge",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("metadata differs: (-got +want)\n%s", diff)
	}
}

func TestGetMetaDataNameSpaces(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":[{"namespace": "test"}]}`)
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

	ns, err := c.GetHostMetaDataNameSpaces(context.Background(), "9rxGOHfVF8F")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ns, []string{"test"}) {
		t.Errorf(`want []string{"test}, got %#v`, ns)
	}
}

func TestPutHostMetaData(t *testing.T) {
	const (
		hostID    = "9rxGOHfVF8F"
		namespace = "testing"
	)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got map[string]interface{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&got); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		want := map[string]interface{}{
			"type":          12345.0,
			"region":        "jp",
			"env":           "staging",
			"instance_type": "c4.xlarge",
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("metadata differs: (-got +want)\n%s", diff)
		}
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: got %s, want %s", http.MethodPut, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"success":true}`)
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

	metadata := struct {
		Type         int    `json:"type"`
		Region       string `json:"region"`
		Env          string `json:"env"`
		InstanceType string `json:"instance_type"`
	}{
		Type:         12345,
		Region:       "jp",
		Env:          "staging",
		InstanceType: "c4.xlarge",
	}
	err = c.PutHostMetaData(context.Background(), hostID, namespace, metadata)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteMetaData(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: got %s, want %s", http.MethodDelete, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"success":true`)
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

	err = c.DeleteHostMetaData(context.Background(), "9rxGOHfVF8F", "test")
	if err != nil {
		t.Error(err)
	}
}
