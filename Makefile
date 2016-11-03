NAME := authz-webhook
REPO := jetstackexperimental/authz-webhook
TAGS := canary

clean:
	rm -f $(NAME)

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(NAME)

image: build
	docker build -t $(REPO):$(TAGS) .

push: image
	true

