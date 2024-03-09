package service_test_cases

import (
	"fmt"
)

// dummyEnum for testing ValidateSlice

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

var ValidEnum = func(de dummyEnum) bool {
	return validEnum[de]
}

type ValidEnumFuncType func(dummyEnum) bool

type ValidateSliceExpects struct {
	ExpectError error
}

type ValidateSliceTestCase struct {
	Slice     []dummyEnum
	ValidEnum ValidEnumFuncType
	Expects   ValidateSliceExpects
}

func ValidateSliceTestCases() map[string]ValidateSliceTestCase {
	testCases := map[string]ValidateSliceTestCase{
		"valid slice": { // all 1, 2, 3 belong to dummyEnum

			Slice:     []dummyEnum{dummyEnum("1"), dummyEnum("2"), dummyEnum("3")},
			ValidEnum: ValidEnum,
			Expects: ValidateSliceExpects{
				ExpectError: nil,
			},
		},
		"invalid slice": { // all 1, 2, 3 belong to dummyEnum

			Slice:     []dummyEnum{dummyEnum("1"), dummyEnum("4"), dummyEnum("3")},
			ValidEnum: ValidEnum,
			Expects: ValidateSliceExpects{
				ExpectError: fmt.Errorf("invalid item in slice of type %T: %v", dummyEnum("4"), dummyEnum("4")),
			},
		},
	}

	return testCases
}

// testCases := []struct {
// 	name           string
// 	slice          []dummyEnum
// 	expectedResult error
// }{

// 	{ // dummyEnum("4") doesn't belong to dummyEnum
// 		name:           "Invalid slice",
// 		slice:          []dummyEnum{dummyEnum("1"), dummyEnum("4"), dummyEnum("3")},
// 		expectedResult: fmt.Errorf("invalid item in slice of type %T: %v", dummyEnum("4"), dummyEnum("4")),
// 	},
// }
