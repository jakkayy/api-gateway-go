package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User Service 1")
	})

	fmt.Println("User service running on port 8001")
	http.ListenAndServe(":8001", nil)
}
