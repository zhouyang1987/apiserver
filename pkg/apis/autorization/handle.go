package autorization

import (
	"net/http"

	r "apiserver/pkg/router"

	"github.com/gorilla/mux"
)

func InstallApi(router *mux.Router) {
	//install user's api handle
	r.RegisterHttpHandler(router, "/users", "CREATE", CreateUser)
	r.RegisterHttpHandler(router, "/users", "GET", GetUsers)
	r.RegisterHttpHandler(router, "/users/{id}", "DELETE", DeleteUser)
	r.RegisterHttpHandler(router, "/users", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/users/{id}", "OPTIONS", Option)

	//install team's api handle
	r.RegisterHttpHandler(router, "/teams", "CREATE", CreateTeam)
	r.RegisterHttpHandler(router, "/teams", "GET", GetTeams)
	r.RegisterHttpHandler(router, "/teams/{id}", "DELETE", DeleteTeam)
	r.RegisterHttpHandler(router, "/teams", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/teams/{id}", "OPTIONS", Option)

	//install permission's api handle
	r.RegisterHttpHandler(router, "/permissions", "CREATE", CreatePermission)
	r.RegisterHttpHandler(router, "/permissions", "GET", GetPermissions)
	r.RegisterHttpHandler(router, "/permissions/{id}", "DELETE", DeletePermission)
	r.RegisterHttpHandler(router, "/permissions", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/permissions/{id}", "OPTIONS", Option)

	//install role's api handle
	r.RegisterHttpHandler(router, "/roles", "CREATE", CreateRole)
	r.RegisterHttpHandler(router, "/roles", "GET", GetRoles)
	r.RegisterHttpHandler(router, "/roles/{id}", "DELETE", DeleteRole)
	r.RegisterHttpHandler(router, "/roles", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/roles/{id}", "OPTIONS", Option)
}

func Option(request *http.Request) (string, interface{}) {
	return r.StatusOK, "ok"
}
