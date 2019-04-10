package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type nodeRouter struct {
	nodeService tfdirectory.NodeService
}

func NewUserRouter(u tfdirectory.NodeService, router *mux.Router) *mux.Router {
	nodeRouter := &nodeRouter{u}

	router.HandleFunc("/", nodeRouter.createNodeHandler).Methods("PUT")
	router.HandleFunc("/", nodeRouter.ListNodeHandler).Methods("GET")
	router.HandleFunc("/{node_id}", nodeRouter.getNodeHandler).Methods("GET")
	return router
}

func (nr *nodeRouter) createNodeHandler(w http.ResponseWriter, r *http.Request) {
	node, err := decodeNode(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = nr.nodeService.Register(context.TODO(), &node)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, err)
}

func (nr *nodeRouter) getNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	nodeID := vars["node_id"]

	node, err := nr.nodeService.GetByID(context.TODO(), nodeID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, node)
}

func (nr *nodeRouter) ListNodeHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := nr.nodeService.Search(context.TODO(), tfdirectory.NodeQuery{})
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, nodes)
}

func decodeNode(r *http.Request) (node tfdirectory.Node, err error) {
	if r.Body == nil {
		return node, errors.New("no request body")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return node, decoder.Decode(&node)
}
