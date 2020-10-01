package data

type DataProvider interface {
	ByCityCountry(city string, country string) Destination
	ByCountry(country string) []Destination
	All() []DestinationList
}
