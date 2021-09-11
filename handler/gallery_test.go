package handler

import (
	"bytes"
	"encoding/json"
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

func NewRequest(t *testing.T, method string, url string, body string) (*http.Response, error) {
	if req, err := http.NewRequest(method, url, bytes.NewBufferString(body)); err != nil {
		t.Fatal(err)
		return nil, nil
	} else {
		return http.DefaultClient.Do(req)
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

	if resp, err := http.Post("http://localhost:8080/arts", "application/json", bytes.NewBufferString(`{"title":"title"}`)); err != nil {
		t.Fatal(err)
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

	if resp, err := NewRequest(t, http.MethodPut, "http://localhost:8080/arts/"+art.Id, `{"title":"title2"}`); err != nil {
		t.Fatal(err)
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

	if resp, err := NewRequest(t, http.MethodDelete, "http://localhost:8080/arts/"+art.Id, ""); err != nil {
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
