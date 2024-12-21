terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  backend "gcs" {
    bucket = "base-420514-terraform"
    prefix = "mods/ng002"
  }

  required_version = ">= 1.10.3"
}
