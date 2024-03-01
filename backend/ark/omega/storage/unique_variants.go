package storage

import (
	"context"
	"github.com/samber/lo"

	"github.com/samber/do"

	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
	variantModel "mods-explore/ark/omega/logic/variant/domain/model"
)

type StoreUniqueVariants struct {
	UniqueDinoID int
	VariantID    int
}

type UniqueVariantsModel struct {
	ID           int `db:"id"`
	UniqueDinoID int `db:"unique_id"`
	VariantID    int `db:"variant_id"`
}

type UniqueVariantsClient struct {
	*Client
}

func NewUniqueVariantsClient(injector *do.Injector) (service.UniqueVariantsCommand, error) {
	return UniqueVariantsClient{
		do.MustInvoke[*Client](injector),
	}, nil
}

func (c UniqueVariantsClient) Insert(ctx context.Context, create service.CreateVariants) (model.UniqueVariantID, error) {
	ids := create.VariantIDs()
	records := lo.Map(ids[:], func(id variantModel.VariantID, _ int) map[string]any {
		return map[string]any{"variant_id": id, "unique_id": create.UniqueDinosaurID()}
	})
	id, err := NamedStore[int](
		ctx,
		c.Client,
		`INSERT INTO unique_variants (variant_id) VALUES (:variant_id, :unique_id) RETURNING id;`,
		records,
	)
	if err != nil {
		return 0, err
	}
	return model.UniqueVariantID(id), err
}

func (c UniqueVariantsClient) Update(ctx context.Context, update service.UpdateVariants) error {
	ids := update.VariantIDs() // メソッドで戻ってきた配列のアドレスを直接sliceに変換できないので、一度アドレスを格納する変数に入れ直す
	records := lo.Map(ids[:], func(id variantModel.VariantID, _ int) map[string]any {
		return map[string]any{"variant_id": id, "unique_id": update.UniqueDinosaurID()}
	})
	_, err := NamedStore[int](
		ctx,
		c.Client,
		`UPDATE unique_variants SET unique_id = :unique_id, variant_id = :variant_id, updated_at = NOW() WHERE id = :id;`,
		records,
	)
	if err != nil {
		return err
	}
	return err
}

func (c UniqueVariantsClient) Delete(ctx context.Context, id model.UniqueVariantID) error {
	return NamedDelete(ctx, c.Client, `DELETE FROM unique_variants WHERE id = :id;`, map[string]any{"id": id})
}
