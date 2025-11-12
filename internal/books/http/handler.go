package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/bkiran6398/library/internal/books/service"
	"github.com/bkiran6398/library/internal/http/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for book operations.
type Handler struct {
	service service.Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service service.Service) Handler {
	return Handler{service: service}
}

func (handler Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := parseListQueryParameters(r)
	books, err := handler.service.List(r.Context(), filter)
	if err != nil {
		response.MapServiceErrorToHTTP(w, err)
		return
	}
	response.JSON(w, http.StatusOK, books)
}

// parseListQueryParameters extracts and parses query parameters for listing books.
func parseListQueryParameters(r *http.Request) domain.ListFilter {
	queryParams := r.URL.Query()

	title := queryParams.Get("title")
	author := queryParams.Get("author")
	isbn := queryParams.Get("isbn")
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	offset, _ := strconv.Atoi(queryParams.Get("offset"))

	var titlePtr, authorPtr, isbnPtr *string
	if title != "" {
		titlePtr = &title
	}
	if author != "" {
		authorPtr = &author
	}
	if isbn != "" {
		isbnPtr = &isbn
	}

	return domain.ListFilter{
		Title:  titlePtr,
		Author: authorPtr,
		ISBN:   isbnPtr,
		Limit:  limit,
		Offset: offset,
	}
}

func (handler Handler) Create(w http.ResponseWriter, r *http.Request) {
	var createRequest domain.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "Invalid JSON body", nil)
		return
	}

	book, err := handler.service.Create(r.Context(), createRequest)
	if err != nil {
		response.MapServiceErrorToHTTP(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, book)
}

func (handler Handler) Get(w http.ResponseWriter, r *http.Request) {
	bookID, err := parseBookIDFromPath(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "Invalid book ID", nil)
		return
	}

	book, err := handler.service.Get(r.Context(), bookID)
	if err != nil {
		response.MapServiceErrorToHTTP(w, err)
		return
	}
	response.JSON(w, http.StatusOK, book)
}

// parseBookIDFromPath extracts and parses the book ID from the request path.
func parseBookIDFromPath(r *http.Request) (uuid.UUID, error) {
	bookIDString := mux.Vars(r)["id"]
	return uuid.Parse(bookIDString)
}

func (handler Handler) Update(w http.ResponseWriter, r *http.Request) {
	bookID, err := parseBookIDFromPath(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "Invalid book ID", nil)
		return
	}

	var updateRequest domain.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "Invalid JSON body", nil)
		return
	}

	book, err := handler.service.Update(r.Context(), bookID, updateRequest)
	if err != nil {
		response.MapServiceErrorToHTTP(w, err)
		return
	}
	response.JSON(w, http.StatusOK, book)
}

func (handler Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bookID, err := parseBookIDFromPath(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "Invalid book ID", nil)
		return
	}

	if err := handler.service.Delete(r.Context(), bookID); err != nil {
		response.MapServiceErrorToHTTP(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
