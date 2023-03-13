package handler

type Client struct{}

func New() Interface {
	return &Client{}
}
