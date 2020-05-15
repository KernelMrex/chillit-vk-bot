package vkdialogs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fileContent string = `{
	"greeting": {
		"text": "test1"
	}
}`

func TestLoader(t *testing.T) {
	dialogs, err := Load(bytes.NewReader([]byte(fileContent)))
	if err != nil {
		t.Errorf("error was not expected: %v", err)
		return
	}

	assert.Equal(t, dialogs["greeting"].Text, "test1", "Two strings must be the same")
}

func TestTextGetter(t *testing.T) {
	dialogs := make(Dialogs)
	dialogs["greeting"] = &Message{
		Text: "test1",
	}

	result, err := dialogs.GetText("greeting")
	assert.Equal(t, result, "test1", "Two strings must be the same")
	assert.NoError(t, err, "Error must be nil")

	_, err = dialogs.GetText("not existing")
	assert.Error(t, err, "Must be error")
}
