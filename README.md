Packer Template post-processor
==============================

Create templated files for post-process images. For example, a Terraform config file that includes the built AMIs.

Usage
-----
Add the post-processor to your packer template:

    {
        "post-processors": [
          {
            "type": "template",
            "template_file": "test.tfvars",
            "output_file": "out.tfvars"
          }
        ]
    }


Installation
------------
Run:

    $ go get github.com/saymedia/packer-post-processor-template
    $ go install github.com/saymedia/packer-post-processor-template

Add the post-processor to ~/.packerconfig:

    {
      "post-processors": {
        "template": "packer-post-processor-template"
      }
    }
