package memkv

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestPutGetDelete(t *testing.T) {
	keys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	values := []string{"golang", "leveldb", "postgres", "redis", "kubernets", "git", "docker", "ollama", "rag", "nginx"}
	l := len(keys)
	rand.Shuffle(l, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	db := New()
	fmt.Println(db.DumpKeys())
	for i := 0; i < l; i++ {
		if err := db.Put([]byte(keys[i]), []byte(values[i])); err != nil {
			log.Fatal(err)
		}
	}
	dks, dvs := db.Dump()
	for i, k := range dks {
		fmt.Println(string(k), string(dvs[i]))
	}
	for i := 0; i < l; i++ {
		if rv, err := db.Get([]byte(keys[i])); err != nil {
			log.Fatal(err)
		} else if !bytes.Equal(rv, []byte(values[i])) {
			log.Fatal("value not match", string(rv), values[i])
		}
	}
	for i := 0; i < l; i++ {
		if err := db.Delete([]byte(keys[i])); err != nil {
			log.Fatal(err)
		}
	}
	if db.Num() != 0 {
		log.Fatal("incorrect number of keys")
	}
}

func TestInitSkiplist(t *testing.T) {
	_ = New()
}

func TestPutGet(t *testing.T) {
	db := New()

	k := []byte("test put_get key")
	v := []byte("test put_get value")
	if err := db.Put(k, v); err != nil {
		log.Fatal(err)
	}
	if rv, err := db.Get(k); err != nil {
		log.Fatal(err)
	} else if !bytes.Equal(v, rv) {
		log.Fatal("value not match", v, rv)
	}
}
