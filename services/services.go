package services

import (
	"fmt"
)

// Service baseline service
type Service struct {
	Name    string
	baseURL string
}

// String for debugging
func (s Service) String() string {
	return fmt.Sprintf("((%v basUrl: %v))", s.Name, s.baseURL)
}

// Make create a new instance of a service
func Make(service string, name string, url string) (Service, error) {
	var s Service
	if service == "radarr" {
		// s = new(*Radarr)
	} else if service == "sonarr" {
		// s = Sonarr{}
		// s = make(Sonarr)
	} else {
		return Service{}, fmt.Errorf("Unknown service %s", service)
	}
	s.Name = name
	s.baseURL = url
	return s, nil
}
