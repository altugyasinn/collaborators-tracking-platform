package controllers

import (
	"collaborators-tracking-platforms/configs"
	"collaborators-tracking-platforms/models"
	"collaborators-tracking-platforms/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dailyRequestsCollection *mongo.Collection = configs.GetCollection(configs.DB, "dailyRequests")
var DRvalidate = validator.New()

func CreateDailyRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var dailyRequests models.DailyRequests
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&dailyRequests); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&dailyRequests); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newDailyRequests := models.DailyRequests{
		Id:                primitive.NewObjectID(),
		CrateRequest:      dailyRequests.CrateRequest,
		MoneyRequest:      dailyRequests.MoneyRequest,
		EmptyCrateRequest: dailyRequests.EmptyCrateRequest,
	}

	result, err := dailyRequestsCollection.InsertOne(ctx, newDailyRequests)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})

}

func GetDailyRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var dailyRequests []models.DailyRequests
	defer cancel()

	results, err := dailyRequestsCollection.Find(ctx, bson.M{})
	fmt.Println(results)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var dailyRequestsToAppend models.DailyRequests
		if err = results.Decode(&dailyRequestsToAppend); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		dailyRequests = append(dailyRequests, dailyRequestsToAppend)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": dailyRequests}},
	)
}

func GetADailyRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	dailyRequestsId := c.Params("dailyRequestsId")
	var OneDailyRequests models.DailyRequests
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(dailyRequestsId)

	err := dailyRequestsCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&OneDailyRequests)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": OneDailyRequests}})

}

func EditADailyRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	dailyRequestsId := c.Params("dailyRequestsId")
	var dailyRequests models.DailyRequests
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(dailyRequestsId)

	//validate the request body
	if err := c.BodyParser(&dailyRequests); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&dailyRequests); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"crateRequest":      dailyRequests.CrateRequest,
		"moneyRequest":      dailyRequests.MoneyRequest,
		"emptyCrateRequest": dailyRequests.EmptyCrateRequest,
		"postScript":        dailyRequests.PostScript}

	result, err := dailyRequestsCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated dailyRequests details``
	var updatedDailyRequests models.DailyRequests
	if result.MatchedCount == 1 {
		err := dailyRequestsCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedDailyRequests)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedDailyRequests}})
}

func DeleteADailyRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	dailyRequestsId := c.Params("dailyRequestsId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(dailyRequestsId)

	result, err := dailyRequestsCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CollaboratorsResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "dailyRequests with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "dailyRequests successfully deleted!"}},
	)
}
