package gorgojo

// Query to Bugzilla on various attributes. It is chainable,
// so you can do:
//   query := client.Query().Summary("crashed").AssignedTo("john")
type Query struct {
	Client   *Client
	QueryMap map[string][]interface{}
}

func NewQuery(client *Client) *Query {
	return &Query{Client: client, QueryMap: make(map[string][]interface{})}
}

func (q *Query) appendQuery(key string, value interface{}) *Query {
	if _, ok := q.QueryMap[key]; !ok {
		q.QueryMap[key] = make([]interface{}, 0)
	}
	// key already exists
	q.QueryMap[key] = append(q.QueryMap[key], value)
	return q
}

// arbitrary field name
func (q *Query) Field(key string, value interface{}) *Query {
	return q.appendQuery(key, value)
}

// The login name of a user that a bug is assigned to.
func (q *Query) AssignedTo(who string) *Query {
	return q.appendQuery("assigned_to", who)
}

//  Searches for substrings in the single-line Summary field on bugs.
// If you specify an array, then bugs whose summaries match any of the passed
// substrings will be returned.
func (q *Query) Summary(what string) *Query {
	return q.appendQuery("summary", what)
}

// The current status of a bug (not including its resolution,
// if it has one, which is a separate field above).
func (q *Query) Status(status string) *Query {
	return q.appendQuery("status", status)
}

// Shortcut for all statuses that keep the bug "open"
func (q *Query) Open() *Query {
	return q.Status("new").Status("assigned").Status("needinfo").Status("reopened")
}

// L3 bugs
func (q *Query) L3() *Query {
	return q.appendQuery("summary", "L3")
}

func (q *Query) Result() ([]Bug, error) {
	return q.Client.Search(q.QueryMap)
}
