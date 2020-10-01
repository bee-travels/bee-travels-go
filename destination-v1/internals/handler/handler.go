package handler

import "github.com/bee-travels/bee-travels-go/destination-v1/internals/data"

type Handler struct {
	db data.DataProvider
}

func New(provider data.DataProvider) *Handler {
	return &Handler{db: provider}
}
