FROM golang:1.22-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/dog_defs

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/dog_defs .
COPY --from=build /app/templates ./templates

CMD ["./dog_defs"]