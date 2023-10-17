package mock

import "sync"

type Sessions struct {
	SessionsList sync.Map
}

var SessionDB = Sessions{
	SessionsList: sync.Map{},
}
