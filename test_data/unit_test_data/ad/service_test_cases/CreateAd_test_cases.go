package service_test_cases

import (
	"fmt"

	"github.com/MarkLai0317/Advertising/ad"
)

type CreateAdExpects struct {
	ExpectError error
}

type CreateAdTestCase struct {
	Ad      ad.Advertisement
	Expects CreateAdExpects
}

func CreateAdTestCases() map[string]CreateAdTestCase {
	testCases := map[string]CreateAdTestCase{
		"create success": {
			Ad: ad.Advertisement{},
			Expects: CreateAdExpects{
				ExpectError: nil,
			},
		},
		"create fail: invalid Advertisement": {
			Ad: ad.Advertisement{},
			Expects: CreateAdExpects{
				ExpectError: fmt.Errorf("invalid Advertisement: %w", fmt.Errorf("invalid field in advertisement")),
			},
		},
		"create fail: DB create fail": {
			Ad: ad.Advertisement{},
			Expects: CreateAdExpects{
				ExpectError: fmt.Errorf("error creating Advertisement in DB: %w", fmt.Errorf("error in DB")),
			},
		},
	}

	return testCases
}
