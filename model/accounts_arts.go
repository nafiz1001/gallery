package model

import (
	"github.com/nafiz1001/gallery-go/dto"
)

type AccountsArtsDB struct {
	accountDB      *AccountDB
	artDB          *ArtDB
	authorIdToArts map[string]map[string]*dto.ArtDto
	artIdToAuthor  map[string]*dto.AccountDto
}

func (db *AccountsArtsDB) Init(accountDB *AccountDB, artDB *ArtDB) error {
	db.accountDB = accountDB
	db.artDB = artDB
	db.authorIdToArts = map[string]map[string]*dto.ArtDto{}
	db.artIdToAuthor = map[string]*dto.AccountDto{}
	return nil
}

func (db *AccountsArtsDB) AddArt(account dto.AccountDto, art dto.ArtDto) (*dto.ArtDto, error) {
	if acc, err := db.accountDB.GetAccount(account.Username); err != nil {
		return nil, err
	} else {
		if a, err := db.artDB.CreateArt(art); err != nil {
			return nil, err
		} else {
			if _, ok := db.authorIdToArts[acc.Id]; !ok {
				db.authorIdToArts[acc.Id] = map[string]*dto.ArtDto{}
			}
			db.authorIdToArts[acc.Id][a.Id] = a
			db.artIdToAuthor[a.Id] = acc

			return a, nil
		}
	}
}

func (db *AccountsArtsDB) IsAuthor(account dto.AccountDto, id string) bool {
	if acc, err := db.accountDB.GetAccount(account.Username); err != nil {
		return false
	} else {
		if author, ok := db.artIdToAuthor[id]; ok {
			return acc.Id == author.Id
		} else {
			return false
		}
	}
}

func (db *AccountsArtsDB) GetArtsByUsername(username string) ([]dto.ArtDto, error) {
	v := []dto.ArtDto{}

	if account, err := db.accountDB.GetAccount(username); err != nil {
		return v, err
	} else {
		if arts, ok := db.authorIdToArts[account.Id]; ok {
			for _, art := range arts {
				v = append(v, *art)
			}
		}

		return v, nil
	}
}

func (db *AccountsArtsDB) DeleteArt(id string) (*dto.ArtDto, error) {

	if art, err := db.artDB.DeleteArt(id); err != nil {
		return nil, err
	} else {
		delete(db.artIdToAuthor, id)
		for authorId := range db.authorIdToArts {
			delete(db.authorIdToArts[authorId], id)
		}

		return art, nil
	}
}
