package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/bkiran6398/library/internal/books/service/mocks"
	intErr "github.com/bkiran6398/library/internal/errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	createRequest := domain.CreateBookRequest{
		Title:       "New Book",
		Author:      "New Author",
		ISBN:        "ISBN-NEW",
		CopiesTotal: 5,
	}

	expectedBook := domain.Book{
		ID:              uuid.New(),
		Title:           "New Book",
		Author:          "New Author",
		ISBN:            "ISBN-NEW",
		CopiesTotal:     5,
		CopiesAvailable: 5,
	}

	mockService.EXPECT().
		Create(gomock.Any(), createRequest).
		Return(expectedBook, nil).
		Times(1)

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Equal(t, expectedBook.ID, result.ID)
	require.Equal(t, expectedBook.Title, result.Title)
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	createRequest := domain.CreateBookRequest{
		Title:       "", // Empty title
		Author:      "Author",
		ISBN:        "ISBN-123",
		CopiesTotal: 5,
	}

	mockService.EXPECT().
		Create(gomock.Any(), createRequest).
		Return(domain.Book{}, intErr.ErrBadRequest).
		Times(1)

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Create_ConflictError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	createRequest := domain.CreateBookRequest{
		Title:       "Title",
		Author:      "Author",
		ISBN:        "ISBN-123",
		CopiesTotal: 5,
	}

	mockService.EXPECT().
		Create(gomock.Any(), createRequest).
		Return(domain.Book{}, intErr.ErrConflict).
		Times(1)

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	require.Equal(t, http.StatusConflict, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "conflict", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()
	expectedBook := domain.Book{
		ID:              bookID,
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	mockService.EXPECT().
		Get(gomock.Any(), bookID).
		Return(expectedBook, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/v1/books/"+bookID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Get(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Equal(t, expectedBook.ID, result.ID)
	require.Equal(t, expectedBook.Title, result.Title)
}

func TestHandler_Get_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/v1/books/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	w := httptest.NewRecorder()

	handler.Get(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()

	mockService.EXPECT().
		Get(gomock.Any(), bookID).
		Return(domain.Book{}, intErr.ErrNotFound).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/v1/books/"+bookID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Get(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "not_found", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "Updated Title",
		Author:          "Updated Author",
		ISBN:            "ISBN-UPDATED",
		CopiesTotal:     10,
		CopiesAvailable: 8,
	}

	updatedBook := domain.Book{
		ID:              bookID,
		Title:           "Updated Title",
		Author:          "Updated Author",
		ISBN:            "ISBN-UPDATED",
		CopiesTotal:     10,
		CopiesAvailable: 8,
	}

	mockService.EXPECT().
		Update(gomock.Any(), bookID, updateRequest).
		Return(updatedBook, nil).
		Times(1)

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/v1/books/"+bookID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Equal(t, updatedBook.Title, result.Title)
	require.Equal(t, updatedBook.Author, result.Author)
}

func TestHandler_Update_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	updateRequest := domain.UpdateBookRequest{
		Title:           "Title",
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/v1/books/invalid-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Update_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/v1/books/"+bookID.String(), bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "Title",
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	mockService.EXPECT().
		Update(gomock.Any(), bookID, updateRequest).
		Return(domain.Book{}, intErr.ErrNotFound).
		Times(1)

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/v1/books/"+bookID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "not_found", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Update_ConflictError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "Title",
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	mockService.EXPECT().
		Update(gomock.Any(), bookID, updateRequest).
		Return(domain.Book{}, intErr.ErrConflict).
		Times(1)

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/v1/books/"+bookID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	require.Equal(t, http.StatusConflict, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "conflict", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()

	mockService.EXPECT().
		Delete(gomock.Any(), bookID).
		Return(nil).
		Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/v1/books/"+bookID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.Bytes())
}

func TestHandler_Delete_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/v1/books/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "bad_request", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	bookID := uuid.New()

	mockService.EXPECT().
		Delete(gomock.Any(), bookID).
		Return(intErr.ErrNotFound).
		Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/v1/books/"+bookID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": bookID.String()})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "not_found", errorResponse["error"].(map[string]interface{})["code"])
}

func TestHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	expectedBooks := []domain.Book{
		{
			ID:              uuid.New(),
			Title:           "Book 1",
			Author:          "Author 1",
			ISBN:            "ISBN-1",
			CopiesTotal:     5,
			CopiesAvailable: 3,
		},
		{
			ID:              uuid.New(),
			Title:           "Book 2",
			Author:          "Author 2",
			ISBN:            "ISBN-2",
			CopiesTotal:     10,
			CopiesAvailable: 8,
		},
	}

	mockService.EXPECT().
		List(gomock.Any(), domain.ListFilter{}).
		Return(expectedBooks, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result []domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expectedBooks[0].Title, result[0].Title)
}

func TestHandler_List_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	titleFilter := "Gatsby"
	authorFilter := "Fitzgerald"
	limit := 10
	offset := 5

	expectedBooks := []domain.Book{
		{
			ID:              uuid.New(),
			Title:           "The Great Gatsby",
			Author:          "F. Scott Fitzgerald",
			ISBN:            "ISBN-1",
			CopiesTotal:     5,
			CopiesAvailable: 3,
		},
	}

	expectedFilter := domain.ListFilter{
		Title:  &titleFilter,
		Author: &authorFilter,
		Limit:  limit,
		Offset: offset,
	}

	mockService.EXPECT().
		List(gomock.Any(), expectedFilter).
		Return(expectedBooks, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/v1/books?title=Gatsby&author=Fitzgerald&limit=10&offset=5", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result []domain.Book
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestHandler_List_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := NewHandler(mockService)

	mockService.EXPECT().
		List(gomock.Any(), domain.ListFilter{}).
		Return(nil, errors.New("database error")).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	require.Equal(t, "internal_error", errorResponse["error"].(map[string]interface{})["code"])
}

func TestParseListQueryParameters(t *testing.T) {
	tests := []struct {
		name           string
		queryString    string
		expectedFilter domain.ListFilter
	}{
		{
			name:        "empty query",
			queryString: "",
			expectedFilter: domain.ListFilter{
				Title:  nil,
				Author: nil,
				ISBN:   nil,
				Limit:  0,
				Offset: 0,
			},
		},
		{
			name:        "with title only",
			queryString: "title=Gatsby",
			expectedFilter: domain.ListFilter{
				Title:  stringPtr("Gatsby"),
				Author: nil,
				ISBN:   nil,
				Limit:  0,
				Offset: 0,
			},
		},
		{
			name:        "with all filters",
			queryString: "title=Gatsby&author=Fitzgerald&isbn=ISBN-123&limit=10&offset=5",
			expectedFilter: domain.ListFilter{
				Title:  stringPtr("Gatsby"),
				Author: stringPtr("Fitzgerald"),
				ISBN:   stringPtr("ISBN-123"),
				Limit:  10,
				Offset: 5,
			},
		},
		{
			name:        "with invalid limit and offset",
			queryString: "limit=abc&offset=xyz",
			expectedFilter: domain.ListFilter{
				Title:  nil,
				Author: nil,
				ISBN:   nil,
				Limit:  0,
				Offset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/books?"+tt.queryString, nil)
			filter := parseListQueryParameters(req)

			require.Equal(t, tt.expectedFilter.Limit, filter.Limit)
			require.Equal(t, tt.expectedFilter.Offset, filter.Offset)

			if tt.expectedFilter.Title == nil {
				require.Nil(t, filter.Title)
			} else {
				require.NotNil(t, filter.Title)
				require.Equal(t, *tt.expectedFilter.Title, *filter.Title)
			}

			if tt.expectedFilter.Author == nil {
				require.Nil(t, filter.Author)
			} else {
				require.NotNil(t, filter.Author)
				require.Equal(t, *tt.expectedFilter.Author, *filter.Author)
			}

			if tt.expectedFilter.ISBN == nil {
				require.Nil(t, filter.ISBN)
			} else {
				require.NotNil(t, filter.ISBN)
				require.Equal(t, *tt.expectedFilter.ISBN, *filter.ISBN)
			}
		})
	}
}

func TestParseBookIDFromPath(t *testing.T) {
	tests := []struct {
		name        string
		bookID      string
		expectError bool
	}{
		{
			name:        "valid UUID",
			bookID:      uuid.New().String(),
			expectError: false,
		},
		{
			name:        "invalid UUID",
			bookID:      "invalid-uuid",
			expectError: true,
		},
		{
			name:        "empty string",
			bookID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/books/"+tt.bookID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			bookID, err := parseBookIDFromPath(req)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.bookID, bookID.String())
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

