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
		Width:          800,
		Height:         600,
		Gravity:        0,
		SizeMultiplier: 10,
		Paused:         true,
	}

	sim.GenerateBalls(n)

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

	http.HandleFunc("/changeballnumber", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		dec := json.NewDecoder(r.Body)

		var data map[string]int
		err := dec.Decode(&data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
		}
		numBalls := data["n"]
		sim.GenerateBalls(numBalls)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/gravity", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		dec := json.NewDecoder(r.Body)

		var data map[string]string
		err := dec.Decode(&data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
		}
		gravity, err := strconv.ParseFloat(data["gravity"], 64)
		if err != nil {
			log.Println("Error parsing gravity:", err)
		}
		sim.Gravity = gravity
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/size", func(w http.ResponseWriter, r *http.Request) {
		sim.Mu.Lock()
		defer sim.Mu.Unlock()
		dec := json.NewDecoder(r.Body)

		var data map[string]string
		err := dec.Decode(&data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
		}
		size, err := strconv.ParseFloat(data["size"], 64)
		if err != nil {
			log.Println("Error parsing size:", err)
		}
		sim.OldSizeMultiplier = sim.SizeMultiplier
		sim.SizeMultiplier = size
		sim.UpdateBallSize()
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
				sim.FPS = 1.0 / duration.Seconds()
			}
			// time.Sleep(16 * time.Millisecond)
			time.Sleep(time.Second / time.Duration(framerate))
		}

	}()

	handler := cors.Default().Handler(http.DefaultServeMux)
	log.Println("Starting server on http://localhost:8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
