package testrequest

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type (
	// A RequestBuilder interface defines methods for constructing a test http.Request.
	RequestBuilder interface {
		// SetHeader sets the request's target.
		SetTarget(target string) RequestBuilder
		// SetMethod sets the request's HTTP method.
		SetMethod(method string) RequestBuilder
		// SetQuery sets the query for the request.
		SetQuery(q url.Values) RequestBuilder
		// SetQueryValue sets the query parameter for the request.
		SetQueryValue(key string, value ...string) RequestBuilder
		// SetHeader sets the header for the request.
		SetHeader(key string, value ...string) RequestBuilder
		// SetContentType sets the request's Content-Type header.
		//
		// Example of a value: application/json; charset=UTF-8.
		//
		// See RFC 7231, Section 3.1.1.5.
		SetContentType(value string) RequestBuilder
		// SetContentTypeWithParam sets the request's Content-Type header with parameter.
		//SetContentTypeWithParam(mimeType, param string) RequestBuilder
		// SetAccept sets the request's Accept header.
		SetAccept(value string) RequestBuilder
		// SetAcceptLanguage sets the request's Accept-Language header.
		SetAcceptLanguage(value string) RequestBuilder
		// SetCookies sets the request's cookie.
		SetCookies(cookie ...*http.Cookie) RequestBuilder
		// SetContext sets the request's context.
		SetContext(context context.Context) RequestBuilder
		// SetContextValue sets the request's context value.
		SetContextValue(key, value interface{}) RequestBuilder
		// SetBody sets the request's body.
		SetBody(reader io.Reader) RequestBuilder
		// SetPostForm sets the request's body as PostForm.
		// If PostForm is nil, it is initialized. The method is set as POST.
		// Content type as application/x-www-form-urlencoded.
		SetPostForm(postForm url.Values) RequestBuilder
		// SetPostFormValue sets field for the request's PostForm.
		// If PostForm is nil, it is initialized. The method is set as POST.
		// Content type is set as application/x-www-form-urlencoded.
		SetPostFormValue(key string, value ...string) RequestBuilder
		// SetJSON sets JSON-encoded data to the request body.
		//
		// If HTTP method is not set or GET or DELETE, the value is set as POST.
		//
		// If Content-Type header is not set, then the value is set as application/json;charset=UTF8.
		SetJSON(data []byte) RequestBuilder
		// SetJSONFromValue converts the value to JSON encoding
		// and sets data to the request body. An error encoding the value will cause a panic.
		//
		// If HTTP method is not set or GET or DELETE, the value is set as POST.
		//
		// If Content-Type header is not set, then the value is set as application/json;charset=UTF8.
		SetJSONFromValue(v interface{}) RequestBuilder
		// SetAuth sets the request's Authorization header.
		// Prefix specifies the authentication scheme.
		SetAuth(prefix, value string) RequestBuilder
		// SetBasicAuth sets Authorization header to use HTTP
		// Basic Authentication with the provided username and password.
		// See RFC 2617, Section 2.
		SetBasicAuth(username, password string) RequestBuilder
		// SetBearerAuth sets the request's Authorization header to use HTTP Bearer Authentication.
		// See RFC 6750, bearer tokens to access OAuth 2.0-protected resources.
		SetBearerAuth(token string) RequestBuilder
		// SetUserAgent sets User-Agent header.
		SetUserAgent(value string) RequestBuilder
		// Request constructs and returns a new incoming server http.Request for testing.
		//
		// If PostForm is initialized, the request body will be set as strings.Reader of its values.
		// Other method SetBody calls will be ignored.
		Request() *http.Request
	}
	requestBuilder struct {
		target   string
		method   string
		headers  map[string][]string
		query    url.Values
		body     io.Reader
		postForm url.Values
		context  context.Context
		cookies  []*http.Cookie
	}
)

var (
	defaultTarget = "https://server.test"
)

//SetDefaultTarget sets default target for RequestBuilder
func SetDefaultTarget(target string) {
	defaultTarget = target
}

