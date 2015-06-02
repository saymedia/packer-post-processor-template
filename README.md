Packer Template post-processor
==============================

Create templated files for post-process images. For example, a Terraform config file that includes the built AMIs.

[![travis build status for packer-post-processor-template](https://travis-ci.org/saymedia/packer-post-processor-template.svg)](https://travis-ci.org/saymedia/packer-post-processor-template)

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

Check out [`amazon_test.json`](./amazon_test.json) and [`docker_test.json`](./docker_test.json) for examples. The Docker test is not particularly useful but it does run quickly for testing purposes.

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

Tests
-----

### STEP 0:

Install packer-post-processor-template as detailed above.

### STEP 1:

Create a file called `amazon_test_variables.json` with the contents similar to (use your own AWS values of course):

    {
        "AWS_VPC_ID": "vpc-2421cc41",
        "AWS_SUBNET_ID": "subnet-49c78e2a",
        "AWS_SG_ID": "sg-4a57f80b",
        "AWS_AMI_ID": "ami-ca9d798e",
        "AWS_REGION": "us-west-1"
    }

### STEP 2:

Make sure Docker is running.

### STEP 3:

Run:

  
    $ export AWS_ACCESS_KEY='<YOUR_AWS_ACCESS_KEY>'
    $ export AWS_SECRET_KEY='<YOUR_AWS_SECRET_KEY>'
    $ ./test.sh
