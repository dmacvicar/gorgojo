package gorgojo

import (
	"github.com/dmacvicar/gorgojo/plugins/novell"
)

type Plugin interface {
	// transforms the bugzilla website url
	TransformSiteUrlHook(url string) (string, error)
	// transforms the bugzilla API endpoint url
	TransformApiUrlHook(url string) (string, error)
}

func AllPlugins() []Plugin {
	return []Plugin{novell.CreatePlugin()}
}
