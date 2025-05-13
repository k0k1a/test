package src

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

type CacheManager struct {
	rw    *sync.RWMutex
	cache map[string]struct{}
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

func (m *CacheManager) Clear() {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.cache = make(map[string]struct{})
}

func (m *CacheManager) HandleHttpCacheAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var keys []string
	b, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(b, &keys)

	ans := make([]bool, 0, len(keys))
	for _, key := range keys {
		isAdd := m.Add(key)
		ans = append(ans, !isAdd)
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(ans)
	w.Write(response)
}

type HttpsServer struct {
	server       *http.Server
	cacheManager *CacheManager
}

func StartHttpsServer(addr string) (*HttpsServer, error) {
	manager := NewCacheManager()
	mux := http.NewServeMux()
	mux.HandleFunc("/cache/add", manager.HandleHttpCacheAdd)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServeTLS("certificate.crt", "private.key"); errors.Is(err, http.ErrServerClosed) {
		} else if err != nil {
			panic(err)
		}
	}()

	return &HttpsServer{
		server:       httpServer,
		cacheManager: manager,
	}, nil
}
