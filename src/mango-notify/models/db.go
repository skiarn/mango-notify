package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

//DBName is generall db file.
const DBName = "mangonotify.db"

//LinesBucketName is bucket for stored file lines that exist in auth file.
const LinesBucketName = "lines"

//OpenDB opens and initializes buckets if not exist already.
func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open(DBName, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(LinesBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

//GetLine line from db. If no exist nil is returned.
func GetLine(db *bolt.DB, lNumber int) (*Line, error) {
	key := []byte(strconv.Itoa(lNumber))
	var lb []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LinesBucketName))
		lb = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if lb != nil {
		var line Line
		err = json.Unmarshal(lb, &line)
		return &line, nil
	}
	return nil, nil

}

//SaveLine line to db.
func SaveLine(db *bolt.DB, line *Line) error {
	key := []byte(strconv.Itoa(line.Number))

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(LinesBucketName))
		if err != nil {
			return err
		}
		fmt.Println("puting, key:", line.Number)
		data, err := json.Marshal(line)
		if err != nil {
			return err
		}
		fmt.Println("puting, value:", data)
		return b.Put(key, data)
	})
}

//GetUnsentLines returns unsent lines from db.
func GetUnsentLines() (Lines, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var unsentLines Lines

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LinesBucketName))
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(LinesBucketName))
		}
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			var l Line
			err := json.Unmarshal(b.Get([]byte(k)), &l)
			if err != nil {
				return err
			}
			if !l.Sent {
				unsentLines = append(unsentLines, l)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return unsentLines, nil
}

//SetLinesAsSent sets the lines to sent in db.
func SetLinesAsSent(ls *Lines) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	for i, l := range *ls {
		fmt.Println("index:", i, "Value:", l)
		l.Sent = true
		SaveLine(db, &l)
	}
	return nil
}
