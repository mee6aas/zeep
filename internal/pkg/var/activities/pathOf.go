package acts

import "path/filepath"

// PathOf returns the path of the resource used by specified username.
func PathOf(username string) (p string, ok bool) {
	if _, ok = activities[username]; !ok {
		return
	}

	p = filepath.Join(rootDirPath, username)

	return
}
