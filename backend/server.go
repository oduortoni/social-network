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

func main() {
	Port := utils.Port(Port)
	srvAddr := fmt.Sprintf("%s:%d", Host, Port)
	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)

	http.HandleFunc("/", controllers.Index)

	http.ListenAndServe(srvAddr, nil)
}
