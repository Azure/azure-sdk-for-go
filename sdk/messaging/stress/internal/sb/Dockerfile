FROM mcr.microsoft.com/oss/go/microsoft/golang:1.23-cbl-mariner2.0 AS build
# you'll need to run this build from the root of the azservicebus module
ENV GOOS=linux 
ENV GOARCH=amd64 
ENV CGO_ENABLED=0

ADD ./azservicebus /src/azservicebus
ADD ./stress /src/stress

WORKDIR /src/stress/internal/sb
RUN go build -o /sb-stress .

WORKDIR /src/stress/internal/sb/tests/benchmarks
RUN go test -o /sb-stress-bench -c

# The first container is just for building the artifacts, and contains all the source, etc...
# That container instance only ever lives on your local machine (or the build machine).
#
# This image pulls from the other container (using COPY --from=build). The resulting image
# will be pushed to our container registry.
FROM mcr.microsoft.com/cbl-mariner/base/core:2.0

COPY --from=build /sb-stress /app/sb-stress
COPY --from=build /sb-stress-bench /app/sb-stress-bench

RUN yum update -y && yum install -y ca-certificates

WORKDIR /app
ENTRYPOINT ["/bin/bash"]
