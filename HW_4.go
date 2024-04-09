package main

import "net/http"

type service struct {
	store map[string]string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("I am alive!"))
	})
	http.ListenAndServe(":8080", mux)
}
