set GOPATH=%cd%

::call get-dependencies.bat

set GOOS=darwin
set GOARCH=386
go build -o analyzePortfolio-darwin-386 analyzer/analyzePortfolio.go 


set GOOS=linux
set GOARCH=386
go build -o analyzePortfolio-linux-386 analyzer/analyzePortfolio.go 


set GOOS=solaris
set GOARCH=amd64
go build -o analyzePortfolio-solaris-amd64 analyzer/analyzePortfolio.go

set GOOS=windows
set GOARCH=386
go build -o analyzePortfolio-windows-386.exe analyzer/analyzePortfolio.go

