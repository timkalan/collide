package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

  "collide/simulation"
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

	http.HandleFunc("/simulation", func(w http.ResponseWriter, r *http.Request) {
		sim.Update(dt)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sim); err != nil {
			log.Println("Error encoding JSON:", err)
		}
	})

	go func() {
		for {
			sim.Update(dt)
			// time.Sleep(16 * time.Millisecond)
			time.Sleep(time.Second / time.Duration(framerate))

		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	log.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
