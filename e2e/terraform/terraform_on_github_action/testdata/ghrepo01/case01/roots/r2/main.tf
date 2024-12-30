provider "google" {
    project = "prj01"
}

module "r2m1" {
    source = "../../commons/r2m1"
}

module "r2m2" {
    source = "../../commons/r2m2"
}
