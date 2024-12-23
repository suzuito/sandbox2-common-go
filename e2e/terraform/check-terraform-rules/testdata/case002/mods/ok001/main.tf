terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  backend "gcs" {
    bucket = "base-420514-terraform"
    prefix = "mods/ok001"
  }

  required_version = ">= 1.10.3"
}

provider "google" {
  project = "base-420514"
}
