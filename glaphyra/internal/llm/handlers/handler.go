package handlers

type Handler interface {
	CallAPI(request interface{}) (interface{}, error)
}
