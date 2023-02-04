package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
)

func main2() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./orm/query",
		OutFile:      "./orm/query/query.go",
		ModelPkgPath: "model",

		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		//FieldNullable: true,
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessray or it will panic
	db, err := gorm.Open(sqlite.Open("file:./assets/db?cache=shared"))
	g.UseDB(db)
	if err != nil {
		log.Fatal(err)
	}

	//for _, m := range g.GenerateAllTable() {
	//	//m.(*check.BaseStruct).StructInfo.PkgPath="you pkg path"
	//	pointer := reflect.ValueOf(m).Elem()
	//	structInfoField := pointer.FieldByName("StructInfo")
	//	pkgPathField := structInfoField.FieldByName("PkgPath")
	//	pkgPathField.SetString("query")
	//	g.ApplyBasic(m)
	//}

	g.ApplyBasic(
		g.GenerateModel("error_code"),
		g.GenerateModel("i18n"),
		g.GenerateModel("lang"),
		g.GenerateModel("from_to_config"),
		g.GenerateModel("admin"),
	)

	//g.GenerateAllTable()

	// apply diy interfaces on structs or table models
	//g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

	// execute the action of code generation
	g.Execute()
}
