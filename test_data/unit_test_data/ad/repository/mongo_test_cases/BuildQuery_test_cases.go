package mongo_test_cases

import (
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"go.mongodb.org/mongo-driver/bson"
)

type BuildQueryTestCase struct {
	Input   BuildQueryInput
	Expects bson.D
}

type BuildQueryInput struct {
	Client *ad.Client
	Now    time.Time
}

func BuildQueryTestCases() map[string]BuildQueryTestCase {
	testCases := map[string]BuildQueryTestCase{
		"input with no missing parameter": {
			Input: BuildQueryInput{
				Client: &ad.Client{
					Age:      20,
					Gender:   ad.GenderType("M"),
					Country:  ad.CountryCode("TW"),
					Platform: ad.PlatformType("ios"),
				},
				Now: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			Expects: bson.D{
				{"conditions.countries", "TW"},
				{"conditions.ageStart", bson.D{{"$lte", 20}}},
				{"conditions.ageEnd", bson.D{{"$gte", 20}}},
				{"conditions.genders", "M"},
				{"conditions.platforms", "ios"},
			},
		},
		"input with missing parameter": {
			Input: BuildQueryInput{
				Client: &ad.Client{
					AgeMissing:      true,
					CountryMissing:  true,
					PlatformMissing: true,
					GenderMissing:   true,
				},
				Now: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			Expects: bson.D{},
		},
	}

	return testCases
}
