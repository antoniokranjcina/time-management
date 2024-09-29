package middleware

import (
	"net/http"
	"time-management/internal/shared/util"
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
)

func RoleMiddleware(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*domain.User)
			if !ok {
				_ = util.WriteJson(w, http.StatusForbidden, util.ApiError{Error: "Forbidden"})
				return
			}

			// Check if user is a SuperAdmin, if yes, grant access
			if user.Role == role.SuperAdmin.String() || user.Role == role.Admin.String() {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user's role matches any of the allowed roles
			for _, allowedRole := range allowedRoles {
				if user.Role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			// If no roles matched, return forbidden
			_ = util.WriteJson(w, http.StatusForbidden, util.ApiError{Error: "Forbidden"})
		})
	}
}
