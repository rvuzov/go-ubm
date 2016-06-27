package bclient

import "fmt"

type (
	Client struct {
		Addr    string
		APIAddr string
	}
)

func NewClient(addr string) Client {
	return Client{
		Addr:    addr,
		APIAddr: fmt.Sprintf("http://%s/", addr),
	}
}
