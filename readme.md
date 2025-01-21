# Tap Waiter ![build](https://github.com/triole/tap-waiter/actions/workflows/build.yaml/badge.svg) ![test](https://github.com/triole/tap-waiter/actions/workflows/test.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Configuration](#configuration)
- [URL Parameters](#url-parameters)
- [Filters](#filters)
  - [Logical Operators](#logical-operators)

<!-- /toc -->

## Synopsis

Tap Waiter offers an http api that serves json objects containing information about toml, yaml, json or markdown files that were indexed from specific folders.

## Configuration

A configuration file is required and defines the listening port and the api endpoint definitions. A simple example can be found below and another in the `testdata` folder.

```go mdox-exec="tail -n +2 conf/conf.yaml"
bind: 0.0.0.0:17777 # bind server to
default_cache_lifetime: 5m # e.g.: 180s, 3m... 0 for no caching

api: # list of api endpoints
  # url where to endpoint is reachable
  all.json:

    # data source, can be folder, file or url
    source: ../testdata

    # method to read the data, to specify http request method
    # not necessary for file reading
    # method: get

    # only detect files which fit the regex filter
    regex_filter: ".+"

    # name of all files that are considered to be sort files, there can be
    # multiple of them, each one refers to its folder and possible substructures
    sort_file_name: sort.yaml

    # files that should not be returned in index
    regex_ignore_list:
      - sort.yaml

    # maximum file size up to which file content will appear in json
    # default is 10K to avoid too big json outputs
    # only relevant if content return is enabled
    max_return_size: 10KB

    # set of return values to add to the final json
    # bool values are by default false and have to be enabled explicitely
    return:
      # an array of the file path split at the path separator
      split_path: false

      # file size in bytes
      size: false

      # file created date
      created: false

      # file modified date
      lastmod: false

      # display file content, note max_return_size
      content: false

      # front matter of markdown files can be split from content
      # set content=false and split_markdown_front_matter=true to
      # have only front matter in the final json
      split_markdown_front_matter: false
```

# Sort Files

Sort files contain a list that determines the `order` in which files in the same or in sub folders shell be returned in the index. The `exclusive` entry is a boolean value that switches between excluding files that are not on the list and displaying them. If true, the index will only contain the exact elements that are in the `order` list even if there are more files. If false, index will consist of the entries of `order` and any other relevant files will follow below.

```go mdox-exec="tail -n +2 testdata/sort.yaml"
exclusive: true
order:
  - filter
  - dump
  - markdown/subfolder2/1.md
```

# URL Parameters

Here are a few URL parameter examples which are hopefully self explanatory. Please keep in mind that special characters have to be url encoded.

```go mdox-exec="sh/display_test_urls.sh"
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/dump [36mrxfilter[0m=.+ [36murl[0m=/all.json
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/mapval/spec.yaml [36mspecs[0m=[map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/markdown/1.md exp:[title1] key:front_matter.title] map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/markdown/1.md exp:[tag1 tag2] key:front_matter.tags] map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/yaml/cpx/data_aip.yaml exp:[Data Services @ AIP] key:title] map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/yaml/cpx/data_aip.yaml exp:[open] key:metadata.access] map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/yaml/cpx/data_aip.yaml exp:[vo IVOA Daiquiri] key:metadata.tags] map[content_file:/home/ole/rolling/golang/projects/tap-waiter/testdata/dump/yaml/cpx/data_aip.yaml exp:[https://data.aip.de] key:metadata.url]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_co_tag1.yaml [36mspecs[0m=[map[exp:[markdown/1.md markdown/subfolder1/1.md markdown/subfolder1/2.md markdown/subfolder1/3.md markdown/subfolder1/4.md] urls:[/all.json?filter=front_matter.tags==tag1]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_eq_ivoa.yaml [36mspecs[0m=[map[exp:[yaml/cpx/data_aip.yaml] urls:[/all.json?filter=metadata.tags==ivoa]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_eq_tag1.yaml [36mspecs[0m=[map[exp:[markdown/subfolder1/1.md] urls:[/all.json?filter=front_matter.tags===tag1]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_nco_tag1.yaml [36mspecs[0m=[map[exp:[markdown/subfolder2/1.md markdown/subfolder2/2.md markdown/subfolder2/3.md] urls:[/all.json?filter=front_matter.tags!=tag1]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_neq_ivoa.yaml [36mspecs[0m=[map[exp:[markdown/subfolder2/3.md markdown/subfolder2/2.md markdown/subfolder2/1.md markdown/subfolder1/4.md markdown/subfolder1/3.md markdown/subfolder1/2.md markdown/subfolder1/1.md yaml/4.yaml markdown/no_front_matter.md markdown/1.md json/list.json html/3.html html/2.html html/1.html] urls:[/all.json?filter=tags!==ivoa&order=desc]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_filter_front_matter_tags_neq_tag1.yaml [36mspecs[0m=[map[exp:[markdown/subfolder2/3.md markdown/subfolder2/2.md markdown/subfolder2/1.md markdown/subfolder1/4.md markdown/subfolder1/3.md markdown/subfolder1/2.md markdown/1.md] urls:[/all.json?filter=front_matter.tags!==tag1&order=desc]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_sortby_front_matter_title.yaml [36mspecs[0m=[map[exp:[markdown/subfolder2/1.md markdown/subfolder1/1.md markdown/subfolder1/3.md markdown/subfolder1/4.md markdown/subfolder1/2.md markdown/1.md markdown/subfolder2/2.md markdown/subfolder2/3.md sort.yaml binary/binary_1k.file binary/binary_5k.file html/1.html html/2.html html/3.html html/more_than_10k.html html/sort.yaml json/list.json markdown/no_front_matter.md markdown/sort.yaml yaml/1.yaml yaml/2.yaml yaml/3.yaml yaml/4.yaml yaml/more_than_10k.yaml yaml/cpx/data_aip.yaml] urls:[/all.json?sortby=front_matter.title]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all_sortby_size.yaml [36mspecs[0m=[map[exp:[yaml/more_than_10k.yaml html/more_than_10k.html binary/binary_5k.file binary/binary_1k.file yaml/1.yaml yaml/4.yaml yaml/2.yaml yaml/3.yaml yaml/cpx/data_aip.yaml html/3.html html/2.html html/1.html markdown/sort.yaml markdown/subfolder1/3.md markdown/subfolder2/1.md markdown/subfolder1/2.md markdown/subfolder2/2.md markdown/subfolder1/4.md markdown/1.md markdown/subfolder2/3.md markdown/subfolder1/1.md markdown/no_front_matter.md json/list.json sort.yaml html/sort.yaml] urls:[/all.json?sortby=size&order=desc]]]
/home/ole/rolling/golang/projects/tap-waiter/testdata/specs/server/all.yaml [36mspecs[0m=[map[exp:[sort.yaml binary/binary_1k.file binary/binary_5k.file html/1.html html/2.html html/3.html html/more_than_10k.html html/sort.yaml json/list.json markdown/1.md markdown/no_front_matter.md markdown/sort.yaml yaml/1.yaml yaml/2.yaml yaml/3.yaml yaml/4.yaml yaml/more_than_10k.yaml markdown/subfolder1/1.md markdown/subfolder1/2.md markdown/subfolder1/3.md markdown/subfolder1/4.md markdown/subfolder2/1.md markdown/subfolder2/2.md markdown/subfolder2/3.md yaml/cpx/data_aip.yaml] urls:[/all.json]]]
```

# Filters

Filters supporting different logical operators are applied by url parameter. A filter contains of a `prefix` defining the part of the document to consider, an `operator` specifying the logical operation and a `suffix` which resembles the filter criterion. The content to be considered can be an array. For example when a markdown file has a list of tags. Therefore the filter criterion can be a comma separated list as well (see above).

## Logical Operators

| op               | returns a document if...                                               |
|------------------|------------------------------------------------------------------------|
| ===              | ...prefix and suffix are exactly equal                                 |
| !==              | ...prefix and suffix are not exactly equal                             |
| ==               | ...prefix contains suffix                                              |
| !=               | ...prefix does not contain suffix                                      |
| ==~              | ...if every regex suffix matches every prefix                          |
| =~               | ...if every regex suffix matches at least one prefix                   |
|                  |                                                                        |
| ** exactly equal | equal regarding size and every single entry                            |
| ** contains      | equal regarding every element of the suffix can be found in the prefix |
