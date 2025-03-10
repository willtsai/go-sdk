package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dapr/go-sdk/actor/api"
	"github.com/dapr/go-sdk/actor/mock"

	"github.com/stretchr/testify/assert"

	"github.com/dapr/go-sdk/service/common"
)

func testTopicFunc(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	if e == nil {
		return false, errors.New("nil content")
	}
	if e.DataContentType != "application/json" {
		return false, fmt.Errorf("invalid content type: %s", e.DataContentType)
	}
	return false, nil
}

func testErrorTopicFunc(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	return true, errors.New("error to cause a retry")
}

func TestEventNilHandler(t *testing.T) {
	s := newServer("", nil)
	sub := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
		Route:      "/",
		Metadata:   map[string]string{},
	}
	err := s.AddTopicEventHandler(sub, nil)
	assert.Errorf(t, err, "expected error adding event handler")
}

func TestEventHandler(t *testing.T) {
	data := `{
		"specversion" : "1.0",
		"type" : "com.github.pull.create",
		"source" : "https://github.com/cloudevents/spec/pull",
		"subject" : "123",
		"id" : "A234-1234-1234",
		"time" : "2018-04-05T17:31:00Z",
		"comexampleextension1" : "value",
		"comexampleothervalue" : 5,
		"datacontenttype" : "application/json",
		"data" : "eyJtZXNzYWdlIjoiaGVsbG8ifQ=="
	}`

	s := newServer("", nil)

	sub := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
		Route:      "/",
		Metadata:   map[string]string{},
	}
	err := s.AddTopicEventHandler(sub, testTopicFunc)
	assert.NoErrorf(t, err, "error adding event handler")

	sub2 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "errors",
		Route:      "/errors",
		Metadata:   map[string]string{},
	}
	err = s.AddTopicEventHandler(sub2, testErrorTopicFunc)
	assert.NoErrorf(t, err, "error adding error event handler")

	s.registerBaseHandler()

	makeEventRequest(t, s, "/", data, http.StatusOK)
	makeEventRequest(t, s, "/", "", http.StatusSeeOther)
	makeEventRequest(t, s, "/", "not JSON", http.StatusSeeOther)
	makeEventRequest(t, s, "/errors", data, http.StatusOK)
}

func TestHealthCheck(t *testing.T) {
	s := newServer("", nil)
	s.registerBaseHandler()
	makeRequest(t, s, "/healthz", "", http.MethodGet, http.StatusOK)
}

func TestActorConfig(t *testing.T) {
	s := newServer("", nil)
	s.registerBaseHandler()
	makeRequest(t, s, "/dapr/config", "", http.MethodGet, http.StatusOK)
}

func TestActorHandler(t *testing.T) {
	reminderReqData, _ := json.Marshal(api.ActorReminderParams{
		Data:    []byte("hello"),
		DueTime: "5s",
		Period:  "5s",
	})

	timerReqData, _ := json.Marshal(api.ActorTimerParam{
		CallBack: "Invoke",
		DueTime:  "5s",
		Period:   "5s",
		Data:     []byte(`"hello"`),
	})

	timerReqDataWithBadCallBackFunction, _ := json.Marshal(api.ActorTimerParam{
		CallBack: "UnexistedFunc",
		DueTime:  "5s",
		Period:   "5s",
		Data:     []byte(`"hello"`),
	})
	s := newServer("", nil)
	s.registerBaseHandler()
	// invoke actor API without target actor defined
	makeRequest(t, s, "/actors/testActorType/testActorID/method/Invoke", "", http.MethodPut, http.StatusNotFound)
	makeRequest(t, s, "/actors/testActorType/testActorID", "", http.MethodDelete, http.StatusNotFound)
	makeRequest(t, s, "/actors/testActorType/testActorID/method/remind/testReminderName", string(reminderReqData), http.MethodPut, http.StatusNotFound)
	makeRequest(t, s, "/actors/testActorType/testActorID/method/timer/testTimerName", string(timerReqData), http.MethodPut, http.StatusNotFound)

	// register test actor factory
	s.RegisterActorImplFactory(mock.ActorImplFactory)

	// invoke actor API with internal error
	makeRequest(t, s, "/actors/testActorType/testActorID/method/remind/testReminderName", `{
"dueTime": "5s",
"period": "5s",
"data": "test data"`, http.MethodPut, http.StatusInternalServerError)
	makeRequest(t, s, "/actors/testActorType/testActorID/method/Invoke", "bad request param", http.MethodPut, http.StatusInternalServerError)
	makeRequest(t, s, "/actors/testActorType/testActorID/method/timer/testTimerName", string(timerReqDataWithBadCallBackFunction), http.MethodPut, http.StatusInternalServerError)

	// invoke actor API with success status
	makeRequestWithExpectedBody(t, s, "/actors/testActorType/testActorID/method/Invoke", `"invoke request"`, http.MethodPut, http.StatusOK, []byte(`"invoke request"`))
	makeRequest(t, s, "/actors/testActorType/testActorID/method/remind/testReminderName", string(reminderReqData), http.MethodPut, http.StatusOK)
	makeRequest(t, s, "/actors/testActorType/testActorID/method/timer/testTimerName", string(timerReqData), http.MethodPut, http.StatusOK)
	makeRequest(t, s, "/actors/testActorType/testActorID", "", http.MethodDelete, http.StatusOK)

	// register not reminder callee actor factory
	s.RegisterActorImplFactory(mock.NotReminderCalleeActorFactory)
	// invoke call reminder to not reminder callee actor type
	makeRequest(t, s, "/actors/testActorNotReminderCalleeType/testActorID/method/remind/testReminderName", string(reminderReqData), http.MethodPut, http.StatusInternalServerError)
}

