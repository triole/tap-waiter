# Tyson Tap ![build](https://github.com/triole/tyson-tap/actions/workflows/build.yaml/badge.svg) ![test](https://github.com/triole/tyson-tap/actions/workflows/test.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Configuration](#configuration)
- [URL Parameters](#url-parameters)
- [Filters](#filters)
  - [Logical Operators](#logical-operators)

<!-- /toc -->

## Synopsis

Tyson Tap offers an http api that serves json objects containing information about toml, yaml, json or markdown files that were indexed from specific folders.

## Configuration

A configuration file is required and defines the listening port and the api endpoint definitions. A simple example can be found below and another in the `testdata` folder.

```go mdox-exec="tail -n +2 conf/conf.yaml"
port: 17777 # port to listen at
api: # list of api endpoints
  # url where to endpoint is reachable
  all.json:

    # folder to be scanned for files
    folder: ../testdata

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
    return_values:
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
/all.json
/all.json?filter=front_matter.tags!=tag1
/all.json?filter=front_matter.tags===tag1
/all.json?filter=front_matter.tags==tag1
/all.json?filter=front_matter.tags!==tag1&order=desc
/all.json?filter=metadata.tags==ivoa
/all.json?filter=tags!==ivoa&order=desc
/all.json?sortby=front_matter.title
/all.json?sortby=size&order=desc
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
