package util

import (
	"context"

	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/filedata"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/otel"
)

func RefreshUserTotalDataSize(ctx context.Context, userValue *ent.User, tx *ent.Tx) error {
	ctx, span := otel.NewSpan(ctx, "refreshUserTotalDataSize")
	defer span.End()
	if totalDataSize, err := tx.User.Query().Where(user.ID(userValue.ID)).QueryFileinfos().
		Aggregate(ent.Sum(filedata.FieldSize)).
		Int(ctx); err != nil {
		return err
	} else {
		if tx != nil {
			return tx.User.UpdateOne(userValue).SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
		}

		return userValue.Update().SetTotalDataSize(int64(totalDataSize)).Exec(ctx)
	}
}
