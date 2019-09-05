package assigns

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Add inserts an assign and return id of inserted assign.
func Add(addr string, username string) (invkID string, c chan interface{}) {
	var (
		err error
		id  uuid.UUID
	)

	for {
		if id, err = uuid.NewRandom(); err != nil {
			err = errors.Wrapf(err, "Failed to create random UUID")
			panic(err)
		}

		invkID = id.String()

		if _, ok := assigns[invkID]; !ok {
			break
		}
	}

	c = make(chan interface{})
	assigns[invkID] = assign{
		id:       invkID,
		assignee: username,
		address:  addr,
		holder:   c,
	}

	return
}
