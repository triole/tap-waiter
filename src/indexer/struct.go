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

func (ji JoinerIndex) Len() int {
	return len(ji)
}

func (ji JoinerIndex) Less(i, j int) bool {
	switch ji[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if ji[i].SortIndex == ji[j].SortIndex {
			return ji[i].Path > ji[j].Path
		}
		return ut.ToFloat(ji[i].SortIndex) < ut.ToFloat(ji[j].SortIndex)
	default:
		if ji[i].SortIndex.(string) == ji[j].SortIndex.(string) {
			return ji[i].Path > ji[j].Path
		}
		return ji[i].SortIndex.(string) < ji[j].SortIndex.(string)
	}
}

func (ji JoinerIndex) Swap(i, j int) {
	ji[i], ji[j] = ji[j], ji[i]
}
