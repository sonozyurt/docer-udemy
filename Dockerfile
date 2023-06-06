From golang:alpine as builder

workdir "usr/src/app"

copy ./go.mod ./go.sum ./
run go mod download && go mod verify

copy . .

run go build -v -o "/usr/src/app/build" ./...

from nginx 
copy --from=builder /usr/src/app/build /usr/share/nginx/html