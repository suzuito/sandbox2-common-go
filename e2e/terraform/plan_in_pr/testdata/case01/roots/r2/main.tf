provider "google" {
    project = "hoge"
}

module "r2m1" {
    source = "../../commons/r2m1"
}

module "r2m2" {
    source = "../../commons/r2m2"
}
