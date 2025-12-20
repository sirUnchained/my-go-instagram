package main

import "net/http"

func (s *server) checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("everything is ok"))
}
