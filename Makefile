# Build and push the image to a local kind cluster regsitry
kind-build:
	docker build -t localhost:5001/joy-generator:latest .
	docker push localhost:5001/joy-generator:latest

fmt:
	goimports --local github.com/nestoca/joy-generator -w .

test:
	@INTERNAL_TESTING=true \
	CATALOG_URL=https://github.com/nestoca/catalog \
	GH_USER=nestobot \
	GH_TOKEN=$(shell gh auth token) \
	REGISTRY=northamerica-northeast1-docker.pkg.dev \
	go test ./... -p 1 -v -race
