From golang:alpine as builder

workdir "usr/src/app"

copy ./go.mod ./go.sum ./
run go mod download && go mod verify

copy . .

run GOOS=linux go build  -v -o /app 


from alpine:latest as deployment
expose 80
workdir /usr/local/bin/myapp
copy --from=builder /app /usr/local/bin/myapp/app
copy ./templates ./templates

cmd ["./app"]

