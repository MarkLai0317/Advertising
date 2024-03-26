package main_test_cases

import (
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/repository"
)

type GetAdvertisementsTestCase struct {
	InputUrl string
	TestData []interface{}
	Expects  GetAdvertisementsExpects
}

type GetAdvertisementsExpects struct {
	ReturnData string
}

func GetAdvertisementsTestCases() map[string]GetAdvertisementsTestCase {

	testCases := map[string]GetAdvertisementsTestCase{
		"query the 8, 2 document": {
			InputUrl: "http://localhost:80/api/v1/ad?offset=0&limit=5&age=24&gender=M&platform=ios&country=TW",
			TestData: []interface{}{
				repository.AdvertisementMongo{
					Title:   "integration test 1:  first doc",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 1:  first doc",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 2: second doc",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  15,
						AgeEnd:    30,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP"), ad.CountryCode("US")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("web")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 3: startAt not match",
					StartAt: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 4: age not match",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  50,
						AgeEnd:    100,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 5: genders not match",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("F")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 6: contries not match",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 7: platform not match",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  1,
						AgeEnd:    50,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("android")},
					},
				},
				repository.AdvertisementMongo{
					Title:   "integration test 8: third doc",
					StartAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndAt:   time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
					Conditions: repository.ConditionsMongo{
						AgeStart:  5,
						AgeEnd:    30,
						Genders:   []ad.GenderType{ad.GenderType("M")},
						Countries: []ad.CountryCode{ad.CountryCode("TW"), ad.CountryCode("JP")},
						Platforms: []ad.PlatformType{ad.PlatformType("ios")},
					},
				},
			},
			Expects: GetAdvertisementsExpects{
				ReturnData: `{
					{
						Title: "integration test 2: second doc",
						EndAt: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Title: "integration test 8: third doc",
						EndAt: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
					},
				}`,
			},
		},
	}

	return testCases
}
