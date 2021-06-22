.PHONY : all clean test cover docker

PACKAGES =	./helper ./config ./src/member/v1/delivery ./src/member/v1/model ./src/member/v1/query ./src/member/v1/usecase \
			./src/services ./src/shared/model \


app-linux: main.go
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $@

test:
	$(foreach pkg, $(PACKAGES), \
	go test $(pkg);)

app-image:
	docker build --build-arg SSH_PRIVATE_KEY=$(cp ~/.ssh/id_rsa .) -t member .