package service_test_cases

import (
	"fmt"

	"github.com/MarkLai0317/Advertising/ad"
)

type AdvertiseExpects struct {
	AdSlice     []ad.Advertisement
	ExpectError error
}

type AdvertiseTestCase struct {
	Client  ad.Client
	Offset  int
	Limit   int
	Expects AdvertiseExpects
}

func AdvertiseTestCases() map[string]AdvertiseTestCase {
	testCases := map[string]AdvertiseTestCase{
		"Advertise success": {
			Client: ad.Client{
				Offset: 1,
				Limit:  1,
			},

			Expects: AdvertiseExpects{
				AdSlice:     []ad.Advertisement{{}, {}},
				ExpectError: nil,
			},
		},
		"invalid client": {
			Client: ad.Client{ // can be arbitraary client  because validate function is mocked
				Age:    -1,
				Offset: -1,
				Limit:  -1,
			},

			Expects: AdvertiseExpects{
				AdSlice:     nil,
				ExpectError: fmt.Errorf("invalid client: %w", fmt.Errorf("invalid client field")),
			},
		},
		"getAdvertisement error": {
			Client: ad.Client{
				Offset: 1,
				Limit:  1,
			},

			Expects: AdvertiseExpects{
				AdSlice:     nil,
				ExpectError: fmt.Errorf("GetAdvertisements Error: %w", fmt.Errorf("DB error")),
			},
		},
	}

	return testCases
}
