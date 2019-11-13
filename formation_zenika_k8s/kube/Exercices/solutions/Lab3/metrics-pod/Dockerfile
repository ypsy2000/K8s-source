FROM golang:1.11 as build
RUN mkdir /app
ADD main.go /app
WORKDIR /app
RUN go mod init server
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /app/server /
EXPOSE 80
CMD ["/server"]