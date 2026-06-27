package functionevaluation

import (
	"fmt"
	"testing"
)

func Test_Lower(t *testing.T) {
	/*
		Tests whether the lower function delivers the correct result
	*/

	testCases := [][]string{{"allLower", "test", "test"}, {"mixed", "TesT", "test"}, {"allUpper", "TEST", "test"}}

	for _, testCase := range testCases {
		t.Run(testCase[0], func(t *testing.T) {
			result, err := lower(map[string]any{"str": testCase[1]})

			if err != nil {
				t.Fatal(err)
			}

			if result != testCase[2] {
				t.Fatalf("Expected result to be %q but was %q instead", testCase[2], result)
			}
		})
	}
}

func Test_Lower_MissingParam(t *testing.T) {
	/*
		Tests whether the lower function returns an error when a parameter is missing
	*/

	_, err := lower(map[string]any{})

	if err == nil {
		t.Fatal("Expected an error as parameter is missing")
	}
}

func Test_Lower_WrongParamType(t *testing.T) {
	/*
		Tests whether the lower function returns an error when a parameter has a wrong type
	*/

	_, err := lower(map[string]any{"str": 4711})

	if err == nil {
		t.Fatal("Expected an error as parameter has wrong type")
	}
}

func Test_Upper(t *testing.T) {
	/*
		Tests whether the upper function delivers the correct result
	*/

	testCases := [][]string{{"allLower", "test", "TEST"}, {"mixed", "TesT", "TEST"}, {"allUpper", "TEST", "TEST"}}

	for _, testCase := range testCases {
		t.Run(testCase[0], func(t *testing.T) {
			result, err := upper(map[string]any{"str": testCase[1]})

			if err != nil {
				t.Fatal(err)
			}

			if result != testCase[2] {
				t.Fatalf("Expected result to be %q but was %q instead", testCase[2], result)
			}
		})
	}
}

func Test_Upper_MissingParam(t *testing.T) {
	/*
		Tests whether the upper function returns an error when a parameter is missing
	*/

	_, err := upper(map[string]any{})

	if err == nil {
		t.Fatal("Expected an error as parameter is missing")
	}
}

func Test_Upper_WrongParamType(t *testing.T) {
	/*
		Tests whether the upper function returns an error when a parameter has a wrong type
	*/

	_, err := upper(map[string]any{"str": 4711})

	if err == nil {
		t.Fatal("Expected an error as parameter has wrong type")
	}
}

func Test_Substr(t *testing.T) {
	/*
		Tests whether the substr function delivers the correct results
	*/

	testCases := []struct {
		str      string
		startIdx uint64
		endIdx   uint64
		expected string
	}{
		{str: "This is a test string", startIdx: 0, endIdx: 4, expected: "This"},
		{str: "This is a test string", startIdx: 5, endIdx: 7, expected: "is"},
		{str: "This is a test string", startIdx: 15, endIdx: 21, expected: "string"},
		{str: "This is a test string", startIdx: 15, endIdx: 22, expected: "string"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test_Substr_Ok_%d", i), func(t *testing.T) {
			result, err := substr(map[string]any{"str": testCase.str, "startIdx": testCase.startIdx, "endIdx": testCase.endIdx})

			if err != nil {
				t.Fatal(err)
			}

			if result != testCase.expected {
				t.Fatalf("Expected result to be %q, but was %q instead", testCase.expected, result)
			}
		})
	}
}

func Test_Substr_MissingParams(t *testing.T) {
	/*
		Tests whether the substr returns an error when parameters are missing
	*/

	testCases := []struct {
		testName string
		paramMap map[string]any
	}{
		{testName: "StrParamMissing", paramMap: map[string]any{}},
		{testName: "StartIdxParamMissing", paramMap: map[string]any{"str": "asdf"}},
		{testName: "EndIdxParamMissing", paramMap: map[string]any{"str": "asdf", "startIdx": 0}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test_Substr_MissingParams_%s", testCase.testName), func(t *testing.T) {
			_, err := substr(testCase.paramMap)

			if err == nil {
				t.Fatal("Expected function to return an error because of missing params")
			}
		})
	}
}

func Test_Substr_StrParamWrongType(t *testing.T) {
	/*
		Tests whether the substr returns an error when the str parameter is of wrong type
	*/

	_, err := substr(map[string]any{"str": 4711, "startIdx": uint64(0), "endIdx": uint64(4)})

	if err == nil {
		t.Fatal("Expected substr function to return an error as the type of the str parameter is wrong")
	}
}

func Test_Substr_StartIdxParamWrongType(t *testing.T) {
	/*
		Tests whether the substr returns an error when the startIdx parameter is of wrong type
	*/

	_, err := substr(map[string]any{"str": "4711", "startIdx": 0, "endIdx": uint64(2)})

	if err == nil {
		t.Fatal("Expected substr function to return an error as the type of the startIdx parameter is wrong")
	}
}

func Test_Substr_EndIdxParamWrongType(t *testing.T) {
	/*
		Tests whether the substr returns an error when the endIdx parameter is of wrong type
	*/

	_, err := substr(map[string]any{"str": "4711", "startIdx": uint64(0), "endIdx": 2})

	if err == nil {
		t.Fatal("Expected substr function to return an error as the type of the endIdx parameter is wrong")
	}
}

func Test_Substr_StartIdxToHigh(t *testing.T) {
	/*
		Tests whether the substr returns an error when the startIdx parameter is greater than the length of the string
	*/

	_, err := substr(map[string]any{"str": "4711", "startIdx": uint64(8), "endIdx": uint64(2)})

	if err == nil {
		t.Fatal("Expected substr function to return an error as the startIdx parameter is greater than the length of the string")
	}
}

func Test_Substr_StartIdxHigherThanEndIdx(t *testing.T) {
	/*
		Tests whether the substr returns an error when the startIdx parameter is greater than the endIdx parameter
	*/

	_, err := substr(map[string]any{"str": "4711", "startIdx": uint64(3), "endIdx": uint64(2)})

	if err == nil {
		t.Fatal("Expected substr function to return an error as the startIdx parameter is greater than the length of the string")
	}
}
