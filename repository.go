package main

import (
	"gopkg.in/doug-martin/goqu.v5"
)

type Repository struct {
	database *goqu.Database
}

func NewModuleRepository(database *goqu.Database) *ModuleRepository {
	return &ModuleRepository{
		tableName:  "module",
		Repository: Repository{database},
	}
}

type ModuleRepository struct {
	Repository
	tableName string
}

func (repo *ModuleRepository) Find(ID uint64) (module Module, err error) {
	_, err = repo.database.From(repo.tableName).Prepared(true).Where(goqu.Ex{
		"id": ID,
	}).ScanStruct(&module)
	return module, err
}

func (repo *ModuleRepository) FindAll() (modules []Module, err error) {
	err = repo.database.From(repo.tableName).Prepared(true).ScanStructs(&modules)
	return modules, err
}

func (repo *ModuleRepository) FindHistory(pkg string) (modules []Module, err error) {
	err = repo.database.From(repo.tableName).Prepared(true).Where(goqu.Ex{
		"package": pkg,
	}).ScanStructs(&modules)
	return modules, err
}

func (repo *ModuleRepository) FindAllByDevice(device string) (modules []Module, err error) {
	err = repo.database.
		From(repo.tableName).
		LeftJoin(goqu.I("device_module_through"), goqu.On(goqu.I("device_module_through.module_id").Eq(goqu.I("module.id")))).
		Prepared(true).
		Where(goqu.Ex{"device_module_through.device_id": device}).
		ScanStructs(&modules)
	return modules, err
}
