#run this locally so that you do not need to restore all the time the external dependencies (go get)
#of course, in a CI/CD environment, you need to change this approach 
export CGO_ENABLED=0
go build -o microservice-rating .
