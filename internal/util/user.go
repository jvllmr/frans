package util

import (
	"context"

	"codeberg.org/jvllmr/frans/internal/ent"
	"codeberg.org/jvllmr/frans/internal/ent/filedata"
	"codeberg.org/jvllmr/frans/internal/ent/user"
	"codeberg.org/jvllmr/frans/internal/otel"
)

func RefreshUserTotalDataSize(ctx context.Context, u *ent.User, tx *ent.Tx) error {
	ctx, span := otel.NewSpan(ctx, "refreshUserTotalDataSize")
	defer span.End()
	filesCount, err := tx.User.Query().QueryFiles().Count(ctx)
	if err != nil {
		return err
	}
	totalDataSize := 0
	if filesCount > 0 {
		totalDataSize, err = tx.User.Query().Where(user.ID(u.ID)).QueryFiles().QueryData().
			Aggregate(ent.Sum(filedata.FieldSize)).
			Int(ctx)
		if err != nil {
			return err
		}
	}

	return tx.User.UpdateOne(u).SetTotalDataSize(int64(totalDataSize)).Exec(ctx)

}
