package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type farmerRouter struct {
	service tfdirectory.FarmerService
}

func NewFarmerRouter(service tfdirectory.FarmerService, router *mux.Router) *mux.Router {
	farmerRouter := &farmerRouter{service}

	router.HandleFunc("", farmerRouter.createFarmerHandler).Methods("POST")
	router.HandleFunc("/{organization}", farmerRouter.updateFarmerHandler).Methods("PUT")
	router.HandleFunc("", farmerRouter.ListFarmerHandler).Methods("GET")
	router.HandleFunc("/{organization}", farmerRouter.getFarmerHandler).Methods("GET")
	return router
}

func (nr *farmerRouter) createFarmerHandler(w http.ResponseWriter, r *http.Request) {
	scope := ScopeFromContext(r.Context())
	if len(scope) < 1 {
		Error(w, http.StatusUnauthorized, "missing required scope")
		return
	}

	farmer, err := decodeFarmer(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !isSubset(scope, []string{
		fmt.Sprintf("user:memberof:%s", farmer.Organization),
	}) {
		Error(w, http.StatusUnauthorized, "missing required scope")
		return
	}

	err = nr.service.Insert(context.TODO(), &farmer)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, err)
}

func (nr *farmerRouter) updateFarmerHandler(w http.ResponseWriter, r *http.Request) {
	scope := ScopeFromContext(r.Context())
	if len(scope) < 1 {
		Error(w, http.StatusUnauthorized, "missing required scope")
		return
	}

	farmer, err := decodeFarmer(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	vars := mux.Vars(r)
	org := vars["organization"]

	if farmer.Organization != org {
		Error(w, http.StatusBadRequest, "cannot modify the organization of a farm")
		return
	}

	if !isSubset(scope, []string{
		fmt.Sprintf("user:memberof:%s", org),
	}) {
		Error(w, http.StatusUnauthorized, "missing required scope")
		return
	}

	err = nr.service.Update(context.TODO(), &farmer)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, farmer)
}

func (nr *farmerRouter) getFarmerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["organization"]

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
