# Readme


## Requirements

* Golang
* gcc (e.g. install from https://github.com/mattn/go-sqlite3/issues/url)

## Installation

### GCC

Follow this: https://github.com/mattn/go-sqlite3/issues/212#issuecomment-273531789

# Building

```
rice embed-go
go build
```
`rice embed-go` generates the `rice-box.go` file containing all the static files (e.g html, css, js). For further information, see 
https://github.com/GeertJohan/go.rice.