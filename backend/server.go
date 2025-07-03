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
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := utils.Port(Port) // Fixed variable shadowing
	srvAddr := fmt.Sprintf("%s:%d", Host, port)
	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.Index)
	
	handler := corsMiddleware(mux)
	http.ListenAndServe(srvAddr, handler)
}