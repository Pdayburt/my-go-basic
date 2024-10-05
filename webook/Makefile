.PHONY: docker mockgen
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s  -o webook .
	@docker rmi -f jack/webbook:v0.0.1
	@docker build -t jack/webbook:v0.0.2 .

mockgen:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/svcmocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/svcmocks/code.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/repomocks/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/repomocks/code.mock.go
	@mockgen -package=redismocks -destination=./internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -source=./internal/service/sms/types.go -package=ratelimitmocks -destination=./internal/service/sms/ratelimitmocks/types.go
	@mockgen -source=./pkg/ratelimit/types.go -package=limitmocks -destination=./pkg/ratelimit/limitmocks/types.go
	@go mod tidy
