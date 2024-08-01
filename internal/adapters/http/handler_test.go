package http

import (
	"bytes"
	"errors"
	"github.com/segmentio/kafka-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/crseat/example-data-pipeline/internal/app"
	"github.com/crseat/example-data-pipeline/internal/domain"
)

// MockProducer is a mock implementation of the ProducerService interface
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) ProcessPostData(postData domain.PostData) error {
	args := m.Called(postData)
	return args.Error(0)
}

func (m *MockProducer) WritePostDataToKafka(postData domain.PostData) error {
	args := m.Called(postData)
	return args.Error(0)
}

func (m *MockProducer) WriteMessageToKafka(message kafka.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestHandler_handleSubmit(t *testing.T) {
	tests := []struct {
		name                string
		requestBody         string
		mockServiceErr      error
		expectedStatus      int
		expectedBody        string
		expectedProcessCall bool
	}{
		{
			name:                "Invalid JSON",
			requestBody:         `{invalid json}`,
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        `{"error":"Invalid JSON"}`,
			expectedProcessCall: false,
		},
		{
			name: "Validation Error",
			requestBody: `{
                    "ip_address": "invalid-ip",
                    "user_agent": "Mozilla/5.0",
                    "referring_url": "http://example.com",
                    "advertiser_id": "123456",
                    "metadata": {
                        "campaign": "summer_sale",
                        "clicks": 120
                    }
                }`,
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        `{"error":"Key: 'PostData.IPAddress' Error:Field validation for 'IPAddress' failed on the 'ipv4' tag"}`,
			expectedProcessCall: false,
		},
		{
			name: "Service Error",
			requestBody: `{
                    "ip_address": "192.168.1.1",
                    "user_agent": "Mozilla/5.0",
                    "referring_url": "http://example.com",
                    "advertiser_id": "123456",
                    "metadata": {
                        "campaign": "summer_sale",
                        "clicks": 120
                    }
                }`,
			mockServiceErr:      errors.New("service error"),
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        `{"error":"Internal Server Error"}`,
			expectedProcessCall: true,
		},
		{
			name: "Successful Request",
			requestBody: `{
                    "ip_address": "192.168.1.1",
                    "user_agent": "Mozilla/5.0",
                    "referring_url": "http://example.com",
                    "advertiser_id": "123456",
                    "metadata": {
                        "campaign": "summer_sale",
                        "clicks": 120
                    }
                }`,
			expectedStatus:      http.StatusNoContent,
			expectedBody:        ``,
			expectedProcessCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			// Register Validator
			e.Validator = &CustomValidator{validator: validator.New()}

			// Mock service
			mockProducer := new(MockProducer)

			// If we're expecting to call the process method, check whether the error (or lack thereof) that we expect is present
			if tt.expectedProcessCall == true {
				if tt.mockServiceErr != nil {
					//mockProducer.On("ProcessPostData", mock.AnythingOfType("domain.PostData")).Return(tt.mockServiceErr)
					mockProducer.On("WritePostDataToKafka", mock.AnythingOfType("domain.PostData")).Return(tt.mockServiceErr)
				} else {
					//mockProducer.On("ProcessPostData", mock.AnythingOfType("domain.PostData")).Return(nil)
					mockProducer.On("WritePostDataToKafka", mock.AnythingOfType("domain.PostData")).Return(nil)
				}
			}

			handler := NewHandler(&app.ProducerService{})
			handler.service.SetProducer(mockProducer)

			req := httptest.NewRequest(http.MethodPost, "/submit", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.handleSubmit(c)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.expectedBody != `` {
					assert.JSONEq(t, tt.expectedBody, rec.Body.String())
				}
			}

			mockProducer.AssertExpectations(t)
		})
	}
}
