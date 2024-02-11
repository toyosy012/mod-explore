package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"mods-explore/ark/omega/logic/variant/domain/model"
	"mods-explore/ark/omega/logic/variant/domain/service"
)

// GroupVariantModel variantsに集約しても良さそうだったがgroups単体で取り扱う可能性があるので分離しておく
type groupVariantModel struct {
	ID   uint8  `db:"id"`
	Name string `db:"name"`
}

type VariantGroupClient struct {
	*sqlx.DB
}

func NewVariantGroupClient(dsn string) (service.VariantGroupRepository, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(timeout); err != nil {
		return nil, fmt.Errorf("connection timeout: %s ", err)
	}

	return VariantGroupClient{
		DB: db,
	}, nil
}

func (v VariantGroupClient) Select(ctx context.Context, id model.VariantGroupID) (*model.VariantGroup, error) {
	query, args, err := v.BindNamed(
		`SELECT id, name FROM groups WHERE id = :id;`,
		map[string]any{"id": id},
	)
	if err != nil {
		return nil, err
	}

	var row groupVariantModel
	if err = v.GetContext(ctx, &row, query, args...); errors.Is(err, sql.ErrNoRows) {
		return nil, service.NotFound
	} else if err != nil {
		return nil, err
	}

	var variant = model.NewVariantGroup(
		model.VariantGroupID(row.ID),
		model.VariantGroupName(row.Name),
	)
	return &variant, nil
}

func (v VariantGroupClient) List(ctx context.Context) (model.VariantGroups, error) {
	query := `SELECT id, name FROM groups;`

	var rows []groupVariantModel
	if err := v.SelectContext(ctx, &rows, query); errors.Is(err, sql.ErrNoRows) {
		return nil, service.NotFound
	} else if err != nil {
		return nil, err
	}

	var results model.VariantGroups
	for _, r := range rows {
		results = append(
			results,
			model.NewVariantGroup(
				model.VariantGroupID(r.ID),
				model.VariantGroupName(r.Name),
			),
		)
	}

	return results, nil
}

func (v VariantGroupClient) Insert(ctx context.Context, create service.CreateVariantGroup) (*model.VariantGroup, error) {
	stmt, err := v.PrepareNamedContext(
		ctx,
		`INSERT INTO groups (name) VALUES (:name) RETURNING id;`,
	)
	if err != nil {
		return nil, err
	}

	var id int
	err = stmt.
		QueryRowxContext(ctx, map[string]any{"name": create.Name()}).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	result, err := v.Select(ctx, model.VariantGroupID(id))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (v VariantGroupClient) Update(ctx context.Context, update service.UpdateVariantGroup) (*model.VariantGroup, error) {
	query, args, err := v.BindNamed(
		`UPDATE groups SET name = :name, updated_at = NOW() WHERE id = :id RETURNING id;`,
		map[string]any{"id": update.ID(), "name": update.Name()},
	)
	if err != nil {
		return nil, err
	}

	_, err = v.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := v.Select(ctx, update.ID())
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (v VariantGroupClient) Delete(ctx context.Context, id model.VariantGroupID) error {
	_, err := v.NamedExecContext(
		ctx,
		`DELETE FROM groups WHERE id = :id;`,
		map[string]any{"id": id},
	)
	return err
}
