package middleware

import (
	"net/http"

	"api-mobile-app/src/api/common/constants"
	errorHandling "api-mobile-app/src/api/common/errorHandling"
	Logging "api-mobile-app/src/api/common/logging"
)

func AllowRoles(handler http.Handler, logger *Logging.Logger, allowed ...string) http.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(constants.UserRoleKey).(string)
		if _, ok := allowedSet[role]; !ok {
			pd := errorHandling.Forbidden("role not permitted for this route", nil)
			logger.LogRequest(r, &pd)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
