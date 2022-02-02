package handlers

import "github.com/julienschmidt/httprouter"

// Handler - interface for creating handlers
type Handler interface {
	Register(router *httprouter.Router)
}
