package protoc_gen_plugin

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init()  {
	generator.RegisterPlugin(&RegisterPlugin{})
}

type RegisterPlugin struct {
	g *generator.Generator
	msgRegisterList []string
}

func (r *RegisterPlugin) Name() string {
	return "RegisterPlugin"
}

func (r *RegisterPlugin) Init(g *generator.Generator) {
	r.g = g
}

func (r *RegisterPlugin) Generate(file *generator.FileDescriptor) {
	for _, desc := range file.MessageType {
		// Don't generate virtual messages for maps.
		if desc.GetOptions().GetMapEntry() {
			continue
		}

		addInitf("%s.RegisterType((*%s)(nil), %q)", g.Pkg["proto"], goTypeName, fullName)
		r.msgRegisterList = append(r.msgRegisterList, fmt.Sprintf("(%v, %v)", 111, generator.CamelCaseSlice(desc.Name)))
	}
}

func (r *RegisterPlugin) GenerateImports(file *generator.FileDescriptor) {

}
