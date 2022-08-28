package urlshort

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

const (
	dbPath     = "bolt.db"
	dbName     = "DB"
	dbPathName = "REDIRECT"
	dbPathKey  = "URL"
)

func storeData(db *bolt.DB, storeData []RedirectPath) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbName)).Bucket([]byte(dbPathName))

		for _, store := range storeData {
			err := b.Put([]byte(store.Path), []byte(store.Url))

			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return nil
}

func getDb() *bolt.DB {
	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		panic(err)
	}

	return db
}

func getUrlByPath(url string) (string, error) {
	db := getDb()
	defer db.Close()

	var foundPath string

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbName)).Bucket([]byte(dbPathName))
		foundPath = string(b.Get([]byte(url)))
		return nil
	})

	if err != nil {
		return "", err
	}

	return foundPath, nil
}

func getRedirectPaths(jsonPath string) ([]RedirectPath, error) {
	jsonData, err := readDataFile(jsonPath)

	if err != nil {
		return nil, err
	}

	var pathData []RedirectPath

	if err := json.Unmarshal(jsonData, &pathData); err != nil {
		return nil, err
	}

	return pathData, nil
}

func SetupDb(jsonPath string) {
	db := getDb()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(dbName))

		if err != nil {
			return err
		}

		_, err = root.CreateBucketIfNotExists([]byte(dbPathName))

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	redirectData, err := getRedirectPaths(jsonPath)

	if err != nil {
		panic(err)
	}

	storeData(db, redirectData)
}
