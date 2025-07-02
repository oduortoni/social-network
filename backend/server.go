package main

import (
	"fmt"
	"net/http"

	"github.com/tajjjjr/social-network/backend/pkg/utils"
	"github.com/tajjjjr/social-network/backend/www/controllers"
)

var (
	Host = "0.0.0.0"
	Port = 9000
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	Port := utils.Port(Port)
	srvAddr := fmt.Sprintf("%s:%d", Host, Port)
	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.Index)

	handler := corsMiddleware(mux)
	http.ListenAndServe(srvAddr, handler)
}
