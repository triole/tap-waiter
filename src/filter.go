package main

func equalSlices(filter, content []string) (r bool) {
	r = false
	if len(filter) == len(content) {
		r = true
		for i := 0; i < len(filter); i++ {
			if filter[i] != content[i] {
				r = false
			}
		}
	}
	return
}

func notEqualSlices(filter, content []string) (r bool) {
	return !equalSlices(filter, content)
}
