go_library(
  name = "repo",
  srcs = ["repo.go"],
  deps = ["//third_party/go:gopkg.in_yaml.v2"],
  visibility = [
    "//pkg/...",
    "//src/...",
  ],
)

go_test(
  name = "repo_test",
  srcs = ["repo_test.go"],
  deps = [
    ":repo", # shorthand for the `//src/repo:repo` rule in the same package
    "//third_party/go:github.com_stretchr_testify",
  ],
)
