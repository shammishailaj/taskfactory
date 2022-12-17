FROM shammishailaj/gobuilder:0.0.6 AS builder

ARG CGO
ARG OS
ARG ARCH

RUN mkdir /go/src/taskfactory

WORKDIR /go/src/taskfactory

COPY . ./

RUN go env

#RUN echo "CGO = "$CGO
#RUN echo "OS = "$OS
#RUN echo "ARCH = "$ARCH
RUN CGO_ENABLED=$CGO GOOS=$OS GOARCH=$ARCH govvv build -mod vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/taskfactory ./cmd/taskfactory/taskfactory.go
RUN cp ./bin/taskfactory /bin/

# See: https://stackoverflow.com/a/52979541/6670698
FROM scratch AS production
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/taskfactory/configs /taskfactory/configs/
COPY --from=builder /go/src/taskfactory/web /taskfactory/web/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/src/taskfactory/bin/taskfactory /taskfactory/taskfactory
WORKDIR /taskfactory
ENV TZ Local
ENTRYPOINT ["/taskfactory/taskfactory","serve"]