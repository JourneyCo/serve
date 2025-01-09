package router

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
	"net/http"
	"serve/api/v1/projects"
	"serve/app/auth0"
	"serve/app/errors"
	"serve/app/middleware"
)

func Route(domain string, audience string) *mux.Router {

	r := mux.NewRouter()
	secureMiddleware := secure.New()
	r.Use(secureMiddleware.Handler)

	//r.HandleFunc("/", home).Methods("GET")
	//r.HandleFunc("/contact", contact).Methods("GET")
	//r.HandleFunc("/api/widgets", apiGetWidgets).Methods("GET")
	//r.HandleFunc("/api/widgets", apiCreateWidget).Methods("POST")
	//r.HandleFunc("/api/widgets/{slug}", apiUpdateWidget).Methods("POST")
	//r.HandleFunc("/api/widgets/{slug}/parts", apiCreateWidgetPart).Methods("POST")
	//r.HandleFunc("/api/widgets/{slug}/parts/{id:[0-9]+}/update", apiUpdateWidgetPart).Methods("POST")
	//r.HandleFunc("/api/widgets/{slug}/parts/{id:[0-9]+}/delete", apiDeleteWidgetPart).Methods("POST")
	//r.HandleFunc("/{slug}", widgetGet).Methods("GET")
	//r.HandleFunc("/{slug}/admin", widgetAdmin).Methods("GET")
	//r.HandleFunc("/{slug}/image", widgetImage).Methods("POST")

	r.HandleFunc("/", errors.NotFoundHandler)
	r.HandleFunc("/api/messages/public", projects.PublicApiHandler)
	r.Handle("/api/messages/protected", auth0.ValidateJWT(audience, domain, http.HandlerFunc(projects.ProtectedApiHandler)))
	r.Handle("/api/messages/admin", auth0.ValidateJWT(audience, domain, http.HandlerFunc(projects.AdminApiHandler)))

	r.Use(middleware.HandleCacheControl)
	return r
}
