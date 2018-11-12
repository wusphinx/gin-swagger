package codegen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wusphinx/gin-swagger/codegen"
)

func TestPrinter(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("package some_package\n", codegen.DeclPackage("some_package"))
	t.Equal("type Test int\n", codegen.DeclType("Test", "int"))
}
