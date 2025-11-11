package establishments

import (
	getInCoordinateRange "api-mobile-app/src/api/establishments/get-in-coordinate-range/orchestrator"
	"net/http"
)

func EstablishmentsRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/get-in-coordinate-range", func(w http.ResponseWriter, r *http.Request) {
		err := getInCoordinateRange.Orchestrate(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return router
}
