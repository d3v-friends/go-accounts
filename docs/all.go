package docs

import (
	"fmt"
)

type KV struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type KVs []*KV

func (x KVs) Find(key string) (res *KV, err error) {
	for _, kv := range x {
		if kv.Key == key {
			res = kv
			return
		}
	}
	err = fmt.Errorf("not found kv: key=%s", key)
	return
}

/* ------------------------------------------------------------------------------------------------------------ */
