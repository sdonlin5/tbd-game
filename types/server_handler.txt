package types



type ServerHandler interface {
	ResponseHandler(r *game.Result)
}
