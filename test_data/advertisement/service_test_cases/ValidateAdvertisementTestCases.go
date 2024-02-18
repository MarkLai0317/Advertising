package service_test

import (
	"fmt"
	"time"

	ad "github.com/MarkLai0317/Advertising/advertisement"
)

type Expects struct {
	ExpectResultAd ad.Advertisement
	ExpectError    error
}

type TestCase struct {
	Ad      ad.Advertisement
	Setup   func()
	Expects Expects
}

func ValidateAdvertisementTestCases(allGenders []ad.GenderType, allCountries []ad.CountryCode, allPlatforms []ad.PlatformType) map[string]TestCase {
	testCases := map[string]TestCase{
		"all valid": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{ad.GenderType("M")},
					Countries: []ad.CountryCode{ad.CountryCode("AA")},
					Platforms: []ad.PlatformType{ad.PlatformType("ios")},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("AA")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
				ExpectError: nil,
			},
		},
		"all valid: StartAt less than Now": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{ad.GenderType("M")},
					Countries: []ad.CountryCode{ad.CountryCode("AA")},
					Platforms: []ad.PlatformType{ad.PlatformType("ios")},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("AA")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
				ExpectError: nil,
			},
		},
		"all valid: conditions set to default if gender, country, platform slices are empty": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{},
					Countries: []ad.CountryCode{},
					Platforms: []ad.PlatformType{},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   allGenders,
						Countries: allCountries,
						Platforms: allPlatforms,
					},
				},
				ExpectError: nil,
			},
		},

		"invalid Title": {
			Ad: ad.Advertisement{
				Title: "",
			},
			Expects: Expects{
				ExpectError: fmt.Errorf("Title cannot be empty"),
			},
		},

		"invalid EndTime: EndAt < Now": {
			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Expects: Expects{
				ExpectError: fmt.Errorf("EndAt cannot be smaller than current Time %s", time.Now().Format("2006-01-02T15:04:05.000Z")),
			},
		},

		"invalid EndTime: EndAt < StartAt": {
			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			Expects: Expects{
				ExpectError: fmt.Errorf("EndAt cannot be smaller than StartAt"),
			},
		},
		"invalid Age": {
			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart: -1,
					AgeEnd:   200,
				},
			},
			Expects: Expects{
				ExpectError: fmt.Errorf("Age should be in range 1 to 100"),
			},
		},

		"invalid gender": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{ad.GenderType("M")},
					Countries: []ad.CountryCode{ad.CountryCode("AA")},
					Platforms: []ad.PlatformType{ad.PlatformType("ios")},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("AA")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
				ExpectError: fmt.Errorf("invalid advertisement : %w", fmt.Errorf("gender slice error")),
			},
		},

		"invalid country": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{ad.GenderType("M")},
					Countries: []ad.CountryCode{ad.CountryCode("AA")},
					Platforms: []ad.PlatformType{ad.PlatformType("ios")},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("AA")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
				ExpectError: fmt.Errorf("invalid advertisement : %w", fmt.Errorf("country slice error")),
			},
		},
		"invalid platform": {

			Ad: ad.Advertisement{
				Title:   "Test Title",
				StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
				EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Conditions: ad.Conditions{
					AgeStart:  1,
					AgeEnd:    100,
					Genders:   []ad.GenderType{ad.GenderType("M")},
					Countries: []ad.CountryCode{ad.CountryCode("AA")},
					Platforms: []ad.PlatformType{ad.PlatformType("ios")},
				},
			},

			Expects: Expects{
				ExpectResultAd: ad.Advertisement{
					Title:   "Test Title",
					StartAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Conditions: ad.Conditions{
						AgeStart:  1,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("AA")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
				ExpectError: fmt.Errorf("invalid advertisement : %w", fmt.Errorf("platform slice error")),
			},
		},
		// Add more test cases as needed
	}
	return testCases
}
