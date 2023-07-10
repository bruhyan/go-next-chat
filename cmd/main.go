package main

import (
	"log"
	"server/db"
	"server/internal/user"
	router "server/router"
)

func main() {
	dbConn, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Failed to initialize database connection %s", err)
	}

	log.Println("Successfully connected to database")

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	router.InitRouter(userHandler)
	_ = router.Start("0.0.0.0:8081")
}
