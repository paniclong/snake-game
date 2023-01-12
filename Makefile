build_snake:
	 go build -o snake cmd/snake/main.go
build_race:
	go build -race -o snake cmd/snake/main.go