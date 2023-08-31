FROM mcr.microsoft.com/oss/go/microsoft/golang:1.18-fips-cbl-mariner2.0 as build
# you'll need to run this build from the root of the azservicebus module
ENV GOOS=linux 
ENV GOARCH=amd64 
ENV CGO_ENABLED=0
ADD . /src
WORKDIR /src/internal/stress
RUN go build -o stress .

# The first container is just for building the artifacts, and contains all the source, etc...
# That container instance only ever lives on your local machine (or the build machine).
#
# This image pulls from the other container (using COPY --from=build). The resulting image
# will be pushed to our container registry.
FROM mcr.microsoft.com/cbl-mariner/base/core:1.0
WORKDIR /app
COPY --from=build /src/internal/stress/stress /app/stress
ENTRYPOINT ["/bin/bash"]