// Builder returns a new RequestBuilder for constructing a test http.Request.
// By default HTTP method is set as GET.
func Builder() RequestBuilder {
	return &requestBuilder{
		target:  defaultTarget,
		method:  http.MethodGet,
		headers: map[string][]string{},
		query:   url.Values{},
		cookies: []*http.Cookie{},
	}
}

func (b *requestBuilder) SetTarget(target string) RequestBuilder {
	b.target = target
	return b
}

func (b *requestBuilder) SetMethod(method string) RequestBuilder {
	b.method = method
	return b
}

func (b *requestBuilder) SetQuery(q url.Values) RequestBuilder {
	b.query = q
	return b
}

func (b *requestBuilder) SetQueryValue(key string, value ...string) RequestBuilder {
	b.query[key] = value
	return b
}

func (b *requestBuilder) SetHeader(key string, value ...string) RequestBuilder {
	b.headers[key] = value
	return b
}

func (b *requestBuilder) SetContentType(value string) RequestBuilder {
	return b.SetHeader("Content-Type", value)
}

func (b *requestBuilder) SetAccept(value string) RequestBuilder {
	return b.SetHeader("Accept", value)
}

func (b *requestBuilder) SetAcceptLanguage(value string) RequestBuilder {
	return b.SetHeader("Accept-Language", value)
}

func (b *requestBuilder) SetCookies(cookie ...*http.Cookie) RequestBuilder {
	b.cookies = cookie
	return b
}

func (b *requestBuilder) SetContext(context context.Context) RequestBuilder {
	b.context = context
	return b
}

func (b *requestBuilder) SetContextValue(key, value interface{}) RequestBuilder {
	if b.context == nil {
		b.context = context.Background()
	}
	b.context = context.WithValue(b.context, key, value)
	return b
}

func (b *requestBuilder) SetBody(reader io.Reader) RequestBuilder {
	b.body = reader
	return b
}

func (b *requestBuilder) SetPostForm(postForm url.Values) RequestBuilder {
	b.postForm = postForm
	_, ok := b.headers["Content-Type"]
	if !ok {
		v := "application/x-www-form-urlencoded;charset=UTF-8"
		b.SetContentType(v)
	}
	return b.SetMethod(http.MethodPost)
}

func (b *requestBuilder) SetPostFormValue(key string, value ...string) RequestBuilder {
	if b.postForm == nil {
		b.SetPostForm(url.Values{})
	}
	b.postForm[key] = value
	return b
}

func (b *requestBuilder) SetJSON(data []byte) RequestBuilder {
	if b.method == http.MethodGet || b.method == http.MethodDelete {
		b.method = http.MethodPost
	}
	_, ok := b.headers["Content-Type"]
	if !ok {
		v := "application/json;charset=UTF-8"
		b.SetContentType(v)
	}
	return b.SetBody(bytes.NewReader(data))
}

func (b *requestBuilder) SetJSONFromValue(v interface{}) RequestBuilder {
	bts, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b.SetJSON(bts)
}

func (b *requestBuilder) SetAuth(prefix, value string) RequestBuilder {
	auth := strings.Join([]string{prefix, value}, " ")
	return b.SetHeader("Authorization", auth)
}

func (b *requestBuilder) SetBasicAuth(username, password string) RequestBuilder {
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return b.SetAuth("basic", auth)
}

func (b *requestBuilder) SetBearerAuth(token string) RequestBuilder {
	return b.SetAuth("bearer", token)
}

func (b *requestBuilder) SetUserAgent(value string) RequestBuilder {
	return b.SetHeader("User-Agent", value)
}

func (b *requestBuilder) Request() *http.Request {
	body := b.body
	if b.postForm != nil {
		body = strings.NewReader(b.postForm.Encode())
	}
	req := httptest.NewRequest(b.method, b.target, body)
	if b.query != nil {
		req.URL.RawQuery = b.query.Encode()
	}
	for key, values := range b.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	for _, c := range b.cookies {
		req.AddCookie(c)
	}
	if b.context != nil {
		return req.WithContext(b.context)
	}
	return req
}
