package activity

import (
	"encoding/json"
	"io/ioutil"
)

// Unmarshal parses activity configuration.
func Unmarshal(act []byte) (a Activity, e error) {
	e = json.Unmarshal(act, &a)
	a = Normalize(a)

	return
}

// UnmarshalFromFile parses activity manifest from given file.
func UnmarshalFromFile(actDesc string) (a Activity, e error) {
	var (
		act []byte
	)

	if act, e = ioutil.ReadFile(actDesc); e != nil {
		return
	}

	a, e = Unmarshal(act)

	return
}
