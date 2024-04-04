# Hands-on Testing

## To work locally with the provider

* modify `~/.terraformrc` on your home directory

```hcl
provider_installation {
  dev_overrides {
    "arielb135/paragon" = "/Users/me/dev/go/path/bin/"
  }
  direct {}
}
```