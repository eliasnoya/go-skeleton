FROM golang:1.20.5-alpine3.18

WORKDIR /tuples

ENV TZ=UTC

COPY . .
RUN go get
# build API
RUN go build -o ./bin/api
# build CLI
RUN go build -o .\cli.exe .\cli\cli.go

EXPOSE 8000

CMD ["/tuples/bin/api"]
