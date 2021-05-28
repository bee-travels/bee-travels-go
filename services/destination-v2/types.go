package main

type Destination struct {
	ID          string   `json:"id"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Population  int      `json:"population"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type Location struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

type Error struct {
	Error string `json:"error"`
}
