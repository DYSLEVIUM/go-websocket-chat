package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chat/config"
)

var PORT string

func init() {
	config.LoadEnv()

	PORT = os.Getenv("PORT")
}

func main() {
	manager := config.NewManager()

	// r := mux.NewRouter()
	// r.HandleFunc("/ws", manager.ServesWS)

	http.HandleFunc("/ws", manager.ServesWS)

	// corsObj := handlers.AllowedOrigins([]string{"*"})
	// corsObj := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	// log.Fatal(http.ListenAndServe(":"+PORT, handlers.CORS(corsObj)(r)))

	log.Fatal(http.ListenAndServe(":"+PORT, nil))

	// app := fiber.New()
	// Routes(app)
	// log.Fatal(app.Listen(":" + PORT))
}
