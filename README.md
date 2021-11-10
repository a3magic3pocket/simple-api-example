# simple-api-example

## 설명

- 간단한 로커 관리 서비스([simple locker](https://isthereanymerch.com))의 api 입니다.

## 의존성 패키지 설치

- ```bash
  go mod tidy
  ```

## 실행

- ```bash
    # 즉시 실행
    go run main.go

    # 빌드 후 실행
    go build -o simple-api
    ./simple-api
  ```

## 테스트

- ```bash
  # 테스트 폴더로 이동
  cd test

  # 모든 테스트 실행
  go test -v

  # 특정 테스트 파일만 실행
  go test main_test.go example1_test.go -v
  ```

## 유의사항

- image tag 변경
  - 새 버전을 배포할 때 simple-api.yml의 image tag 값을 변경해야합니다.  
    (ex. image: a3magic3pocket/simple-api:0.0.4 -> a3magic3pocket/simple-api:0.0.5)
  - registry server(지금은 dockerhub)에 동일한 tag가 존재할 경우, 배포 시 테스트에 실패합니다.
- kubernetes에서 sqlite3 디렉토리는 pv(persistent volume)로 mount 됨
  - kubernetes에서 simple-api-example가 배포될 때에는 sqlite3가 pv으로 mount 됩니다.
  - 따라서 pv, pvc를 먼저 배포 후 simpe-api deployment를 배포해야합니다.

## API 문서

- 배포된 문서
  - [api 문서](https://api.isthereanymerch.com/swagger/index.html)
- 로컬에서 문서 보는 법

  - ```bash
    # controller에서 swag annotation을 수정

    # 문서 생성
    swag init

    # 서버 실행
    go run main go

    # localhost:8080/swagger/index.html에 접속하면 볼 수 있음
    ```

## 관련 레포

- [github - simple-web-example](https://github.com/a3magic3pocket/simple-web-example)
- [github - simple-manifest](https://github.com/a3magic3pocket/simple-manifest)
- [github - simple-image-tag-extraction](https://github.com/a3magic3pocket/simple-image-tag-extraction)
- [docker hub - simple-web](https://hub.docker.com/r/a3magic3pocket/simple-web/tags)
- [docker hub - simple-api](https://hub.docker.com/r/a3magic3pocket/simple-api/tags)
