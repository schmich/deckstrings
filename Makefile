test:
	go test

doc:
	@echo http://localhost:8888/pkg/github.com/schmich/deckstrings/
	godoc -http :8888

.PHONY: test doc
