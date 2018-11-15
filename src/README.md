1) How to run

Go to the directory containing main and do

1) go build
2) go run main.go

There are 2 API's,

a) http://localhost:8000/cert/{domain name}
This returns the certificate if it is already present and if not creates a new one and returns
b)http://localhost:8000/certs
It returns all the certs that are in the cert store at the moment.
