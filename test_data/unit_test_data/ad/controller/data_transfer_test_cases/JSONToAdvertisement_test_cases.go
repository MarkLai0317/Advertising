package data_transfer_test_cases

import (
	"bytes"
	"net/http"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
)

type JSONToAdvertisementExpects struct {
	ResultAd    ad.Advertisement
	ExpectError error
}

type JSONToAdvertisementTestCase struct {
	Req     *http.Request
	Expects JSONToAdvertisementExpects
}

func JSONToAdvertisementTestCases() map[string]*JSONToAdvertisementTestCase {
	testCases := map[string]*JSONToAdvertisementTestCase{
		"all json attribute filled": {
			Req: newRequest([]byte(`
				{
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
					
				}
			`)),
			Expects: JSONToAdvertisementExpects{
				ResultAd: ad.Advertisement{
					Title:   "test 1",
					StartAt: parseTimeJSONToAdvertisement("2023-12-10T03:00:00.000Z"),
					EndAt:   parseTimeJSONToAdvertisement("2023-12-31T16:00:00.000Z"),
					Conditions: ad.Conditions{
						AgeStart:  20,
						AgeEnd:    30,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("android"), ad.PlatformType("ios")},
					},
				},
				ExpectError: nil,
			},
		},
		"empty json": {
			Req: newRequest([]byte(`{}`)),
			Expects: JSONToAdvertisementExpects{
				ResultAd:    ad.Advertisement{},
				ExpectError: nil,
			},
		},
	}
	return testCases
}

func newRequest(jsonPayload []byte) *http.Request {
	req, _ := http.NewRequest("POST", "/your-endpoint", bytes.NewBuffer(jsonPayload))

	// Set Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	return req
}

func parseTimeJSONToAdvertisement(timestr string) time.Time {
	resultTime, _ := time.Parse("2006-01-02T15:04:05.000Z", timestr)
	return resultTime
}
