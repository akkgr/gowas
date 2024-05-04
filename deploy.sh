cd src
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o ../artifacts/bootstrap main.go
cd  ../artifacts
zip myFunction.zip bootstrap
aws lambda create-function --function-name myFunction \
--runtime provided.al2023 --handler bootstrap \
--architectures arm64 \
--role arn:aws:iam::761059477267:role/lambdaRole \
--zip-file fileb://myFunction.zip