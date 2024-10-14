
build:
	@echo "Building..."
	@docker build -t concrete .

run:
	@echo "Running..."
	@docker run -it -v ./imgs:/imgs concrete:latest /bin/bash