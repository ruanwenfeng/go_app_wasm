build:
	@GOARCH=wasm GOOS=js go build -o app.wasm ./app
	@go build -o demo ./server
	@go build -o dicom ./testdicom
	@go build -o downloadcsv ./downloadcsv
	@go build -o testchan ./testchan

dicom:build
	./dicom

csv:build
	./downloadcsv

run: build
	PORT=8001 ./demo

testchan: build
	./testchan

deploy: build
	gcloud app deploy --project=murlok