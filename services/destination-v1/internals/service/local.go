package service

import "github.com/bee-travels/bee-travels-go/services/destination-v1/internals/data"

type LocalDB struct {
	destination []data.Destination
}

func NewLocalDB(destination []data.Destination) LocalDB {
	return LocalDB{destination: destination}
}

func (l LocalDB) ByCityCountry(city, country string) data.Destination {
	var res data.Destination
	for _, destination := range l.destination {
		if normalize(destination.City) == city && normalize(destination.Country) == country {
			res = destination
			break
		}
	}
	return res
}

func (l LocalDB) ByCountry(country string) []data.Destination {
	res := make([]data.Destination, 0)

	for _, destination := range l.destination {
		if normalize(destination.Country) == country {
			res = append(res, destination)
		}
	}
	return res
}

func (l LocalDB) All() []data.DestinationList {
	res := make([]data.DestinationList, len(l.destination))
	for i, destination := range l.destination {
		res[i] = data.DestinationList{
			Country: destination.Country,
			City:    destination.City,
		}
	}
	return res
}
