![Image alt](https://github.com/redagain/go-testrequest/raw/master/go-testrequest.png)
[![License](https://img.shields.io/github/license/redagain/go-testrequest)](LICENSE)
[![Build Status](https://travis-ci.com/redagain/go-testrequest.svg?branch=master)](https://travis-ci.com/redagain/go-testrequest)
[![Go Reference](https://pkg.go.dev/badge/github.com/redagain/go-testrequest.svg)](https://pkg.go.dev/github.com/redagain/go-testrequest)

Simple helpers for testing http-handlers and other handlers

## Install

```bash
go get github.com/redagain/go-testrequest
```

## Helpers

* Test request builder
* No-op http.ResponseWriter

## Usage example

```go
func TestHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				w: testrequest.NopResponseWriter(),
				r: testrequest.Builder().SetJSONFromValue(
					map[string]interface{}{
						"title":   "The Go Programming Language",
						"isbn":    "978-0134190440",
						"authors": []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
					}).
					SetAccept("application/json").
					SetBearerAuth("4mwbMA6zq9Nxf3XzLk9n01MJX57jdjMAdfYCaJu44vEUJVfdVLF9").
					Request(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {    
```

## Examples

- [simpleapi](https://github.com/redagain/go-testrequest/tree/master/examples/simpleapi) - examples of using helpers to test a simple API
- [clientcredentials](https://github.com/redagain/go-testrequest/tree/master/examples/clientcredentials) - example of testing a function to extract client credentials from a http.Request
