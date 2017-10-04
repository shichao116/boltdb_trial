package main

import (
    "fmt"
    "time"
    "log"
    "encoding/json"
    "github.com/boltdb/bolt"
)

func main() {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create bucket
    db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        return nil
    })

    // Add key, value pairs
    db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        if err := b.Put([]byte("answer"), []byte("42")); err != nil {
            return err
        }
        if err := b.Put([]byte("question1"), []byte("what")); err != nil {
            return err
        }
        if err := b.Put([]byte("question2"), []byte("where")); err != nil {
            return err
        }
        if err := b.Put([]byte("question3"), []byte("who")); err != nil {
            return err
        }
        if err := b.Put([]byte("question4"), []byte("how")); err != nil {
            return err
        }
        return nil
    })

    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        v := b.Get([]byte("answer"))
        fmt.Printf("The answer is: %s\n", v)
        return nil
    })

    //Iterating over keys
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        c := b.Cursor()

        for k,v := c.First(); k != nil; k,v = c.Next() {
            fmt.Printf("key=%s, value=%s\n", k, v)
        }
        return nil
    })

    /* Accessing variables from BoltDB
    Create an empty instance of the variable outside of the function
    scope and use it in your program later. This ensures that data is
    returned in the correct order and that there are no unexpected
    results*/
    var val []byte
    err = db.Batch(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        val = b.Get([]byte("question1"))
        return nil
    })

    if err != nil {
        return
    }
    fmt.Println("Got value %v", val)

    /* Nested Buckets */
    type Data struct {
        Name string `json:"name"`
    }

    JSONResult := Data{}

    db.Update(func(tx *bolt.Tx) error {
        w, err := tx.CreateBucketIfNotExists([]byte("Primary Bucket"))
        if err != nil {
            fmt.Println(err)
            return err
        }
        x, err := w.CreateBucketIfNotExists([]byte("Secondary Bucket"))
        if err != nil {
            return err
        }
        if err := x.Put([]byte("name"), []byte(`{"name": "chao"}`)); err != nil {
            return err
        }
        return nil
    })
    db.View(func(tx *bolt.Tx) error {
        x := tx.Bucket([]byte("Primary Bucket")).Bucket([]byte("Secondary Bucket"))
        JSON := x.Get([]byte("name"))
        fmt.Println(JSON)
        err := json.Unmarshal(JSON, &JSONResult)
        if err != nil {
            fmt.Println(err)
        }
        return nil
    })
    fmt.Printf("%+v\n", JSONResult)

    go func() {
        put := func(tx *bolt.Tx) error {
            now := time.Now().Format("15:04:05.000")
            fmt.Printf("updating at\t%s\n", now)
            bucket, err := tx.CreateBucketIfNotExists([]byte("bukkit"))
            if err != nil {
                return err
            }
            if err := bucket.Put([]byte("clock"), []byte(now)); err != nil {
                return err
            }
            return nil
        }
        for {
            if err := db.Update(put); err != nil {
                log.Fatal(err)
            }
            time.Sleep(100 * time.Millisecond)
        }
    }()

    time.Sleep(500 * time.Millisecond)

    var result []byte
    get := func(tx *bolt.Tx) error {
        fmt.Println("reading...")
        time.Sleep(500 * time.Millisecond)
        bucket := tx.Bucket([]byte("bukkit"))
        clock := bucket.Get([]byte("clock"))
        fmt.Println("observing %s\n", clock)
        result = make([]byte, len(clock))
        copy(result, clock)
        return nil
    }
    if err := db.View(get); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("result %s\n", result)
}
