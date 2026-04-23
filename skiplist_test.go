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
	values := []string{"golang", "leveldb", "postgres", "redis", "kubernets", "git", "docker", "ollama", "rag", "nginx", "mysql", "java", "rocm", "ubuntu", "npm", "ansible", "bitcoin", "xfs", "awk", "sed", "curl", "tmux", "tr", "tee", "df", "apt", "systemctl", "ssh", "less", "vim", "touch", "mv", "cd", "cp", "scp", "jump", "box", "kafka", "ln"}
	keys = make([]string, len(values))
	for i := range values {
		keys[i] = fmt.Sprintf("%03d", i)
	}
	l := len(keys)
	rand.Shuffle(l, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	sl := &skipList{
		rnd:       rand.New(rand.NewSource(0xabacabed)),
		maxHeight: 1,
	}
	db := sl
	sl.print()
	//fmt.Println(db.DumpKeys())
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
	sl.print()
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
