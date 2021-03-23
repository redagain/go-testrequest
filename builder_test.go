package testrequest

import (
	"bytes"
	"context"
	"encoding/json"

	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func Test_InitBuilder(t *testing.T) {
	want := &builder{
		target:  defaultTarget,
		method:  http.MethodGet,
		headers: map[string][]string{},
		query:   url.Values{},
		cookies: []*http.Cookie{},
	}
	if got := Builder(); !reflect.DeepEqual(got, want) {
		t.Errorf("Builder() = %v, want %v", got, want)
	}
}

func Test_builder_SetTarget(t *testing.T) {
	want := "http://example.com"
	req := Builder().SetTarget(want).Request()
	got := req.RequestURI
	if got != want {
		t.Errorf("SetTarget() = %v, wantTarget %v", got, want)
	}
}

func Test_builder_SetMethod(t *testing.T) {
	want := http.MethodPatch
	req := Builder().SetMethod(want).Request()
	got := req.Method
	if got != want {
		t.Errorf("SetMethod() = %v, wantMethod%v", got, want)
	}
}

func Test_builder_SetQuery(t *testing.T) {
	req := Builder().SetQuery(url.Values{"test": {"test"}}).Request()
	want := url.Values{"test": {"test"}}
	got := req.URL.Query()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetQuery() = %v, wantQuery%v", got, want)
		return
	}
}

func Test_builder_SetQueryValue(t *testing.T) {
	key := "test"
	req := Builder().SetQueryValue(key, "test").Request()
	got := req.URL.Query()
	want := url.Values{"test": {"test"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetQueryValue() = %v, wantQuery%v", got, want)
		return
	}
}

func Test_builder_SetQueryValues(t *testing.T) {
	key := "test"
	req := Builder().SetQueryValue(key, "test1", "test2").Request()
	got := req.URL.Query()
	want := url.Values{"test": {"test1", "test2"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetQueryValue() = %v, wantQuery%v", got, want)
		return
	}
}

func Test_builder_SetHeader(t *testing.T) {
	req := Builder().SetHeader("test", "test").Request()
	want := http.Header{}
	want.Set("test", "test")
	got := req.Header
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetHeader() = %v, wantHeader %v", got, want)
	}
}

func Test_builder_SetHeaderValues(t *testing.T) {
	req := Builder().SetHeader("test", "test1", "test2").Request()
	want := http.Header{}
	want.Add("test", "test1")
	want.Add("test", "test2")
	got := req.Header
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetHeader() = %v, wantHeader %v", got, want)
	}
}

func Test_builder_SetContentType(t *testing.T) {
	req := Builder().SetContentType("application/json").Request()
	want := "application/json"
	got := req.Header.Get("Content-Type")
	if got != want {
		t.Errorf("SetContentType() = %v, want %v", got, want)
	}
}

func Test_builder_SetContentTypeWithParam(t *testing.T) {
	req := Builder().SetContentTypeWithParam("application/json", CharsetUTF8).Request()
	want := "application/json;charset=UTF-8"
	got := req.Header.Get("Content-Type")
	if got != want {
		t.Errorf("SetContentType() = %v, want %v", got, want)
	}
}

func Test_builder_SetCookie(t *testing.T) {
	req := Builder().SetCookie(&http.Cookie{Name: "test", Value: "test"}).Request()
	want := []*http.Cookie{{
		Name:  "test",
		Value: "test",
	}}
	got := req.Cookies()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetCookie() = %v, want %v", got, want)
	}
}

func Test_builder_SetCookies(t *testing.T) {
	req := Builder().
		SetCookie(&http.Cookie{Name: "test1", Value: "test1"}, &http.Cookie{Name: "test2", Value: "test2"}).Request()
	want := []*http.Cookie{{
		Name:  "test1",
		Value: "test1",
	},
		{
			Name:  "test2",
			Value: "test2",
		}}
	got := req.Cookies()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetCookie() = %v, want %v", got, want)
	}
}

func Test_builder_SetContext(t *testing.T) {
	req := Builder().SetContext(context.Background()).Request()
	got := req.Context()
	want := context.Background()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetContext() = %v, want %v", got, want)
	}
}

