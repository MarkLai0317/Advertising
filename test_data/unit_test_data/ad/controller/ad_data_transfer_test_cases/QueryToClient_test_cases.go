package ad_data_transfer_test_cases

import (
	"fmt"
	"net/http"

	"github.com/MarkLai0317/Advertising/ad"
)

type QueryToClientExpects struct {
	Client      *ad.Client
	ExpectError error
}

type QueryToClientTestCase struct {
	Req     *http.Request
	Expects QueryToClientExpects
}

func QueryToClientTestCases() map[string]QueryToClientTestCase {
	testCases := map[string]QueryToClientTestCase{
		"valid Query": {
			Req: newQuery("1", "M", "TW", "ios", "2", "1", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client: &ad.Client{
					Age:             1,
					Gender:          ad.GenderType("M"),
					Country:         ad.CountryCode("TW"),
					Platform:        ad.PlatformType("ios"),
					Offset:          2,
					Limit:           1,
					AgeMissing:      false,
					GenderMissing:   false,
					CountryMissing:  false,
					PlatformMissing: false,
				},

				ExpectError: nil,
			},
		},
		"condition empty string": {
			Req: newQuery("", "", "", "", "2", "1", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client: &ad.Client{
					Age:             0,
					Gender:          ad.GenderType(""),
					Country:         ad.CountryCode(""),
					Platform:        ad.PlatformType(""),
					Offset:          2,
					Limit:           1,
					AgeMissing:      true,
					GenderMissing:   true,
					CountryMissing:  true,
					PlatformMissing: true,
				},
				ExpectError: nil,
			},
		},
		"offset or limit empty": {
			Req: newQuery("1", "M", "TW", "ios", "", "", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client: &ad.Client{
					Age:             1,
					Gender:          ad.GenderType("M"),
					Country:         ad.CountryCode("TW"),
					Platform:        ad.PlatformType("ios"),
					Offset:          0,
					Limit:           5,
					AgeMissing:      false,
					GenderMissing:   false,
					CountryMissing:  false,
					PlatformMissing: false,
				},

				ExpectError: nil,
			},
		},
		"condition not provide": {
			Req: newQuery("", "", "", "", "2", "1", false, false, false, false, true, true),
			Expects: QueryToClientExpects{
				Client: &ad.Client{
					Offset:          2,
					Limit:           1,
					AgeMissing:      true,
					GenderMissing:   true,
					CountryMissing:  true,
					PlatformMissing: true,
				},
				ExpectError: nil,
			},
		},
		"offset and limit not provide": {
			Req: newQuery("1", "M", "TW", "ios", "", "", true, true, true, true, false, false),
			Expects: QueryToClientExpects{
				Client: &ad.Client{
					Age:      1,
					Gender:   ad.GenderType("M"),
					Country:  ad.CountryCode("TW"),
					Platform: ad.PlatformType("ios"),
					Offset:   0,
					Limit:    5,
				},

				ExpectError: nil,
			},
		},
		"age format error": {
			Req: newQuery("1age", "M", "TW", "ios", "2", "1", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client:      nil,
				ExpectError: fmt.Errorf("age format incorrect"),
			},
		},
		"offset format error": {
			Req: newQuery("1", "M", "TW", "ios", "off2", "1", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client:      nil,
				ExpectError: fmt.Errorf("offset format incorrect"),
			},
		},
		"limit format error": {
			Req: newQuery("1", "M", "TW", "ios", "2", "limit1", true, true, true, true, true, true),
			Expects: QueryToClientExpects{
				Client:      nil,
				ExpectError: fmt.Errorf("limit format incorrect"),
			},
		},
	}
	return testCases

}

// use*  indicate having the key in query or not
func newQuery(age string, gender string, country string, platform string, offset string, limit string,
	useAge bool, useGender bool, useCountry bool, usePlatform bool, useOffset bool, useLimit bool) *http.Request {
	req, _ := http.NewRequest("GET", "http://test.com", nil)

	// Add query parameters
	query := req.URL.Query()
	if useAge {
		query.Add("age", age)
	}
	if useGender {
		query.Add("gender", gender)
	}
	if useCountry {
		query.Add("country", country)
	}
	if usePlatform {
		query.Add("platform", platform)
	}
	if useOffset {
		query.Add("offset", offset)
	}
	if useLimit {
		query.Add("limit", limit)
	}
	req.URL.RawQuery = query.Encode()
	return req
}
