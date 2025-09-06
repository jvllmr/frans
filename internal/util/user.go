package util

import (
	"context"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
)

func RefreshUserTotalDataSize(ctx context.Context, userValue *ent.User, tx *ent.Tx) error {
	if totalDataSize, err := userValue.QueryTickets().
		QueryFiles().
		Aggregate(ent.Sum(file.FieldSize)).
		Int(ctx); err != nil {
		return err
	} else {
		if tx != nil {
			tx.User.UpdateOne(userValue).SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
		}

		return userValue.Update().SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
	}
}
