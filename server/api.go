package server

import "net/http"

const (
	post = "POST"
	get  = "GET"
)

//endpoint struct have all endpoint related information
type endpoint struct {
	method  string
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
}

type API interface {
	apiEndpoint() []endpoint
}

//Api struct have all the endpoints in it
type Api struct {
	endpoints []endpoint
}

//apiEndpoint returnes all endpoints of API
func (a *Api) apiEndpoint() []endpoint {
	return a.endpoints
}

//NewAPI returns a new Api
func NewAPI() *Api {
	api := &Api{}
	api.endpoints = getEndpoints()
	return api
}

//getEndpoints gives all the endpoints for application
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
