package service_test_cases

import (
	"fmt"

	"github.com/MarkLai0317/Advertising/ad"
)

type ValidateClientExpects struct {
	ExpectError error
}

type ValidateClientTestCase struct {
	Client  ad.Client
	Expects ValidateClientExpects
}

func ValidateClientTestCases() map[string]ValidateClientTestCase {
	testCases := map[string]ValidateClientTestCase{
		"valid Client": {
			Client: ad.Client{
				Age:      2,
				Gender:   ad.GenderType("valid gender"),
				Country:  ad.CountryCode("valid country"),
				Platform: ad.PlatformType("valid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: nil,
			},
		},
		"valid Client: missing conditions": {
			Client: ad.Client{
				Offset:          0,
				Limit:           1,
				AgeMissing:      true,
				GenderMissing:   true,
				CountryMissing:  true,
				PlatformMissing: true,
			},
			Expects: ValidateClientExpects{
				ExpectError: nil,
			},
		},
		"invalid Client: invalid Age < 1": {
			Client: ad.Client{
				Age:      0,
				Gender:   ad.GenderType("valid gender"),
				Country:  ad.CountryCode("valid country"),
				Platform: ad.PlatformType("valid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: fmt.Errorf("age should be in range 1 - 100"),
			},
		},
		"invalid Client: invalid Age > 100": {
			Client: ad.Client{
				Age:      200,
				Gender:   ad.GenderType("valid gender"),
				Country:  ad.CountryCode("valid country"),
				Platform: ad.PlatformType("valid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: fmt.Errorf("age should be in range 1 - 100"),
			},
		},
		"invalid Client: invalid Gender": {
			Client: ad.Client{
				Age:      1,
				Gender:   ad.GenderType("invalid gender"),
				Country:  ad.CountryCode("valid country"),
				Platform: ad.PlatformType("valid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: fmt.Errorf("invalid Gender type"),
			},
		},
		"invalid Client: invalid Country": {
			Client: ad.Client{
				Age:      1,
				Gender:   ad.GenderType("valid gender"),
				Country:  ad.CountryCode("invalid country"),
				Platform: ad.PlatformType("valid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: fmt.Errorf("invalid Country type"),
			},
		},
		"invalid Client: invalid Platform": {
			Client: ad.Client{
				Age:      1,
				Gender:   ad.GenderType("valid gender"),
				Country:  ad.CountryCode("valid country"),
				Platform: ad.PlatformType("invalid platform"),
				Offset:   0,
				Limit:    1,
			},
			Expects: ValidateClientExpects{
				ExpectError: fmt.Errorf("invalid Platform type"),
			},
		},
	}
	return testCases
}
