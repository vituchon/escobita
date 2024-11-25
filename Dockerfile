FROM golang:1.20 as builder

WORKDIR /app

COPY . .

EXPOSE 9090

CMD ["./escobita"]
