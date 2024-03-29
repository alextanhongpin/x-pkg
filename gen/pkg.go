package gen

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// SkipCurrentPackagePath returns "" if the package is referring to an entity
// that belongs to itself, e.g. package foo calling foo.Foo.
func SkipCurrentPackagePath(pkgPath string, fieldPkgPath string) string {
	if pkgPath == fieldPkgPath {
		return ""
	}
	return fieldPkgPath
}

// packagePath returns the github package path from any given path,
// e.g. path/to/github.com/your-repo/your-pkg returns github.com/your-repo/your-pkg
// If your package is not hosted on github, you may need to override $PKG to
// set the prefix of your package.
func packagePath(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		base := filepath.Base(path)
		path = path[:len(path)-len(base)]
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	pkg := os.Getenv("PKG")
	if pkg == "" {
		pkg = "github.com"
	}
	idx := strings.Index(path, pkg)
	return path[idx:]
}

// packageName returns the base package name.
func packageName(path string) string {
	return filepath.Base(packagePath(path))
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		log.Fatalf("failed to load package: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	return pkgs[0]
}
