package server

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"

    h "go_blog/internal/transport/http"
    "go_blog/internal/repository"
    "go_blog/internal/service"
    user "go_blog/internal/domain/user"
    post "go_blog/internal/domain/post"
    comment "go_blog/internal/domain/comment"
)

type memUserRepo struct{ seq uint; byID map[uint]*user.User; byEmail map[string]*user.User }
func (m *memUserRepo) Create(u *user.User) error { m.seq++; u.ID=m.seq; m.byID[u.ID]=u; m.byEmail[u.Email]=u; return nil }
func (m *memUserRepo) GetByEmail(email string) (*user.User, error) { if u:=m.byEmail[email]; u!=nil { return u,nil }; return nil, gormErr }
func (m *memUserRepo) GetByID(id uint) (*user.User, error) { if u:=m.byID[id]; u!=nil { return u,nil }; return nil, gormErr }

type memPostRepo struct{ seq uint; byID map[uint]*post.Post; list []*post.Post }
func (m *memPostRepo) Create(p *post.Post) error { m.seq++; p.ID=m.seq; m.byID[p.ID]=p; m.list=append(m.list,p); return nil }
func (m *memPostRepo) Update(p *post.Post) error { m.byID[p.ID]=p; return nil }
func (m *memPostRepo) Delete(id uint) error { delete(m.byID,id); return nil }
func (m *memPostRepo) FindByID(id uint) (*post.Post, error) { if p:=m.byID[id]; p!=nil { return p,nil }; return nil, gormErr }
func (m *memPostRepo) List(filter repository.PostListFilter, page, size int) ([]post.Post, int64, error) {
    var out []post.Post
    for _, p := range m.list { out = append(out, *p) }
    return out, int64(len(out)), nil
}
func (m *memPostRepo) IncrementViews(id uint) error { if p:=m.byID[id]; p!=nil { p.Views++; return nil }; return gormErr }

type memCommentRepo struct{ seq uint; byID map[uint]*comment.Comment; byPost map[uint][]comment.Comment }
func (m *memCommentRepo) Create(c *comment.Comment) error { m.seq++; c.ID=m.seq; m.byID[c.ID]=c; m.byPost[c.PostID]=append(m.byPost[c.PostID], *c); return nil }
func (m *memCommentRepo) Update(c *comment.Comment) error { m.byID[c.ID]=c; return nil }
func (m *memCommentRepo) Delete(id uint) error { delete(m.byID,id); return nil }
func (m *memCommentRepo) FindByID(id uint) (*comment.Comment, error) { if c:=m.byID[id]; c!=nil { return c,nil }; return nil, gormErr }
func (m *memCommentRepo) ListByPost(postID uint) ([]comment.Comment, error) { return m.byPost[postID], nil }

var gormErr = errNotFound{}
type errNotFound struct{}
func (e errNotFound) Error() string { return "not found" }

func setupTestRouter(t *testing.T) *gin.Engine {
    gin.SetMode(gin.TestMode)
    userRepo := &memUserRepo{byID: map[uint]*user.User{}, byEmail: map[string]*user.User{}}
    userSvc := service.NewUserService(userRepo, []byte("testsecret"), time.Hour)
    postRepo := &memPostRepo{byID: map[uint]*post.Post{}}
    postSvc := service.NewPostService(postRepo)
    commentRepo := &memCommentRepo{byID: map[uint]*comment.Comment{}, byPost: map[uint][]comment.Comment{}}
    commentSvc := service.NewCommentService(commentRepo)

    r := gin.New()
    api := r.Group("/api/v1")
    authH := NewAuthHandler(userSvc)
    api.POST("/auth/register", authH.Register)
    api.POST("/auth/login", authH.Login)
    authMW := NewAuthMiddleware([]byte("testsecret"))
    apiAuth := api.Group("")
    apiAuth.Use(authMW)
    apiAuth.GET("/auth/me", authH.Me)

    postH := NewPostHandler(postSvc)
    apiAuth.POST("/posts", postH.Create)
    api.GET("/posts", postH.List)
    api.GET("/posts/:id", postH.Get)

    commentH := NewCommentHandler(commentSvc)
    apiAuth.POST("/posts/:id/comments", commentH.Create)
    api.GET("/posts/:id/comments", commentH.ListTree)
    return r
}

func TestAuthPostCommentFlow(t *testing.T) {
    r := setupTestRouter(t)

    rr := httptest.NewRecorder()
    reqBody := h.RegisterRequest{Username: "u1", Email: "u1@example.com", Password: "secret"}
    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("register failed: %d", rr.Code) }

    rr = httptest.NewRecorder()
    lb := h.LoginRequest{Email: "u1@example.com", Password: "secret"}
    b, _ = json.Marshal(lb)
    req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("login failed: %d", rr.Code) }
    var loginResp map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &loginResp)
    data := loginResp["data"].(map[string]interface{})
    token := data["token"].(string)

    rr = httptest.NewRecorder()
    pb := h.PostCreateRequest{Title: "Hello World", Content: "content", Status: "draft"}
    b, _ = json.Marshal(pb)
    req = httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("create post failed: %d", rr.Code) }

    rr = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&page_size=10", nil)
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("list posts failed: %d", rr.Code) }

    rr = httptest.NewRecorder()
    cb := h.CommentCreateRequest{Content: "nice"}
    b, _ = json.Marshal(cb)
    req = httptest.NewRequest(http.MethodPost, "/api/v1/posts/1/comments", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("create comment failed: %d", rr.Code) }

    rr = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodGet, "/api/v1/posts/1/comments", nil)
    r.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("list comments failed: %d", rr.Code) }
}