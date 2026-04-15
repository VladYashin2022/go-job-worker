package handler

import (
	"encoding/json"
	"go-job-worker/internal/service"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	jobsService *service.JobsService
}

func NewHandler(jbsrv *service.JobsService) *Handler {
	return &Handler{
		jobsService: jbsrv,
	}
}

func (h *Handler) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/jobs", h.handler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetJob(w, r)
	case http.MethodPost:
		h.handleCreateJob(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GET
func (h *Handler) handleGetJob(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		//Get All Jobs
		jobs := h.jobsService.GetAllJobs()

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(jobs)
		if err != nil {
			log.Println("encoder error:", err)
			return
		}

	} else {
		//Get Job by ID
		idJob, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		job, err := h.jobsService.GetJob(idJob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(job)
		if err != nil {
			log.Println("encoder error:", err)
			return
		}

	}
}

// POST
func (h *Handler) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req requestJob
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Decode error", http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.Payload == "" {
		http.Error(w, "Empty parameter in request", http.StatusBadRequest)
		return
	}

	job, err := h.jobsService.CreateJob(req.Type, req.Payload)
	if err != nil {
		http.Error(w, "Create task error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		log.Println("Encoder error")
		return
	}
}

type requestJob struct {
	Type    string
	Payload string
}
