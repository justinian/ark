package main

import (
	"log"
	"net/http"
)

type statusSaver struct {
	s int
	w http.ResponseWriter
}

func (s *statusSaver) Status() int                 { return s.s }
func (s *statusSaver) Header() http.Header         { return s.w.Header() }
func (s *statusSaver) Write(b []byte) (int, error) { return s.w.Write(b) }
func (s *statusSaver) WriteHeader(c int)           { s.s = c; s.w.WriteHeader(c) }

func loggingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := statusSaver{s: 200, w: w}
		h.ServeHTTP(&s, r)
		log.Printf("%21s %3d%7s %s", r.RemoteAddr, s.Status(), r.Method, r.URL)
	})
}
