package ad_data_transfer_test_cases

import (
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/controller"
)

type UnmarshalJSONExpects struct {
	AdJson      controller.AdvertisementJSON
	ExpectError error
}

type UnmarshalJSONTestCase struct {
	InputJson string
	Expects   UnmarshalJSONExpects
}

func parseTimeUnmarshalJSON(timestr string) time.Time {
	resultTime, _ := time.Parse("2006-01-02T15:04:05.000Z", timestr)
	return resultTime
}

func UnmarsalJSONTestCases() map[string]UnmarshalJSONTestCase {
	testCases := map[string]UnmarshalJSONTestCase{
		"valid JSON all attribute filleed": {
			InputJson: `
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
			`,
			Expects: UnmarshalJSONExpects{
				AdJson: controller.AdvertisementJSON{
					Title:   "test 1",
					StartAt: parseTimeUnmarshalJSON("2023-12-10T03:00:00.000Z"),
					EndAt:   parseTimeUnmarshalJSON("2023-12-31T16:00:00.000Z"),
					Conditions: controller.ConditionsJSON{
						AgeStart:  20,
						AgeEnd:    30,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("android"), ad.PlatformType("ios")},
					},
				},
			},
		},
		"valid empty JSON": {
			InputJson: `{}`,
			Expects: UnmarshalJSONExpects{
				AdJson: controller.AdvertisementJSON{
					Title:   "",
					StartAt: time.Time{},
					EndAt:   time.Time{},
					Conditions: controller.ConditionsJSON{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{},
						Countries: []ad.CountryCode{},
						Platforms: []ad.PlatformType{},
					},
				},
			},
		},
	}

	return testCases

}
