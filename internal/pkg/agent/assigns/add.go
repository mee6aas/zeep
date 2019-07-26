package assigns

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Add inserts an assign and return id of inserted assign.
func Add() (id string, c chan interface{}) {
	var (
		err error
		uid uuid.UUID
	)

	for {
		if uid, err = uuid.NewRandom(); err != nil {
			err = errors.Wrapf(err, "Failed to create random UUID")
			panic(err)
		}

		id = uid.String()

		if _, ok := assigns[id]; !ok {
			break
		}
	}

	c = make(chan interface{})
	assigns[id] = c

	return
}
