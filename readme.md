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
  all.json:
    folder: ../testdata
    rxfilter: ".+"
    max_return_size: 10K
    return_values:
      size: false
      file_created: false
      file_lastmod: false
      content: false
      split_markdown_front_matter: false

  list.json:
    folder: ../testdata
    rxfilter: ".+"
    return_values:
      size: false
      file_created: false
      file_lastmod: false
      content: false
      split_markdown_front_matter: false

  md_body_only.json:
    folder: ../testdata
    rxfilter: "\\.(md)$"
    max_return_size: 10K
    return_values:
      size: false
      file_created: false
      file_lastmod: false
      content: true
      split_markdown_front_matter: false

  md_frontmatter_only.json:
    folder: ../testdata
    rxfilter: "\\.(md)$"
    max_return_size: 10K
    return_values:
      size: false
      file_created: false
      file_lastmod: false
      content: false
      split_markdown_front_matter: true

  yaml_dates.json:
    folder: ../testdata
    rxfilter: "\\.(ya?ml)$"
    max_return_size: 1K
    return_values:
      file_created: true
      file_lastmod: true

  yaml_sizes_created_and_content.json:
    folder: ../testdata
    rxfilter: "\\.(ya?ml)$"
    max_return_size: 500B
    return_values:
      file_created: true
      size: true
      content: true
```
