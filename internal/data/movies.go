package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Tittle    string    `json:"tittle"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
