package main

import "log"

func main() {
	config := Config{
		Port:        "8080",
		DatabaseURL: "postgres://postgres-url",
	}
	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}
