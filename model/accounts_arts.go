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

func (db *AccountsArtsDB) AddArt(account *dto.AccountDto, art *dto.ArtDto) (*dto.ArtDto, error) {
	if acc, err := db.accountDB.GetAccount(account.Username); err == nil {
		if a, err := db.artDB.StoreArt(*art); err != nil {
			return nil, err
		} else {
			if _, ok := db.authorIdToArts[acc.Id]; !ok {
				db.authorIdToArts[acc.Id] = map[string]*dto.ArtDto{}
			}
			db.authorIdToArts[acc.Id][a.Id] = a
			db.artIdToAuthor[a.Id] = acc

			return a, nil
		}
	} else {
		return nil, err
	}
}

func (db *AccountsArtsDB) IsAuthor(account *dto.AccountDto, id string) bool {
	if acc, err := db.accountDB.GetAccount(account.Username); err == nil {
		if author, ok := db.artIdToAuthor[id]; ok {
			return acc.Id == author.Id
		} else {
			return false
		}
	} else {
		return false
	}
}

func (db *AccountsArtsDB) GetArtsByUsername(username string) ([]*dto.ArtDto, error) {
	if account, err := db.accountDB.GetAccount(username); err == nil {
		v := []*dto.ArtDto{}

		if arts, ok := db.authorIdToArts[account.Id]; ok {
			for _, art := range arts {
				v = append(v, art)
			}
		}

		return v, nil
	} else {
		return nil, err
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
