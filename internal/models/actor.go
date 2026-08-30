package models

// actor structs
type Actor struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birth_date"`
	MovieIDs  []int64 `json:"movie_ids"`
}

type ActorSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
