package main

import (
	"embed"
	"flag"
	"go/ast"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"golang.org/x/tools/go/packages"
)

//go:embed template
var fs embed.FS

func main() {
	var module, template, output string
	flag.StringVar(&module, "module", "", "module name")
	flag.StringVar(&template, "template", "", "template file name")
	flag.StringVar(&output, "output", "./generated.go", "output file name")

	flag.Parse()

	if err := write(template, output, "gofmt", module, GetDataList(module)); err != nil {
		log.Fatal(err)
	}
}

type Data struct {
	ClientFactory   string
	RPCType         string
	OriginArgsType  string
	OriginReplyType string
	RPCFactory      string
	ReplyGetter     string
}

func GetDataList(moduleName string) []Data {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax} // TODO: 現sampleでは NeedTypes を入れないと、Syntaxが取れない
	pkgs, err := packages.Load(cfg, moduleName)
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
							rep := regexp.MustCompile(`ServiceClient$`)
							if rep.Match([]byte(t.Name.Name)) {
								for _, m := range tt.Methods.List {
									// 第２引数が Args であることを決め打ち
									// 第１返値が Reply であることを決め打ち
									var originArgsType, originReplyType, replyGetter string
									switch args := m.Type.(*ast.FuncType).Params.List[1].Type.(*ast.StarExpr).X; args.(type) {
									case *ast.SelectorExpr:
										originArgsType = "emptypb.Empty"
									default:
										originArgsType = "rpc." + args.(*ast.Ident).Name
									}

									switch reply := m.Type.(*ast.FuncType).Results.List[0].Type.(*ast.StarExpr).X; reply.(type) {
									case *ast.SelectorExpr:
										originReplyType = "emptypb.Empty"
										replyGetter = "Get" + rep.ReplaceAllString(t.Name.Name, "") + m.Names[0].Name + "Reply"
									default:
										originReplyType = "rpc." + reply.(*ast.Ident).Name
										replyGetter = "Get" + reply.(*ast.Ident).Name
									}
									dataList = append(dataList, Data{
										ClientFactory:   "rpc.New" + t.Name.Name,
										RPCType:         m.Names[0].Name,
										OriginArgsType:  originArgsType,
										OriginReplyType: originReplyType,
										RPCFactory:      "New" + m.Names[0].Name,
										ReplyGetter:     replyGetter,
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
func write(templfile, output, fmt string, moduleName string, data []Data) error {
	force := true
	funcMap := template.FuncMap{
		//"title":   strings.Title,
		//"toUpper": strings.ToUpper,
		//"toLower": strings.ToLower,
		//"split":   strings.Split,
	}
	var templ *template.Template
	if templfile == "" {
		templ = template.Must(template.ParseFS(fs, "template/generated.go.txt")).Funcs(funcMap)
	} else {
		templ = template.Must(template.New(filepath.Base(templfile)).Funcs(funcMap).ParseFiles(templfile))
	}

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
		ModuleName string
		DataList   []Data
	}{
		ModuleName: moduleName,
		DataList:   data,
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
