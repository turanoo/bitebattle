package tests

import (
	"testing"

	"github.com/turanoo/bitebattle/internal/restaurant"
)

func TestNewService(t *testing.T) {
	svc := restaurant.NewService()
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestSearchRestaurants(t *testing.T) {
	svc := restaurant.NewService()
	_, err := svc.SearchRestaurants("pizza", "37.7749,-122.4194", "1000")
	if err == nil {
		t.Error("expected error due to missing Google Places API key")
	}
}
