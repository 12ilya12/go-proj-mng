FROM golang:alpine  
WORKDIR /app  
COPY go.mod ./  
RUN go mod download  
COPY . .  
RUN go build -o main .  
EXPOSE 8001
CMD ["./main"]  