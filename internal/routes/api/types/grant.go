package apiTypes

type RequestedGrantParam struct {
	ID string `uri:"grantId" binding:"required,uuid"`
}
