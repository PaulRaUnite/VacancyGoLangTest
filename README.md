# Go Lang vacancy test

This is my realization of the task below for Go programmer vacancy at [Geeks.Team](https://geeks.team).

## About current realization

Requires: `Golang 1.8+`

## Problem

Using Gin framework <https://github.com/gin-gonic/gin> create a web server with a handler `/checkText`.
Handler will listen for `POST` request with such `JSON`:
```json
{
   "Site":["https://google.com","https://yahoo.com"],
   "SearchText":"Google"
}
```

Request structure:
```go
type Request struct {
    Site []string // Slice of strings: https://blog.golang.org/go-slices-usage-and-internals
    SearchText string
}
```

After getting request handler must get the body content of each website mentioned in `Site` variable (this is list of urls) and search in it for a `SearchText`. You can use default Go http client to get the websites body content.
* If the requested `SearchText` was found on the page of any `Site` url, webserver must return `JSON` with the url of the site at which text was found.

Response example:
```json
{
    "FoundAtSite":"https://google.com"
}
```

Response structure:
```go
type Response struct {
    FoundAtSite string
}
```

* If text was not found return `HTTP Code 204 No Content`.