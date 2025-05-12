package src

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func NewTLSServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", tls)
	server := &http.Server{
		Addr:    "8848",
		Handler: mux,
	}

	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		panic(err)
	}
}

var cacheManager = NewCacheManager()

type CacheManager struct {
	cache map[string]struct{}
	rw    *sync.RWMutex
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		cache: make(map[string]struct{}),
		rw:    new(sync.RWMutex),
	}
}

func (m *CacheManager) Add(key string) bool {
	m.rw.Lock()
	defer m.rw.Unlock()

	var isAdded bool
	if ok := m.Contains(key); ok {
		isAdded = false
	} else {
		m.cache[key] = struct{}{}
		isAdded = true
	}
	return isAdded
}

func (m *CacheManager) Del(key string) {}

func (m *CacheManager) Contains(key string) bool {
	m.rw.RLock()
	defer m.rw.RUnlock()
	_, ok := m.cache[key]
	return ok
}

func tls(w http.ResponseWriter, r *http.Request) {
	body, err := r.GetBody()
	if err != nil {
		return
	}

	b, _ := io.ReadAll(body)
	args := new(struct {
		Keys []string `json:"keys"`
	})
	_ = json.Unmarshal(b, args)

	ans := make([]bool, 0, len(args.Keys))
	for _, key := range args.Keys {
		isAdd := cacheManager.Add(key)
		ans = append(ans, !isAdd)
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(map[string]interface{}{
		"is_add": ans,
	})
	w.Write(response)
}
