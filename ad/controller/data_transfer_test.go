package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"testing"

	"github.com/MarkLai0317/Advertising/ad"
	"github.com/MarkLai0317/Advertising/ad/controller"

	"github.com/MarkLai0317/Advertising/test_data/unit_test_data/ad/controller/data_transfer_test_cases"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/suite"
)

type DataTransferUnitTestSuite struct {
	suite.Suite
}

func TestDataTransferUnitTestSuite(t *testing.T) {
	suite.Run(t, &DataTransferUnitTestSuite{})
}

func (uts *DataTransferUnitTestSuite) TestUnmarshalJSON() {

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// use stub decouple the test of UnmarshalJSON with defaultAdvertisementJSON
	// test won't fail if we change default value for budsiness requirment
	patches.ApplyFunc(controller.DefaultAdvertisementJSON, func() *controller.AdvertisementJSON {
		return &controller.AdvertisementJSON{
			Conditions: controller.ConditionsJSON{
				AgeStart:  1,
				AgeEnd:    100,
				Genders:   []ad.GenderType{},
				Countries: []ad.CountryCode{},
				Platforms: []ad.PlatformType{},
			},
		}
	})
	testCases := data_transfer_test_cases.UnmarsalJSONTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {
			adJson := controller.AdvertisementJSON{}
			reader := strings.NewReader(tc.InputJson)
			req, _ := http.NewRequest("POST", "/endpoint", reader)

			err := json.NewDecoder(req.Body).Decode(&adJson)
			uts.Equal(nil, err, "error should be nil")
			uts.Equal(tc.Expects.AdJson, adJson, "adJson not equal")
		})
	}
}

func (uts *DataTransferUnitTestSuite) TestJSONToAdvertisement() {

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	//tempAd := &controller.AdvertisementJSON{}
	patches.ApplyMethod(reflect.TypeOf(&controller.AdvertisementJSON{}), "UnmarshalJSON", func(adJson *controller.AdvertisementJSON, text []byte) error {
		// Call the original json.Unmarshal
		type AdAlias controller.AdvertisementJSON
		ad := AdAlias{}
		if err := json.Unmarshal(text, &ad); err != nil {
			return fmt.Errorf("JSON unmarshal fail: %w", err)
		}
		*adJson = controller.AdvertisementJSON(ad)
		return nil
	})
	// use stub decouple the test of UnmarshalJSON with defaultAdvertisementJSON
	// test won't fail if we change default value for budsiness requirment
	testCases := data_transfer_test_cases.JSONToAdvertisementTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {
			dt := controller.AdDataTransferer{}
			resultAd, err := dt.JSONToAdvertisement(tc.Req)
			uts.Equal(tc.Expects.ExpectError, err, "error not equal")

			uts.Equal(&tc.Expects.ResultAd, resultAd, "result ad not equal")

		})
	}
}

func (uts *DataTransferUnitTestSuite) TestQueryToClient() {
	testCases := data_transfer_test_cases.QueryToClientTestCases()

	for name, tc := range testCases {
		uts.Run(name, func() {
			df := controller.AdDataTransferer{}
			client, err := df.QueryToClient(tc.Req)
			uts.Equal(tc.Expects.Client, client, "client not equal")
			uts.Equal(tc.Expects.ExpectError, err)

		})
	}
}

func (uts *DataTransferUnitTestSuite) TestAdvertisementSliceToJSON() {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	//tempAd := &controller.AdvertisementJSON{}
	patches.ApplyMethod(reflect.TypeOf(controller.ResponseTime{}), "MarshalJSON", func(rt controller.ResponseTime) ([]byte, error) {
		t := time.Time(rt)
		formattedTime := t.Format("2006-01-02T15:04:05.000Z") // Adjust the layout according to your requirement
		return []byte(fmt.Sprintf(`"%s"`, formattedTime)), nil
	})

	testCases := data_transfer_test_cases.AdvertisementSliceToJSONTestCases()
	for name, tc := range testCases {
		uts.Run(name, func() {
			dt := controller.AdDataTransferer{}
			adResponse, err := dt.AdvertisementSliceToJSON(tc.AdSlice)
			uts.Equal(nil, err, "should not be error")

			// unmashal []byte of json string to map for ignoring indentation of json string
			var expectResponse, actualResponse map[string]interface{}

			err = json.Unmarshal(tc.Expects.AdResponse, &expectResponse)
			uts.Equal(nil, err, "expect response err should be nil: check if test file is correct")

			err = json.Unmarshal(adResponse, &actualResponse)
			uts.Equal(nil, err, "ad response should be nil")

			// compare unmashal json
			uts.Equal(expectResponse, actualResponse, "should be equal")

		})
	}
}
