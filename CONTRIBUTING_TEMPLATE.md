<!-- use this template to generate the contributor docs with the following command: `$ lingo run docs --template CONTRIBUTING_TEMPLATE.md  --output CONTRIBUTING.md` -->
# Contributing Changes

gofpdf is a global community effort and you are invited to make it even better. If you have implemented a new feature or corrected a problem, please consider contributing your change to the project. A contribution that does not directly pertain to the core functionality of gofpdf should be placed in its own directory directly beneath the contrib directory.

Here are guidelines for making submissions. Your change should

* be compatible with the MIT License
* be properly documented
* be formatted with go fmt
* include an example in [fpdf_test.go](https://github.com/jung-kurt/gofpdf/blob/master/fpdf_test.go) if appropriate
* conform to the standards of [golint](https://github.com/golang/lint) and [go vet](https://godoc.org/golang.org/x/tools/cmd/vet), that is, golint . and go vet . should not generate any warnings
* not diminish [test coverage](https://blog.golang.org/cover)

[Pull requests](https://help.github.com/articles/using-pull-requests) work nicely as a means of contributing your changes.

# Code Review Comments and Effective Go Guidelines

{{range .}}
## {{.title}}
{{.body}}
{{end}}