package swagger

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/wusphinx/gin-swagger/codegen"
	"github.com/wusphinx/gin-swagger/program"
)

func NewEnumGenerator(packagePath string, registerEnumMethod string) *EnumGenerator {
	prog := program.NewProgram(packagePath, true)

	return &EnumGenerator{
		PackagePath:        packagePath,
		RegisterEnumMethod: registerEnumMethod,
		Program:            prog,
	}
}

type Option struct {
	Value string
	Label string
}

type Enum struct {
	Name        string
	Type        string
	Pathname    string
	PackageName string
	Values      []Option
}

type EnumGenerator struct {
	PackagePath        string
	RegisterEnumMethod string
	Program            *program.Program
	Enums              []Enum
}

func (g *EnumGenerator) addEnum(name string, tpe string, pkg *types.Package, enumOptions []program.Option) {
	if g.Enums == nil {
		g.Enums = []Enum{}
	}

	options := []Option{}

	for _, option := range enumOptions {
		switch option.Value.(type) {
		case string:
			options = append(options, Option{
				Value: option.Value.(string),
				Label: option.Label,
			})
		}
	}

	if len(options) > 0 {
		g.Enums = append(g.Enums, Enum{
			Name:        name,
			Type:        tpe,
			Pathname:    pkg.Path(),
			PackageName: pkg.Name(),
			Values:      options,
		})
	}
}

func (g *EnumGenerator) Scan() {
	for _, pkgInfo := range g.Program.AllPackages {
		if program.IsSubPackageOf(g.PackagePath, pkgInfo.Pkg.Path()) {
			for id, def := range pkgInfo.Defs {
				doc := program.GetTextFromCommentGroup(g.Program.CommentGroupFor(id))

				_, hasEnum := ParseEnum(doc)

				if hasEnum {
					options := g.Program.GetEnumOptionsByType(id)
					if def.Name() != "Bool" {
						g.addEnum(
							def.Name(),
							def.Type().Underlying().String(),
							def.Pkg(),
							options,
						)
					}
				}
			}
		}

	}
}

// HasElem  especially, if src is empty, still return true
func HasElem(src []string, x string) bool {
	if len(src) == 0 {
		return true
	}
	for _, i := range src {
		if i == x {
			return true
		}
	}
	return false
}

func (g *EnumGenerator) Output(src ...string) {
	g.Scan()

	if len(g.Enums) < 1 {
		return
	}

	for _, enum := range g.Enums {
		if HasElem(src, enum.Name) == false {
			continue
		}

		relPath, _ := filepath.Rel(g.PackagePath, enum.Pathname)

		name := strings.Replace(codegen.ToLowerSnakeCase(enum.Name), "_", " ", -1)

		imports := []string{"errors", "strings"}

		registerEnumImport, registerEnumMethod := program.ParsePkgExpose(g.RegisterEnumMethod)

		if registerEnumImport != "" {
			imports = append(imports, registerEnumImport)
		}

		blocks := []string{
			codegen.DeclPackage(enum.PackageName),
			codegen.DeclImports(imports...),
			codegen.DeclVar("Invalid"+enum.Name, `errors.New("invalid `+name+`")`),
			ParseEnumStringify(enum),
			ParseEnumLabel(enum),
			ParseEnumParser(enum),
			ParseEnumLabelParser(enum),
			ParseEnumJSONMarshal(enum),
		}

		if registerEnumMethod != "" {
			blocks = append(blocks, ParseRegisterEnum(registerEnumMethod, enum))
		}

		codegen.GenerateGoFile(
			codegen.JoinWithSlash(relPath, codegen.ToLowerSnakeCase(enum.Name)+".go"),
			strings.Join(blocks, "\n\n"),
		)
	}
}

func ParseRegisterEnum(registerEnumMethod string, enum Enum) string {
	codes := `func init () {
	`

	for _, option := range enum.Values {
		codes += registerEnumMethod + `(` + strings.Join([]string{
			codegen.WithQuotes(enum.Name),
			codegen.WithQuotes(option.Value),
			codegen.WithQuotes(option.Label),
		}, ", ") + `)
		`
	}

	codes += `}`
	return codes
}

func ParseEnumLabelParser(enum Enum) string {
	firstLine := codegen.TemplateRender(`func Parse{{ .Name }}FromLabelString(s string) ({{ .Name }}, error) {
	switch s {`)(enum)

	var lines = []string{
		firstLine,
	}

	prefix := codegen.ToUpperSnakeCase(enum.Name)

	lines = append(lines, codegen.DeclCase(codegen.WithQuotes("")))
	lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"_UNKNOWN", "nil")))

	for _, option := range enum.Values {
		lines = append(lines, codegen.DeclCase(codegen.WithQuotes(option.Label)))
		lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"__"+option.Value, "nil")))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"_UNKNOWN", codegen.TemplateRender(`Invalid{{ .Name }}`)(enum))))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseEnumParser(enum Enum) string {
	firstLine := codegen.TemplateRender(`func Parse{{ .Name }}FromString(s string) ({{ .Name }}, error) {
	switch s {`)(enum)

	var lines = []string{
		firstLine,
	}

	prefix := codegen.ToUpperSnakeCase(enum.Name)

	lines = append(lines, codegen.DeclCase(codegen.WithQuotes("")))
	lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"_UNKNOWN", "nil")))

	for _, option := range enum.Values {
		lines = append(lines, codegen.DeclCase(codegen.WithQuotes(option.Value)))
		lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"__"+option.Value, "nil")))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(codegen.JoinWithComma(prefix+"_UNKNOWN", codegen.TemplateRender(`Invalid{{ .Name }}`)(enum))))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseEnumStringify(enum Enum) string {
	firstLine := codegen.TemplateRender(`func (v {{ .Name }}) String() string {
	switch v {`)(enum)

	var lines = []string{
		firstLine,
	}

	prefix := codegen.ToUpperSnakeCase(enum.Name)

	lines = append(lines, codegen.DeclCase(prefix+"_UNKNOWN"))
	lines = append(lines, codegen.DeclReturn(codegen.WithQuotes("")))

	for _, option := range enum.Values {
		lines = append(lines, codegen.DeclCase(prefix+"__"+option.Value))
		lines = append(lines, codegen.DeclReturn(codegen.WithQuotes(option.Value)))
	}

	lines = append(lines, `}
	return "UNKNOWN"
	}`)

	return strings.Join(lines, "\n")
}

func ParseEnumLabel(enum Enum) string {
	firstLine := codegen.TemplateRender(`func (v {{ .Name }}) Label() string {
	switch v {`)(enum)

	var lines = []string{
		firstLine,
	}

	prefix := codegen.ToUpperSnakeCase(enum.Name)

	lines = append(lines, codegen.DeclCase(prefix+"_UNKNOWN"))
	lines = append(lines, codegen.DeclReturn(codegen.WithQuotes("")))

	for _, option := range enum.Values {
		lines = append(lines, codegen.DeclCase(prefix+"__"+option.Value))
		lines = append(lines, codegen.DeclReturn(codegen.WithQuotes(option.Label)))
	}

	lines = append(lines, `}
	return "UNKNOWN"
	}`)

	return strings.Join(lines, "\n")
}

func ParseEnumJSONMarshal(enum Enum) string {
	return codegen.TemplateRender(`
func (v {{ .Name }}) MarshalJSON() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, Invalid{{ .Name }}
	}
	return []byte("\"" + str + "\""), nil
}

func (v *{{ .Name }}) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(strings.ToUpper(string(data)), "\"")
	*v, err = Parse{{ .Name }}FromString(s)
	return
}`)(enum)
}
