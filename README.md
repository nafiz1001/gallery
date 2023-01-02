# Gallery

A basic CRUD server for art gallery written in Go. 

## Usage

Download all dependencies
```
$ go get ./...
```

Run server
```
$ go run cmd/main.go 
2023/01/02 12:00:10 Listening to localhost:8080
```

Run tests
```
$ go clean -testcache
$ go test -v ./...
?       github.com/nafiz1001/gallery-go/dto     [no test files]
=== RUN   TestGallery
=== RUN   TestGallery/No_art_at_the_beginning
=== RUN   TestGallery/Successfully_create_first_account
=== RUN   TestGallery/First_account_created_exists
=== RUN   TestGallery/Don't_create_art_because_basic_auth_is_missing
=== RUN   TestGallery/Successfully_create_art
=== RUN   TestGallery/The_new_art_exists_with_valid_information
=== RUN   TestGallery/Fail_to_update_art_because_of_invalid_credential
=== RUN   TestGallery/Successfully_update_existing_art
=== RUN   TestGallery/Fail_to_delete_art_because_of_invalid_credential
=== RUN   TestGallery/Successfully_delete_art
--- PASS: TestGallery (0.52s)
    --- PASS: TestGallery/No_art_at_the_beginning (0.00s)
    --- PASS: TestGallery/Successfully_create_first_account (0.00s)
    --- PASS: TestGallery/First_account_created_exists (0.00s)
    --- PASS: TestGallery/Don't_create_art_because_basic_auth_is_missing (0.00s)
    --- PASS: TestGallery/Successfully_create_art (0.00s)
    --- PASS: TestGallery/The_new_art_exists_with_valid_information (0.00s)
    --- PASS: TestGallery/Fail_to_update_art_because_of_invalid_credential (0.00s)
    --- PASS: TestGallery/Successfully_update_existing_art (0.00s)
    --- PASS: TestGallery/Fail_to_delete_art_because_of_invalid_credential (0.00s)
    --- PASS: TestGallery/Successfully_delete_art (0.00s)
PASS
ok      github.com/nafiz1001/gallery-go/handler 0.520s
=== RUN   TestAccountDBInit
--- PASS: TestAccountDBInit (0.00s)
=== RUN   TestCreateAccount
--- PASS: TestCreateAccount (0.00s)
=== RUN   TestGetAccountById
--- PASS: TestGetAccountById (0.00s)
=== RUN   TestGetAccountByUsername
--- PASS: TestGetAccountByUsername (0.00s)
=== RUN   TestArtDBInit
--- PASS: TestArtDBInit (0.00s)
=== RUN   TestCreateArt
--- PASS: TestCreateArt (0.00s)
=== RUN   TestGetArt
--- PASS: TestGetArt (0.00s)
=== RUN   TestGetArts
--- PASS: TestGetArts (0.00s)
=== RUN   TestUpdateArt
--- PASS: TestUpdateArt (0.00s)
=== RUN   TestDeleteArt
--- PASS: TestDeleteArt (0.00s)
PASS
ok      github.com/nafiz1001/gallery-go/model   0.012s
```
