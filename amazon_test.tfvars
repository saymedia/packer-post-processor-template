{{$arr := split .Artifact.Id ":"}}
variable "images" {
    default = {
        {{index $arr 0}} = "{{index $arr 1}}"
    }
}