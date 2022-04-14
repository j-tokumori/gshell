package main

import (
	"go/ast"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"golang.org/x/tools/go/packages"
	//"github.com/NamcoBandaiStudios/prism/cmd/util/cmdflag"
)

func main() {
	//flg := flag.NewFlag()
	//flg.Parse()

	if err := write("cmd/gengshell/template/generated.go.txt",
		"test/generated.go", "", GetDataList()); err != nil {
		log.Fatal(err)
	}
}

type Data struct {
	ServiceClient string
	RPC           string
	Args          string
	Reply         string
}

func GetDataList() []Data {
	cfg := &packages.Config{Mode: packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "github.com/j-tokumori/gshell/cmd/test/api")
	if err != nil {
		panic(err)
	}

	dataList := make([]Data, 0)
	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			ast.Inspect(syntax, func(n ast.Node) bool {
				switch t := n.(type) {
				case *ast.TypeSpec:
					if t.Name.IsExported() {
						switch tt := t.Type.(type) {
						case *ast.InterfaceType:
							if regexp.MustCompile(`Client$`).Match([]byte(t.Name.Name)) {
								for _, m := range tt.Methods.List {
									dataList = append(dataList, Data{
										t.Name.Name,
										m.Names[0].Name,
										m.Type.(*ast.FuncType).Params.List[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name,  // 第２引数が Args であることを決め打ち
										m.Type.(*ast.FuncType).Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name, // 第１返値が Reply であることを決め打ち
									})
								}
							}
						}
					}
				}
				return true
			})
		}
	}
	return dataList
}

// write テンプレートファイルから書き出し
func write(templfile, output, fmt string, data []Data) error {
	force := true
	funcMap := template.FuncMap{
		//"title":   strings.Title,
		//"toUpper": strings.ToUpper,
		//"toLower": strings.ToLower,
		//"split":   strings.Split,
	}
	templ := template.Must(template.New(filepath.Base(templfile)).Funcs(funcMap).ParseFiles(templfile))

	cmd := exec.Command(fmt)
	file, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cmd.Stderr = os.Stderr
	if output == "-" {
		cmd.Stdout = os.Stdout
	} else {
		if force || !isExist(output) {
			cmd.Stdout, err = os.Create(output)
			if err != nil {
				return err
			}
			log.Println("Writing to", output)
		} else {
			log.Println("Already generated to", output)
		}
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	dataList := struct {
		DataList []Data
	}{
		DataList: data,
	}
	err = templ.Execute(file, dataList)
	if err != nil {
		return err
	}

	file.Close()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
