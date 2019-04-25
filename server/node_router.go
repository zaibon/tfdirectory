package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type nodeRouter struct {
	nodeService tfdirectory.NodeService
}

func NewUserRouter(u tfdirectory.NodeService, router *mux.Router) *mux.Router {
	nodeRouter := &nodeRouter{u}

	router.HandleFunc("", nodeRouter.createNodeHandler).Methods("PUT")
	router.HandleFunc("", nodeRouter.ListNodeHandler).Methods("GET")
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

func queryStr(r *http.Request, target, ifNil string) string {
	q, ok := r.URL.Query()[target]
	if ok {
		return q[0]
	}
	return ifNil
}

func queryFloat64(r *http.Request, target string, ifNil float64) (float64, error) {
	q, ok := r.URL.Query()[target]
	if ok {
		f, err := strconv.ParseFloat(q[0], 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}
	return ifNil, nil
}

func (nr *nodeRouter) ListNodeHandler(w http.ResponseWriter, r *http.Request) {
	cru, err := queryFloat64(r, "cru", 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("cru format error:%v", err), http.StatusBadRequest)
	}
	mru, err := queryFloat64(r, "mru", 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("mru format error:%v", err), http.StatusBadRequest)
	}
	hru, err := queryFloat64(r, "hru", 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("hru format error:%v", err), http.StatusBadRequest)
	}
	sru, err := queryFloat64(r, "sru", 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("sru format error:%v", err), http.StatusBadRequest)
	}

	q := tfdirectory.NodeQuery{
		Farmer: queryStr(r, "farmer", ""),
		Location: tfdirectory.Location{
			Country: queryStr(r, "country", ""),
		},
		Resource: tfdirectory.Resource{
			CRU: cru,
			MRU: mru,
			SRU: sru,
			HRU: hru,
		},
	}
	log.Println(q)
	nodes, err := nr.nodeService.Search(context.TODO(), q)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	log.Println(len(nodes))

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
