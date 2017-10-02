package main

import (
    "fmt"
    "time"
    "log"
    "github.com/boltdb/bolt"
)

func main() {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    go func() {
        put := func(tx *bolt.Tx) error {
            now := time.Now().Format("15:04:05.000")
            fmt.Printf("updating at\t%s\n", now)
            bucket := tx.Bucket([]byte("bukkit"))
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
        bucket := tx.Bucket([]byte("bukit"))
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
