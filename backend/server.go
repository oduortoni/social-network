package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"log"

	"github.com/tajjjjr/social-network/backend/pkg/db/sqlite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tajjjjr/social-network/backend/pkg/db/sqlite"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
	"github.com/tajjjjr/social-network/backend/www/controllers"
)

var (
	Host = "0.0.0.0"
	Port = 9000
)

func main() {
	// Initialize DB and run migrations
	db, err := sqlite.Migration()
	if err != nil {
		panic(fmt.Sprintf("DB migration failed: %v", err))
	}
	defer db.Close()

	Port := utils.Port(Port)
	srvAddr := fmt.Sprintf("%s:%d", Host, Port)
	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)

	http.HandleFunc("/", controllers.Index)

	http.ListenAndServe(srvAddr, nil)
}
