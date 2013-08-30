package main

import (
	"flag"
	"fmt"
	"github.com/theplant/goapigen/parser"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var pkg = flag.String("pkg", "", "Put a full go package path like 'github.com/theplant/qortexapi', make sure you did 'go get github.com/theplant/qortexapi'")
var lang = flag.String("lang", "javascript", "put language like 'javascript', 'objc', 'java'")
var outdir = flag.String("outdir", ".", "the dir to output the generated source code")
var impl = flag.String("impl", "", "implementation package like 'github.com/theplant/qortex/services'")
var prefix = flag.String("prefix", "", "the prefix of structs and services")
var javapackage = flag.String("java_package", "", "the package of generated java source code")

func main() {
	flag.Parse()

	buildpkg, err := build.Default.Import(*pkg, "", 0)

	if err != nil {
		die(err)
	}

	apis := parser.Parse(buildpkg.Dir, *prefix)

	switch *lang {
	case "javascript":
		printjavascript(*outdir, apis)
	case "server":
		printserver(*outdir, apis, *pkg, *impl)
	case "objc":
		printobjc(*outdir, apis)
	case "java":
		printjava(*outdir, apis, *javapackage)
	}

}

func die(message interface{}) {
	fmt.Println(message)
	os.Exit(0)
}

func dieIf(message interface{}, condition bool) {
	if !condition {
		return
	}
	die(message)
}

func codeTemplate() (tpl *template.Template) {
	tpl = template.New("")
	tpl = tpl.Funcs(template.FuncMap{
		"title":       strings.Title,
		"snake":       snake,
		"downcase":    strings.ToLower,
		"dotlastname": dotLastName,
	})
	tpl = template.Must(tpl.Parse(Templates))
	return
}

func snake(name string) (r string) {
	r = strings.ToLower(name[:1]) + name[1:]
	return
}

func dotLastName(pkg string) (r string) {
	names := strings.Split(pkg, "/")
	r = names[len(names)-1]
	return
}

func printserver(dir string, apiset *parser.APISet, apipkg string, impl string) {
	if impl == "" {
		die("must use -impl=your.package/full/path to give implementation package")
	}

	apiset.ServerImports = []string{
		"time",
		"encoding/json",
		apipkg,
		impl,
		"net/http",
		"github.com/sunfmin/govalidations",
	}
	apiset.ImplPkg = impl

	tpl := codeTemplate()

	p := filepath.Join(dir, apiset.Name+"httpimpl", "gen.go")
	os.Mkdir(filepath.Dir(p), 0755)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	err = tpl.ExecuteTemplate(f, "httpserver", apiset)
	if err != nil {
		panic(err)
	}
}

func printjava(dir string, apiset *parser.APISet, javapackage string) {
	dieIf("must use -java_package=com.qortex.android to give java package", javapackage == "")
	filedir := filepath.Join(dir, strings.Replace(javapackage, ".", "/", -1))
	err1 := os.MkdirAll(filedir, 0666)
	dieIf(err1, err1 != nil)
	tpl := codeTemplate()

	for _, dataobj := range apiset.DataObjects {
		writeSingleJavaFile(tpl, filedir, javapackage, "java/dataobject", dataobj.Name, dataobj)
	}

	for _, inf := range apiset.Interfaces {
		data := make(map[string]interface{})
		data["Prefix"] = apiset.Prefix
		data["Interface"] = inf
		data["PkgName"] = strings.Title(apiset.Name)
		writeSingleJavaFile(tpl, filedir, javapackage, "java/interface", inf.Name, data)
	}
	writeSingleJavaFile(tpl, filedir, javapackage, "java/remote_error", "RemoteError", nil)
	writeSingleJavaFile(tpl, filedir, javapackage, "java/validated", "Validated", nil)
}

func writeSingleJavaFile(tpl *template.Template, filedir string, javapackage string, templateName string, name string, data interface{}) {
	jfile, err := os.Create(filepath.Join(filedir, strings.Title(name)+".java"))
	defer jfile.Close()
	dieIf(err, err != nil)
	fmt.Fprintf(jfile, "package %s;\n\n", javapackage)
	err = tpl.ExecuteTemplate(jfile, templateName, data)
	dieIf(err, err != nil)
}

func printobjc(dir string, apiset *parser.APISet) {
	tpl := codeTemplate()
	hfile, err1 := os.Create(filepath.Join(dir, strings.Title(apiset.Name)+".h"))
	dieIf(err1, err1 != nil)
	mfile, err2 := os.Create(filepath.Join(dir, strings.Title(apiset.Name)+".m"))
	dieIf(err2, err2 != nil)
	err3 := tpl.ExecuteTemplate(hfile, "objc/h", apiset)
	dieIf(err3, err3 != nil)
	err4 := tpl.ExecuteTemplate(mfile, "objc/m", apiset)
	dieIf(err4, err4 != nil)
}

func printjavascript(dir string, apiset *parser.APISet) {
	tpl := codeTemplate()
	p := filepath.Join(dir, apiset.Name+".js")
	f, err := os.Create(p)
	dieIf(err, err != nil)
	err = tpl.ExecuteTemplate(f, "javascript/interfaces", apiset)
	dieIf(err, err != nil)
}
