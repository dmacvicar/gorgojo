package gorgojo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientBadUrl(t *testing.T) {
	assert := assert.New(t)

	_, err := NewClient("hfds%ttps://b%ugzilla.gnome.org")
	assert.NotNil(err)
}

func TestClientVersion(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	version, err := client.Version()
	assert.Nil(err)
	assert.Equal("4.4.9", version)
}

func TestClientSearch(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	query := map[string][]interface{}{"summary": []interface{}{"Kolab"}}
	bugs, err := client.Search(query)

	assert.Equal(11, len(bugs))
}

func TestClientQuery(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	bugs, err := client.Query().Summary("Kolab").Result()
	assert.Equal(11, len(bugs))
}
