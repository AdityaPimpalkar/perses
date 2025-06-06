FROM alpine AS build-env
RUN apk add --update --no-cache mailcap
RUN mkdir /perses
RUN mkdir /plugins

FROM gcr.io/distroless/static-debian12

LABEL maintainer="The Perses Authors <perses-team@googlegroups.com>"

USER nobody

COPY --chown=nobody:nobody perses                            /bin/perses
COPY --chown=nobody:nobody percli                            /bin/percli
COPY --chown=nobody:nobody LICENSE                           /LICENSE
COPY --chown=nobody:nobody plugins-archive/                  /etc/perses/plugins-archive/
COPY --chown=nobody:nobody docs/examples/config.docker.yaml  /etc/perses/config.yaml
COPY --from=build-env --chown=nobody:nobody                  /perses         /perses
COPY --from=build-env --chown=nobody:nobody                  /plugins        /etc/perses/plugins
COPY --from=build-env --chown=nobody:nobody                  /etc/mime.types /etc/mime.types

WORKDIR /perses

EXPOSE     8080
VOLUME     ["/perses"]
ENTRYPOINT [ "/bin/perses" ]
CMD        ["--config=/etc/perses/config.yaml", \
            "--log.level=error"]
