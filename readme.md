# Tyson Tap ![example workflow](https://github.com/triole/tyson-tap/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Conf](#conf)
- [```go mdox-exec=&quot;cat conf/conf.yaml&quot;](#go-mdox-execcat-confconfyaml)

<!-- /toc -->

## Synopsis

Tyson Tap offers an http api that serves json objects containing information about toml, yaml, json or markdown files that were indexed from specific folders.

## Conf

A configuration file is required and defines the listening port and the api definition. Here is a very early stage config example....

```go mdox-exec="cat conf/conf.yaml"
---
port: 8080
api:
  test.json:
    folder: ../testdata
    rxfilter: "\\.md$"
    no_file_content: false
    no_file_metadata: false

test_no_content.json:
  folder: ../testdata
  rxfilter: "\\.(md|json|toml|yaml)$"
  no_file_content: true
  no_file_metadata: false

test_no_content_no_metadata.json:
  folder: ../testdata
  rxfilter: "\\.(md|json|toml|yaml)$"
  no_file_content: true
  no_file_metadata: true
```

The settings above configure two api endpoints listening on `/test.json` and `/testslim.json` that serve json data indexed from the mentioned folders.
