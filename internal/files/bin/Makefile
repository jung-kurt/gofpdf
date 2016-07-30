TRG= ../files.go
JSON= ../../../font/calligra.json
Z= ../../../font/calligra.z

${TRG} : ${JSON} ${Z} ./bin
	echo "package files" >  ${TRG}
	echo "" >> ${TRG}
	echo "// CalligraJson is embedded byte slice for calligra.json" >> ${TRG}
	echo "var CalligraJson = []byte{" >> ${TRG}
	./bin < ${JSON} >> ${TRG}
	echo "}" >> ${TRG}
	echo "" >> ${TRG}
	echo "// CalligraZ is embedded byte slice for calligra.z" >> ${TRG}
	echo "var CalligraZ = []byte{" >> ${TRG}
	./bin < ${Z} >> ${TRG}
	echo "}" >> ${TRG}
	gofmt -s -w ${TRG}

./bin : bin.go
	go build -v

clean :
	rm -f ./bin ${TRG}

