FROM alpine:3.4
RUN apk add --update bash
COPY dockerci /
CMD ["./dockerci"]