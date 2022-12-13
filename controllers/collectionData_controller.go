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

var collectionDataCollection *mongo.Collection = configs.GetCollection(configs.DB, "collectionData")
var CDvalidate = validator.New()

func CreateCollectionData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var collectionData models.CollectionData
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&collectionData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&collectionData); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newcollectionData := models.CollectionData{
		CrateNumber:        collectionData.CrateNumber,
		WeightOfCrop:       collectionData.WeightOfCrop,
		WeightOfOliveOil:   collectionData.WeightOfOliveOil,
		WeightOfNum12:      collectionData.WeightOfNum12,
		WeightOfNum13:      collectionData.WeightOfNum13,
		WeightOfNum14:      collectionData.WeightOfNum14,
		WeightOfNum15:      collectionData.WeightOfNum15,
		WeightOfNum16:      collectionData.WeightOfNum16,
		WeightOfNum17:      collectionData.WeightOfNum17,
		WaitingCrateOfCrop: collectionData.WaitingCrateOfCrop,
	}

	result, err := collectionDataCollection.InsertOne(ctx, newcollectionData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})

}

func GetCollectionData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var collectionData []models.CollectionData
	defer cancel()

	results, err := collectionDataCollection.Find(ctx, bson.M{})
	fmt.Println(results)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var collectionDataToAppend models.CollectionData
		if err = results.Decode(&collectionDataToAppend); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		collectionData = append(collectionData, collectionDataToAppend)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": collectionData}},
	)
}

func GetACollectionData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collectionDataId := c.Params("collectionDataId")
	var OnecollectionData models.CollectionData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(collectionDataId)

	err := collectionDataCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&OnecollectionData)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": OnecollectionData}})

}

func EditACollectionData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collectionDataId := c.Params("collectionDataId")
	var collectionData models.CollectionData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(collectionDataId)

	//validate the request body
	if err := c.BodyParser(&collectionData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&collectionData); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CollaboratorsResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"crateNumber":        collectionData.CrateNumber,
		"weightOfCrop":       collectionData.WeightOfCrop,
		"weightOfOliveOil":   collectionData.WeightOfOliveOil,
		"weightOfNum12":      collectionData.WeightOfNum12,
		"weightOfNum13":      collectionData.WeightOfNum13,
		"weightOfNum14":      collectionData.WeightOfNum14,
		"weightOfNum15":      collectionData.WeightOfNum15,
		"weightOfNum16":      collectionData.WeightOfNum16,
		"weightOfNum17":      collectionData.WeightOfNum17,
		"waitingCrateOfCrop": collectionData.WaitingCrateOfCrop}

	result, err := collectionDataCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated collectionData details``
	var updatedcollectionData models.CollectionData
	if result.MatchedCount == 1 {
		err := collectionDataCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedcollectionData)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedcollectionData}})
}

func DeleteACollectionData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collectionDataId := c.Params("collectionDataId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(collectionDataId)

	result, err := collectionDataCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CollaboratorsResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CollaboratorsResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "collectionData with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CollaboratorsResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "collectionData successfully deleted!"}},
	)
}
