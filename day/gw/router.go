package gw

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}
