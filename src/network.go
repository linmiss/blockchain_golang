package main

import (
	"net/http"
	"time"
)

func run() error {
	mux := makeMuxRouter()
	httpAddress := "8080"
	s := &http.Server{
		Addr:           ":" + httpAddress,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, //2^20
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
