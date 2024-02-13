package sentry

import (
	"fmt"
	"github.com/getkevin/terraform-provider-sentry/internal/acctest"
	sentry "github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSentryOrganizationDataSource_basic(t *testing.T) {
	var organization sentry.Organization

	rn := "data.sentry_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryOrganizationExists(rn, &organization),
					resource.TestCheckResourceAttrSet(rn, "name"),
					resource.TestCheckResourceAttr(rn, "slug", acctest.TestOrganization),
					resource.TestCheckResourceAttrWith(rn, "internal_id", func(v string) error {
						want := sentry.StringValue(organization.ID)
						if v != want {
							return fmt.Errorf("got organization ID %s; want %s", v, want)
						}
						return nil
					}),
				),
			},
		},
	})
}

var testAccSentryOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, acctest.TestOrganization)
