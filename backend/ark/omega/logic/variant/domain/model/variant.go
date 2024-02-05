package model

type VariantID int

type GroupName string

type Name string

type Variant struct {
	id    VariantID
	group GroupName
	name  Name
}

type Variants []Variant

func NewVariant(id VariantID, group GroupName, name Name) Variant {
	return Variant{
		id:    id,
		group: group,
		name:  name,
	}
}

func (v Variant) ID() VariantID    { return v.id }
func (v Variant) Group() GroupName { return v.group }
func (v Variant) Name() Name       { return v.name }
