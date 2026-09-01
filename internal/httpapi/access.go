package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/auth"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const sessionCookieName = "stockpilot_session"

type principalContextKey struct{}

type accessHandler struct {
	next         http.Handler
	auth         *auth.Service
	users        repository.Access
	secureCookie bool
}

func WithAccess(next http.Handler, authService *auth.Service, users repository.Access, secureCookie bool) http.Handler {
	return &accessHandler{next: next, auth: authService, users: users, secureCookie: secureCookie}
}

func (a *accessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/meta" || r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
		a.next.ServeHTTP(w, r)
		return
	}
	if r.URL.Path == "/api/v1/auth/login" && r.Method == http.MethodPost {
		a.login(w, r)
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
		return
	}
	principal, err := a.auth.Resolve(r.Context(), cookie.Value)
	if err != nil {
		a.clearCookie(w)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), principalContextKey{}, principal))

	if isMutation(r.Method) && r.Header.Get("X-StockPilot-CSRF") != "1" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "CSRF confirmation header is required"})
		return
	}

	switch {
	case r.URL.Path == "/api/v1/auth/logout" && r.Method == http.MethodPost:
		a.logout(w, r, cookie.Value)
		return
	case r.URL.Path == "/api/v1/auth/me" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, principal.User)
		return
	case strings.HasPrefix(r.URL.Path, "/api/v1/users"):
		if !principal.User.Role.Can(domain.PermissionUsersManage) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "insufficient permission"})
			return
		}
		a.usersAPI(w, r, principal)
		return
	}

	permission, ok := permissionFor(r)
	if !ok || !principal.User.Role.Can(permission) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "insufficient permission"})
		return
	}
	a.next.ServeHTTP(w, r)
}

func (a *accessHandler) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	result, err := a.auth.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}
	maxAge := int(time.Until(result.ExpiresAt).Seconds())
	if maxAge < 1 {
		maxAge = 1
	}
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieName, Value: result.Token, Path: "/", HttpOnly: true, Secure: a.secureCookie,
		SameSite: http.SameSiteStrictMode, Expires: result.ExpiresAt, MaxAge: maxAge,
	})
	writeJSON(w, http.StatusOK, map[string]any{"user": result.User, "expiresAt": result.ExpiresAt})
}

func (a *accessHandler) logout(w http.ResponseWriter, r *http.Request, rawToken string) {
	if err := a.auth.Logout(r.Context(), rawToken); err != nil {
		writeDomainError(w, err)
		return
	}
	a.clearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (a *accessHandler) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: a.secureCookie,
		SameSite: http.SameSiteStrictMode, MaxAge: -1, Expires: time.Unix(1, 0),
	})
}

func (a *accessHandler) usersAPI(w http.ResponseWriter, r *http.Request, principal auth.Principal) {
	switch {
	case r.URL.Path == "/api/v1/users" && r.Method == http.MethodGet:
		users, err := a.users.ListUsers(r.Context())
		if err != nil {
			writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": users})
	case r.URL.Path == "/api/v1/users" && r.Method == http.MethodPost:
		var body struct {
			Email       string      `json:"email"`
			DisplayName string      `json:"displayName"`
			Role        domain.Role `json:"role"`
			Password    string      `json:"password"`
		}
		if !decodeJSON(w, r, &body) {
			return
		}
		user, err := a.auth.CreateUser(r.Context(), body.Email, body.DisplayName, body.Role, body.Password)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, user)
	case strings.HasSuffix(r.URL.Path, "/role") && r.Method == http.MethodPut:
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/users/"), "/role")
		var body struct {
			Role domain.Role `json:"role"`
		}
		if !decodeJSON(w, r, &body) {
			return
		}
		if id == principal.User.ID && body.Role != domain.RoleAdmin {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "administrators cannot remove their own administrator role"})
			return
		}
		if err := a.users.UpdateUserRole(r.Context(), id, body.Role); err != nil {
			writeDomainError(w, err)
			return
		}
		user, err := a.users.GetUser(r.Context(), id)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, user)
	case strings.HasSuffix(r.URL.Path, "/active") && r.Method == http.MethodPut:
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/users/"), "/active")
		var body struct {
			Active bool `json:"active"`
		}
		if !decodeJSON(w, r, &body) {
			return
		}
		if id == principal.User.ID && !body.Active {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "administrators cannot deactivate their own account"})
			return
		}
		if err := a.users.SetUserActive(r.Context(), id, body.Active); err != nil {
			writeDomainError(w, err)
			return
		}
		user, err := a.users.GetUser(r.Context(), id)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, user)
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func permissionFor(r *http.Request) (domain.Permission, bool) {
	path := r.URL.Path
	read := r.Method == http.MethodGet || r.Method == http.MethodHead
	switch {
	case strings.HasPrefix(path, "/api/v1/categories"), strings.HasPrefix(path, "/api/v1/suppliers"), strings.HasPrefix(path, "/api/v1/products"):
		if read {
			return domain.PermissionCatalogRead, true
		}
		return domain.PermissionCatalogWrite, true
	case strings.HasPrefix(path, "/api/v1/warehouses"), strings.HasPrefix(path, "/api/v1/locations"), strings.HasPrefix(path, "/api/v1/lots"), strings.HasPrefix(path, "/api/v1/inventory"):
		if read {
			return domain.PermissionInventoryRead, true
		}
		return domain.PermissionInventoryWrite, true
	case strings.HasPrefix(path, "/api/v1/orders"):
		if read {
			return domain.PermissionOrdersRead, true
		}
		return domain.PermissionOrdersWrite, true
	case strings.HasPrefix(path, "/api/v1/reports"):
		return domain.PermissionReportsRead, read
	case strings.HasPrefix(path, "/api/v1/audit"):
		return domain.PermissionAuditRead, read
	default:
		return "", false
	}
}

func isMutation(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func PrincipalFromContext(ctx context.Context) (auth.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(auth.Principal)
	return principal, ok
}
