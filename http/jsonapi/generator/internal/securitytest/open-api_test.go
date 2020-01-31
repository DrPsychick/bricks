// Code generated by github.com/pace/bricks DO NOT EDIT.
package securitytest

import (
	"context"
	mux "github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	oauth2 "github.com/pace/bricks/http/oauth2"
	apikey "github.com/pace/bricks/http/security/apikey"
	errors "github.com/pace/bricks/maintenance/errors"
	metrics "github.com/pace/bricks/maintenance/metric/jsonapi"
	"net/http"
)

type AuthorizationBackend interface {
	AuthorizeOAuth2(r *http.Request, w http.ResponseWriter, scope string) (context.Context, bool)
	CanAuthorizeOAuth2(r *http.Request) bool
	AuthorizeProfileKey(r *http.Request, w http.ResponseWriter) (context.Context, bool)
	CanAuthorizeProfileKey(r *http.Request) bool
	Init(cfgOAuth2 *oauth2.Config, cfgProfileKey *apikey.Config)
}

var cfgOAuth2 = &oauth2.Config{
	AuthorizationCode: &oauth2.Flow{
		AuthorizationURL: "https://id.pace.cloud/oauth2/authorize",
		RefreshURL:       "https://id.pace.cloud/oauth2/token",
		Scopes:           map[string]string{"anything": "test"},
		TokenURL:         "https://id.pace.cloud/oauth2/token",
	},
	Description: "",
}
var cfgProfileKey = &apikey.Config{
	Description: "prefix with \"Bearer \"",
	In:          "header",
	Name:        "Authorization",
}

/*
GetTestHandler handles request/response marshaling and validation for
 Get /beta/test
*/
func GetTestHandler(service Service, authBackend AuthorizationBackend) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer errors.HandleRequest("GetTestHandler", w, r)

		var ctx context.Context
		var ok bool
		if authBackend.CanAuthorizeOAuth2(r) {

			ctx, ok = authBackend.AuthorizeOAuth2(r, w, "anything")
			if !ok {
				return
			}
		} else if authBackend.CanAuthorizeProfileKey(r) {

			ctx, ok = authBackend.AuthorizeProfileKey(r, w)
			if !ok {
				return
			}
		} else {
			http.Error(w, "Authorization Error", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(ctx)

		// Trace the service function handler execution
		handlerSpan, ctx := opentracing.StartSpanFromContext(r.Context(), "GetTestHandler")
		defer handlerSpan.Finish()

		// Setup context, response writer and request type
		writer := getTestResponseWriter{
			ResponseWriter: metrics.NewMetric("securitytest", "/beta/test", w, r),
		}
		request := GetTestRequest{
			Request: r.WithContext(ctx),
		}

		// Scan and validate incoming request parameters

		// Invoke service that implements the business logic
		err := service.GetTest(ctx, &writer, &request)
		if err != nil {
			errors.HandleError(err, "GetTestHandler", w, r)
		}
	})
}

/*
GetTestResponseWriter is a standard http.ResponseWriter extended with methods
to generate the respective responses easily
*/
type GetTestResponseWriter interface {
	http.ResponseWriter
	OK()
}
type getTestResponseWriter struct {
	http.ResponseWriter
}

// OK responds with empty response (HTTP code 200)
func (w *getTestResponseWriter) OK() {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(200)
}

/*
GetTestRequest is a standard http.Request extended with the
un-marshaled content object
*/
type GetTestRequest struct {
	Request *http.Request `valid:"-"`
}

// Service interface for all handlers
type Service interface {
	// GetTest Test
	GetTest(context.Context, GetTestResponseWriter, *GetTestRequest) error
}

/*
Router implements: PACE Payment API

Welcome to the PACE Payment API documentation.
This API is responsible for managing payment methods for users as well as authorizing payments on behalf of PACE services.
*/
func Router(service Service, authBackend AuthorizationBackend) *mux.Router {
	authBackend.Init(cfgOAuth2, cfgProfileKey)
	router := mux.NewRouter()
	// Subrouter s1 - Path: /pay
	s1 := router.PathPrefix("/pay").Subrouter()
	s1.Methods("GET").Path("/beta/test").Handler(GetTestHandler(service, authBackend)).Name("GetTest")
	return router
}