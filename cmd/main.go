package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vickon16/third-backend-tutorial/cmd/db"
	"github.com/vickon16/third-backend-tutorial/cmd/project"
	"github.com/vickon16/third-backend-tutorial/cmd/sqlc"
	"github.com/vickon16/third-backend-tutorial/cmd/user"
	"github.com/vickon16/third-backend-tutorial/cmd/utils"
	"github.com/vickon16/third-backend-tutorial/cmd/validator"
)

func main() {
	router := mux.NewRouter()

	conn, err := db.NewMySQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	validator.RegisterCustomValidators()

	db := sqlc.New(conn)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, 200, map[string]string{
			"message": "All good",
		})
	})

	user.NewHandler(db).RegisterRoutes(router)
	project.NewHandler(db).RegisterRoutes(router)

	log.Println("Server is running on port 4000")
	if err = http.ListenAndServe(":"+utils.Configs.PORT, router); err != nil {
		log.Fatal(err)
	}
}
