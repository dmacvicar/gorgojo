package gorgojo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryMap(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	query := client.Query().Summary("KDE").AssignedTo("duncan")

	expected := map[string][]interface{}{
		"summary":     []interface{}{"KDE"},
		"assigned_to": []interface{}{"duncan"},
	}
	AssertDeepEqual(t, expected, query.QueryMap)
}

func TestQueryMapMultiple(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	query := client.Query().Summary("KDE").Open()

	expected := map[string][]interface{}{
		"summary": []interface{}{"KDE"},
		"status":  []interface{}{"new", "assigned", "needinfo", "reopened", "confirmed"},
	}

	AssertDeepEqual(t, expected, query.QueryMap)
}
