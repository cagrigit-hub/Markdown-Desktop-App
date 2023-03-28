BINARY_NAME=Markdown.app
APP_NAME=Markdown
VERSION=1.0.0

build:
	rm -rf $(BINARY_NAME)
	rm -f MARKDOWN-APP
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

run:
	go run .


clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf $(BINARY_NAME)
	@echo "Clean up complete."

test:
	go test -v ./...