package user

type RequestWellKnownName struct {
	Name string `json:"name" path:"name" query:"name"`
}
