package model

type VariantID int

type Name string

type Variant struct {
	id    VariantID
	group VariantGroupName
	name  Name
}

type Variants []Variant

func NewVariant(id VariantID, group VariantGroupName, name Name) Variant {
	return Variant{
		id:    id,
		group: group,
		name:  name,
	}
}

func (v Variant) ID() VariantID           { return v.id }
func (v Variant) Group() VariantGroupName { return v.group }
func (v Variant) Name() Name              { return v.name }

type VariantGroup struct {
	id   VariantGroupID
	name VariantGroupName
}

func NewVariantGroup(id VariantGroupID, name VariantGroupName) VariantGroup {
	return VariantGroup{id, name}
}

type VariantGroupID uint
type VariantGroupName string

func (g VariantGroup) ID() VariantGroupID     { return g.id }
func (g VariantGroup) Name() VariantGroupName { return g.name }
