cd cmd
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o ../artifacts/bootstrap main.go
cd  ../artifacts
zip myFunction.zip bootstrap
aws lambda update-function-code --function-name myFunction \
--zip-file fileb://myFunction.zip