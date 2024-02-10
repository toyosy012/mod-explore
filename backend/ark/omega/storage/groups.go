package storage

// groupModel variantsに集約しても良さそうだったがgroups単体で取り扱う可能性があるので分離しておく
type groupModel struct {
	ID   uint8  `db:"id"`
	Name string `db:"name"`
}

// TODO Implement CRUD methods
