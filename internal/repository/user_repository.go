package repository

type UserRepository interface {
	create()
	getById()
	getAll()
	updateById()
	updateAll()
	deleteById()
	deleteAll()
}
