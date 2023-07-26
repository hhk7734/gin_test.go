## 초기 세팅

asdf 설치

https://asdf-vm.com/guide/getting-started.html

```shell
make init
```

## Development

```shell
alias log2jq="jq -R -r '. as \$line | try fromjson catch \$line'"
```

```shell
go run ./cmd/server 2>&1 | log2jq
```

### Dependency

새 패키지 추가

```shell
make init
```

새 의존성 주입 추가

```shell
make wire
```
