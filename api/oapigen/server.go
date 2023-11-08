// Package oapigen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package oapigen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a blob
	// (POST /v1/blobs)
	CreateBlob(w http.ResponseWriter, r *http.Request)
	// Retrieve a blob
	// (GET /v1/blobs/{id})
	GetBlobByID(w http.ResponseWriter, r *http.Request, id string)
	// Retrieve health status
	// (GET /v1/health)
	GetHealth(w http.ResponseWriter, r *http.Request)
	// Retrieve all prompts
	// (GET /v1/prompts)
	GetPrompts(w http.ResponseWriter, r *http.Request, params GetPromptsParams)
	// Create a prompt
	// (POST /v1/prompts)
	CreatePrompt(w http.ResponseWriter, r *http.Request)
	// Retrieve a prompt
	// (GET /v1/prompts/{id})
	GetPromptByID(w http.ResponseWriter, r *http.Request, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// CreateBlob operation middleware
func (siw *ServerInterfaceWrapper) CreateBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateBlob(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetBlobByID operation middleware
func (siw *ServerInterfaceWrapper) GetBlobByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBlobByID(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealth(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetPrompts operation middleware
func (siw *ServerInterfaceWrapper) GetPrompts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPromptsParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPrompts(w, r, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreatePrompt operation middleware
func (siw *ServerInterfaceWrapper) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreatePrompt(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetPromptByID operation middleware
func (siw *ServerInterfaceWrapper) GetPromptByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPromptByID(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/blobs", wrapper.CreateBlob)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/blobs/{id}", wrapper.GetBlobByID)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/health", wrapper.GetHealth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/prompts", wrapper.GetPrompts)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/prompts", wrapper.CreatePrompt)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/prompts/{id}", wrapper.GetPromptByID)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}