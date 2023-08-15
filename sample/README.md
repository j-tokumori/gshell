proto から go コード生成
```bash
$ protoc --go_out=./grpc --go_opt=paths=source_relative \
--go-grpc_out=./grpc --go-grpc_opt=paths=source_relative \
sample.proto
```

gshell コード生成
```bash
$ go run ../cmd/gengshell/... -module="github.com/j-tokumori/gshell/sample/grpc" -output=generated.go
```

サンプルサーバ起動
```bash
$ cd server
$ go run main.go
```

grpcurl でチェック
```bash
$ grpcurl -plaintext localhost:8080 list
$ grpcurl -plaintext -d '{"name": "hsaki"}' localhost:8080 service.SampleService.Hello
```

gshell 起動
```bash
$ go run ./main.go
```
