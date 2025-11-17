FROM golang:1.25.4-alpine

WORKDIR /app

RUN useradd --create-home gopher \
  && chown -R gopher:gopher /app

USER gopher

COPY --chown=gopher:gopher . .