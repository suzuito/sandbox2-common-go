provider "google" {
    project = "hoge"
}

module "r3m1" {
    source = "../../commons/r3m1"
}

module "r3m2" {
    source = "../../commons/r3m2"
}
