
# All targets in this package that do not specify otherwise will be visible to the entire repo.
# This is useful for third-party packages that anyone can use them.
package(default_visibility = ['PUBLIC'])

go_module(
  name = "github.com_stretchr_testify",
  install = ["..."],
  module = "github.com/stretchr/testify",
  version = "v1.7.0",
  deps = [
    ":github.com_davecgh_go-spew",
    ":github.com_pmezard_go-difflib",
    ":github.com_stretchr_objx",
    ":gopkg.in_yaml.v3",
  ],
)

go_module(
  name = "gopkg.in_yaml.v2",
  install = ["..."],
  module = "gopkg.in/yaml.v2",
  version = "v2.4.0",
  deps = [":gopkg.in_check.v1"]
)

go_module(
  name = "github.com_davecgh_go-spew",
  install = ["..."],
  module = "github.com/davecgh/go-spew",
  version = "v1.1.1",
)

go_module(
  name = "github.com_pmezard_go-difflib",
  install = ["..."],
  module = "github.com/pmezard/go-difflib",
  version = "v1.0.0",
)

go_mod_download(
  name = "github.com_stretchr_objx_download",
  module = "github.com/stretchr/objx",
  version = "v0.3.0",
)

go_module(
  name = "github.com_stretchr_objx",
  download = ":github.com_stretchr_objx_download",
  install = ["..."],
  module = "github.com/stretchr/objx",
  strip = ["objx_test"],
)

go_module(
  name = "gopkg.in_yaml.v3",
  install = ["..."],
  module = "gopkg.in/yaml.v3",
  version = "496545a6307b2a7d7a710fd516e5e16e8ab62dbc",
  deps = [":gopkg.in_check.v1"]
)

go_module(
  name = "gopkg.in_check.v1",
  install = ["..."],
  module = "gopkg.in/check.v1",
  version = "10cb98267c6cb43ea9cd6793f29ff4089c306974",
  deps = [
    ":github.com_kr_pretty",
  ]
)

go_module(
  name = "github.com_kr_pretty",
  install = ["..."],
  module = "github.com/kr/pretty",
  version = "v0.3.0",
  deps = [
    ":github.com_kr_text",
    ":github.com_rogpeppe_go-internal",
  ]
)

go_module(
  name = "github.com_kr_text",
  install = ["..."],
  module = "github.com/kr/text",
  version = "v0.2.0",
  deps = [":github.com_creack_pty"]
)

go_module(
  name = "github.com_rogpeppe_go-internal",
  install = ["..."],
  module = "github.com/rogpeppe/go-internal",
  version = "v1.8.0",
  deps = [
    ":github.com_pkg_diff",
    ":gopkg.in_errgo.v2",
  ]
)

go_module(
  name = "github.com_creack_pty",
  install = ["..."],
  module = "github.com/creack/pty",
  version = "v1.1.17",
)

go_module(
  name = "github.com_pkg_diff",
  install = ["..."],
  module = "github.com/pkg/diff",
  version = "20ebb0f2a09e612109b224b32f79370409108bcc",
)

go_mod_download(
  name = "gopkg.in_errgo.v2_download",
  module = "gopkg.in/errgo.v2",
  version = "v2.1.0",
)

go_module(
  name = "gopkg.in_errgo.v2",
  download = ":gopkg.in_errgo.v2_download",
  install = ["..."],
  module = "gopkg.in/errgo.v2",
  strip = [
    "errors/errors_test",
    "fmt/errors/errors_test",
  ],
)
