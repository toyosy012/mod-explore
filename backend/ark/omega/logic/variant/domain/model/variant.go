package model

type GroupName string

type Name string

type Variant struct {
	group GroupName
	name  Name
}

func NewVariant(group GroupName, name Name) Variant {
	return Variant{
		group: group,
		name:  name,
	}
}

func (v Variant) Group() GroupName { return v.group }
func (v Variant) Name() Name       { return v.name }
