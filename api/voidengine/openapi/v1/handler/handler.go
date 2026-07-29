package handler

import "BlackHole/internal/voidengine/service"

type Handler struct {
	services *service.Services
}

func New(services *service.Services) *Handler {
	return &Handler{services: services}
}
