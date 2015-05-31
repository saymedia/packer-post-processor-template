package template

import (
	"bufio"
	// "encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	// ui.Message(fmt.Sprintf("Creating artifact: %s", artifact))
	log.Printf("[DEBUG] Creating artifact: %+v\n", artifact)
	log.Printf("[DEBUG] Creating artifact: %+v\n", artifact.BuilderId())
	log.Printf("[DEBUG] Creating artifact: %+v\n", artifact.Id())
	log.Printf("[DEBUG] Creating artifact: %+v\n", artifact.String())

	if artifact.BuilderId() == "amazon-ebs" {
		// log.Printf("[DEBUG] Artifact.BuilderId: %+v\n", artifact.BuilderId)
		// log.Printf("[DEBUG] Artifact.Amis: %+v\n", artifact.Amis)
		// log.Printf("[DEBUG] Artifact.Amis[\"us-west-1\"]: %+v\n", artifact.Amis["us-west-1"])
		// log.Printf("[DEBUG] Artifact.BuilderIdValue: %+v\n", artifact.BuilderIdValue)
	}

	renderErr := p.renderTemplate(tmpl, artifact, output)
	if renderErr != nil { // could not render template
		return nil, false, renderErr
	}
	log.Printf("[DEBUG] Template file: %s\n", output)
	return artifact, false, nil

	// // Write our Vagrantfile
	// var customVagrantfile string
	// if config.VagrantfileTemplate != "" {
	//     ui.Message(fmt.Sprintf("Using custom Vagrantfile: %s", config.VagrantfileTemplate))
	//     customBytes, err := ioutil.ReadFile(config.VagrantfileTemplate)
	//     if err != nil {
	//         return nil, false, err
	//     }

	//     customVagrantfile = string(customBytes)
	// }

	// f, err := os.Create(filepath.Join(dir, "Vagrantfile"))
	// if err != nil {
	//     return nil, false, err
	// }

	// t := template.Must(template.New("root").Parse(boxVagrantfileContents))
	// err = t.Execute(f, &vagrantfileTemplate{
	//     ProviderVagrantfile: vagrantfile,
	//     CustomVagrantfile:   customVagrantfile,
	// })
	// f.Close()
	// if err != nil {
	//     return nil, false, err
	// }
}

func (p *PostProcessor) renderTemplate(inputFile string, artifact packer.Artifact, outputFile string) (_err error) {

	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Template file read failed: %s", err)
	}
	templateStr := string(buf)

	f, err := os.Create(outputFile)
	if err != nil { // could not render template
		return err
	}
	templateWriter := bufio.NewWriter(f)

	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"join":  strings.Join,
		"split": strings.Split,
		"splitA": func(s string) string {
			out := strings.Split(s, ":")
			if len(out) > 0 {
				return out[0]
			}
			return ""
		},
		"splitB": func(s string) string {
			out := strings.Split(s, ":")
			if len(out) > 1 {
				return out[1]
			}
			return ""
		},
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
