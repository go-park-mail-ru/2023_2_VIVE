package nullTypes

import (
	// "encoding/json"
	// "fmt"
	// "fmt"
	"encoding/json"
	"testing"
)

type testStruct struct {
	NullInt NullInt32  `json:"int32"`
	NullStr NullString `json:"str"`
	PStr    *string    `json:"pstr"`
}

var (
	str = "string pointer"

	validData = testStruct{
		NullInt: NewNullInt(1, true),
		NullStr: NewNullString("test string", true),
		PStr:    &str,
	}
	validJson = `{"int32":1, "str":"test string"}`

	nullData = testStruct{
		NullInt: NewNullInt(0, false),
		NullStr: NewNullString("", false),
		PStr:    nil,
	}
	nullJson = `{"int32":"null", "str":"null"}`
)

var testCases = []struct {
	data         testStruct
	expectedJson string
}{
	{
		data:         validData,
		expectedJson: validJson,
	},
	{
		data:         nullData,
		expectedJson: nullJson,
	},
}

func TestNullMarshal(t *testing.T) {
	for _, testCase := range testCases {
		actualJson, err := json.Marshal(testCase.data)
		// fmt.Printf("%s", actualJson)
		if err != nil {
			t.Errorf("unexpected error while marshaling data: %s", err)
		}

		actualJsonStr := string(actualJson)
		if actualJsonStr != testCase.expectedJson {
			t.Errorf("wrong answer while marshaling\n\texpected: %s\n\tgot: %s", testCase.expectedJson, actualJsonStr)
		}
	}
}
