package storage

// import (
// 	"log"

// 	"github.com/nutsdb/nutsdb"
// )

// type Storage struct {
// 	*nutsdb.DB
// 	Dir string
// }

// func init() {
// 	db, err := nutsdb.Open(
// 		nutsdb.DefaultOptions,
// 		nutsdb.WithDir("/tmp/nutsdb"),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func createBucket() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			return tx.NewBucket(nutsdb.DataStructureBTree, bucket)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func delete() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			return tx.Delete(bucket, key)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func put() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			val := []byte("val1")
// 			return tx.Put(bucket, key, val, 0)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }
// func put2() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			val := []byte("val2")
// 			return tx.Put(bucket, key, val, 0)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func put3() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name2")
// 			val := []byte("val2")
// 			return tx.PutIfNotExists(bucket, key, val, 0)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func put4() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name2")
// 			val := []byte("val2")
// 			return tx.PutIfExists(bucket, key, val, 100)
// 		}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func read() {
// 	if err := db.View(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			value, err := tx.Get(bucket, key)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println("val:", string(value))

// 			return nil
// 		}); err != nil {
// 		log.Println(err)
// 	}
// }

// func valueLen() {
// 	if err := db.View(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			value, err := tx.ValueLen(bucket, key)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println("value length:", value)

// 			return nil
// 		}); err != nil {
// 		log.Println(err)
// 	}
// }

// func getSet() {
// 	if err := db.Update(
// 		func(tx *nutsdb.Tx) error {
// 			key := []byte("name1")
// 			val := []byte("val3")
// 			oldValue, err := tx.GetSet(bucket, key, val)
// 			if err != nil {
// 				return err
// 			}

// 			fmt.Println("old value :", string(oldValue))

// 			return nil
// 		}); err != nil {
// 		log.Println(err)
// 	}
// }

// func getTTL() {
// 	if err := db.View(func(tx *nutsdb.Tx) error {
// 		key := []byte("name2")
// 		ttl, err := tx.GetTTL(bucket, key)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("ttl :", ttl)
// 		return nil
// 	}); err != nil {
// 		log.Println(err)
// 	}
// }

// func persist() {
// 	if err := db.Update(func(tx *nutsdb.Tx) error {
// 		key := []byte("name2")
// 		return tx.Persist(bucket, key)
// 	}); err != nil {
// 		log.Println(err)
// 	}
// }
