package main

func readDataFile(filename string, basePath string, chin chan string, chout chan tJoinerEntry) {
	chin <- filename
	je := tJoinerEntry{
		Path: filename,
	}
	chout <- je
	_ = <-chin
}
