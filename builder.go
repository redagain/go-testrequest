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

// Returns charset=UTF-8
const CharsetUTF8 = "charset=UTF-8"

var (
	defaultTarget = "https://server.test"
)

//SetDefaultTarget sets default target for Builder
func SetDefaultTarget(target string) {
	defaultTarget = target
}

type (
	builder struct {
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

// Builder returns a new Builder for constructing a test http.Request, HTTP method is set as GET.
func Builder() *builder {
	return &builder{
		target:  defaultTarget,
		method:  http.MethodGet,
		headers: map[string][]string{},
		query:   url.Values{},
		cookies: []*http.Cookie{},
	}
}

// SetHeader sets the request's target.
func (b *builder) SetTarget(target string) *builder {
	b.target = target
	return b
}

// SetMethod sets the request's HTTP method.
func (b *builder) SetMethod(method string) *builder {
	b.method = method
	return b
}

// SetQuery sets the query for the request.
func (b *builder) SetQuery(q url.Values) *builder {
	b.query = q
	return b
}

// SetQueryValue sets the query parameter for the request.
func (b *builder) SetQueryValue(key string, value ...string) *builder {
	b.query[key] = value
	return b
}

// SetHeader sets the header for the request.
func (b *builder) SetHeader(key string, value ...string) *builder {
	b.headers[key] = value
	return b
}

// SetContentTypeWithParam sets the request's Content-Type header with parameter.
func (b *builder) SetContentTypeWithParam(mimeType, param string) *builder {
	return b.SetContentType(strings.Join([]string{mimeType, param}, ";"))
}

// SetContentType sets the request's Content-Type header.
//
// Example of a value: application/json; charset=UTF-8.
//
// See RFC 7231, Section 3.1.1.5.
func (b *builder) SetContentType(value string) *builder {
	return b.SetHeader("Content-Type", value)
}

// SetAccept sets the request's Accept header.
func (b *builder) SetAccept(value string) *builder {
	return b.SetHeader("Accept", value)
}

// SetAcceptLanguage sets the request's Accept-Language header.
func (b *builder) SetAcceptLanguage(value string) *builder {
	return b.SetHeader("Accept-Language", value)
}

// SetCookie sets the request's cookie.
func (b *builder) SetCookie(cookie ...*http.Cookie) *builder {
	b.cookies = cookie
	return b
}

// SetContext sets the request's context.
func (b *builder) SetContext(context context.Context) *builder {
	b.context = context
	return b
}

// SetContextValue sets the request's context value.
func (b *builder) SetContextValue(key, value interface{}) *builder {
	if b.context == nil {
		b.context = context.Background()
	}
	b.context = context.WithValue(b.context, key, value)
	return b
}

// SetBody sets the request's body.
func (b *builder) SetBody(reader io.Reader) *builder {
	b.body = reader
	return b
}

// SetPostForm sets the request's body as PostForm.
// If PostForm is nil, it is initialized. The method is set as POST.
// Content type as application/x-www-form-urlencoded.
func (b *builder) SetPostForm(postForm url.Values) *builder {
	b.postForm = postForm
	return b.SetMethod(http.MethodPost).SetContentTypeWithParam("application/x-www-form-urlencoded", CharsetUTF8)
}

// SetPostFormValue sets field for the request's PostForm.
// If PostForm is nil, it is initialized. The method is set as POST.
// Content type is set as application/x-www-form-urlencoded.
func (b *builder) SetPostFormValue(key string, value ...string) *builder {
	if b.postForm == nil {
		b.SetPostForm(url.Values{})
	}
	b.postForm[key] = value
	return b
}

// SetJSON sets JSON-encoded data to the request body.
//
// If HTTP method is not set or GET or DELETE, the value is set as POST.
//
// If Content-Type header is not set, then the value is set as application/json;charset=UTF8.
func (b *builder) SetJSON(data []byte) *builder {
	if b.method == http.MethodGet || b.method == http.MethodDelete {
		b.method = http.MethodPost
	}
	_, ok := b.headers["Content-Type"]
	if !ok {
		b.SetContentTypeWithParam("application/json", CharsetUTF8)
	}
	return b.SetBody(bytes.NewReader(data))
}

// SetJSONFromValue converts the value to JSON encoding
// and sets data to the request body. An error encoding the value will cause a panic.
//
// If HTTP method is not set or GET or DELETE, the value is set as POST.
//
// If Content-Type header is not set, then the value is set as application/json;charset=UTF8.
func (b *builder) SetJSONFromValue(v interface{}) *builder {
	bts, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b.SetJSON(bts)
}

/*func (b *builder) SetMultipartFormDataFromFile(path string, options map[string]string) *builder {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	for k, v := range options {
		_ = writer.WriteField(k, v)
	}
	contentType := writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	return b.SetMethod(http.MethodPost).SetContentType(contentType).SetBody(buffer)
}*/

// SetAuth sets the request's Authorization header.
// Prefix specifies the authentication scheme.
func (b *builder) SetAuth(prefix, value string) *builder {
	auth := strings.Join([]string{prefix, value}, " ")
	return b.SetHeader("Authorization", auth)
}

// SetBasicAuth sets Authorization header to use HTTP
// Basic Authentication with the provided username and password.
// See RFC 2617, Section 2.
func (b *builder) SetBasicAuth(username, password string) *builder {
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return b.SetAuth("basic", auth)
}

// SetBearerAuth sets the request's Authorization header to use HTTP Bearer Authentication.
// See RFC 6750, bearer tokens to access OAuth 2.0-protected resources.
func (b *builder) SetBearerAuth(token string) *builder {
	return b.SetAuth("bearer", token)
}

// SetUserAgent sets User-Agent header.
func (b *builder) SetUserAgent(value string) *builder {
	return b.SetHeader("User-Agent", value)
}

// Request constructs and returns a new incoming server http.Request for testing.
//
// If PostForm is initialized, the request body will be set as strings.Reader of its values.
// Other method SetBody calls will be ignored.
func (b *builder) Request() *http.Request {
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
