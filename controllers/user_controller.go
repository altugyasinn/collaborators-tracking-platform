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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.User{
		Id:           primitive.NewObjectID(),
		Name:         user.Name,
		CompanyName:  user.CompanyName,
		Location:     user.Location,
		PhoneNumber:  user.PhoneNumber,
		SecondNumber: user.SecondNumber,
		IDNumber:     user.IDNumber,
		LastTonnage:  user.LastTonnage,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})

}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	fmt.Println(results)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var userToAppend models.User
		if err = results.Decode(&userToAppend); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, userToAppend)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var OneUser models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&OneUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": OneUser}})

}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"name":         user.Name,
		"companyName":  user.CompanyName,
		"location":     user.Location,
		"phoneNumber":  user.PhoneNumber,
		"secondNumber": user.SecondNumber,
		"idNumber":     user.IDNumber,
		"lastTonnage":  user.LastTonnage}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details``
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CollaboratorsResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}
