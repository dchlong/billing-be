package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/dchlong/billing-be/internal/presenter"
	"github.com/dchlong/billing-be/internal/services"
	"github.com/dchlong/billing-be/pkg/logger"
)

type HandlerTestSuite struct {
	suite.Suite
	httpHandler http.Handler
	billService *services.MockBillService
}

func (s *HandlerTestSuite) SetupSuite() {
	router := gin.Default()
	mockCtrl := gomock.NewController(s.T())
	s.billService = services.NewMockBillService(mockCtrl)
	ilogger, _, err := logger.ProvideLogger()
	s.Require().NoError(err)
	h := &handler{
		billService: s.billService,
		ilogger:     ilogger,
	}
	h.Register(router)
	s.httpHandler = router
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestCreateHistoryCall() {
	const (
		testUserName     = "hlong"
		testCallDuration = int64(1000)
	)
	expectedResp := &presenter.CreateCallHistoryResponse{
		ID: 1,
	}

	s.billService.EXPECT().CreateCallHistory(gomock.Any(), testUserName, testCallDuration).Return(expectedResp, nil)
	path := fmt.Sprintf("/mobile/%s/call", testUserName)
	body := fmt.Sprintf(`{"call_duration":%d}`, testCallDuration)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusOK)
	var resp presenter.CreateCallHistoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().Equal(expectedResp.ID, resp.ID)
}

func (s *HandlerTestSuite) TestCreateHistoryCall_Error() {
	const (
		testUserName     = "hlong"
		testCallDuration = int64(1000)
	)

	s.billService.EXPECT().CreateCallHistory(gomock.Any(), testUserName, testCallDuration).Return(nil, errors.New(""))
	path := fmt.Sprintf("/mobile/%s/call", testUserName)
	body := fmt.Sprintf(`{"call_duration":%d}`, testCallDuration)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusInternalServerError)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("internal_server_error", resp.Error.Code)
}

func (s *HandlerTestSuite) TestCreateHistoryCall_InvalidUserName() {
	const (
		testUserName     = "hlong111112222233333444445555566666777778888899999"
		testCallDuration = int64(1000)
	)

	path := fmt.Sprintf("/mobile/%s/call", testUserName)
	body := fmt.Sprintf(`{"call_duration":%d}`, testCallDuration)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_user_name", resp.Error.Code)
}

func (s *HandlerTestSuite) TestCreateHistoryCall_EmptyUserName() {
	const (
		testCallDuration = int64(1000)
	)

	path := `/mobile//call`
	body := fmt.Sprintf(`{"call_duration":%d}`, testCallDuration)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_user_name", resp.Error.Code)
}

func (s *HandlerTestSuite) TestCreateHistoryCall_InvalidCallDuration() {
	const (
		testUserName     = "hlong"
		testCallDuration = int64(-1000)
	)

	path := fmt.Sprintf("/mobile/%s/call", testUserName)
	body := fmt.Sprintf(`{"call_duration":%d}`, testCallDuration)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_input", resp.Error.Code)
}

func (s *HandlerTestSuite) TestCreateHistoryCall_EmptyCallDuration() {
	const (
		testUserName = "hlong"
	)

	path := fmt.Sprintf("/mobile/%s/call", testUserName)
	body := `{}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_input", resp.Error.Code)
}

func (s *HandlerTestSuite) TestGetBill() {
	const (
		testUserName = "hlong"
	)

	expectedResp := &presenter.GetBillResponse{
		CallCount:  3,
		BlockCount: 5,
	}

	s.billService.EXPECT().GetBill(gomock.Any(), testUserName).Return(expectedResp, nil)
	path := fmt.Sprintf("/mobile/%s/billing", testUserName)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusOK)
	var resp presenter.GetBillResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().Equal(expectedResp.CallCount, resp.CallCount)
	s.Require().Equal(expectedResp.BlockCount, resp.BlockCount)
}

func (s *HandlerTestSuite) TestGetBill_Error() {
	const (
		testUserName = "hlong"
	)

	s.billService.EXPECT().GetBill(gomock.Any(), testUserName).Return(nil, errors.New(""))
	path := fmt.Sprintf("/mobile/%s/billing", testUserName)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusInternalServerError)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("internal_server_error", resp.Error.Code)
}

func (s *HandlerTestSuite) TestGetBill_InvalidUserName() {
	const (
		testUserName = "hlong111112222233333444445555566666777778888899999"
	)

	path := fmt.Sprintf("/mobile/%s/billing", testUserName)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_user_name", resp.Error.Code)
}

func (s *HandlerTestSuite) TestGetBill_EmptyUserName() {
	path := "/mobile//billing"
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	w := httptest.NewRecorder()
	s.httpHandler.ServeHTTP(w, req)
	s.Require().Equal(w.Code, http.StatusBadRequest)
	var resp presenter.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)
	s.Require().NotNil(resp.Error)
	s.Require().Equal("invalid_user_name", resp.Error.Code)
}
