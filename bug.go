package gorgojo

import (
	"time"
)

type Bug struct {
	// The unique numeric id of this bug.
	Id int `xmlrpc:"id"`
	// The summary of this bug.
	Summary string `xmlrpc:"summary"`
	// When the bug was created
	CreationTime time.Time `xmlrpc:"creation_time"`
	// The login name of the user to whom the bug is assigned.
	AssignedTo string `xmlrpc:"assigned_to"`
	// he name of the current component of this bug.
	Component string `xmlrpc:"component"`
	// When the bug was last changed.
	LastChangeTime time.Time `xmlrpc:"last_change_time"`
	// The current severity of the bug.
	Severity string `xmlrpc:"severity"`
	// The current status of the bug.
	Status string `xmlrpc:"status"`
}
