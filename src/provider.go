package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider initial used in main.go
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"example_server": resourceServer(),
		},
	}
}
