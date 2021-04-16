package mockserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xicoalmeida/mockserver-go/model"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// Verifier
func None() *model.VerificationTimes {
	return &model.VerificationTimes{
		AtLeast: 0,
		AtMost:  0,
	}
}

func Once() *model.VerificationTimes {
	return &model.VerificationTimes{
		AtLeast: 1,
		AtMost:  1,
	}
}

func Twice() *model.VerificationTimes {
	return &model.VerificationTimes{
		AtLeast: 2,
		AtMost:  2,
	}
}

//Fluent API to help create a expectation request
type Request struct {
	method         string
	path           string
	header         map[string]string
	body           interface{}
	queryParams    map[string]interface{}
	pathParameters map[string][]string
}

func ARequest() *Request {
	return &Request{}
}

func (r *Request) WithMethod(m string) *Request {
	r.method = m
	return r
}
func (r *Request) WithPath(p string) *Request {
	r.path = p
	return r
}
func (r *Request) WithPathVariable(p map[string][]string) *Request {
	r.pathParameters = p
	return r
}
func (r *Request) WithHeader(h map[string]string) *Request {
	r.header = h
	return r
}
func (r *Request) WithBody(b interface{}) *Request {
	r.body = b
	return r
}
func (r *Request) WithQueryParams(qp map[string]interface{}) *Request {
	r.queryParams = qp
	return r
}

//Fluent API to help create a expectation response
type Response struct {
	body       interface{}
	statusCode int
}

func AResponse() *Response {
	return &Response{}
}

func (r *Response) WithStatusCode(code int) *Response {
	r.statusCode = code
	return r
}
func (r *Response) WithBody(body interface{}) *Response {
	r.body = body
	return r
}

//Fluent API to help create a expectation
type Expectation struct {
	request  *Request
	response *Response
}

func NewExpectation() *Expectation {
	return &Expectation{}
}

func (ms *Expectation) When(request *Request) *Expectation {
	ms.request = request
	return ms
}
func (ms *Expectation) Respond(response *Response) *Expectation {
	ms.response = response
	return ms
}
func (ms *Expectation) buildExpectation() *model.ExpectationRequest {

	body := &model.Body{
		ContentType: "application/json",
		Type:        "JSON",
		MatchType:   "STRICT",
		Json:        ms.request.body,
	}

	req := model.HttpRequest{
		Method:          ms.request.method,
		Path:            ms.request.path,
		QueryParameters: ms.request.queryParams,
		Body:            body,
	}

	if len(ms.request.pathParameters) != 0 {
		req.PathParameters = ms.request.pathParameters
	}

	respBody, err := json.Marshal(ms.response.body)

	if err != nil {
		return nil
	}

	e := &model.ExpectationRequest{
		HttpRequest: req,
		HttpResponse: model.HttpResponse{
			StatusCode: ms.response.statusCode,
			Body:       string(respBody),
		},
	}
	return e
}

//Fluent API to help create verify expectations
type Verifier struct {
	times       *model.VerificationTimes
	expectation *Expectation
	retry       int
}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Expectation(expectation *Expectation) *Verifier {
	v.expectation = expectation
	return v
}
func (v *Verifier) ToBeCalled(times *model.VerificationTimes) *Verifier {
	v.times = times
	return v
}
func (v *Verifier) Retry(retry int) *Verifier {
	v.retry = retry
	return v
}
func (v *Verifier) buildVerifier() *model.VerifierRequest {
	reqBody, err := json.Marshal(v.expectation.request.body)
	if err != nil {
		return nil
	}

	body := &model.Body{
		ContentType: "application/json",
		Type:        "JSON",
		MatchType:   "STRICT",
		Json:        string(reqBody),
	}

	er := model.HttpRequest{
		Method:          v.expectation.request.method,
		Path:            v.expectation.request.path,
		QueryParameters: v.expectation.request.queryParams,
		Body:            body,
	}

	if len(v.expectation.request.pathParameters) != 0 {
		er.PathParameters = v.expectation.request.pathParameters
	}

	return &model.VerifierRequest{
		HttpRequest: er,
		Times: model.VerificationTimes{
			AtLeast: v.times.AtLeast,
			AtMost:  v.times.AtMost,
		},
	}
}

//Fluent API to manage mockserver behaviors
type MockServerClient struct {
	address string
	client  *http.Client
}

func (ms *MockServerClient) Expect(t *testing.T, expectation *Expectation) *Expectation {
	uri := fmt.Sprintf("http://%s/mockserver/expectation", ms.address)

	exReq := expectation.buildExpectation()

	if exReq == nil {
		t.Fatal("Unable to create expectation")
	}

	resp, err := ms.execute(exReq, uri)

	if err != nil {
		t.Fatal(err)
	}

	switch statusCode := resp.StatusCode; statusCode {

	case 400:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		t.Fatalf("Incorrent request format: %v", string(body))
	case 406:
		t.Fatalf("MockServer unknown error. Details: %v", resp)
	case 201:
		t.Log("Expectation has been created")
		return expectation
	default:
		t.Fatalf("Unknown error. Check response: %v ", resp)
	}
	return nil
}

func (ms *MockServerClient) Verify(t *testing.T, v *Verifier) {
	for i := 1; i <= v.retry; i++ {

		t.Logf("Attempting %v of %v", i, v.retry)

		uri := fmt.Sprintf("http://%s/mockserver/verify", ms.address)

		verifier := v.buildVerifier()

		if verifier == nil {
			t.Fatal("Unable to create verifier")
		}

		resp, err := ms.execute(verifier, uri)

		if err != nil {
			t.Fatal(err)
		}

		switch statusCode := resp.StatusCode; statusCode {

		case 403:
			if i == v.retry {
				t.Errorf("MockServer unknown error. Details: %v", resp)
				return
			}
		case 400:
			if i == v.retry {
				t.Errorf("Incorrent request format: %v", resp)
				return
			}
		case 406:
			if i == v.retry {
				body, _ := ioutil.ReadAll(resp.Body)
				t.Errorf("Expectation was not met. Api not received specified numbers of time. Error response body: %v", string(body))
				return
			}
		default:
			if statusCode == 202 {
				t.Log("Expectation has been met")
				return
			}
		}
		if i == v.retry {
			t.Errorf("Unknown error. Check response: %v ", resp)
		}
		time.Sleep(1 * time.Second)
	}

}

func (ms *MockServerClient) ResetExpectations() error {
	uri := fmt.Sprintf("http://%s/mockserver/reset", ms.address)

	req, err := http.NewRequest(http.MethodPut, uri, nil)

	if err != nil {
		return fmt.Errorf("reset expectations fail. Error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	result, err := ms.client.Do(req)

	if err != nil {
		return err

	}
	if 200 != result.StatusCode {
		return fmt.Errorf("Unable to reset expectations. Error: %w.", err)
	}
	return nil
}

func (ms *MockServerClient) execute(request interface{}, uri string) (*http.Response, error) {

	if request == nil {
		return nil, errors.New("request not provided")
	}

	p, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("Json marshal error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(p))

	if err != nil {
		return nil, fmt.Errorf("expectation request fail. Error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return ms.client.Do(req)
}

func NewMockServerClient(address string, httpClient *http.Client) *MockServerClient {
	return &MockServerClient{
		address: address,
		client:  httpClient,
	}
}
