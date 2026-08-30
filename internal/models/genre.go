package models

// genre structs
type Genre struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	MovieIDs []int64 `json:"movie_ids"`
}

type GenreSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
