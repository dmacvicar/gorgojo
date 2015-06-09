package novell

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"gopkg.in/ini.v1"
)

const (
	oscSection = "https://api.opensuse.org"
)

var OscrcPath string = os.Getenv("HOME") + "/.oscrc"

func CreatePlugin() novellPlugin {
	return novellPlugin{}
}

type novellPlugin struct {
}

func (p novellPlugin) TransformSiteUrlHook(url string) (string, error) {
	if url == "bnc" {
		return "https://bugzilla.novell.com", nil
	}
	if url == "bsc" {
		return "https://bugzilla.suse.com", nil
	}
	if url == "boo" {
		return "https://bugzilla.opensuse.org", nil
	}
	return url, nil
}

func (p novellPlugin) TransformApiUrlHook(urlStr string) (string, error) {
	urlStr = strings.Replace(urlStr, "/bugzilla.suse.com", "/apibugzilla.novell.com", 1)
	urlStr = strings.Replace(urlStr, "/bugzilla.novell.com", "/apibugzilla.novell.com", 1)

	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr, err
	}

	// add the credentials
	if _, err := os.Stat(OscrcPath); os.IsNotExist(err) {
		// no problem if there are not
		return urlStr, nil
	}

	user, pass, err := ReadOscCredentials(OscrcPath)
	if err != nil {
		return urlStr, nil
	}
	u.User = url.UserPassword(user, pass)
	return u.String(), nil
}


// returns the $HOME/.oscrc credentials
func ParseOscCredentials(reader io.Reader) (string, string, error) {

	username := ""
	password := ""

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return username, password, err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return username, password, err
	}

	section, err := cfg.GetSection(oscSection)
		if err != nil {
		return username, password, err
	}

	key, err := section.GetKey("user")
	if err != nil {
		return username, password, err
	}
	username = key.String()

	key, err = section.GetKey("pass")
	if err != nil {
		return username, password, err
	}
	password = key.String()

	return username, password, nil
}

// Reads the credentials file
func ReadOscCredentials(path string) (string, string, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	return ParseOscCredentials(file)
}
