FROM golang as build
# you'll need to run this build from the root of the repo
ENV GOOS=linux 
ENV GOARCH=amd64 
ENV CGO_ENABLED=0
ADD . /src
WORKDIR /src/internal/stress
RUN go build -o stress .

FROM alpine
RUN mkdir -p /app
COPY --from=build /src/internal/stress/stress /app/stress
WORKDIR /app
ENTRYPOINT ["./stress"]