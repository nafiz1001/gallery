package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/nafiz1001/gallery-go/dto"

	_ "github.com/mattn/go-sqlite3"
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

func CheckError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestGallery(t *testing.T) {
	go func() {
		// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
		// psqlconn, ok := os.LookupEnv("DATABASE_URL")
		// if !ok {
		// 	log.Fatal("expected env var DATABASE_URL")
		// }

		// // open database
		// db, err := sql.Open("postgres", psqlconn)

		db, err := sql.Open("sqlite3", "file::memory:?cache=shared&_foreign_keys=true")
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(-1)
		CheckError(t, err)

		// close database
		defer db.Close()

		// check db
		err = db.Ping()
		CheckError(t, err)

		fmt.Println("Connected!")

		h := GalleryHandler{}
		err = h.Init(db)
		CheckError(t, err)

		srv := &http.Server{
			Handler: h,
			Addr:    "localhost:8080",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	time.Sleep(500 * time.Millisecond)

	var art dto.ArtDto
	var account dto.AccountDto

	if arts := GetArts(t); len(arts) != 0 {
		t.Fatalf("%v length is not 0", arts)
	}

	if resp, err := NewRequest(t, http.MethodPost, "http://localhost:8080/accounts", `{"username":"good", "password":"good"}`, "", ""); err != nil {
		t.Fatal(err)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
			t.Fatal(err)
		} else {
			if account.Username != "good" {
				t.Fatalf("%s is not equal to 'good'", account.Username)
			}
		}
	}

	if resp, err := NewRequest(t, http.MethodGet, fmt.Sprintf("http://localhost:8080/accounts/%d", account.Id), "", "", ""); err != nil {
		t.Fatal(err)
	} else {
		var tmp dto.AccountDto
		if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
			t.Fatal(err)
		} else {
			if tmp.Id != account.Id {
				t.Fatalf("%d is not equal to %d", tmp.Id, account.Id)
			}
		}
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
			if art.AuthorId != account.Id {
				t.Fatalf("the authorId of response (%v) is not equal to '%d'", art, account.Id)
			}
		}
	}

	if arts := GetArts(t); arts[0].Id != art.Id {
		t.Fatalf("the response (%v) does not have art with id %d", arts, art.Id)
	} else if arts[0].AuthorId != account.Id {
		t.Fatalf("the authorId of response (%v) is not equal to '%d'", arts[0].AuthorId, account.Id)
	}

	if _, err := NewRequest(t, http.MethodPut, fmt.Sprintf("http://localhost:8080/arts/%d", art.Id), `{"title":"title2"}`, "good", "bad"); err == nil {
		t.Fatal(err)
	}

	if resp, err := NewRequest(t, http.MethodPut, fmt.Sprintf("http://localhost:8080/arts/%d", art.Id), `{"title":"title2"}`, "good", "good"); err != nil {
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

	if _, err := NewRequest(t, http.MethodDelete, fmt.Sprintf("http://localhost:8080/arts/%d", art.Id), "", "good", "bad"); err == nil {
		t.Fatal(err)
	}

	if arts := GetArts(t); len(arts) != 1 {
		t.Fatalf("%v length is not 1", arts)
	}

	if resp, err := NewRequest(t, http.MethodDelete, fmt.Sprintf("http://localhost:8080/arts/%d", art.Id), "", "good", "good"); err != nil {
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
