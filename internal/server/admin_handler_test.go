package server

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

type fakeAdminSvc struct{}
func (f fakeAdminSvc) Stats() (int64, int64, int64, error) { return 1, 2, 3, nil }

func TestAdminStats(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    ah := NewAdminHandler(fakeAdminSvc{})
    // stub auth: set role=admin
    r.GET("/api/v1/admin/stats", func(c *gin.Context) { c.Set("role", "admin") }, AdminOnly(), ah.Stats)

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/stats", nil)
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
}