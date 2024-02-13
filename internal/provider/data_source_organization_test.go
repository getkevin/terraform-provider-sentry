package provider

import (
	"fmt"
	"github.com/getkevin/terraform-provider-sentry/internal/acctest"
)

var testAccOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, acctest.TestOrganization)
