package postgres

import (
	"context"
	"errors"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"

	"github.com/jackc/pgx/v5/pgconn"
)

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{txManager: txManager}
}

func (r *PgRepository) AddUser(ctx context.Context, req dto.UserDto) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
	insert into users(username, password) values($1, $2)
	`, req.UserName, req.UserPassword)

	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) GetUserByUsername(ctx context.Context, username string) (*dto.UserDto, error) {
	tx := r.txManager.GetQueryEngine(ctx)
	user := &dto.UserDto{UserName: username}
	err := tx.QueryRow(ctx, `
	select password, id from users where username = $1
	`, username).Scan(&user.UserPassword, &user.UserId)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PgRepository) AddRefreshSession(ctx context.Context, req dto.RefreshSessionDto) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
	insert into refreshSessions(user_id, refresh_token, fingerprint, ip, expires_in, created_at) values($1, $2, $3, $4, $5, $6)
	`, req.UserId, req.RefreshToken, req.Fingerprint, req.Ip, req.Expires, req.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) AddAsset(ctx context.Context, req dto.AssetDto) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
	insert into assets(user_id, name, descr, price) values($1, $2, $3, $4)
	`, req.UserId, req.AssetName, req.AssetDescr, req.AssetPrice)

	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) DeleteAsset(ctx context.Context, req dto.DeleteAssetDto) error {
	tx := r.txManager.GetQueryEngine(ctx)
	x, err := tx.Exec(ctx, `
	update assets set price = -1 where user_id = $1 and name = $2
	`, req.UserId, req.AssetName)
	if err != nil {
		return err
	}
	if x.RowsAffected() == 0 {
		return domain.ErrNoSuchAsset
	}
	return nil
}

func (r *PgRepository) GetAssetInfo(ctx context.Context, assetName string) (*dto.AssetDto, error) {
	var asset dto.AssetDto
	tx := r.txManager.GetQueryEngine(ctx)
	err := tx.QueryRow(ctx, `
	select id, user_id, price from assets where name = $1 and price != -1
	`, assetName).Scan(&asset.AssetId, &asset.UserId, &asset.AssetPrice)

	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *PgRepository) Refresh(ctx context.Context, req dto.UpdateRefreshDto) error {
	tx := r.txManager.GetQueryEngine(ctx)
	var check bool
	err := tx.QueryRow(ctx, `
	select exists (select 1 from refreshSessions WHERE user_id = $1 and refresh_token = $2)
	`, req.UserId, req.PriviousRefresh).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return domain.ErrNoSuchSession
	}

	_, err = tx.Exec(ctx, `
	update refreshSessions set refresh_token = $1 where user_id = $2 and refresh_token = $3
	`, req.RefreshToken, req.UserId, req.PriviousRefresh)
	if err != nil {
		return err
	}
	return nil
}

func UnwrapPgCode(err error) string {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code
		}
	}
	return ""
}
