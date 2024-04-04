# Hands-on Testing

## Getting Started

* modify `~/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "arielb135/paragon" = "/Users/me/dev/go/path/bin/"
  }
  direct {}
}
```

* Build the binary, by running `make` in the root of the git repo.
* Spin up a copy of the paragon service. `docker container run --rm --name paragon-service -p 8080:8080 arielb135/paragon-service`
* Run `terraform init`
* Run whatever other terraform commands you want.

## Cleaning up

* Stop the paragon service. `docker container stop paragon-service`
* Comment out the `dev_overrides` in `~/.terraformrc`.
