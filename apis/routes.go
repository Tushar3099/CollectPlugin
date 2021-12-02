package apis

import (
	"net/http"

	"github.com/Tushar3099/CollectPlugin/services"
)

type Route struct {
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}
 
type Routes []Route
 
var routes = Routes{
    Route{
        "GET",
        "/",
        services.Index,
    },
}
