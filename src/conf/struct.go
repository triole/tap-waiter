package conf

import (
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

type Conf struct {
	FileName string
	Threads  int
	Port     int
	API      map[string]Endpoint
	Util     util.Util
	Lg       logseal.Logseal
}

type ConfContent struct {
	Port int                 `yaml:"port"`
	API  map[string]Endpoint `yaml:"api"`
}

type Endpoint struct {
	Source             string `yaml:"source"`
	SourceType         string
	URL                string   `yaml:"url"`
	RxFilter           string   `yaml:"regex_filter"`
	SortFileName       string   `yaml:"sort_file_name"`
	IgnoreList         []string `yaml:"regex_ignore_list"`
	MaxReturnSize      string   `yaml:"max_return_size"`
	MaxReturnSizeBytes uint64
	ReturnValues       ReturnValues `yaml:"return_values"`
}

type ReturnValues struct {
	SplitPath                bool `yaml:"split_path"`
	Metadata                 bool `yaml:"metadata"`
	Content                  bool `yaml:"content"`
	Size                     bool `yaml:"size"`
	LastMod                  bool `yaml:"lastmod"`
	Created                  bool `yaml:"created"`
	SplitMarkdownFrontMatter bool `yaml:"split_markdown_front_matter"`
}
