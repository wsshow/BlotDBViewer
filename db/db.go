package db

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

type Visiter struct {
	db *bolt.DB
}

type Options struct {
	Timeout time.Duration
}

func Open(dbPath string, options *Options) (*Visiter, error) {
	var err error
	visiter := &Visiter{}
	visiter.db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: options.Timeout})
	return visiter, err
}

func (visiter *Visiter) Close() error {
	return visiter.db.Close()
}

func (visiter *Visiter) Buckets() ([]string, error) {
	var buckets []string
	err := visiter.db.View(func(tx *bolt.Tx) error {
		err := tx.Cursor().Bucket().ForEachBucket(func(k []byte) error {
			buckets = append(buckets, string(k))
			return nil
		})
		return err
	})
	return buckets, err
}

func (visiter *Visiter) DataFromBucket(bucketName string) (map[string]string, error) {
	result := make(map[string]string)
	err := visiter.db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(bucketName)).ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})
		return err
	})
	return result, err
}

func (visiter *Visiter) CreateBucket(bucketName string) error {
	return visiter.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

func (visiter *Visiter) DeleteBucket(bucketName string) error {
	return visiter.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

func (visiter *Visiter) Set(bucketName string, key string, value []byte) error {
	return visiter.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put([]byte(key), value)
	})
}

func (visiter *Visiter) Get(bucketName string, key string) ([]byte, error) {
	var value []byte
	err := visiter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		value = bucket.Get([]byte(key))
		return nil
	})
	return value, err
}

func (visiter *Visiter) Delete(bucketName string, key string) error {
	return visiter.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Delete([]byte(key))
	})
}
