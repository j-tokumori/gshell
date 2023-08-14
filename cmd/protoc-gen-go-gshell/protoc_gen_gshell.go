package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"log"
	"strings"
	"text/template"

	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

//go:embed template.gohtml
var fs embed.FS

// TemplateData is template passing structure
type TemplateData struct {
	PbPackageAlias string
	GoPackage      string
	DataList       []rpcData
}

// rpcData is RPC information
type rpcData struct {
	ServiceName     string
	MethodName      string
	OriginArgsType  string
	OriginReplyType string
}

// GShellPlugin is definition file generation plugin for gshell
type GShellPlugin struct{}

const pbPackageAlias = "rpc"

// Generate implements the protokit.Generator interface
func (u GShellPlugin) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	files := protokit.ParseCodeGenRequest(req)
	var resp plugin.CodeGeneratorResponse
	setSupportedFeaturesOnCodeGeneratorResponse(&resp)
	templateData := TemplateData{
		PbPackageAlias: pbPackageAlias,
		GoPackage:      files[0].GetOptions().GetGoPackage(),
	}
	for _, f := range files {
		messageMap := make(map[string]*protokit.Descriptor)
		for _, m := range f.GetMessages() {
			messageMap[m.GetName()] = m
		}
		for _, service := range f.GetServices() {
			for _, method := range service.GetMethods() {
				var data rpcData
				data.ServiceName = service.GetName()
				data.MethodName = method.GetName()
				data.OriginArgsType = originalTypeName(method.GetInputType(), f.GetPackage())
				data.OriginReplyType = originalTypeName(method.GetOutputType(), f.GetPackage())
				templateData.DataList = append(templateData.DataList, data)
			}
		}
	}

	buf := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFS(fs, "template.gohtml"))
	if err := tmpl.Execute(buf, templateData); err != nil {
		log.Fatalf("%+v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String("generated.go"),
		Content: proto.String(string(formatted)),
	})
	return &resp, nil
}

func originalTypeName(t, p string) string {
	if t == ".google.protobuf.Empty" {
		return "emptypb.Empty"
	}
	s := strings.TrimPrefix(t, fmt.Sprintf(".%s.", p))
	return fmt.Sprintf("%s.%s", pbPackageAlias, s)
}

func setSupportedFeaturesOnCodeGeneratorResponse(resp *plugin.CodeGeneratorResponse) {
	// Enable support for optional keyword in proto3.
	sf := uint64(plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp.SupportedFeatures = &sf
}

func main() {
	if err := protokit.RunPlugin(new(GShellPlugin)); err != nil {
		log.Fatalf("%+v", err)
	}
}
