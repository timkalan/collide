package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"collide/simulation"

	"github.com/rs/cors"
)

func main() {

	sim := &simulation.Simulation{
		Width:          800,
		Height:         600,
		Gravity:        0,
		SizeMultiplier: 10,
		Paused:         true,
	}

	sim.GenerateBalls(100)

	dt := 1.0
  velocity := -0.001

	pi := &simulation.Pillision{
		Width:            800,
		Height:           600,
		WeightMultiplier: 1,
		NumCollisions:    0,
		Paused:           true,
		SmallSquare: simulation.Square{
			TopLeft:     simulation.Point{X: 100, Y: 250},
			BottomRight: simulation.Point{X: 200, Y: 350},
			Velocity:    0,
			Weight:      1,
		},
		BigSquare: simulation.Square{
			TopLeft:     simulation.Point{X: 300, Y: 200},
			BottomRight: simulation.Point{X: 500, Y: 400},
			Velocity:    velocity,
			Weight:      100,
		},
	}

	pi.SmallSquare.Width = pi.SmallSquare.BottomRight.X - pi.SmallSquare.TopLeft.X
	pi.BigSquare.Width = pi.BigSquare.BottomRight.X - pi.BigSquare.TopLeft.X
	pi.BigSquare.Weight = math.Pow(pi.BigSquare.Weight, pi.WeightMultiplier)

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
			sim.Reset()
			return
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
			sim.Reset()
			return
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
			sim.Reset()
			return
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
			sim.Reset()
		}
	})

	http.HandleFunc("/pausepi", func(w http.ResponseWriter, r *http.Request) {
		pi.Mu.Lock()
		defer pi.Mu.Unlock()
		pi.Paused = true
		pi.Reset(velocity)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/resumepi", func(w http.ResponseWriter, r *http.Request) {
		pi.Mu.Lock()
		defer pi.Mu.Unlock()
		pi.Paused = false
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/weight", func(w http.ResponseWriter, r *http.Request) {
		pi.Mu.Lock()
		defer pi.Mu.Unlock()
		dec := json.NewDecoder(r.Body)

		var data map[string]string
		err := dec.Decode(&data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			pi.Reset(velocity)
			return
		}
		weight, err := strconv.ParseFloat(data["weight"], 64)
		if err != nil {
			log.Println("Error parsing weight:", err)
		}
		pi.BigSquare.Weight = math.Pow(pi.BigSquare.Weight, 1/pi.WeightMultiplier)
		pi.WeightMultiplier = weight
		pi.BigSquare.Weight = math.Pow(pi.BigSquare.Weight, pi.WeightMultiplier)
		sim.UpdateBallSize()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/pillision", func(w http.ResponseWriter, r *http.Request) {
		pi.Mu.Lock()
		defer pi.Mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pi); err != nil {
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
			time.Sleep(16 * time.Millisecond)
		}

	}()

	go func() {
		for {
			if !pi.Paused {
				pi.Update()
			}
			time.Sleep(time.Microsecond)
		}

	}()

	handler := cors.Default().Handler(http.DefaultServeMux)
	log.Println("Starting server on http://localhost:8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
