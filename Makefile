# Build and push the image to a local kind cluster regsitry
kind-build:
	docker build -t localhost:5001/joy-generator:latest .
	docker push localhost:5001/joy-generator:latest
