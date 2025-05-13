package src

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func NewTLSServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cache/add", cacheAdd)
	server := &http.Server{
		Addr:    ":8848",
		Handler: mux,
	}

	if err := server.ListenAndServeTLS("certificate.crt", "private.key"); err != nil {
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
	if _, ok := m.cache[key]; ok {
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

func cacheAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, _ := io.ReadAll(r.Body)
	keys := []string{}
	_ = json.Unmarshal(b, &keys)

	ans := make([]bool, 0, len(keys))
	for _, key := range keys {
		isAdd := cacheManager.Add(key)
		ans = append(ans, !isAdd)
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(ans)
	w.Write(response)
}
