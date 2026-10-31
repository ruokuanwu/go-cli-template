package service

import (
	"fmt"
	"strings"

	"example.com/go-cli-template/internal/model"
	"github.com/samber/oops"
)

type GreetingRepository interface {
	Get() (model.Greeting, error)
}

type Greeting struct {
	repository GreetingRepository
}

func NewGreeting(repository GreetingRepository) *Greeting {
	return &Greeting{repository: repository}
}

func (s *Greeting) Greet(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", oops.Errorf("name must not be empty")
	}
	greeting, err := s.repository.Get()
	if err != nil {
		return "", oops.Wrapf(err, "get greeting")
	}
	return fmt.Sprintf("%s, %s!", greeting.Prefix, name), nil
}
