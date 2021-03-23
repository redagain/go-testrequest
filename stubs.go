package testrequest

import "net/http"

type nopResponseWriter struct {
	header http.Header
}

//NopResponseWriter returns a http.ResponseWriter with a no-op Header, Write and WriteHeader methods.
//This is a simple stub for testing.
func NopResponseWriter() http.ResponseWriter {
	return &nopResponseWriter{header: http.Header{}}
}

func (w *nopResponseWriter) Header() http.Header { return w.header }

func (w *nopResponseWriter) Write([]byte) (int, error) { return 0, nil }

func (w *nopResponseWriter) WriteHeader(int) {}
