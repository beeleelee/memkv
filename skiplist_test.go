package memkv

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

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
	for i := 0; i < l; i++ {
		if err := db.Put([]byte(keys[i]), []byte(values[i])); err != nil {
			log.Fatal(err)
		}
	}
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
