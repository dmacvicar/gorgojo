package gorgojo

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestQueryMap(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	query := client.Query().Summary("KDE").AssignedTo("duncan")

	expected := map[string]interface{}{
		"summary":     "KDE",
		"assigned_to": "duncan",
	}
	if !reflect.DeepEqual(expected, query.QueryMap) {
		t.Fail()
	}
}

func TestQueryMapMultiple(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	query := client.Query().Summary("KDE").Open()

	expected := map[string]interface{}{
		"summary": "KDE",
		"status":  []string{},
	}

	AssertDeepEqual(t, expected, query.QueryMap)
}
