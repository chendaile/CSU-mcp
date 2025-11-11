package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"

	"campusapp/internal/api/app"
	"campusapp/internal/api/controllers"
	_ "campusapp/internal/api/routers"
	"campusapp/internal/mcp/httpserver"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	root, err := findProjectRoot(filepath.Dir(file))
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(root); err != nil {
		panic(err)
	}
	if err := app.Initialize(); err != nil {
		panic(err)
	}
	beego.ErrorController(&controllers.ErrorController{})
}

func TestAPIHomePage(t *testing.T) {
	Convey("GET / renders the branded API landing page", t, func() {
		resp := performAPIRequest(http.MethodGet, "/")
		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Body.Len(), ShouldBeGreaterThan, 0)
		So(resp.Header().Get("Content-Type"), ShouldContainSubstring, "text/html")
	})
}

func TestAPIUnknownRouteReturns404(t *testing.T) {
	Convey("Unknown routes should respond with a JSON 404 payload", t, func() {
		resp := performAPIRequestWithHeaders(http.MethodGet, "/does-not-exist", map[string]string{
			"Accept": "application/json",
		})
		So(resp.Code, ShouldEqual, http.StatusNotFound)
		So(resp.Header().Get("Content-Type"), ShouldContainSubstring, "application/json")

		var payload struct {
			StateCode int    `json:"StateCode"`
			Error     string `json:"Error"`
		}
		So(json.Unmarshal(resp.Body.Bytes(), &payload), ShouldBeNil)
		So(payload.StateCode, ShouldEqual, 404)
		So(payload.Error, ShouldEqual, "api not found")
	})
}

func TestMCPLandingPage(t *testing.T) {
	Convey("MCP proxy landing page should reuse the shared template", t, func() {
		handler := newMCPTestHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Host = "mcp.local:13000"
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(resp.Header().Get("Content-Type"), ShouldContainSubstring, "text/html")
		So(resp.Body.String(), ShouldContainSubstring, "MCP 服务")
	})
}

func TestMCPDelegatesNonRootRequests(t *testing.T) {
	Convey("Non-root requests should pass through to the MCP handler", t, func() {
		var invoked bool
		handler := newMCPTestHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			invoked = true
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte("ok"))
		}))

		req := httptest.NewRequest(http.MethodPost, "/mcp/stream", strings.NewReader("{}"))
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		So(invoked, ShouldBeTrue)
		So(resp.Code, ShouldEqual, http.StatusAccepted)
		So(resp.Body.String(), ShouldEqual, "ok")
	})
}

func performAPIRequest(method, path string) *httptest.ResponseRecorder {
	return performAPIRequestWithHeaders(method, path, nil)
}

func performAPIRequestWithHeaders(method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(resp, req)
	return resp
}

func newMCPTestHandler(next http.Handler) http.Handler {
	return httpserver.LoggingMiddleware(
		httpserver.LandingMiddleware(
			next,
			"CSU MCP Proxy",
			"http://localhost:12000",
			":13000",
		),
	)
}

func findProjectRoot(start string) (string, error) {
	dir := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", start)
		}
		dir = parent
	}
}
