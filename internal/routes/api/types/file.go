package apiTypes

type RequestedFileParam struct {
	ID string `uri:"fileId" binding:"required,uuid"`
}
