package data

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

type DestinationList struct {
	Country string `json:"country"`
	City    string `json:"city"`
}