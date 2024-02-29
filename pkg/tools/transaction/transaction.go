package transaction

import (
	"github.com/jackc/pgx/v4"

	"arch/pkg/type/context"
	log "arch/pkg/type/logger"
)

func Finish(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return log.ErrorWithContext(ctx, rollbackErr)
		}
		return err
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return log.ErrorWithContext(ctx, commitErr)
		}
		return nil
	}
}
