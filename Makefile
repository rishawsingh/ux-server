.PHONY: lint dep-graph
lint:
	golangci-lint run --fix

dep-graph:
	godepgraph -s ./cmd | dot -Tpng -o "$RANDOM.png"
