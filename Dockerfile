FROM alpine:3.4
RUN apk add --update bash
COPY main /
CMD ["./main"]