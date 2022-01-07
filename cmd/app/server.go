package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Tursunkhuja/crud/pkg/customers"
	"github.com/gorilla/mux"
)

//Server ..............
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
}

//NewServer: Create new Server
func NewServer(mux *mux.Router, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customersSvc}
}

//function for launching handlers through mux
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomer).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleSaveCustomer).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveCustomerById).Methods(DELETE)
	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	//s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	//s.mux.HandleFunc("/customers.save", s.handleSaveCustomer)
	//s.mux.HandleFunc("/customers.removeById", s.handleRemoveCustomerById)

	s.mux.HandleFunc("/customers/active", s.handleGetCustomer).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockCustomerById).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockCustomerById).Methods(DELETE)

	//s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	//s.mux.HandleFunc("/customers.blockById", s.handleBlockCustomerById)
	//s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockCustomerById)
}
func (s *Server) handleGetCustomer(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if idParam == "active" {
		s.handleGetAllActiveCustomers(writer, request)
	}
	if idParam != "active" {
		s.handleGetCustomerByID(writer, request)
	}
}

func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println("parse error 36. idParam=" + idParam)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.customerSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.GetAll(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetAllActiveCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.GetAllActive(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	data, err := json.Marshal(items)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err = s.customerSvc.Save(request.Context(), item)

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {

		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {

		return
	}

}
func (s *Server) handleRemoveCustomerById(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(idParam, 10, 60)
	if err != nil {
		log.Println("parse error 141")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if id == 0 {

		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	id2, err := s.customerSvc.RemoveByID(request.Context(), id)
	if err != nil {

		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(id2)
	if err != nil {

		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}

}

func (s *Server) handleBlockCustomerById(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(idParam, 10, 60)
	if err != nil {
		log.Println("parse error 176")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	item, err := s.customerSvc.BlockByID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {

		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {

		return
	}
}

func (s *Server) handleUnBlockCustomerById(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(idParam, 10, 60)
	if err != nil {
		log.Println("parse error 209")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	item, err := s.customerSvc.UnblockByID(request.Context(), id)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
