package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Record struct {
	Version   string    `json:"version"`
	Name      string    `json:"name"`
	UUID      string    `json:"uuid"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Course    float64   `json:"course"`
	Speed     float64   `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

func Generate(w io.Writer, count int) error {
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < count; i++ {
		// Generate a unique UUID with mrn prefix
		id := uuid.New()
		mrnUUID := fmt.Sprintf("urn:mrn:signalk:uuid:%s", id.String())
		
		// San Francisco Bay area coordinates
		baseLat := 37.78
		baseLong := -122.38
		
		record := Record{
			Version:   "1.0.0",
			Name:      fmt.Sprintf("Boat %d", rand.Intn(100)+1),
			UUID:      mrnUUID,
			Latitude:  baseLat + (rand.Float64()-0.5)*0.02,
			Longitude: baseLong + (rand.Float64()-0.5)*0.02,
			Altitude:  0.0,
			Course:    rand.Float64() * 360,
			Speed:     float64(rand.Intn(15) + 1),
			Timestamp: time.Now(),
		}
		
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}
		
		if _, err := fmt.Fprintln(w, string(data)); err != nil {
			return err
		}
	}
	
	return nil
}