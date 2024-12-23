package computil

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	baseImportPath = ``
	baseRegexpStr  = ``
	baseRegexp     = regexp.MustCompile(baseRegexpStr)
	testFileRegexp = regexp.MustCompile(`.*_test\.go$`)
	windowsFix     = regexp.MustCompile(`\\`)

	// GenesisLibs is the name of the packages within the genesis standard library
	GenesisLibs = map[string]bool{
		"crypto":   true,
		"encoding": true,
		"exec":     true,
		"file":     true,
		"net":      true,
		"os":       true,
		"rand":     true,
		"requests": true,
		"time":     true,
	}

	// InstalledGoPackages holds a cache of all currently installed golang libraries
	InstalledGoPackages = GatherInstalledGoPackages()
)

func regexpForModule(mod ...string) *regexp.Regexp {
	if runtime.GOOS == "windows" {
		return regexp.MustCompile(windowsFix.ReplaceAllString(filepath.Join(append([]string{baseImportPath}, mod...)...), `/`))
	}
	return regexp.MustCompile(filepath.Join(append([]string{baseRegexpStr}, mod...)...))
}

// List packages on workDir.
// workDir is required for module mode. If the workDir is not under module, then it will fallback to GOPATH mode.
func list(opts PkgOptions) (map[string]Pkg, error) {
	pkgs := make(map[string]Pkg)

	if opts.WorkDir == "" {
		// force on GOPATH mode
		// fmt.Println("FORCING GOPATH MODE")
		for _, srcDir := range build.Default.SrcDirs() {
			// fmt.Printf("SOURCE DIR: %v\n", srcDir)
			err := collectPkgs(srcDir, opts.WorkDir, opts.NoVendor, pkgs)
			if err != nil {
				return nil, err
			}
		}
		return pkgs, nil
	}

	mods, err := listMods(opts.WorkDir)
	if err != nil {
		// GOPATH mode
		for _, srcDir := range build.Default.SrcDirs() {
			err = collectPkgs(srcDir, opts.WorkDir, opts.NoVendor, pkgs)
			if err != nil {
				return nil, err
			}
		}
		return pkgs, nil
	}

	// Module mode
	if err = collectPkgs(filepath.Join(build.Default.GOROOT, "src"), opts.WorkDir, false, pkgs); err != nil {
		return nil, err
	}

	for _, m := range mods {
		err = collectModPkgs(m, pkgs)
		if err != nil {
			return nil, err
		}
	}

	return pkgs, nil
}

func SetImportPath(path string) {
	baseImportPath = path
	var specialChars = []string{".", "\\", "+", "*", "?", "|", "{", "}", "(", ")", "[", "]", "^", "$"}
	for _, char := range specialChars {
		path = strings.ReplaceAll(path, char, "\\"+char)
	}
	baseRegexpStr = path
}

// GatherInstalledGoPackages retrieves a list of all installed go packages in the context of current GOPATH and GOROOT
func GatherInstalledGoPackages() map[string]Pkg {
	goPackages, err := List(PkgOptions{NoVendor: true})
	if err != nil {
		panic(err)
	}
	for p, _ := range goPackages {
		if !strings.Contains(p, "/usr/lib/go-1.23/") {
			fmt.Printf(p + "\n")
			fmt.Printf(goPackages[p].ImportPath + "\n")
		}
	}
	if runtime.GOOS == "windows" {
		pathFix := regexp.MustCompile(`\\`)
		newMap := map[string]Pkg{}
		for n, p := range goPackages {
			newMap[pathFix.ReplaceAllString(n, `/`)] = p
		}
		return newMap
	}
	return goPackages
}

// SourceFileIsTest determines if the given source file is named after the test convention
func SourceFileIsTest(src string) bool {
	return testFileRegexp.MatchString(src)
}

// ResolveGoPath attempts to resolve the current user's GOPATH
func ResolveGoPath() string {
	gp := os.Getenv("GOPATH")
	if gp != "" {
		return gp
	}
	u, err := user.Current()
	if err != nil {
		// really shouldn't happen
		panic(err)
	}
	return filepath.Join(u.HomeDir, "go")
}

// ResolveEngineDir attempts to resolve the absolute path of the genesis engine directory
func ResolveEngineDir() (targetDir string, err error) {
	dirMatch := regexpForModule("engine")
	for name, pkg := range InstalledGoPackages {
		if !dirMatch.MatchString(name) {
			continue
		}
		targetDir = pkg.Dir
	}
	if targetDir == "" {
		return targetDir, fmt.Errorf("coult not locate the genesis engine directory")
	}
	return targetDir, nil
}

// ResolveStandardLibraryDir attempts to resolve the absolute path of the specified standard library package
func ResolveStandardLibraryDir(pkg string) (*Pkg, error) {
	dirMatch := regexpForModule("stdlib", pkg)
	for name, gpkg := range InstalledGoPackages {
		if !dirMatch.MatchString(name) {
			continue
		}
		return &gpkg, nil
	}
	return nil, fmt.Errorf("could not locate standard library package %s", pkg)
}

func ResolveGlobalImport(pkg string) (*Pkg, error) {
	for _, gpkg := range InstalledGoPackages {
		if gpkg.ImportPath == pkg {
			return &gpkg, nil
		}
	}
	return nil, errors.New("could not locate gopackage of that name")
}
