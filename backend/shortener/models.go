package shortener

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/dineshappavoo/basex"
)

// Index to mask ID.
const index = 4242

// Shortener is the struct to be stored in db.
type Shortener struct {
	ID       int
	ShortURL string
	LongURL  string
	Views    int
	LastView time.Time
}

// Insert new shortner to db.
func (s *Shortener) Insert() error {
	id, _ := store.bucket.NextSequence()
	s.ID = int(id)

	shorturl, err := basex.Encode(strconv.Itoa(s.ID + index))
	if err != nil {
		return err
	}
	s.ShortURL = shorturl

	// Marshal shortner data into bytes.
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Persist bytes to users bucket.
	if err := store.bucket.Put(itob(int(id)), buf); err != nil {
		return err
	}
	return nil

}

// Retrieve a shortener from db.
func (s *Shortener) Read() error {
	v := store.bucket.Get(itob(s.ID))
	if v == nil {
		return errors.New("key not found")
	}
	return json.Unmarshal(v, &s)
}

// Update a shortener.
func (s *Shortener) Update() error {
	// Marshal shortner data into bytes.
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Persist bytes to users bucket.
	if err := store.bucket.Put(itob(s.ID), buf); err != nil {
		return err
	}
	return nil
}

// SearchByURL get by longURL.
func (s *Shortener) SearchByURL() error {
	c := store.bucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		ts := &Shortener{}
		if err := json.Unmarshal(v, &ts); err != nil {
			return err
		}

		if ts.LongURL == s.LongURL {
			s.ShortURL = ts.ShortURL
			return nil
		}
	}

	return nil
}

func shortToID(shorturl string) (int, error) {
	id, err := basex.Decode(shorturl)
	if err != nil {
		return 0, err
	}

	intid, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	if (intid - index) < 0 {
		return 0, errors.New("invalid shorturl")
	}
	return intid - index, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
