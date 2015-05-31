variable "images" {
    default = {
        {{index .Artifact.Id}}
    }
}