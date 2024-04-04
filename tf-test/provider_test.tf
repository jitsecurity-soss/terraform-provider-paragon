terraform {
  required_providers {
    paragon = {
      source = "arielb135/paragon"
    }
  }
}

# Configure the connection details for the Inventory service
provider "paragon" {
  host = "127.0.0.1"
  port = "8080"
}

#Create new Inventory item
resource "paragon_item" "example" {
  name = "Jones Extreme Sour Cherry Warhead Soda"
  tag  = "USD:2.99"
}
