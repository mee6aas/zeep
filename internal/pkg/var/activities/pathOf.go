package acts

import "path/filepath"

// PathOf returns directory path of resource used by given username.
func PathOf(username string) (p string, ok bool) {
	if _, ok = activities[username]; !ok {
		return
	}

	p = filepath.Join(rootDirPath, username)

	return
}
