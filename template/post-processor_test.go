package template

import (
	"testing"
	// "github.com/hashicorp/packer/packer"
)

func TestPostProcessorConfigure(t *testing.T) {

	var p PostProcessor
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}

	if p.config.OutputFile == "" {
		t.Fatal("should have OutputFile")
	}
	if p.config.TemplateFile == "" {
		t.Fatal("should have TemplateFile")
	}
}

func TestPostProcessor(t *testing.T) {
	var p PostProcessor
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func validDefaults() map[string]interface{} {
	return map[string]interface{}{
		"template_file": "amazon_test.tfvars",
		"output_file":   "amazon_test_out.tfvars",
	}
}
