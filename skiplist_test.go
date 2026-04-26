package memkv

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestContainsFind(t *testing.T) {
	keys := []string{"drink", "sport", "season", "fruit", "vegetable", "city"}
	vals := []string{"milk", "football", "winter", "apple", "lettuce", "shanghai"}
	db := New()
	for i, key := range keys {
		if err := db.Put([]byte(key), []byte(vals[i])); err != nil {
			log.Fatal(err)
		}
	}
	if contains := db.Contains([]byte("culture")); contains {
		log.Fatal("Should note contains key: culture")
	}
	for _, key := range keys {
		if contains := db.Contains([]byte(key)); !contains {
			log.Fatal("Should contains key: ", key)
		}
	}
	if rk, v, err := db.Find([]byte("season")); err != nil || string(rk) != "season" || string(v) != "winter" {
		log.Fatal("Should find key season")
	}
	if rk, v, err := db.Find([]byte("culture")); err != nil || string(rk) != "drink" || string(v) != "milk" {
		log.Fatal("Should find key drink greater than culture")
	}
	if rk, v, err := db.Find([]byte("music")); err != nil || string(rk) != "season" || string(v) != "winter" {
		log.Fatal("Should find key season greater than music")
	}
	if rk, v, err := db.Find([]byte("zoo")); err == nil || rk != nil || v != nil {
		log.Fatal("Should not find key greater than zoo")
	}
}
func TestPutGetDelete(t *testing.T) {
	var keys []string
	vals := []string{"golang", "leveldb", "postgres", "redis", "kubernets", "git", "docker", "ollama", "rag", "nginx", "mysql", "java", "rocm", "ubuntu", "npm", "ansible", "bitcoin", "xfs", "awk", "sed", "curl", "tmux", "tr", "tee", "df", "apt", "systemctl", "ssh", "less", "vim", "touch", "mv", "cd", "cp", "scp", "jump", "box", "kafka", "ln", "cat", "tar", "which", "alias", "date", "echo", "cut", "ping", "uptime", "ps", "top"}
	var values []string
	for i := 0; i < 10; i++ {
		values = append(values, vals...)
	}
	keys = make([]string, len(values))
	for i := range values {
		keys[i] = fmt.Sprintf("%03d", i)
	}
	l := len(keys)
	rand.Shuffle(l, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	sl := newSkipList()
	db := sl
	// put empty value first
	for i := 0; i < l; i++ {
		if err := db.Put([]byte(keys[i]), nil); err != nil {
			log.Fatal(err)
		}
	}
	// sl.print()
	for i := 0; i < l; i++ {
		if err := db.Put([]byte(keys[i]), []byte(values[i])); err != nil {
			log.Fatal(err)
		}
	}
	//sl.print()
	for i := 0; i < l; i++ {
		if rv, err := db.Get([]byte(keys[i])); err != nil {
			log.Fatal(err)
		} else if !bytes.Equal(rv, []byte(values[i])) {
			log.Fatal("value not match", string(rv), values[i])
		}
	}
	// sl.print()
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
