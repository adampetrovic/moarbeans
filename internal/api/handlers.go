package api

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/adampetrovic/moarbeans/internal/woodroaster"
)

type Server struct {
	db     *database.DB
	client *woodroaster.Client
}

func NewServer(db *database.DB) *Server {
	return &Server{
		db:     db,
		client: woodroaster.NewClient(),
	}
}

func (s *Server) handleRequestLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.client.RequestMagicLink(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetNextOrderDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := s.db.GetLatestSession()
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if session == nil {
		http.Error(w, "No valid session found", http.StatusUnauthorized)
		return
	}

	nextDate, err := s.client.GetNextOrderDate(session.SessionToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"next_order_date": nextDate.Format("2006-01-02"),
	})
}

func (s *Server) handleSetNextOrderDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Date string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	session, err := s.db.GetLatestSession()
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if session == nil {
		http.Error(w, "No valid session found", http.StatusUnauthorized)
		return
	}

	if err := s.client.SetNextOrderDate(session.SessionToken, date); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 