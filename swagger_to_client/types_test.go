package swagger_to_client_test

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/wusphinx/gin-swagger/swagger_to_client"
)

func TestGetTypeFromSchema(tt *testing.T) {
	t := assert.New(tt)

	{
		schemaWithRef := spec.RefSchema("Test")
		typeName, _ := swagger_to_client.GetTypeFromSchema(*schemaWithRef)
		t.Equal("Test", typeName)
	}

	{
		schemaWithItems := spec.ArrayProperty(spec.StringProperty())
		typeName, _ := swagger_to_client.GetTypeFromSchema(*schemaWithItems)

		t.Equal("[]string", typeName)
	}
}
