FROM mcr.microsoft.com/oss/go/microsoft/golang:1.18 as build
# you'll need to run this build from the root of the repo
ENV GOOS=linux 
ENV GOARCH=amd64 
ENV CGO_ENABLED=0
ADD . /src
WORKDIR /src/internal/eh/stress
RUN go build -o stress .

FROM mcr.microsoft.com/cbl-mariner/base/core:1.0
WORKDIR /app
COPY --from=build /src/internal/eh/stress/stress /app/stress
ENTRYPOINT ["./stress"]
