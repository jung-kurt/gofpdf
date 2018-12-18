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


## Comment First Word as Subject
Doc comments work best as complete sentences, which allow a wide variety of automated presentations.
The first sentence should be a one-sentence summary that starts with the name being declared.


## Package Comment
Every package should have a package comment, a block comment preceding the package clause. 
For multi-file packages, the package comment only needs to be present in one file, and any one will do. 
The package comment should introduce the package and provide information relevant to the package as a 
whole. It will appear first on the godoc page and should set up the detailed documentation that follows.


## Single Method Interface Name
By convention, one-method interfaces are named by the method name plus an -er suffix 
or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.

There are a number of such names and it's productive to honor them and the function names they capture. 
Read, Write, Close, Flush, String and so on have canonical signatures and meanings. To avoid confusion, 
don't give your method one of those names unless it has the same signature and meaning. Conversely, 
if your type implements a method with the same meaning as a method on a well-known type, give it the 
same name and signature; call your string-converter method String not ToString.


## Avoid Annotations in Comments
Comments do not need extra formatting such as banners of stars. The generated output
may not even be presented in a fixed-width font, so don't depend on spacing for alignment—godoc, 
like gofmt, takes care of that. The comments are uninterpreted plain text, so HTML and other 
annotations such as _this_ will reproduce verbatim and should not be used. One adjustment godoc 
does do is to display indented text in a fixed-width font, suitable for program snippets. 
The package comment for the fmt package uses this to good effect.


## Context as First Argument
Values of the context.Context type carry security credentials, tracing information, 
deadlines, and cancellation signals across API and process boundaries. Go programs 
pass Contexts explicitly along the entire function call chain from incoming RPCs 
and HTTP requests to outgoing requests.

Most functions that use a Context should accept it as their first parameter.


## Do Not Discard Errors
Do not discard errors using _ variables. If a function returns an error, 
check it to make sure the function succeeded. Handle the error, return it, or, 
in truly exceptional situations, panic.


## Go Error Format
Error strings should not be capitalized (unless beginning with proper nouns 
or acronyms) or end with punctuation, since they are usually printed following
other context. That is, use fmt.Errorf("something bad") not fmt.Errorf("Something bad"),
so that log.Printf("Reading %s: %v", filename, err) formats without a spurious 
capital letter mid-message. This does not apply to logging, which is implicitly
line-oriented and not combined inside other messages.


## Use Crypto Rand
Do not use package math/rand to generate keys, even 
throwaway ones. Unseeded, the generator is completely predictable. 
Seeded with time.Nanoseconds(), there are just a few bits of entropy. 
Instead, use crypto/rand's Reader, and if you need text, print to 
hexadecimal or base64


## Avoid Renaming Imports
Avoid renaming imports except to avoid a name collision; good package names
should not require renaming. In the event of collision, prefer to rename the
most local or project-specific import.

