package main

import (
	"fmt"
	"log"

	"lms-vue-go/backend/config"
	"lms-vue-go/backend/routes"
)

func main() {
	// Initialize database connection
	dbConfig := config.DefaultConfig()
	err := config.InitDB(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer config.CloseDB()

	// Menggunakan router yang sudah dibuat
	r := routes.SetupRouter()

	// Jalankan server pada port 3001
	fmt.Println("Server berjalan pada http://localhost:3001")
	err = r.Run(":3001")
	if err != nil {
		log.Fatal("Error menjalankan server:", err)
	}
}
