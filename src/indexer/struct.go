package indexer

import "tap-waiter/src/conf"

type DataSources struct {
	Paths  []string
	Type   string
	Params Params
}

type TapIndex []TapEntry

type TapEntry struct {
	Path      string      `json:"path"`
	FullPath  string      `json:"-"`
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
	Response  string
	Endpoint  conf.Endpoint
	Sources   []string
	Method    string
	Type      string
	Filter    FilterParams
	SortBy    string
	Ascending bool
}

type FilterParams struct {
	Prefix   string
	Operator string
	Suffix   []string
	Errors   []error
	Enabled  bool
}

func (ti TapIndex) Len() int {
	return len(ti)
}

func (ti TapIndex) Less(i, j int) bool {
	switch ti[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if ti[i].SortIndex == ti[j].SortIndex {
			return ti[i].Path > ti[j].Path
		}
		return ut.ToFloat(ti[i].SortIndex) < ut.ToFloat(ti[j].SortIndex)
	default:
		if ti[i].SortIndex.(string) == ti[j].SortIndex.(string) {
			return ti[i].Path > ti[j].Path
		}
		return ti[i].SortIndex.(string) < ti[j].SortIndex.(string)
	}
}

func (ti TapIndex) Swap(i, j int) {
	ti[i], ti[j] = ti[j], ti[i]
}
