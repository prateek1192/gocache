package fcache

import (
  "fmt"
  "sync"
  "crypto/sha1"
  "time"
  "log"
)

type Cache map[string]*CacheShard
var cache_life time.Duration
type CacheShard struct {
  items map[string][]byte
  startTime time.Time
  lock *sync.RWMutex
  defined bool
}

func New(expiration_time time.Duration) Cache {
  cache_life = expiration_time
  c := make(Cache, 256)
  for i := 0; i < 256; i++ {
    c[fmt.Sprintf("%02x", i)] = &CacheShard{
      items: make(map[string][]byte, 2048),
      lock: new(sync.RWMutex),
    }
  }
  go expire_and_renew_element(c)
  return c
}

func (c Cache) Get(key string) []byte {
  shard := c.GetShard(key)
  shard.lock.RLock()
  defer shard.lock.RUnlock()
  return shard.items[key]
}

func (c Cache) Set(key string, data []byte) {
  shard := c.GetShard(key)
  shard.lock.Lock()
  defer shard.lock.Unlock()
  shard.items[key] = data
  shard.startTime = time.Now()
  shard.defined = true
}

func (c Cache) GetShard(key string) (shard *CacheShard) {
  hasher := sha1.New()
  hasher.Write([]byte(key))
  shardKey := fmt.Sprintf("%x", hasher.Sum(nil))[0:2]
  return c[shardKey]
}

func expire_and_renew_element(c Cache) {
	log.Printf("Starting expire and renew")
	for _, shard := range c {
		if time.Now().Sub(shard.startTime) > cache_life && shard.defined == true {
			log.Printf("Renewing element")
			shard.startTime = time.Now()
		}
	}
   time.Sleep(10 * time.Second)
   expire_and_renew_element(c)
}
