package shortener

import (
	"time"

	"github.com/boltdb/bolt"
)

var store struct {
	db     *bolt.DB
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

// ConnectDB open the db and create the bucket.
func ConnectDB(bucket string) error {
	db, err := bolt.Open("/data/boltdb.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}

	b := tx.Bucket([]byte(bucket))

	store.db = db
	store.tx = tx
	store.bucket = b

	return nil
}

func CloseDB() {
	store.db.Close()
}
