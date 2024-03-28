//go:build unit
// +build unit

package ad_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/stretchr/testify/suite"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/mock"

	"github.com/MarkLai0317/Advertising/test_data/unit_test_data/ad/service_test_cases"

	advertisementMock "github.com/MarkLai0317/Advertising/mocks/ad"
)

type UnitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) TestCreateAd() {

	mockRepo := advertisementMock.NewRepository(uts.T())

	testCases := service_test_cases.CreateAdTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {
			patches := gomonkey.NewPatches()
			defer patches.Reset()
			if name == "create success" {

				patches.ApplyFunc(ad.ValidateAdvertisement, func(*ad.Advertisement, ad.Validator) error {
					return nil
				})
				mockRepo.EXPECT().CreateAdvertisement(mock.Anything).Return(nil).Once()

			} else if name == "create fail: invalid Advertisement" {

				patches.ApplyFunc(ad.ValidateAdvertisement, func(*ad.Advertisement, ad.Validator) error {
					return fmt.Errorf("invalid field in advertisement")
				})
				mockRepo.AssertNotCalled(uts.T(), "ValidateCountrySlice", mock.Anything)

			} else {

				patches.ApplyFunc(ad.ValidateAdvertisement, func(*ad.Advertisement, ad.Validator) error {
					return nil
				})
				mockRepo.EXPECT().CreateAdvertisement(mock.Anything).Return(fmt.Errorf("error in DB")).Once()

			}

			service := ad.NewService(mockRepo)
			err := service.CreateAd(&tc.Ad)

			uts.Equal(tc.Expects.ExpectError, err)

		})
	}

}

func (uts *UnitTestSuite) TestValidateAdvertisement() {

	// init gomonkey for stub
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	allGenders := []ad.GenderType{ad.GenderType("Male"), ad.GenderType("Female")}
	allCountries := []ad.CountryCode{ad.CountryCode("AA"), ad.CountryCode("BB")}
	allPlatforms := []ad.PlatformType{ad.PlatformType("ios"), ad.PlatformType("android")}

	// stub AllGenders
	patches.ApplyFunc(ad.AllGenders, func() []ad.GenderType {
		copied := make([]ad.GenderType, len(allGenders))
		copy(copied, allGenders)
		return copied
	})

	// stub AllCountries
	patches.ApplyFunc(ad.AllCountries, func() []ad.CountryCode {
		copied := make([]ad.CountryCode, len(allCountries))
		copy(copied, allCountries)
		return copied
	})

	// stub AllPlatforms
	patches.ApplyFunc(ad.AllPlatforms, func() []ad.PlatformType {
		copied := make([]ad.PlatformType, len(allPlatforms))
		copy(copied, allPlatforms)
		return copied
	})

	// stub time.Now()
	patches.ApplyFunc(time.Now, func() time.Time {
		fakeTime := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
		return fakeTime
	})

	testCases := service_test_cases.ValidateAdvertisementTestCases(allGenders, allCountries, allPlatforms)

	for name, tc := range testCases {
		uts.Run(name, func() {
			// Create a new mock Validator
			mockValidator := advertisementMock.NewValidator(uts.T())

			// Set up mocks' expectations based on the test case
			// test case invalid gender slice
			if name == "invalid gender" {

				mockValidator.EXPECT().ValidateGenderSlice(mock.Anything).Return(fmt.Errorf("gender slice error")).Once()
				mockValidator.AssertNotCalled(uts.T(), "ValidateCountrySlice", mock.Anything)
				mockValidator.AssertNotCalled(uts.T(), "ValidatePlatformSlice", mock.Anything)

				// test case invalid country slice
			} else if name == "invalid country" {

				mockValidator.EXPECT().ValidateGenderSlice(mock.Anything).Return(nil).Once()
				mockValidator.EXPECT().ValidateCountrySlice(mock.Anything).Return(fmt.Errorf("country slice error")).Once()
				mockValidator.AssertNotCalled(uts.T(), "ValidatePlatformSlice", mock.Anything)

				// test case invalid platform slice
			} else if name == "invalid platform" {

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
			mockValidator.AssertExpectations(uts.T())
			// Assert based on the expected outcome
			if tc.Expects.ExpectError != nil {
				uts.Equal(tc.Expects.ExpectError, err)
			} else {
				uts.Equal(tc.Expects.ExpectResultAd, tc.Ad)
			}
		})
	}
}

func (uts *UnitTestSuite) TestValidateSlice() {

	// dummy enum for testing the general function ValidateSlice

	testCases := service_test_cases.ValidateSliceTestCases()

	for name, tc := range testCases {

		uts.Run(name, func() {

			err := ad.ValidateSlice(tc.Slice, tc.ValidEnum)

			uts.Equal(tc.Expects.ExpectError, err)
		})

	}
}

func (uts *UnitTestSuite) TestAdvertise() {

	// init gomonkey for stub
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// stub time.Now()
	patches.ApplyFunc(time.Now, func() time.Time {
		return time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	})

	// stub ad.VaalidateClient()
	patches.ApplyFunc(ad.ValidateClient, func(client *ad.Client) error {
		if client.Age < 0 {
			return fmt.Errorf("invalid client field")
		}
		return nil
	})

	testCases := service_test_cases.AdvertiseTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {

			mockRepo := advertisementMock.NewRepository(uts.T())
			service := ad.NewService(mockRepo)

			// mock Repo
			if name == "getAdvertisement error" {
				mockRepo.On("GetAdvertisements",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil, fmt.Errorf("DB error"))
			} else if name == "invalid client" || name == "invalid offset" || name == "invalid limit" {
				mockRepo.AssertNotCalled(uts.T(), "GetAdvertisements")
			} else {
				mockRepo.On("GetAdvertisements",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(tc.Expects.AdSlice, nil).Once()
			}

			adSlice, err := service.Advertise(&tc.Client)
			mockRepo.AssertExpectations(uts.T())

			uts.Equal(tc.Expects.AdSlice, adSlice)
			uts.Equal(tc.Expects.ExpectError, err)

		})

	}

}

func (uts *UnitTestSuite) TestValidateClient() {

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch AllGenders
	patches.ApplyFunc(ad.ValidGender, func(gender ad.GenderType) bool {
		return gender == ad.GenderType("valid gender")

	})

	// Patch AllCountries
	patches.ApplyFunc(ad.ValidCountry, func(country ad.CountryCode) bool {
		return country == ad.CountryCode("valid country")
	})

	// Patch AllPlatforms
	patches.ApplyFunc(ad.ValidPlatform, func(platform ad.PlatformType) bool {
		return platform == ad.PlatformType("valid platform")
	})

	// Patch time.Now()

	testCases := service_test_cases.ValidateClientTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {

			err := ad.ValidateClient(&tc.Client)
			uts.Equal(tc.Expects.ExpectError, err)
		})
	}
}
