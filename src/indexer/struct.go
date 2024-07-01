package indexer

import "tyson-tap/src/conf"

type JoinerIndex []JoinerEntry

type JoinerEntry struct {
	Path      string      `json:"path"`
	SplitPath []string    `json:"split_path,omitempty"`
	Size      uint64      `json:"size,omitempty"`
	LastMod   int64       `json:"lastmod,omitempty"`
	Created   int64       `json:"created,omitempty"`
	Content   FileContent `json:"content,omitempty"`
	SortIndex interface{} `json:"-"`
}

type FileContent struct {
	Body        interface{} `json:"body,omitempty"`
	FrontMatter interface{} `json:"front_matter,omitempty"`
	Error       error       `json:"-"`
}

type Params struct {
	Endpoint  conf.Endpoint
	Filter    FilterParams
	SortBy    string
	Ascending bool
	Threads   int
}

type FilterParams struct {
	Prefix   string
	Operator string
	Suffix   []string
	Errors   []error
	Enabled  bool
}
