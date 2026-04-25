package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Order Service")
	})

	fmt.Println("Order service running on port 8002")
	http.ListenAndServe(":8002", nil)
}
