package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/do"
	"github.com/samber/lo"

	"mods-explore/ark/omega/logic/creature/domain/model"
	"mods-explore/ark/omega/logic/creature/domain/service"
	"mods-explore/ark/omega/logic/creature/usecase"
	variant "mods-explore/ark/omega/logic/variant/domain/model"
)

type UniqueQueryModel struct {
	UniqueID         int            `db:"unique_id"`
	UniqueName       string         `db:"unique_name"`
	HealthMultiplier float32        `db:"health_multiplier"`
	DamageMultiplier float32        `db:"damage_multiplier"`
	BaseID           int            `db:"base_id"`
	BaseName         string         `db:"base_name"`
	BaseHealth       uint           `db:"base_health"`
	BaseMelee        int            `db:"base_melee"`
	UniqueVariants   UniqueVariants `db:"unique_variants"`
}

type UniqueVariant struct {
	VariantID   int    `db:"variant_id" json:"variant_id"`
	VariantName string `db:"variant_name" json:"variant_name"`
	GroupName   string `db:"group_name" json:"group_name"`
}

type UniqueVariants [2]UniqueVariant

func (v *UniqueVariants) Scan(value any) error {
	switch vv := value.(type) {
	case []byte:
		return json.Unmarshal(vv, v)
	case string:
		return json.Unmarshal([]byte(vv), v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (v UniqueQueryModel) toResponseCreature() (*service.ResponseCreature, error) {
	variants := v.UniqueVariants
	vs := lo.Map(variants[:], func(v UniqueVariant, _ int) model.DinosaurVariant {
		return model.NewDinosaurVariant(
			variant.NewVariant(
				variant.VariantID(v.VariantID),
				variant.VariantGroupName(v.GroupName),
				variant.Name(v.VariantName)),
			model.VariantDescriptions{},
		)
	})
	health, err := model.NewHealth(v.BaseHealth)
	if err != nil {
		return nil, err
	}
	healthMultiplier, err := model.NewUniqueMultiplier[model.Health](model.StatusMultiplier(v.HealthMultiplier))
	if err != nil {
		return nil, err
	}
	damageMultiplier, err := model.NewUniqueMultiplier[model.Melee](model.StatusMultiplier(v.DamageMultiplier))
	if err != nil {
		return nil, err
	}

	return &service.ResponseCreature{
		ResponseDinosaur: service.NewResponseDinosaur(
			model.DinosaurID(v.BaseID),
			model.DinosaurName(v.BaseName),
			health,
			model.Melee(v.BaseMelee),
		),
		ResponseVariants: service.NewResponseVariants(([2]model.DinosaurVariant)(vs)),
		ResponseUnique: service.NewResponseUnique(
			model.UniqueDinosaurID(v.UniqueID),
			model.UniqueName(v.UniqueName),
			*healthMultiplier,
			*damageMultiplier,
		),
	}, nil
}

type UniqueQueryRepo struct {
	*Client
}

func NewUniqueQueryRepo(injector *do.Injector) (usecase.UniqueQueryRepository, error) {
	return &UniqueQueryRepo{
		do.MustInvoke[*Client](injector),
	}, nil
}

func (r UniqueQueryRepo) Select(ctx context.Context, id model.UniqueDinosaurID) (*service.ResponseCreature, error) {
	row, err := NamedGet[UniqueQueryModel](
		ctx,
		r.Client,
		`SELECT
    				u.id as unique_id, u.name as unique_name,
    				u.health_multiplier, u.damage_multiplier,
    				d.id as base_id, d.name as base_name,
    				d.health as base_health, d.melee as base_melee,
					JSONB_AGG(JSONB_BUILD_OBJECT(
							'variant_id', v.id,
							'variant_name', v.name, 
							'group_name', g.name
					    )
					) as unique_variants
				FROM uniques as u 
				    JOIN dinosaurs as d ON u.dinosaur_id = d.id 
				    JOIN unique_variants as uv ON u.id = uv.unique_id 
				    JOIN variants as v ON uv.variant_id = v.id 
				    JOIN groups as g ON g.id = v.group_id 
				WHERE u.id = :id GROUP BY u.id, d.id;`,
		map[string]any{"id": id},
	)
	if err != nil {
		return nil, err
	}

	unique, err := row.toResponseCreature()
	if err != nil {
		return nil, err
	}

	return unique, nil
}

func (r UniqueQueryRepo) List(ctx context.Context) (service.ResponseCreatures, error) {
	rowsList, err := Select[UniqueQueryModel](
		ctx,
		r.Client,
		`SELECT
					u.id as unique_id, u.name as unique_name,
					u.health_multiplier, u.damage_multiplier,
					d.id as base_id, d.name as base_name,
					d.health as base_health, d.melee as base_melee,
					JSONB_AGG(JSONB_BUILD_OBJECT(
							'variant_id', v.id,
							'variant_name', v.name, 
							'group_name', g.name
					    )
					) as unique_variants
				FROM uniques as u
				    JOIN dinosaurs as d ON u.dinosaur_id = d.id
				    JOIN unique_variants as uv ON u.id = uv.unique_id
				    JOIN variants as v ON uv.variant_id = v.id
				    JOIN groups as g ON g.id = v.group_id
				GROUP BY u.id, d.id;`,
	)
	if err != nil {
		return nil, err
	}

	var response service.ResponseCreatures
	for _, rows := range rowsList {
		resp, err := rows.toResponseCreature()
		if err != nil {
			return nil, err
		}
		response = append(response, *resp)
	}
	return response, nil
}

type UniqueModel struct {
	ID               int     `db:"id"`
	Name             string  `db:"name"`
	HealthMultiplier float32 `db:"health_multiplier"`
	DamageMultiplier float32 `db:"damage_multiplier"`
}

type UniqueCommandRepo struct {
	*Client
}

func NewUniqueCommandRepo(injector *do.Injector) (service.UniqueCommandRepository, error) {
	return UniqueCommandRepo{
		do.MustInvoke[*Client](injector),
	}, nil
}

func (r UniqueCommandRepo) Insert(ctx context.Context, create service.CreateUniqueDinosaur) (model.UniqueDinosaurID, error) {
	id, err := NamedStore[int](
		ctx,
		r.Client,
		`INSERT INTO uniques (dinosaur_id, name, health_multiplier, damage_multiplier)
			VALUES (:dinosaur_id, :name, :health_multiplier, :damage_multiplier)
			RETURNING id;`,
		map[string]any{
			"dinosaur_id": create.DinosaurID(), "name": create.Name(),
			"health_multiplier": create.HealthMultiplier().Value(),
			"damage_multiplier": create.DamageMultiplier().Value(),
		},
	)
	if err != nil {
		return model.UniqueDinosaurID(0), err
	}
	return model.UniqueDinosaurID(id), nil
}

func (r UniqueCommandRepo) Update(ctx context.Context, update service.UpdateUniqueDinosaur) error {
	_, err := NamedStore[int](
		ctx,
		r.Client,
		`UPDATE uniques 
			SET dinosaur_id = :dinosaur_id, name = :name, 
			    health_multiplier = :health_multiplier, damage_multiplier = :damage_multiplier, updated_at = NOW() 
			WHERE id = :id;`,
		map[string]any{
			"id": update.ID(), "dinosaur_id": update.DinosaurID(), "name": update.Name(),
			"health_multiplier": update.HealthMultiplier().Value(),
			"damage_multiplier": update.DamageMultiplier().Value(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r UniqueCommandRepo) Delete(ctx context.Context, id model.UniqueDinosaurID) error {
	return NamedDelete(ctx, r.Client, `DELETE FROM uniques WHERE id = :id;`, map[string]any{"id": id})
}