func makeRequest(t *testing.T, s *Server, route, data, method string, expectedStatusCode int) {
	req, err := http.NewRequest(method, route, strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request: %s", data)
	testRequest(t, s, req, expectedStatusCode)
}

func makeRequestWithExpectedBody(t *testing.T, s *Server, route, data, method string, expectedStatusCode int, expectedBody []byte) {
	req, err := http.NewRequest(method, route, strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request: %s", data)
	testRequestWithResponseBody(t, s, req, expectedStatusCode, expectedBody)
}

func makeEventRequest(t *testing.T, s *Server, route, data string, expectedStatusCode int) {
	req, err := http.NewRequest(http.MethodPost, route, strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request: %s", data)
	req.Header.Set("Content-Type", "application/json")
	testRequest(t, s, req, expectedStatusCode)
}

func TestAddingInvalidEventHandlers(t *testing.T) {
	s := newServer("", nil)
	err := s.AddTopicEventHandler(nil, testTopicFunc)
	assert.Errorf(t, err, "expected error adding no sub event handler")

	sub := &common.Subscription{Metadata: map[string]string{}}
	err = s.AddTopicEventHandler(sub, testTopicFunc)
	assert.Errorf(t, err, "expected error adding empty sub event handler")

	sub.Topic = "test"
	err = s.AddTopicEventHandler(sub, testTopicFunc)
	assert.Errorf(t, err, "expected error adding sub without component event handler")

	sub.PubsubName = "messages"
	err = s.AddTopicEventHandler(sub, testTopicFunc)
	assert.Errorf(t, err, "expected error adding sub without route event handler")
}

func TestRawPayloadDecode(t *testing.T) {
	testRawTopicFunc := func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		if e.DataContentType != "application/octet-stream" {
			err = fmt.Errorf("invalid content type: %s", e.DataContentType)
		}
		if e.DataBase64 != "eyJtZXNzYWdlIjoiaGVsbG8ifQ==" {
			err = errors.New("error decode data_base64")
		}
		if err != nil {
			assert.NoErrorf(t, err, "error rawPayload decode")
		}
		return
	}

	const rawData = `{
		"datacontenttype" : "application/octet-stream",
		"data_base64" : "eyJtZXNzYWdlIjoiaGVsbG8ifQ=="
	}`

	s := newServer("", nil)

	sub3 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "testRaw",
		Route:      "/raw",
		Metadata: map[string]string{
			"rawPayload": "true",
		},
	}
	err := s.AddTopicEventHandler(sub3, testRawTopicFunc)
	assert.NoErrorf(t, err, "error adding raw event handler")

	s.registerBaseHandler()
	makeEventRequest(t, s, "/raw", rawData, http.StatusOK)
}