func Test_builder_SetContextValue(t *testing.T) {
	key := "test"
	value := "test"
	req := Builder().SetContextValue(key, value).Request()
	want := context.WithValue(context.Background(), key, value)
	got := req.Context()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetContextValue() = %v, want %v", got, want)
		return
	}
	gotValue, ok := req.Context().Value(key).(string)
	if !ok && gotValue != value {
		t.Errorf("SetContextValue() = %v, wantValue %v", gotValue, value)
	}
}

func Test_builder_SetUserAgent(t *testing.T) {
	req := Builder().SetUserAgent("test").Request()
	want := "test"
	got := req.UserAgent()
	if got != want {
		t.Errorf("SetUserAgent() = %v, want %v", got, want)
	}
}

func Test_builder_SetAuth(t *testing.T) {
	req := Builder().SetAuth("bearer", "test").Request()
	want := "bearer test"
	got := req.Header.Get("Authorization")
	if got != want {
		t.Errorf("SetAuth() = %v, want %v", got, want)
	}
}

func Test_builder_SetBasicAuth(t *testing.T) {
	req := Builder().SetBasicAuth("admin", "p@ssw0rd").Request()
	wantUsername := "admin"
	wantPassword := "p@ssw0rd"
	wantOk := true
	gotUsername, gotPassword, gotOk := req.BasicAuth()
	if gotOk != wantOk {
		t.Errorf("SetBasicAuth() = %v, wantOk %v", gotOk, wantOk)
		return
	}
	if gotUsername != wantUsername {
		t.Errorf("SetBasicAuth() = %v, wantUsername %v", gotUsername, wantUsername)
		return
	}
	if gotPassword != wantPassword {
		t.Errorf("SetBasicAuth() = %v, wantPassword %v", gotPassword, wantPassword)
	}
}

func Test_builder_SetBody(t *testing.T) {
	req := Builder().SetBody(strings.NewReader("test")).Request()
	want := []byte("test")
	got, _ := io.ReadAll(req.Body)
	if !bytes.Equal(got, want) {
		t.Errorf("SetBody() = %v, want %v", got, want)
	}
}

func Test_builder_SetPostForm(t *testing.T) {
	req := Builder().SetPostForm(url.Values{"test": {"test"}}).Request()
	want := url.Values{"test": {"test"}}
	req.ParseForm()
	got := req.PostForm
	if req.Method != http.MethodPost {
		t.Errorf("SetPostForm() wantMethod %v", http.MethodPost)
		return
	}
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		t.Errorf("SetPostForm() wantContentType %v", "application/x-www-form-urlencoded")
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetPostForm() = %v, want %v", got, want)
	}
}

func Test_builder_SetPostFormValue(t *testing.T) {
	req := Builder().SetPostFormValue("test", "test").Request()
	want := url.Values{"test": {"test"}}
	req.ParseForm()
	got := req.PostForm
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetPostFormValue() = %v, want %v", got, want)
	}
}

func Test_builder_SetJSON(t *testing.T) {
	req := Builder().SetJSON([]byte("{}")).Request()
	want := []byte("{}")
	got, _ := io.ReadAll(req.Body)
	if req.Method != http.MethodPost {
		t.Errorf("SetJSON() wantMethod %v", http.MethodPost)
		return
	}
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		t.Errorf("SetJSON() wantContentType %v", "application/json")
		return
	}
	if !bytes.Equal(got, want) {
		t.Errorf("SetJSON() = %v, want %v", got, want)
	}
}

func Test_builder_SetJSONWithPreset(t *testing.T) {
	req := Builder().SetMethod(http.MethodPut).
		SetContentType("application/json").
		SetJSON([]byte("{}")).
		Request()
	if req.Method != http.MethodPut {
		t.Errorf("SetJSON() wantMethod %v", http.MethodPut)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("SetJSON() wantContentType %v", "application/json")
	}
}

func Test_builder_SetJSONFromValue(t *testing.T) {
	type value struct {
		Value string `json:"value"`
	}
	want := value{Value: "test"}
	req := Builder().SetJSONFromValue(want).Request()
	got := value{}
	json.NewDecoder(req.Body).Decode(&got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SetJSONFromValue() = %v, want %v", got, want)
	}
}
