package routes

import (
	controller "collaborators-tracking-platforms/controllers"

	"github.com/gofiber/fiber/v2"
)

func CollaboratorsRoute(app *fiber.App) {

	// User Endpoints
	app.Post("/user", controller.CreateUser)
	app.Get("/user/:userId", controller.GetAUser)
	app.Put("/user/:userId", controller.EditAUser)
	app.Delete("/user/:userId", controller.DeleteAUser)
	app.Get("/users", controller.GetUsers)

	//Daily Requests Endpoints
	app.Post("/dailyRequest", controller.CreateDailyRequests)
	app.Get("/dailyRequest/:dailyRequestId", controller.GetADailyRequests)
	app.Put("/dailyRequest/:dailyRequestId", controller.EditADailyRequests)
	app.Delete("/dailyRequest/:dailyRequestId", controller.DeleteADailyRequests)
	app.Get("/dailyRequests", controller.GetDailyRequests)

	// Collection Data Endpoints
	app.Post("/collectionData", controller.CreateCollectionData)
	app.Get("/collectionData/:collectionDataId", controller.GetACollectionData)
	app.Put("/collectionData/:collectionDataId", controller.EditACollectionData)
	app.Delete("/collectionData/:collectionDataId", controller.DeleteACollectionData)
	app.Get("/collectionDatas", controller.GetCollectionData)

}
