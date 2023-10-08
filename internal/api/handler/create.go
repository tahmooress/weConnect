package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tahmooress/weConnect-task/internal/api/dto"
	"github.com/tahmooress/weConnect-task/internal/service"
)

type Handler struct {
	usecase service.UseCase
}

func New(usecase service.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.Statistics
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		record := dto.DtoToEntity(request)

		id, err := h.usecase.Create(r.Context(), &record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		request.ID = id

		response, err := json.Marshal(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("id is empty"))
		}

		err := h.usecase.Delete(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var response dto.Statistics
		response.ID = id

		resp, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("id is empty"))
		}

		record, err := h.usecase.GetByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		response := dto.EntityToDto(*record)

		resp, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p, l = 0, 10

		page := r.URL.Query().Get("page")
		if page != "" {
			i, err := strconv.Atoi(page)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			p = i
		}

		limit := r.URL.Query().Get("limit")
		if limit != "" {
			i, err := strconv.Atoi(page)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			l = i
		}

		records, err := h.usecase.GetAll(r.Context(), int64(p), int64(l))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		response := make([]dto.Statistics, len(records))

		for i := range records {
			response[i] = dto.EntityToDto(records[i])
		}

		resp, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
