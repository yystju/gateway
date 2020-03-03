FROM ubuntu
EXPOSE 1234
WORKDIR /app
COPY gateway gateway
COPY config.toml config.toml
CMD ./gateway
