package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User Service 2")
	})

	fmt.Println("User service running on port 8003")
	http.ListenAndServe(":8003", nil)
}
