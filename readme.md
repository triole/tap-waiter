# Tyson Tap ![example workflow](https://github.com/triole/tyson-tap/actions/workflows/build.yaml/badge.svg)

- [Synopsis](#synopsis)
- [Configuration](#configuration)

## Synopsis

Tyson Tap offers an http api that serves json objects containing information about toml, yaml, json or markdown files that were indexed from specific folders.

## Configuration

A configuration file is required and defines the listening port and the api endpoint definitions. A simple example can be found below and another in the `testdata` folder.

```go mdox-exec="cat conf/conf.yaml | cat"
---
# bool values are by default false, if not mentioned in the config
# all return values have to be explicitely enabled

# port to listen at
port: 17777

# list of api endpoints
api:
  # url where to endpoint is reachable
  all.json:

    # folder to be scanned for files
    folder: ../testdata

    # only detect files which fit the regex filter
    rxfilter: ".+"

    # maximum file size up to which file content will appear in json
    # default is 10K to avoid too big json outputs
    # only relevant if content return is enabled
    max_return_size: 10KB

    # set of return values to add to the final json
    return_values:
      # an array of the file path split at the path separator
      split_path: false

      # file size in bytes
      size: false

      # file created date
      file_created: false

      # file modified date
      file_lastmod: false

      # display file content, note max_return_size
      content: false

      # front matter of markdown files can be split from content
      # set content=false and split_markdown_front_matter=true to
      # have only front matter in the final json
      split_markdown_front_matter: false
```
