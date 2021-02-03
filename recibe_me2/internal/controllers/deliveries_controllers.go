package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetDelivery returns a Delivery
func GetDelivery(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deliveryID := params["id"]

	if !bson.IsObjectIdHex(deliveryID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	objectID := bson.ObjectIdHex(deliveryID)

	result := models.DeliveryModel{}

	err := helpers.DeliveriesCollection.Find(objectID).One(&result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// ListDeliveries returns a List of Deliveries
func ListDeliveries(w http.ResponseWriter, r *http.Request) {
	var results []models.DeliveryModel

	err := helpers.DeliveriesCollection.Find(nil).Sort("-_id").All(&results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// Rate a Delivery Service
func Rate(responseWriter http.ResponseWriter, request *http.Request) {

	var result map[string]interface{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&result)
	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	if result["rating"] == nil {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El campo 'rating' es obligatorio.", nil)
		return
	}

	fmt.Printf("%v\n", result["rating"])

	rating, err := strconv.Atoi(fmt.Sprintf("%v", result["rating"]))

	if err != nil {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Calificación Inválida: Debe ser un número entero.", nil)
		return
	}

	// rating := int(result["rating"].(float64))

	if rating < 1 || rating > 5 {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Calificación Inválida: Debe ser un número entero entre 1 y 5.", nil)
		return
	}

	params := mux.Vars(request)
	deliveryID := params["id"]

	if !bson.IsObjectIdHex(deliveryID) {
		helpers.Response(responseWriter, http.StatusNotFound, constants.ERR_INVALID_DATA, "El ID es inválido.", nil)
		return
	}

	oid := bson.ObjectIdHex(deliveryID)

	delivery := models.DeliveryModel{}

	err = helpers.DeliveriesCollection.FindId(oid).One(&delivery)

	if err != nil {
		helpers.Response(responseWriter, http.StatusNotFound, constants.ERR_NOT_FOUND, err, nil)
		return
	}

	err = helpers.DeliveriesCollection.UpdateId(oid, bson.M{"$set": bson.M{"rating": rating, "rated": true}})

	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_INTERNAL_ERROR, err, nil)
		return
	}

	delivery.Rating = int64(rating)
	delivery.Rated = true

	helpers.Response(responseWriter, http.StatusNotFound, constants.SUCCESS, nil, delivery)
}
