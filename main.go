package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/outrigdev/outrig"
)

type Config struct {
	MaxMemoryMB int
	DebugMode   bool
}

// Global state that we can monitor with Outrig
var (
	// Configurable variables
	config Config

	// Memory allocation tracking
	memoryAllocated []byte
	memoryMutex     sync.Mutex

	// Request counter
	requestCount int
	requestMutex sync.Mutex
)

func main() {
	outrigConfig := outrig.DefaultConfig()
	outrigConfig.AppName = "OutrigDemo"
	_, err := outrig.Init(outrigConfig)
	if err != nil {
		slog.Error("Failed to initialize Outrig", "error", err)
		return
	}
	defer outrig.AppDone()
	// Initialize default config
	config.MaxMemoryMB = 100
	config.DebugMode = false

	outrig.WatchFunc("app.config", func() Config {
		return config
	})
	outrig.WatchFunc("app.memory", func() int {
		return len(memoryAllocated)
	})

	outrig.WatchCounterSync("app.request_count", &requestMutex, &requestCount)
	// Set up HTTP server
	http.HandleFunc("/config", handleConfig)
	http.HandleFunc("/memory", handleMemory)
	http.HandleFunc("/stats", handleStats)

	// Start background goroutine for logging
	go backgroundLogger()

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

func backgroundLogger() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		requestMutex.Lock()
		count := requestCount
		requestMutex.Unlock()

		memoryMutex.Lock()
		memSize := len(memoryAllocated)
		memoryMutex.Unlock()

		slog.Info("Background stats",
			"request_count", count,
			"memory_allocated_mb", memSize/1024/1024,
			"debug_mode", config.DebugMode)
	}
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	requestMutex.Lock()
	requestCount++
	requestMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig struct {
		MaxMemoryMB int  `json:"max_memory_mb"`
		DebugMode   bool `json:"debug_mode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.MaxMemoryMB = newConfig.MaxMemoryMB
	config.DebugMode = newConfig.DebugMode

	slog.Info("Config updated", "config", config)
	w.WriteHeader(http.StatusOK)
}

func handleMemory(w http.ResponseWriter, r *http.Request) {
	requestMutex.Lock()
	requestCount++
	requestMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Action string `json:"action"` // "allocate" or "release"
		SizeMB int    `json:"size_mb"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	memoryMutex.Lock()
	defer memoryMutex.Unlock()

	switch req.Action {
	case "allocate":
		if req.SizeMB > config.MaxMemoryMB {
			http.Error(w, "Requested size exceeds maximum", http.StatusBadRequest)
			return
		}
		memoryAllocated = make([]byte, req.SizeMB*1024*1024)
		slog.Info("Memory allocated", "size_mb", req.SizeMB)
	case "release":
		memoryAllocated = nil
		slog.Info("Memory released")
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	requestMutex.Lock()
	requestCount++
	requestMutex.Unlock()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	memoryMutex.Lock()
	memSize := len(memoryAllocated)
	memoryMutex.Unlock()

	stats := struct {
		RequestCount    int  `json:"request_count"`
		MemoryAllocated int  `json:"memory_allocated_mb"`
		DebugMode       bool `json:"debug_mode"`
	}{
		RequestCount:    requestCount,
		MemoryAllocated: memSize / 1024 / 1024,
		DebugMode:       config.DebugMode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
