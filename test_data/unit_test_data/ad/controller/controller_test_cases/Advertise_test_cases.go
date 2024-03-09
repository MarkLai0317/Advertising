package controller_test_cases

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/MarkLai0317/Advertising/ad"
)

// define the expected input, the expected behavior of mocks, and the expected output(expects)
type AdvertiseTestCase struct {
	Input   AdvertiseInput
	Mocks   AdvertiseMocks
	Expects AdvertiseExpects
}

type AdvertiseInput struct {
	Request        *http.Request
	ResponseWriter *httptest.ResponseRecorder
}

type AdvertiseMocks struct {
	QueryToClient            QueryToClientMock
	ServiceAdvertise         ServiceAdvertiseMock
	AdvertisementSliceToJSON AdvertisementSliceToJSONMock
}

type QueryToClientMock struct {
	Called       bool
	InputRequest *http.Request
	ReturnClient *ad.Client
	ReturnErr    error
}

type ServiceAdvertiseMock struct {
	Called        bool
	InputClient   *ad.Client
	ReturnAdSlice []ad.Advertisement
	ReturnErr     error
}

type AdvertisementSliceToJSONMock struct {
	Called           bool
	InputAdSlice     []ad.Advertisement
	ReturnAdResponse []byte
	ReturnErr        error
}

type AdvertiseExpects struct {
	StatusCode   int
	ResponseBody []byte
	ErrorMessage string
}

func AdvertiseTestCases() map[string]AdvertiseTestCase {
	testCases := map[string]AdvertiseTestCase{
		"Get Ad success": {
			Input: AdvertiseInput{
				Request:        newAdvertiseRequest(),
				ResponseWriter: httptest.NewRecorder(),
			},
			Mocks: AdvertiseMocks{
				QueryToClient: QueryToClientMock{
					Called:       true,
					InputRequest: newAdvertiseRequest(),
					ReturnClient: &ad.Client{
						Age:      1,
						Gender:   ad.GenderType("M"),
						Country:  ad.CountryCode("TW"),
						Platform: ad.PlatformType("ios"),
						Offset:   0,
						Limit:    2,
					},
					ReturnErr: nil,
				},
				ServiceAdvertise: ServiceAdvertiseMock{
					Called: true,
					ReturnAdSlice: []ad.Advertisement{
						{},
					},
					ReturnErr: nil,
				},
				AdvertisementSliceToJSON: AdvertisementSliceToJSONMock{
					Called: true,
					ReturnAdResponse: []byte(`{
						"items":[
							{
								"title": "test1",
								"endAt": "2024-01-03T00:00:00.000Z"
							},
							{
								"title": "test2",
								"endAt": "2024-02-03T00:00:00.000Z"
							}
						]
					}`),
					ReturnErr: nil,
				},
			},
			Expects: AdvertiseExpects{
				StatusCode: http.StatusOK,
				ResponseBody: []byte(`
					{
						"items":[
							{
								"title": "test1",
								"endAt": "2024-01-03T00:00:00.000Z"
							},
							{
								"title": "test2",
								"endAt": "2024-02-03T00:00:00.000Z"
							}
						]
					}
				`),
				ErrorMessage: "",
			},
		},
	}
	return testCases
}

func newAdvertiseRequest() *http.Request {
	req := httptest.NewRequest("GET", "/ad?age=1&gender=M&country=TW&platform=ios&offset=0&limit=1", bytes.NewBufferString(""))
	return req
}

// ageStr := queryValues.Get("age")
// genderStr := queryValues.Get("gender")
// countryStr := queryValues.Get("country")
// platformStr := queryValues.Get("platform")
// offsetStr := queryValues.Get("offset")
// limitStr := queryValues.Get("limit")

// func (c *Controller) Advertise(resp http.ResponseWriter, req *http.Request) {

// 	resp.Header().Set("Content-Type", "application/json")
// 	client, offset, limit, err := c.DataTransferer.QueryToClient(req)
// 	if err != nil {
// 		// Use http.Error to send the error message back to the client
// 		http.Error(resp, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	adSlice, err := c.AdService.Advertise(client, offset, limit)
// 	if err != nil {
// 		http.Error(resp, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	adResponse, err := c.DataTransferer.AdvertisementSliceToJSON(adSlice)
// 	if err != nil {
// 		http.Error(resp, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	_, err = resp.Write(adResponse)
// 	if err != nil {
// 		http.Error(resp, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp.WriteHeader(http.StatusOK)

// }
