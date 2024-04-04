package controller_test_cases

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// func (c *Controller) CreateAdvertisement(resp http.ResponseWriter, req *http.Request) {

// 	resp.Header().Set("Content-Type", "application/json")
// 	newAd, err := c.DataTransferer.JSONToAdvertisement(req)
// 	if err != nil {
// 		resp.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(resp).Encode(ServiceError{Message: "Error decode req.Body"})
// 		return
// 	}

// 	err = c.AdService.CreateAd(newAd)

// 	if err != nil {
// 		resp.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(resp).Encode(ServiceError{Message: "Error creating advertisement"})
// 		return
// 	}

// 	resp.WriteHeader(http.StatusOK)

// }

type CreateAdvertisementExpects struct {
	StatusCode   int
	ErrorMessage string
}

type CreateAdvertisementTestCase struct {
	Request        *http.Request
	ResponseWriter *httptest.ResponseRecorder
	JsonToAdCalled bool
	JsonToAdError  error
	CreateAdCalled bool
	CreateAdError  error

	Expects CreateAdvertisementExpects
}

func CreateAdvertisementTestCases() map[string]CreateAdvertisementTestCase {
	jsonData := `{
		"title": "test 1",
		"startAt": "2023-12-10T03:00:00.000Z", 
		"endAt": "2023-12-31T16:00:00.000Z", 
		"conditions":{
			"ageStart": 20,
			"ageEnd": 30,
			"gender": ["M"],
			"country": ["TW", "JP"], 
			"platform": ["android", "ios"]
		} 
	}`

	testCases := map[string]CreateAdvertisementTestCase{
		"successful creation": {
			Request:        newCreateAdvertisementRequest("POST", "/ad", jsonData),
			ResponseWriter: httptest.NewRecorder(),
			JsonToAdCalled: true,
			CreateAdCalled: true,
			Expects: CreateAdvertisementExpects{
				StatusCode:   http.StatusOK,
				ErrorMessage: "",
			},
		},
		"error JSON to advertisement": {
			Request:        newCreateAdvertisementRequest("POST", "/ad", jsonData),
			ResponseWriter: httptest.NewRecorder(),
			JsonToAdCalled: true,
			JsonToAdError:  fmt.Errorf("error decoding JSON"),
			Expects: CreateAdvertisementExpects{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error JSON to advertisement:",
			},
		},
		"error creating advertisement": {
			Request:        newCreateAdvertisementRequest("POST", "/ad", jsonData),
			ResponseWriter: httptest.NewRecorder(),
			JsonToAdCalled: true,
			CreateAdCalled: true,
			CreateAdError:  fmt.Errorf("error creating advertisement"),
			Expects: CreateAdvertisementExpects{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "Error creating advertisement:",
			},
		},
	}

	return testCases
}

func newCreateAdvertisementRequest(method string, url string, body string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
