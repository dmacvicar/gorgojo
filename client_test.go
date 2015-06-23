package gorgojo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestClientNamedQueryToUrl(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.gnome.org")
	assert.Nil(err)

	url := client.namedQueryToUrl("My Bugs")
	assert.Equal("https://bugzilla.gnome.org/buglist.cgi?cmdtype=runnamed&ctype=atom&namedcmd=My+Bugs", url)
}

// this test is ran manually, it depends on user data
func disabledTestClientFetchNamedQuery(t *testing.T) {
	assert := assert.New(t)

	client, err := NewClient("https://bugzilla.suse.com")
	assert.Nil(err)

	bugs, err := client.FetchNamedQuery("My Bugs")
	assert.Nil(err)

	assert.Equal(8, len(bugs))
}

var atomXml = `
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Bugzilla Bugs</title>
  <link rel="alternate" type="text/html"
        href="http://bugzilla.foobar.com/buglist.cgi?foo=bar"/>
  <updated>2015-06-10T07:09:28Z</updated>
  <id>http://bugzilla.foobar.com/buglist.cgi?foo=bar</id>
  <entry>
    <title>One bug</title>
    <link rel="alternate" type="text/html"
          href="http://bugzilla.foobar.com/show_bug.cgi?id=1234"/>
    <id>http://bugzilla.foobar.com/show_bug.cgi?id=1234</id>
    <author>
      <name>John Doe</name>
    </author>
    <updated>2015-06-10T07:09:28Z</updated>
    <summary type="html">
      A bad bug
    </summary>
  </entry>
  <entry>
    <title>Two bug</title>
    <link rel="alternate" type="text/html"
          href="http://bugzilla.foobar.com/show_bug.cgi?id=12345"/>
    <id>http://bugzilla.foobar.com/show_bug.cgi?id=12345</id>
    <author>
      <name>Joe Bob</name>
    </author>
    <updated>2015-05-20T09:43:40Z</updated>
    <summary type="html">
      bad bug
    </summary>
  </entry>
</feed>`

func TestClientParseAtomFeed(t *testing.T) {
	assert := assert.New(t)

	ids, err := parseAtomFeed(strings.NewReader(atomXml))
	assert.Nil(err)

	assert.Equal(2, len(ids))
	assert.Equal(1234, ids[0])
	assert.Equal(12345, ids[1])
}
