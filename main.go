package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/hub"
	"github.com/olabanji12-ojo/church-backend/middleware"
	"github.com/olabanji12-ojo/church-backend/routes"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Connect to Databases & External Services
	db := database.ConnectDB()
	database.ConnectRedis()
	config.InitFirebase()

	// 3. Create and Start the Global WebSocket Hub
	globalHub := hub.NewHub()
	go globalHub.Run()

	// 4. Create Router & Initialize Routes
	mainRouter := mux.NewRouter()
	routes.InitRoutes(mainRouter, db, globalHub)

	// 5. Middleware Chain with Negroni (Matching Car Wash App)
	n := negroni.New()
	
	// Recovery middleware
	n.Use(negroni.NewRecovery())
	
	// CORS middleware
	n.Use(middleware.Cors())
	
	// Secure headers middleware
	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		secureMiddleware := middleware.Secure()
		secureMiddleware.HandlerFuncWithNext(w, r, next)
	}))

	// Attach router
	n.UseHandler(mainRouter)

	// 6. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Println("🚀 Church-Match API Listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, n); err != nil {
		logrus.Fatal("Server failed to start: ", err)
	}
}
