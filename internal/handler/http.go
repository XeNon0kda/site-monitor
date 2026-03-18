package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	
	"site-monitor/internal/service"
)

type Handler struct {
	monitor *service.Monitor
}

func NewHandler(monitor *service.Monitor) *Handler {
	return &Handler{monitor: monitor}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/sites", h.addSite).Methods("POST")
	r.HandleFunc("/api/sites", h.getSites).Methods("GET")
	r.HandleFunc("/api/sites/{id}", h.removeSite).Methods("DELETE")
	r.HandleFunc("/ws", h.wsHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
}

func (h *Handler) addSite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}
	site, err := h.monitor.AddSite(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

func (h *Handler) getSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.monitor.GetSites(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func (h *Handler) removeSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.monitor.RemoveSite(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}