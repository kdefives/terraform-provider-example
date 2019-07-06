package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// example.Server represents a concrete Go type that represents an API resource
func TestAccExampleServer_basic(t *testing.T) {
	var serverBefore, serverAfter example.Server
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExampleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceExists("example_server.foo", &serverBefore),
				),
			},
			{
				Config: testAccExampleResource_removedPolicy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceExists("example_server.foo", &serverAfter),
				),
			},
		},
	})
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("EXAMPLE_KEY"); v == "" {
		t.Fatal("EXAMPLE_KEY must be set for acceptance tests")
		if v := os.Getenv("EXAMPLE_SECRET"); v == "" {
			t.Fatal("EXAMPLE_SECRET must be set for acceptance tests")
		}
	}
}

// testAccCheckExampleResourceDestroy verifies the Server
// has been destroyed
func testAccCheckExampleResourceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	conn := testAccProvider.Meta().(*ExampleClient)

	// loop through the resources in state, verifying each server
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "example_server" {
			continue
		}

		// Retrieve our server by referencing it's state ID for API lookup
		request := &example.DescribeServers{
			IDs: []string{rs.Primary.ID},
		}

		response, err := conn.DescribeServers(request)
		if err == nil {
			if len(response.Servers) > 0 && *response.Servers[0].ID == rs.Primary.ID {
				return fmt.Errorf("Server (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		// If the error is equivalent to 404 not found, the server is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "Server not found") {
			return err
		}
	}
	return nil
}
