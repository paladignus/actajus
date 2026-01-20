// Package handler
package handler

// type MockEnterpriseService struct {
// 	expectedError error
// }
//
// func (m *MockEnterpriseService) Execute(ctx context.Context, input dto.EnterpriseInput) error {
// 	return m.expectedError
// }
//
// func TestNewEnterpriseSevice(t *testing.T) {
// 	service := &MockEnterpriseService{
// 		expectedError: nil,
// 	}
// 	logger := &spy.Logger{}
//
// 	t.Run("should sucessfully request", func(t *testing.T) {
// 		requestBody := `{"registered_by": 1, "name": "Actajus", "trade_name": "Actajus sistema jurídico", "cnpj": "10.123.456/0001-00"}`
// 		sut := NewEnterprise(service, logger)
// 		req := httptest.NewRequest(http.MethodPost, "/enterprise", bytes.NewBufferString(requestBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		logger.On("Info", req.Context(), "create enterprise successful", "cnpj", "10.123.456/0001-00", "registered_by", 1, "method", req.Method, "url", req.URL.Path)
// 		sut.Create(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
// 	})
//
// 	t.Run("should an error for invalid request body", func(t *testing.T) {
// 		mockService := &MockEnterpriseService{}
// 		sut := NewEnterprise(mockService, logger)
// 		req := httptest.NewRequest(http.MethodPost, "/enterprise", bytes.NewBufferString("{invalid json"))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		logger.On("Warn", req.Context(), "failed to decode request body for enterprise create", "error", mock.Anything)
// 		sut.Create(w, req)
// 		if w.Code != http.StatusBadRequest {
// 			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
// 		}
// 	})
//
// 	t.Run("should return error if registered_by is invalid", func(t *testing.T) {
// 		mockService := &MockEnterpriseService{
// 			expectedError: exception.ErrInvalidRegisteredBy,
// 		}
// 		logger := &spy.Logger{}
// 		sut := NewEnterprise(mockService, logger)
// 		requestBody := `{"registered_by": 1, "name": "Actajus", "trade_name": "Actajus sistema jurídico", "cnpj": "10.123.456/0001-00"}`
// 		req := httptest.NewRequest(http.MethodPost, "/enterprise", bytes.NewBufferString(requestBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		logger.On(
// 			"Warn",
// 			req.Context(),
// 			"create enterprise failed",
// 			"error",
// 			mockService.expectedError,
// 			"cnpj",
// 			"10.123.456/0001-00",
// 			"method",
// 			req.Method,
// 			"url",
// 			req.URL.Path,
// 		)
// 		sut.Create(w, req)
// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 	})
//
// 	t.Run("should return error if registered_by is invalid", func(t *testing.T) {
// 		mockService := &MockEnterpriseService{
// 			expectedError: fmt.Errorf("failed to create enterprise"),
// 		}
// 		logger := &spy.Logger{}
// 		sut := NewEnterprise(mockService, logger)
// 		requestBody := `{"registered_by": 1, "name": "Actajus", "trade_name": "Actajus sistema jurídico", "cnpj": "10.123.456/0001-00"}`
// 		req := httptest.NewRequest(http.MethodPost, "/enterprise", bytes.NewBufferString(requestBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		logger.On(
// 			"Error",
// 			req.Context(),
// 			"create enterprise failed with server error",
// 			"error",
// 			mockService.expectedError,
// 			"cnpj",
// 			"10.123.456/0001-00",
// 			"method",
// 			req.Method,
// 			"url",
// 			req.URL.Path,
// 		)
// 		sut.Create(w, req)
// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	})
// }
