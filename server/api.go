package server

import "net/http"

const (
	post = "POST"
	get  = "GET"
)

type endpoint struct {
	method  string
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
}

type API interface {
	apiEndpoint() []endpoint
}

type Api struct {
	endpoints []endpoint
}

func (a *Api) apiEndpoint() []endpoint {
	return a.endpoints
}

func NewAPI() *Api {
	api := &Api{}
	api.endpoints = getEndpoints()
	return api
}

func getEndpoints() []endpoint {
	return []endpoint{
		{
			method:  post,
			path:    "/createuser",
			handler: HandleCreateUser,
		},
		{
			method:  post,
			path:    "/login",
			handler: HandleLogin,
		},
		{
			method:  post,
			path:    "/comment",
			handler: HandleComment,
		},
		{
			method:  post,
			path:    "/subcomment",
			handler: HandleCreateSubcomment,
		},
		{
			method:  post,
			path:    "/reaction",
			handler: HandleCreateReaction,
		},
		{
			method:  get,
			path:    "/wall",
			handler: HandleGetWall,
		},
	}
}
