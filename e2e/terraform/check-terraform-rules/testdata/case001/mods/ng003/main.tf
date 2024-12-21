terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  backend "gcs" {
    bucket = "hoge-terraform"
    prefix = "mods/ng003"
  }

  required_version = ">= 1.10.3"
}

provider "google" {
  project = "base-420514"
}
