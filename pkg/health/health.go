package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

const (
	StatusDOWN = "DOWN"
	StatusUP   = "UP"
)

type (
	response struct {
		Status string  `json:"status"`           // UP or DOWN
		Checks []Check `json:"checks,omitempty"` // List of checks. Most likely external services
		Data   any     `json:"data,omitempty"`
	}

	Check struct {
		Name     string `json:"name"`              // Name of the external service being checked
		Status   string `json:"status"`            // UP or DOWN
		Critical bool   `json:"critical"`          // If true, the service is considered down
		Message  string `json:"message,omitempty"` // Optional message. Could be used for errors
	}

	// Checker function should perform a check of some sort of external service
	Checker func(ctx context.Context) Check

	system struct {
		Memory memory `json:"memory"`
		CPU    cpu    `json:"-"` // TODO: implement
	}

	memory struct {
		Used int64 `json:"used"` // bytes used on the heap
		Free int64 `json:"free"` // bytes free on the heap
	}

	cpu struct {
		Used int64 `json:"used"` // cpu used in nanoseconds
		Free int64 `json:"free"` // cpu free in nanoseconds
	}
)

func NewHTTPHandler(serviceName string, checkers []Checker) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc(
		"/health/ping",
		handlePing([]byte(fmt.Sprintf("pong from %s", serviceName))),
	)
	router.HandleFunc("/health/ready", handleReady(checkers))
	router.HandleFunc("/health/live", handleLive)

	return router
}

func handlePing(msg []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(msg)
	}
}

func handleReady(checkers []Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := StatusUP
		checks := make([]Check, len(checkers))
		for i, checker := range checkers {
			check := checker(r.Context())
			checks[i] = check
			if check.Status != StatusUP {
				if check.Critical {
					status = StatusDOWN
				}
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		resp, err := json.Marshal(response{
			Status: status,
			Checks: checks,
		})
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("could not marshal response err: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func handleLive(w http.ResponseWriter, r *http.Request) {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	resp, err := json.Marshal(response{
		Status: "UP",
		Data: system{
			Memory: memory{
				Used: int64(memStats.Alloc),
				Free: int64(memStats.Sys - memStats.Alloc),
			},
		},
	})
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("could not marshal response err: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
