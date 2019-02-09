package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git/sirupsen/logrus"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dbs, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		logrus.Fatal(err)

	}
	defer dbs.Close()

	// 自动迁移模式
	dbs.AutoMigrate(&Product{})

	// 创建
	dbs.Create(&Product{Code: "L1212", Price: 1000})

	// 读取
	var product Product
	dbs.First(&product, 1)                   // 查询id为1的product
	dbs.First(&product, "code = ?", "L1212") // 查询code为l1212的product

	// 更新 - 更新product的price为2000
	dbs.Model(&product).Update("Price", 2000)

	// 删除 - 删除product
	dbs.Delete(&product)

	//https: //www.cnblogs.com/yxdz-hit/p/8536094.html
	//https://www.jianshu.com/p/ee494d459f2c

	//https://jasperxu.github.io/gorm-zh/database.html#dbc
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}

	//2.写数据库
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))

		var err error
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("firstBucket"))
			if err != nil {
				fmt.Println("createBucket failed!", err)
				os.Exit(1)
			}
		}

		bucket.Put([]byte("aaaa"), []byte("AASDD!"))
		bucket.Put([]byte("bbbbc"), []byte("HelloItcast!"))

		bucket.Put([]byte("aaadd"), []byte("HelloWorld!"))
		bucket.Put([]byte("aaabc"), []byte("HelloItcast!"))
		return nil
	})

	//3.读取数据库
	var value []byte

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("firstBucket"))
		if bucket == nil {
			fmt.Println("Bucket is nil!")
			os.Exit(1)
		}

		value = bucket.Get([]byte("aaaa"))
		fmt.Println("aaaa => ", string(value))
		value = bucket.Get([]byte("bbbb"))
		fmt.Println("bbbb => ", string(value))

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("firstBucket"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})

	//prefix search
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte("firstBucket")).Cursor()

		prefix := []byte("aaa")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("M: key=%s, value=%s\n", k, v)
		}

		return nil
	})

	//foreach
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("firstBucket"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

	// Grab the initial stats.
	prev := db.Stats()

	for {
		// Wait for 10s.
		time.Sleep(10 * time.Second)

		// Grab the current stats and diff them.
		stats := db.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		json.NewEncoder(os.Stderr).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}

	//read only
	/*
		dba, erra := bolt.Open("my.db", 0666, &bolt.Options{ReadOnly: true})
		if erra != nil {
			log.Fatal(erra)
		}
	*/

}
