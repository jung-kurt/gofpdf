TRG= ../files.go
JSON= ../../../font/calligra.json
Z= ../../../font/calligra.z

${TRG} : ${JSON} ${Z} ./bin
	echo "package files" >  ${TRG}
	echo "" >> ${TRG}
	echo "var CalligraJson = []byte{" >> ${TRG}
	./bin < ${JSON} >> ${TRG}
	echo "}" >> ${TRG}
	echo "" >> ${TRG}
	echo "var CalligraZ = []byte{" >> ${TRG}
	./bin < ${Z} >> ${TRG}
	echo "}" >> ${TRG}
	gofmt -w ${TRG}

./bin : bin.go
	go build -v

clean :
	rm -f ./bin ${TRG}

