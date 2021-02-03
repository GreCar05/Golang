package server

import (
	"net/http"
	ctr "recibe_me/internal/controllers"

	"github.com/gorilla/mux"
)

// Route is a Route type
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// Routes is a array of Route
type Routes []Route

// NewRouter ...
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Pattern).
			Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ctr.Index,
		// mid.Authenticate(ctr.Index),
	},
	Route{
		"Issues",
		"POST",
		"/issues",
		ctr.IssueAdd,
	},
	// Route{
	// 	"name",
	// 	"method",
	// 	"/endpoint",
	// 	ctr.anything,
	// },
	// Route{
	// 	"listDeliveries",
	// 	"GET",
	// 	"/deliveries",
	// 	ctr.listDeliveries,
	// },
	// Route{
	// 	"getDelivery",
	// 	"GET",
	// 	"/deliveries/{id}",
	// 	ctr.getDelivery,
	// },
	Route{
		"deliveries.ratings.store",
		"POST",
		"/deliveries/{id}/ratings",
		ctr.Rate,
	},
	Route{
		"Signup",
		"POST",
		"/signup",
		ctr.SignUp,
	},
	Route{
		"Login",
		"POST",
		"/login",
		ctr.Login,
	},
	Route{
		"ResendVerificationCode",
		"GET",
		"/resend-verification-code/{userId}",
		ctr.ResendVerificationCode,
	},
	Route{
		"VerificationAccount",
		"GET",
		"/verification-account/{verificationCode}/{userId}",
		ctr.VerificationAccount,
	},
}
