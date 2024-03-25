//go:build unit
// +build unit

package repository_test

import (
	"testing"

	"github.com/MarkLai0317/Advertising/ad/repository"
	"github.com/MarkLai0317/Advertising/test_data/unit_test_data/ad/repository/mongo_test_cases"
	"github.com/stretchr/testify/suite"
)

type MongoUnitTestSuite struct {
	suite.Suite
}

func TestMongoUnitTestSuite(t *testing.T) {
	suite.Run(t, &MongoUnitTestSuite{})
}

func (uts *MongoUnitTestSuite) TestBuildQuery() {
	testCases := mongo_test_cases.BuildQueryTestCases()

	for _, tc := range testCases {

		query := repository.ExportBuildMongoQuey(tc.Input.Client, tc.Input.Now)
		uts.Equal(tc.Expects, query)
	}

}
