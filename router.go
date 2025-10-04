package http

type Response struct {
	Status      int
	StatusText  string
	ContentType string
	Body        []byte
}

type HandlerFunc func(*Request) Response

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]HandlerFunc)}
}
