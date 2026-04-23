package memkv

import (
	"bytes"
	"log"
	"testing"
)

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
