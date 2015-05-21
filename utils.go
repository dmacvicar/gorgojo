package gorgojo

import (
	//"fmt"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func AssertDeepEqual(t *testing.T, expected interface{}, actual interface{}) {
	assert := assert.New(t)

	if !reflect.DeepEqual(expected, actual) {

		expectedJson, err := json.MarshalIndent(expected, "", "  ")
		assert.Nil(err)

		actualJson, err := json.MarshalIndent(actual, "", "  ")
		assert.Nil(err)

		t.Log(string(expectedJson))
		t.Log(string(actualJson))
		t.Fail()
	}
}
