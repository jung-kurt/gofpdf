doc : README.md doc.go doc/index.html.ok

test :
	go test -v

cov :
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

check :
	golint .
	go vet -all .
	gofmt -s -l .

%.html.ok : %.html
	tidy -quiet -output /dev/null $<
	touch $@

doc/body.md README.md doc.go : document.md
	lua doc/doc.lua
	gofmt -s -w doc.go

doc/index.html : doc/hdr.html doc/body.html doc/ftr.html
	cat doc/hdr.html doc/body.html doc/ftr.html > $@

doc/body.html : doc/body.md
	markdown -f +links,+image,+smarty,+ext,+divquote -o $@ $<

clean :
	rm -f coverage.html coverage doc/*.ok doc/body.md README.md doc.go doc/index.html doc/body.html
