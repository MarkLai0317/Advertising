package advertisement_test

import (
	"fmt"
	"testing"
	"time"

	ad "github.com/MarkLai0317/Advertising/advertisement"
	"github.com/stretchr/testify/suite"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/mock"

	service_test "github.com/MarkLai0317/Advertising/test_data/advertisement/service_test_cases"

	advertisementMock "github.com/MarkLai0317/Advertising/mocks/advertisement"
)

type UnitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) TestTrue() {

}

func (uts *UnitTestSuite) TestValidateSlice() {

	// dummy enum for testing the general function ValidateSlice
	type dummyEnum string
	const (
		dummyEnum1 dummyEnum = "1"
		dummyEnum2 dummyEnum = "2"
		dummyEnum3 dummyEnum = "3"
	)

	var validEnum = map[dummyEnum]bool{
		dummyEnum1: true,
		dummyEnum2: true,
		dummyEnum3: true,
	}

	ValidEnum := func(de dummyEnum) bool {

		return validEnum[de]
	}

	testCases := []struct {
		name           string
		slice          []dummyEnum
		expectedResult error
	}{
		{ // all 1, 2, 3 belong to dummyEnum
			name:           "Valid slice",
			slice:          []dummyEnum{dummyEnum("1"), dummyEnum("2"), dummyEnum("3")},
			expectedResult: nil,
		},
		{ // dummyEnum("4") doesn't belong to dummyEnum
			name:           "Invalid slice",
			slice:          []dummyEnum{dummyEnum("1"), dummyEnum("4"), dummyEnum("3")},
			expectedResult: fmt.Errorf("invalid item in slice of type %T: %v", dummyEnum("4"), dummyEnum("4")),
		},
	}

	for _, testCase := range testCases {
		uts.Equal(testCase.expectedResult, ad.ValidateSlice(testCase.slice, ValidEnum), fmt.Sprintf("fail at test case %s", testCase.name))
	}
}

func (uts *UnitTestSuite) TestValidateAdvertisement() {

	// Patch functions with incrementing counters
	patches := gomonkey.NewPatches()

	allGenders := []ad.GenderType{ad.GenderType("Male"), ad.GenderType("Female")}
	allCountries := []ad.CountryCode{ad.CountryCode("AA"), ad.CountryCode("BB")}
	allPlatforms := []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")}

	// Patch AllGenders
	patches.ApplyFunc(ad.AllGenders, func() []ad.GenderType {
		copied := make([]ad.GenderType, len(allGenders))
		copy(copied, allGenders)
		return copied
	})

	// Patch AllCountries
	patches.ApplyFunc(ad.AllCountries, func() []ad.CountryCode {
		copied := make([]ad.CountryCode, len(allCountries))
		copy(copied, allCountries)
		return copied
	})

	// Patch AllPlatforms
	patches.ApplyFunc(ad.AllPlatforms, func() []ad.PlatformType {
		copied := make([]ad.PlatformType, len(allPlatforms))
		copy(copied, allPlatforms)
		return copied
	})

	// Patch time.Now()

	patches.ApplyFunc(time.Now, func() time.Time {
		fakeTime := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
		return fakeTime
	})

	defer patches.Reset()

	// type ad struct {
	// 	Title      string
	// 	StartAt    time.Time
	// 	EndAt      time.Time
	// 	Conditions Conditions
	// }

	// type Expects struct {
	// 	expectResultAd ad.Advertisement
	// 	expectError    error
	// }

	// type TestCase struct {
	// 	ad      ad.Advertisement
	// 	setup   func()
	// 	expects Expects
	// }

	// AgeStart  int
	// AgeEnd    int
	// Genders   []GenderType
	// Countries []CountryCode
	// Platforms []PlatformType
	testCases := service_test.ValidateAdvertisementTestCases(allGenders, allCountries, allPlatforms)

	for name, tc := range testCases {
		uts.Run(name, func() {
			// Create a new mock Validator
			mockValidator := advertisementMock.NewValidator(uts.T())

			// Set up mocks' expectations based on the test case

			// test case invalid gender slice
			if tc.Expects.ExpectError != nil && tc.Expects.ExpectError.Error() == fmt.Errorf("invalid advertisement : gender slice error").Error() {

				fmt.Printf("testestesteset")
				mockValidator.EXPECT().ValidateGenderSlice(mock.Anything).Return(fmt.Errorf("gender slice error")).Once()
				mockValidator.AssertNotCalled(uts.T(), "ValidateCountrySlice", mock.Anything)
				mockValidator.AssertNotCalled(uts.T(), "ValidatePlatformSlice", mock.Anything)

				// test case invalid country slice
			} else if tc.Expects.ExpectError != nil && tc.Expects.ExpectError.Error() == fmt.Errorf("invalid advertisement : country slice error").Error() {
				mockValidator.EXPECT().ValidateGenderSlice(mock.Anything).Return(nil).Once()
				mockValidator.EXPECT().ValidateCountrySlice(mock.Anything).Return(fmt.Errorf("country slice error")).Once()
				mockValidator.AssertNotCalled(uts.T(), "ValidatePlatformSlice", mock.Anything)

				// test case invalid platform slice
			} else if tc.Expects.ExpectError != nil && tc.Expects.ExpectError.Error() == fmt.Errorf("invalid advertisement : platform slice error").Error() {

				mockValidator.EXPECT().ValidateGenderSlice(mock.Anything).Return(nil).Once()
				mockValidator.EXPECT().ValidateCountrySlice(mock.Anything).Return(nil).Once()
				mockValidator.EXPECT().ValidatePlatformSlice(mock.Anything).Return(fmt.Errorf("platform slice error")).Once()

				// set expectation based on if the slices are empty
			} else {
				if len(tc.Ad.Conditions.Genders) == 0 { // if empty: will not execute validateSlice
					mockValidator.AssertNotCalled(uts.T(), "ValidateGenderSlice", mock.Anything)
				} else {
					mockValidator.EXPECT().ValidateGenderSlice(tc.Ad.Conditions.Genders).Return(nil).Once()
				}

				if len(tc.Ad.Conditions.Countries) == 0 {
					mockValidator.AssertNotCalled(uts.T(), "ValidateCountrySlice", mock.Anything)
				} else {
					mockValidator.EXPECT().ValidateCountrySlice(tc.Ad.Conditions.Countries).Return(nil).Once()
				}

				if len(tc.Ad.Conditions.Platforms) == 0 {
					mockValidator.AssertNotCalled(uts.T(), "ValidatePlatformSlice", mock.Anything)
				} else {
					mockValidator.EXPECT().ValidatePlatformSlice(tc.Ad.Conditions.Platforms).Return(nil).Once()
				}
			}

			// Call the function under test
			err := ad.ValidateAdvertisement(&tc.Ad, mockValidator)

			// Assert based on the expected outcome
			if tc.Expects.ExpectError != nil {

				//uts.NotEqual(err, nil)
				uts.Equal(tc.Expects.ExpectError, err)
			} else {
				uts.Equal(tc.Expects.ExpectResultAd, tc.Ad)
			}
		})
	}
}
