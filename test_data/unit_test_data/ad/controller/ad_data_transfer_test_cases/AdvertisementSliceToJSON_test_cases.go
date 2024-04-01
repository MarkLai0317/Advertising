package ad_data_transfer_test_cases

import (
	"time"

	"github.com/MarkLai0317/Advertising/ad"
)

type AdvertisementSliceToJSONExpects struct {
	AdResponse []byte
}

type AdvertisementSliceToJSONTestCase struct {
	AdSlice []ad.Advertisement
	Expects AdvertisementSliceToJSONExpects
}

func AdvertisementSliceToJSONTestCases() map[string]AdvertisementSliceToJSONTestCase {
	testCases := map[string]AdvertisementSliceToJSONTestCase{
		"ValidSlice": {
			AdSlice: []ad.Advertisement{
				{
					Title: "test1",
					EndAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "test2",
					EndAt: time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			Expects: AdvertisementSliceToJSONExpects{
				AdResponse: []byte(`
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
			},
		},
	}
	return testCases
}
