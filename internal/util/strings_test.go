package util

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagsFromString(t *testing.T) {
	// given
	str := "auth-{name}-{age}"
	// when
	tags := GetTagsFromString(str)
	// then
	assert.Equal(t, []string{"name", "age"}, tags)
}

func TestReplaceTags(t *testing.T) {
	// given
	str := "auth-{name}-{age}"
	tags := map[string]interface{}{
		"name": "taste",
		"age":  "25",
	}
	JSON, _ := json.Marshal(tags)
	// when
	newStr, err := ReplaceTags(context.Background(), str, JSON)
	// then
	assert.Nil(t, err)
	assert.Equal(t, "auth-taste-25", newStr)
}

func TestReplaceTagsError(t *testing.T) {
	// given
	str := "auth-{name}-{age}"
	tags := map[string]interface{}{
		"age": "25",
	}
	JSON, _ := json.Marshal(tags)
	// when
	newStr, err := ReplaceTags(context.Background(), str, JSON)
	// then
	assert.NotNil(t, err)
	assert.Error(t, err, "Tag not found")
	assert.Equal(t, "", newStr)
}
