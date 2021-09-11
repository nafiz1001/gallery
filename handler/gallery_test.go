package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nafiz1001/gallery-go/dto"
)

func GetArts(t *testing.T) []dto.ArtDto {
	var arr []dto.ArtDto

	if resp, err := http.Get("http://localhost:8080/arts"); err != nil {
		t.Fatal(err)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
			t.Fatal(err)
		}
	}

	return arr
}

func NewRequest(t *testing.T, method string, url string, body string, username string, password string) (*http.Response, error) {
	if req, err := http.NewRequest(method, url, bytes.NewBufferString(body)); err != nil {
		t.Fatal(err)
		return nil, nil
	} else {
		req.Header.Set("Content-Type", "application/json")
		if len(username) > 0 {
			req.SetBasicAuth(username, password)
		}
		if resp, _ := http.DefaultClient.Do(req); 400 <= resp.StatusCode && resp.StatusCode <= 599 {
			b, _ := io.ReadAll(resp.Body)
			return resp, fmt.Errorf("%s: %s", resp.Status, string(b))
		} else {
			return resp, nil
		}
	}
}

func TestGallery(t *testing.T) {
	h := GalleryHandler{}
	h.Init()

	srv := &http.Server{
		Handler: h,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go srv.ListenAndServe()

	time.Sleep(500 * time.Millisecond)

	var art dto.ArtDto

	if arts := GetArts(t); len(arts) != 0 {
		t.Fatalf("%v length is not 0", arts)
	}

	if _, err := http.Post("http://localhost:8080/accounts", "application/json", bytes.NewBufferString(`{"username":"bad", "password":"bad"}`)); err != nil {
		t.Fatal(err)
	}

	if _, err := http.Post("http://localhost:8080/accounts", "application/json", bytes.NewBufferString(`{"username":"good", "password":"good"}`)); err != nil {
		t.Fatal(err)
	}

	if resp, err := NewRequest(t, http.MethodPost, "http://localhost:8080/arts", `{"title":"title"}`, "", ""); err == nil {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected to fail uploading because basic auth is missing\n%s", string(b))
	}

	if resp, err := NewRequest(t, http.MethodPost, "http://localhost:8080/arts", `{"title":"title"}`, "good", "good"); err != nil {
		t.Fatalf("%s", err)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&art); err != nil {
			t.Fatal(err)
		} else {
			if art.Title != "title" {
				t.Fatalf("the title of response (%v) is not equal to 'title'", art)
			}
		}
	}

	if arts := GetArts(t); arts[0].Id != art.Id {
		t.Fatalf("the response (%v) does not have art with id %s", arts, art.Id)
	}

	if _, err := NewRequest(t, http.MethodPut, "http://localhost:8080/arts/"+art.Id, `{"title":"title2"}`, "bad", "bad"); err == nil {
		t.Fatal(err)
	}

	if resp, err := NewRequest(t, http.MethodPut, "http://localhost:8080/arts/"+art.Id, `{"title":"title2"}`, "good", "good"); err != nil {
		t.Fatalf("%s", err)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&art); err != nil {
			t.Fatal(err)
		} else {
			if art.Title != "title2" {
				t.Fatalf("the title of response (%v) is not equal to 'title2'", art)
			}
		}
	}

	if arts := GetArts(t); arts[0].Title != art.Title {
		t.Fatalf("the response (%v) does not have art with %s", arts, art.Title)
	}

	if _, err := NewRequest(t, http.MethodDelete, "http://localhost:8080/arts/"+art.Id, "", "bad", "bad"); err == nil {
		t.Fatal(err)
	}

	if arts := GetArts(t); len(arts) != 1 {
		t.Fatalf("%v length is not 1", arts)
	}

	if resp, err := NewRequest(t, http.MethodDelete, "http://localhost:8080/arts/"+art.Id, "", "good", "good"); err != nil {
		t.Fatal(err)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&art); err != nil {
			t.Fatal(err)
		}
	}

	if arts := GetArts(t); len(arts) != 0 {
		t.Fatalf("%v length is not 0", arts)
	}
}
