package optional

import (
	"encoding/json"
	"testing"

	"github.com/eduardolat/openroutergo/internal/assert"
)

func TestOptionalGenericType(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		// Test unset value
		unset := &Optional[string]{IsSet: false}
		data, err := unset.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "null", string(data))

		// Test set value
		set := &Optional[string]{IsSet: true, Value: "hello"}
		data, err = set.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `"hello"`, string(data))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		// Test null value
		var opt Optional[int]
		err := opt.UnmarshalJSON([]byte("null"))
		assert.NoError(t, err)
		assert.False(t, opt.IsSet)

		// Test valid value
		err = opt.UnmarshalJSON([]byte("42"))
		assert.NoError(t, err)
		assert.True(t, opt.IsSet)
		assert.Equal(t, 42, opt.Value)
	})
}

func TestDerivedTypes(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		// Marshal then unmarshal to verify full cycle
		original := String{IsSet: true, Value: "test"}
		data, err := json.Marshal(&original)
		assert.NoError(t, err)

		// Unmarshal
		var result String
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)
		assert.True(t, result.IsSet)
		assert.Equal(t, "test", result.Value)
	})

	t.Run("Int", func(t *testing.T) {
		// Test int with null value
		var num Int
		err := json.Unmarshal([]byte("null"), &num)
		assert.NoError(t, err)
		assert.False(t, num.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("99"), &num)
		assert.NoError(t, err)
		assert.True(t, num.IsSet)
		assert.Equal(t, 99, num.Value)

		// Marshal
		data, err := json.Marshal(&num)
		assert.NoError(t, err)
		assert.Equal(t, `99`, string(data))
	})

	t.Run("Bool", func(t *testing.T) {
		// Testing marshal/unmarshal false value
		original := Bool{IsSet: true, Value: false}
		data, err := json.Marshal(original)
		assert.NoError(t, err)

		var result Bool
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)
		assert.True(t, result.IsSet)
		assert.False(t, result.Value)
	})

	t.Run("Float64", func(t *testing.T) {
		// Test marshal
		original := Float64{IsSet: true, Value: 3.14}
		data, err := json.Marshal(original)
		assert.NoError(t, err)
		assert.Equal(t, `3.14`, string(data))

		// Test float64 with null value
		var num Float64
		err = json.Unmarshal([]byte("null"), &num)
		assert.NoError(t, err)
		assert.False(t, num.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("3.14"), &num)
		assert.NoError(t, err)
		assert.True(t, num.IsSet)
		assert.Equal(t, 3.14, num.Value)
	})

	t.Run("Any", func(t *testing.T) {
		// Test anyVar with null value
		var anyVar Any
		err := json.Unmarshal([]byte("null"), &anyVar)
		assert.NoError(t, err)
		assert.False(t, anyVar.IsSet)

		// Test anyVar with string
		err = json.Unmarshal([]byte(`"test"`), &anyVar)
		assert.NoError(t, err)
		assert.True(t, anyVar.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("{\"key\": \"value\"}"), &anyVar)
		assert.NoError(t, err)
		assert.True(t, anyVar.IsSet)
		anyMap := anyVar.Value.(map[string]interface{})
		assert.Equal(t, "value", anyMap["key"])

		// Marshal test
		data, err := json.Marshal(&anyVar)
		assert.NoError(t, err)
		assert.Equal(t, `{"key":"value"}`, string(data))
	})

	t.Run("MapStringAny", func(t *testing.T) {
		// Test map[string]any with null value
		var mapAny MapStringAny
		err := json.Unmarshal([]byte("null"), &mapAny)
		assert.NoError(t, err)
		assert.False(t, mapAny.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("{\"key\": \"value\"}"), &mapAny)
		assert.NoError(t, err)
		assert.True(t, mapAny.IsSet)
		mapAnyMap := mapAny.Value
		assert.Equal(t, "value", mapAnyMap["key"])

		// Marshal test
		data, err := json.Marshal(&mapAny)
		assert.NoError(t, err)
		assert.Equal(t, `{"key":"value"}`, string(data))
	})

	t.Run("MapStringString", func(t *testing.T) {
		// Test map[string]string with null value
		var mapString MapStringString
		err := json.Unmarshal([]byte("null"), &mapString)
		assert.NoError(t, err)
		assert.False(t, mapString.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("{\"key\": \"value\"}"), &mapString)
		assert.NoError(t, err)
		assert.True(t, mapString.IsSet)
		mapStringMap := mapString.Value
		assert.Equal(t, "value", mapStringMap["key"])

		// Marshal test
		data, err := json.Marshal(&mapString)
		assert.NoError(t, err)
		assert.Equal(t, `{"key":"value"}`, string(data))
	})

	t.Run("MapIntInt", func(t *testing.T) {
		// Test map[int]int with null value
		var mapInt MapIntInt
		err := json.Unmarshal([]byte("null"), &mapInt)
		assert.NoError(t, err)
		assert.False(t, mapInt.IsSet)

		// Test with value
		err = json.Unmarshal([]byte("{\"1\": 2, \"3\": 4}"), &mapInt)
		assert.NoError(t, err)
		assert.True(t, mapInt.IsSet)
		mapIntMap := mapInt.Value
		assert.Equal(t, 2, mapIntMap[1])
		assert.Equal(t, 4, mapIntMap[3])

		// Marshal test
		data, err := json.Marshal(&mapInt)
		assert.NoError(t, err)
		assert.Equal(t, `{"1":2,"3":4}`, string(data))
	})
}

func TestEmptyStringBehavior(t *testing.T) {
	// This test is critical because empty string isn't valid JSON
	t.Run("Direct vs json.Marshal", func(t *testing.T) {
		opt := &Optional[string]{IsSet: false}

		// Direct call returns empty string
		direct, err := opt.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "null", string(direct))

		// json.Marshal might handle it differently
		_, err = json.Marshal(opt)
		assert.NoError(t, err)
	})
}

func TestRealWorldUsage(t *testing.T) {
	// Test with a realistic struct
	type User struct {
		Name     string `json:"name"`
		Age      Int    `json:"age"`
		Email    String `json:"email,omitempty"`
		Verified Bool   `json:"verified,omitempty"`
	}

	t.Run("Complete serialization cycle", func(t *testing.T) {
		original := User{
			Name:     "Jane",
			Age:      Int{IsSet: true, Value: 30},
			Email:    String{IsSet: true, Value: "jane@example.com"},
			Verified: Bool{IsSet: false},
		}

		data, err := json.Marshal(original)
		assert.NoError(t, err)

		var result User
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		assert.Equal(t, original.Name, result.Name)
		assert.Equal(t, original.Age.Value, result.Age.Value)
		assert.True(t, result.Age.IsSet)
		assert.Equal(t, original.Email.Value, result.Email.Value)
		assert.True(t, result.Email.IsSet)
		assert.False(t, result.Verified.IsSet)
	})
}
