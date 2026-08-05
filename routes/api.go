package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/hub"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"github.com/olabanji12-ojo/church-backend/services"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitAuthController sets up the DI chain for Auth
func InitAuthController(db *mongo.Database) *controllers.AuthController {
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	return controllers.NewAuthController(authService)
}

// InitUserController sets up the DI chain for Users
func InitUserController(db *mongo.Database) *controllers.UserController {
	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	profileService := services.NewProfileService(userRepo)
	notifService := services.NewFirebaseNotificationService()
	swipeService := services.NewSwipeService(userRepo, matchRepo, messageRepo, notifService)

	return controllers.NewUserController(profileService, swipeService)
}

// InitMatchController sets up the DI chain for Matches
func InitMatchController(db *mongo.Database) *controllers.MatchController {
	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	notifService := services.NewFirebaseNotificationService()
	swipeService := services.NewSwipeService(userRepo, matchRepo, messageRepo, notifService)

	return controllers.NewMatchController(swipeService)
}

// InitPrayerController sets up the DI chain for Prayers
func InitPrayerController(db *mongo.Database) *controllers.PrayerController {
	prayerRepo := repositories.NewPrayerRepository(db)
	prayerService := services.NewPrayerService(prayerRepo)

	return controllers.NewPrayerController(prayerService)
}

// InitReportController sets up the DI chain for Reports
func InitReportController(db *mongo.Database) *controllers.ReportController {
	reportRepo := repositories.NewReportRepository(db)
	userRepo := repositories.NewUserRepository(db)
	reportService := services.NewReportService(reportRepo, userRepo)

	return controllers.NewReportController(reportService)
}

// InitChatController sets up the WebSocket Hub and Chat dependencies
func InitChatController(db *mongo.Database, globalHub *hub.Hub) *controllers.ChatController {
	matchRepo := repositories.NewMatchRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	userRepo := repositories.NewUserRepository(db)
	notifService := services.NewFirebaseNotificationService()
	messageService := services.NewMessageService(messageRepo, matchRepo, userRepo, notifService)

	return controllers.NewChatController(globalHub, messageService)
}

// InitRoutes wires everything to the Gorilla Mux Router
func InitRoutes(router *mux.Router, db *mongo.Database, globalHub *hub.Hub) {
	// Base health check
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok", "app": "Church-Match API"}`))
	}).Methods("GET")

	// Initialize Controllers
	authController := InitAuthController(db)
	userController := InitUserController(db)
	matchController := InitMatchController(db)
	prayerController := InitPrayerController(db)
	reportController := InitReportController(db)
	chatController := InitChatController(db, globalHub)

	// Mount Routes
	AuthRoutes(router, authController)
	UserRoutes(router, userController)
	DiscoverRoutes(router, userController)
	MatchRoutes(router, matchController)
	PrayerRoutes(router, prayerController)
	ReportRoutes(router, reportController)
	ChatRoutes(router, chatController)
}
