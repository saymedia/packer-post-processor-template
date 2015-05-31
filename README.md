Packer Terraform post-processor
=============================

Run Terraform scripts for post-process images

Usage
-----
Add the post-processor to your packer template:

    {
        "post-processors": [
          {
            "type": "terraform",
            "template_file": "test.tfvars",
            "output_file": "out.tfvars"
          }
        ]
    }

Available configuration options:



Installation
------------
Run:

    $ go get github.com/saymedia/packer-post-processor-terraform
    $ go install github.com/saymedia/packer-post-processor-terraform

Add the post-processor to ~/.packerconfig:

    {
      "post-processors": {
        "terraform": "packer-post-processor-terraform"
      }
    }
