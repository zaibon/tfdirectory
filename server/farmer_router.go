package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type farmerRouter struct {
	service tfdirectory.FarmerService
}

func NewFarmerRouter(service tfdirectory.FarmerService, router *mux.Router) *mux.Router {
	farmerRouter := &farmerRouter{service}

	router.HandleFunc("", farmerRouter.createFarmerHandler).Methods("PUT")
	router.HandleFunc("", farmerRouter.ListFarmerHandler).Methods("GET")
	router.HandleFunc("/{node_id}", farmerRouter.getFarmerHandler).Methods("GET")
	return router
}

func (nr *farmerRouter) createFarmerHandler(w http.ResponseWriter, r *http.Request) {
	farmer, err := decodeFarmer(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = nr.service.Insert(context.TODO(), &farmer)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, err)
}

func (nr *farmerRouter) getFarmerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["iyo_organization"]

	farmer, err := nr.service.GetByID(context.TODO(), org)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, farmer)
}

func (nr *farmerRouter) ListFarmerHandler(w http.ResponseWriter, r *http.Request) {
	farmers, err := nr.service.List(context.TODO())
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, farmers)
}

func decodeFarmer(r *http.Request) (farmer tfdirectory.Farmer, err error) {
	if r.Body == nil {
		return farmer, errors.New("no request body")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return farmer, decoder.Decode(&farmer)
}
