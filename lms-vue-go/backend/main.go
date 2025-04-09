package main

import (
	"fmt"
	"lms-vue-go/backend/routes"
)

func main() {
	// Menggunakan router yang sudah dibuat
	r := routes.SetupRouter()

	// Jalankan server pada port 3000
	fmt.Println("Server berjalan pada http://localhost:3000")
	r.Run(":3000")
} 