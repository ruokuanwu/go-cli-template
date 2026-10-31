package repository

import "example.com/go-cli-template/internal/model"

type Greeting struct{}

func (Greeting) Get() (model.Greeting, error) {
	return model.Greeting{Prefix: "Hello"}, nil
}
