package template

import (
	"bufio"
	// "encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

// BuilderId is the name of this post-processor in the logs
const BuilderID = "packer.post-processor.template"

// Config options
type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	TemplateFile string `mapstructure:"template_file"`
	OutputFile   string `mapstructure:"output_file"`

	ctx interpolate.Context
}

// PostProcessor master
type PostProcessor struct {
	config Config
	t      *template.Template
}

// OutputFileTemplate is for use within the templates
type OutputFileTemplate struct {
	ArtifactID string
	BuildName  string
	Provider   string
}

// Configure sets up the config options to be used later
func (p *PostProcessor) Configure(raws ...interface{}) error {

	// Configure
	err := config.Decode(&p.config,
		&config.DecodeOpts{
			Interpolate: true,
			InterpolateFilter: &interpolate.RenderFilter{
				Exclude: []string{
					"output",
				},
			},
		},
		raws...)
	if err != nil {
		return err
	}

	// apply default, var not required
	if p.config.OutputFile == "" {
		p.config.OutputFile = "out.tfvars"
	}

	required := map[string]*string{
		"template_file": &p.config.TemplateFile,
	}

	var errs *packer.MultiError
	for key, ptr := range required {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("%s must be set", key))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

// PostProcess applies the artifact to the templates
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	tmpl := p.config.TemplateFile
	output := p.config.OutputFile

	ui.Message(fmt.Sprintf("Creating artifact: %s", artifact))
	log.Printf("[DEBUG] Creating artifact: %+v\n", artifact)

	renderErr := p.renderTemplate(tmpl, artifact, output)
	if renderErr != nil { // could not render template
		return nil, false, renderErr
	}
	return artifact, false, nil
}

func (p *PostProcessor) renderTemplate(inputFile string, artifact packer.Artifact, outputFile string) (_err error) {

	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Template file read failed: %s", err)
	}
	templateStr := string(buf)

	// Create a new file if already exists
	origFile := outputFile
	err = nil
	for i := 2; i < 100; i = i + 1 {
		if i >= 100 {
			// this is not the same 100 as the line above due to the loop
			log.Printf("[DEBUG] Too many files exist. Giving up after %s", outputFile)
			break
		}
		_, err = os.Stat(outputFile)
		if err != nil {
			break
		}
		log.Printf("[DEBUG] Output file exists. Trying %s", outputFile)
		outputFile = strings.Replace(origFile, ".", "_"+strconv.Itoa(i)+".", 1)
	}

	log.Printf("[DEBUG] Template file: %s\n", outputFile)
	f, err := os.Create(outputFile)
	if err != nil { // could not render template
		return err
	}
	templateWriter := bufio.NewWriter(f)

	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"join":  strings.Join,
		"split": strings.Split,
	}).Parse(templateStr)
	if err != nil {
		return err
	}

	type TemplatePage struct {
		Artifact packer.Artifact
	}
	var thePage = TemplatePage{}
	thePage.Artifact = artifact
	if err := t.Execute(templateWriter, thePage); err != nil {
		return err
	}
	templateWriter.Flush()

	return nil
}
