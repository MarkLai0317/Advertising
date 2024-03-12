//go:build unit
// +build unit

package controller_test

import (
	"encoding/json"
	"testing"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/controller"
	"github.com/MarkLai0317/Advertising/test_data/unit_test_data/ad/controller/controller_test_cases"
	"github.com/stretchr/testify/suite"

	advertisementMock "github.com/MarkLai0317/Advertising/mocks/ad"
	controllerMock "github.com/MarkLai0317/Advertising/mocks/ad/controller"

	"github.com/stretchr/testify/mock"
)

type ControllerUnitTestSuite struct {
	suite.Suite
}

func TestControllerUnitTestSuite(t *testing.T) {
	suite.Run(t, &ControllerUnitTestSuite{})
}

func (uts *ControllerUnitTestSuite) TestCreateAdvertisement() {
	testCases := controller_test_cases.CreateAdvertisementTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {

			//init mock
			dataTransferMock := controllerMock.NewDataTransferer(uts.T())
			if tc.JsonToAdCalled {
				dataTransferMock.EXPECT().JSONToAdvertisement(tc.Request).Return(&ad.Advertisement{}, tc.JsonToAdError)
			}

			// EXPECT CreateAd to be called with anything (mock.Anything)
			serviceMock := advertisementMock.NewUseCase(uts.T())

			if tc.CreateAdCalled {
				serviceMock.EXPECT().CreateAd(mock.Anything).Return(tc.CreateAdError)
			}

			// init controller
			adController := controller.NewAdvertisementController(serviceMock, dataTransferMock)

			// init input param
			req := tc.Request
			resp := tc.ResponseWriter

			// run the tested function
			adController.CreateAdvertisement(resp, req)

			// check return status code
			uts.Equal(tc.Expects.StatusCode, resp.Code, "response status code not correct")

			if tc.JsonToAdError != nil {
				var serviceErr controller.ServiceError
				err := json.NewDecoder(resp.Body).Decode(&serviceErr)
				uts.NoError(err)
				uts.Contains(serviceErr.Message, tc.Expects.ErrorMessage, "error message not correct")
			} else if tc.CreateAdError != nil {
				var serviceErr controller.ServiceError
				err := json.NewDecoder(resp.Body).Decode(&serviceErr)
				uts.NoError(err)
				uts.Contains(serviceErr.Message, tc.Expects.ErrorMessage, "error message not correct")
			}

		})
	}
}

func (uts *ControllerUnitTestSuite) TestAdvertise() {
	testCases := controller_test_cases.AdvertiseTestCases()
	for name, tc := range testCases {
		uts.Run(name, func() {

			// init mocks
			dataTransferMock := controllerMock.NewDataTransferer(uts.T())
			serviceMock := advertisementMock.NewUseCase(uts.T())

			// set mocks expected param and return
			if tc.Mocks.QueryToClient.Called {
				dataTransferMock.EXPECT().QueryToClient(tc.Input.Request).Return(
					tc.Mocks.QueryToClient.ReturnClient,
					tc.Mocks.QueryToClient.ReturnErr,
				)
			}

			if tc.Mocks.ServiceAdvertise.Called {
				serviceMock.EXPECT().Advertise(tc.Mocks.QueryToClient.ReturnClient).Return(
					tc.Mocks.ServiceAdvertise.ReturnAdSlice,
					tc.Mocks.ServiceAdvertise.ReturnErr,
				)
			}

			if tc.Mocks.AdvertisementSliceToJSON.Called {
				dataTransferMock.EXPECT().AdvertisementSliceToJSON(tc.Mocks.ServiceAdvertise.ReturnAdSlice).Return(
					tc.Mocks.AdvertisementSliceToJSON.ReturnAdResponse,
					tc.Mocks.AdvertisementSliceToJSON.ReturnErr,
				)
			}

			// init the struct which impliment the method that we want to test
			adController := controller.NewAdvertisementController(serviceMock, dataTransferMock)

			// init input param
			req := tc.Input.Request
			resp := tc.Input.ResponseWriter

			// call the tested method
			adController.Advertise(resp, req)
			// check statuscode and error message if any
			uts.Equal(tc.Expects.StatusCode, resp.Code, "response status code not correct")

			// when return error messeage
			if tc.Expects.ErrorMessage != "" {
				// expect error message
				serviceErr := controller.ServiceError{}
				err := json.NewDecoder(resp.Body).Decode(&serviceErr)
				uts.NoError(err)
				uts.Contains(serviceErr.Message, tc.Expects.ErrorMessage, "error message not correct")

				// when no error and return advertisement list
			} else {
				// expect correct ad list
				respBodyObj := map[string]interface{}{}
				expectObj := map[string]interface{}{}

				err := json.NewDecoder(resp.Body).Decode(&respBodyObj)
				uts.NoError(err, "decode resp body should not be error")

				err = json.Unmarshal(tc.Expects.ResponseBody, &expectObj)
				uts.NoError(err, "unmarshal expect body should not be error. check test cases")

				uts.Equal(expectObj, respBodyObj, "body returned incorrect")

			}

		})
	}

}
