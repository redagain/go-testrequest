package testrequest

import "net/http"

type nopResponseWriter struct{}

//NopResponseWriter returns a http.ResponseWriter with a no-op Header, Write and WriteHeader methods.
//This is a simple stub for testing.
func NopResponseWriter() http.ResponseWriter { return &nopResponseWriter{} }

func (w *nopResponseWriter) Header() http.Header { return http.Header{} }

func (w *nopResponseWriter) Write([]byte) (int, error) { return 0, nil }

func (w *nopResponseWriter) WriteHeader(int) {}
