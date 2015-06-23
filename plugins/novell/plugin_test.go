package novell

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransformApiUrlHook(t *testing.T) {
	assert := assert.New(t)

	n := novellPlugin{}

	// make sure oscrc is not read from the users directory
	OscrcPath = "/dummy"

	newUrl, err := n.TransformApiUrlHook("http://bugzilla.novell.com")
	assert.Nil(err)
	assert.Equal("http://apibugzilla.novell.com", newUrl)

	newUrl, err = n.TransformApiUrlHook("http://bugzilla.suse.com")
	assert.Nil(err)
	assert.Equal("http://apibugzilla.novell.com", newUrl)

	newUrl, err = n.TransformApiUrlHook("http://bugzilla.gnome.org")
	assert.Nil(err)
	assert.Equal("http://bugzilla.gnome.org", newUrl)
}

func TestTransformSiteUrlHook(t *testing.T) {
	assert := assert.New(t)

	n := novellPlugin{}
	newUrl, err := n.TransformSiteUrlHook("bnc")
	assert.Nil(err)
	assert.Equal("https://bugzilla.novell.com", newUrl)

	newUrl, err = n.TransformSiteUrlHook("bsc")
	assert.Nil(err)
	assert.Equal("https://bugzilla.suse.com", newUrl)

	newUrl, err = n.TransformSiteUrlHook("boo")
	assert.Nil(err)
	assert.Equal("https://bugzilla.opensuse.org", newUrl)

}

const (
	oscrc = `
# Force using of keyring for this API
#keyring = 1
trusted_prj=openSUSE:Factory SUSE:SLE-11:SP3 devel:languages:ruby:backports openSUSE:13.2 devel:languages:python Java:packages openSUSE:13.1 devel:languages:go Virtualization SUSE:SLE-12:GA Virtualization:containers

[https://api.opensuse.org]
user=foo
pass=bar
trusted_prj=SUSE:SLE-11:GA SUSE:SLE-11-SP3:GA SUSE:SLE-11-SP1:GA SUSE:SLE-11-SP3:Update SUSE:SLE-11-SP1:Update SUSE:SLE-11:Update SUSE:SLE-11-SP2:GA SUSE:SLE-11-SP2:Update SUSE:SLE-12:GA SUSE:SLE-12:Update openSUSE.org:openSUSE:Factory
`
)

func TestCredentialsParsing(t *testing.T) {
	assert := assert.New(t)

	data := strings.NewReader(oscrc)
	user, pass, err := ParseOscCredentials(data)

	assert.Nil(err)
	assert.Equal("foo", user)
	assert.Equal("bar", pass)
}
