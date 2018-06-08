REM go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
REM go get gopkg.in/urfave/cli.v1

set GOOS=linux
del drycontacts-iptor
go build -o drycontacts-iptor cmd/drycontacts-iptor/drycontacts-iptor.go
del aws-adeunis-drycontacts.zip
build-lambda-zip -o aws-adeunis-drycontacts.zip drycontacts-iptor
REM "C:\Program Files\7-Zip\7z.exe" a aws-adeunis-drycontacts.zip cmd\ frame\ -xr!*test*