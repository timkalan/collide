package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"collide/simulation"

	"github.com/rs/cors"
)

func main() {
	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Error getting ball number", err)
	}

	sim := &simulation.Simulation{
		Width:   800,
		Height:  600,
		Gravity: 0,
		Balls:   simulation.GenerateBalls(n),
	}

	framerate, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Error parsing framerate:", err)
	}
	// dt := 1.0 / float64(framerate)
	dt := 1.0

	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		sim.Paused = true
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/resume", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		sim.Paused = false
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/startwar", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		sim.War = true
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/stopwar", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		sim.War = false
    sim.RandomizeColors()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/simulation", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sim); err != nil {
			log.Println("Error encoding JSON:", err)
		}
	})

	go func() {
		for {
			if !sim.Paused {
				start := time.Now()
				sim.Update(dt)
				duration := time.Since(start)
				log.Println("Frametime:", duration)
			}
			// time.Sleep(16 * time.Millisecond)
			time.Sleep(time.Second / time.Duration(framerate))
		}

	}()

	// http.Handle("/", http.FileServer(http.Dir("./frontend")))
	handler := cors.Default().Handler(http.DefaultServeMux)
	log.Println("Starting server on http://localhost:8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
