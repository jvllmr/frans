package tasks

import (
	"context"
	"log/slog"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
)

func TicketsLifecycleTask(db *ent.Client, ts services.TicketService) {
	tickets := db.Ticket.Query().WithFiles().AllX(context.Background())
	deletedCount := 0

	for _, ticketValue := range tickets {
		if ts.ShouldDeleteTicket(ticketValue) {
			err := db.Ticket.DeleteOne(ticketValue).Exec(context.Background())
			if err != nil {
				slog.Error("Could not delete ticket", "err", err)
				continue
			}
			deletedCount += 1
		}

	}
	slog.Info("Deleted tickets", "count", deletedCount)
}
