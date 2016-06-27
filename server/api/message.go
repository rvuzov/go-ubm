package api

type (
	APIMessage interface {
		Receive() interface{}
	}
)
