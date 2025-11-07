package util

import (
	"context"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/filedata"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/otel"
)

func RefreshUserTotalDataSize(ctx context.Context, u *ent.User, tx *ent.Tx) error {
	ctx, span := otel.NewSpan(ctx, "refreshUserTotalDataSize")
	defer span.End()
	if totalDataSize, err := tx.User.Query().Where(user.ID(u.ID)).QueryFiles().QueryData().
		Aggregate(ent.Sum(filedata.FieldSize)).
		Int(ctx); err != nil {
		return err
	} else {
		if tx != nil {
			return tx.User.UpdateOne(u).SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
		}

		return u.Update().SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
	}
}
