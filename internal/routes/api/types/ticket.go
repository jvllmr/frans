package apiTypes

type RequestedTicketParam struct {
	ID string `uri:"ticketId" binding:"required,uuid"`
}
