package gen

import (
	"flag"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alextanhongpin/pkg/stringcase"
)

// StructField for the example below.
//type Foo struct {
//  Name sql.NullString `json:"name"
//}
type StructField struct {
	Name string `example:"Name"`
	// Useful when the output directory doesn't match the existing ones.
	PkgPath  string `example:"github.com/alextanhongpin/go-codegen/test"`
	PkgName  string `example:"test"`
	Exported bool   `example:"true"`
	Tag      string `example:"build:'-'"` // To ignore builder.
	*Field
}

type Field struct {
	Type         string `example:"NullString"`
	PkgPath      string `example:"database/sql"`
	IsPointer    bool
	IsCollection bool // Whether it's an array or slice.
	IsMap        bool
	MapKey       *Field
	MapValue     *Field
}

// NewField recursively checks for the field type.
func NewField(typ types.Type) *Field {
	var isPointer, isCollection, isMap bool
	var fieldPkgPath, fieldType string
	var mapKey, mapValue *Field

	switch t := typ.(type) {
	case *types.Pointer:
		isPointer = true
		typ = t.Elem()
	case *types.Slice:
		isCollection = true
		typ = t.Elem()
	case *types.Array:
		isCollection = true
		typ = t.Elem()
	case *types.Map:
		isMap = true
		mapKey = NewField(t.Key())
		mapValue = NewField(t.Elem())
	}

	// In case the slice or array is pointer, we take the elem again.
	switch t := typ.(type) {
	case *types.Pointer:
		isPointer = true
		typ = t.Elem()
	}

	switch t := typ.(type) {
	case *types.Named:
		obj := t.Obj()
		fieldPkgPath = obj.Pkg().Path()
		fieldType = obj.Name()
	default:
		fieldType = t.String()
	}

	return &Field{
		Type:         fieldType,
		PkgPath:      fieldPkgPath,
		IsCollection: isCollection,
		IsPointer:    isPointer,
		IsMap:        isMap,
		MapKey:       mapKey,
		MapValue:     mapValue,
	}
}

type Option struct {
	In         string
	Out        string
	PkgName    string
	PkgPath    string
	StructName string
	Fields     []StructField
}

type Generator func(opt Option) error

func New(fn Generator) error {
	structPtr := flag.String("type", "", "the target struct name")
	inPtr := flag.String("in", os.Getenv("GOFILE"), "the input file, defaults to the file with the go:generate comment")
	outPtr := flag.String("out", "", "the output directory")
	flag.Parse()

	in := fullPath(*inPtr)

	// Allows -type=Foo,Bar
	structNames := strings.Split(*structPtr, ",")
	for _, structName := range structNames {
		var out string
		if o := *outPtr; o == "" {
			// Foo becomes foo.go
			fileName := stringcase.SnakeCase(structName) + ".go"

			// foo.go becomes foo_gen.go
			genFileName := safeAddSuffixToFileName(fileName, "_gen")

			// path/to/main.go becomes path/to/foo_gen.go
			out = safeAddFileName(filepath.Dir(in), genFileName)
		} else {
			out = fullPath(o)
		}

		pkg := loadPackage(packagePath(in)) // github.com/your-github-username/your-pkg.
		pkgPath := pkg.PkgPath              // Specify the config packages.NeedName to get this value.
		pkgName := pkg.Name                 // main

		obj := pkg.Types.Scope().Lookup(structName)
		if obj == nil {
			log.Fatalf("struct %s not found", structName)
		}

		// Check if it is a declared typed.
		if _, ok := obj.(*types.TypeName); !ok {
			log.Fatalf("%v is not a named type", obj)
		}

		// Check if the type is a struct.
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			log.Fatalf("%v is not a struct", obj)
		}

		fields := extractFields(structType)
		if err := fn(Option{
			PkgName:    pkgName,
			PkgPath:    pkgPath,
			Out:        out,
			In:         in,
			StructName: structName,
			Fields:     fields,
		}); err != nil {
			return err
		}
	}
	return nil
}

func extractFields(structType *types.Struct) []StructField {
	fields := make([]StructField, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)

		fields[i] = StructField{
			Name:     field.Name(),
			PkgPath:  field.Pkg().Path(),
			Exported: field.Exported(),
			Field:    NewField(field.Type()),
			Tag:      tag,
		}
	}
	return fields
}
