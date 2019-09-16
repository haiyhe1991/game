package visitors

import "bytes"

// Visitor Gateway to service connection object
type Visitor struct {
	id   int
	sock int32
	data *bytes.Buffer
}
