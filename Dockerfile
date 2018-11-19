FROM alpine:3.8
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY ./microservice-rating /app/
WORKDIR /app
CMD ["./microservice-rating"]