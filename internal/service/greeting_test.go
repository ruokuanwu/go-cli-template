package service

import (
	"errors"
	"testing"

	"example.com/go-cli-template/internal/model"
)

type stubGreetings struct {
	err   error
	calls int
}

func (s *stubGreetings) Get() (model.Greeting, error) {
	s.calls++
	return model.Greeting{Prefix: "Hello"}, s.err
}

func TestGreeting(t *testing.T) {
	repo := &stubGreetings{}
	svc := NewGreeting(repo)
	message, err := svc.Greet("Go")
	if err != nil || message != "Hello, Go!" {
		t.Fatalf("message=%q err=%v", message, err)
	}
	if _, err := svc.Greet("  "); err == nil || repo.calls != 1 {
		t.Fatalf("empty name: err=%v calls=%d", err, repo.calls)
	}
	want := errors.New("repository unavailable")
	repo.err = want
	if _, err := svc.Greet("Go"); !errors.Is(err, want) {
		t.Fatalf("error=%v, want wrapped %v", err, want)
	}
}
