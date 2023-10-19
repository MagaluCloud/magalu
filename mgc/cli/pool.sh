set -xe
go run main.go object-storage buckets create --name "test"
go run main.go object-storage objects upload --src test.txt --dst s3://test/a
go run main.go object-storage objects upload --src test.txt --dst s3://test/b
go run main.go object-storage objects upload --src test.txt --dst s3://test/c
go run main.go object-storage objects upload --src test.txt --dst s3://test/d
go run main.go object-storage objects upload --src test.txt --dst s3://test/e
go run main.go object-storage objects upload --src test.txt --dst s3://test/f
go run main.go object-storage objects upload --src test.txt --dst s3://test/g
go run main.go object-storage objects upload --src test.txt --dst s3://test/h
go run main.go object-storage objects upload --src test.txt --dst s3://test/i
go run main.go object-storage objects upload --src test.txt --dst s3://test/j
go run main.go object-storage objects upload --src test.txt --dst s3://test/k
go run main.go object-storage objects upload --src test.txt --dst s3://test/l
go run main.go object-storage objects upload --src test.txt --dst s3://test/m
go run main.go object-storage objects upload --src test.txt --dst s3://test/n
go run main.go object-storage objects upload --src test.txt --dst s3://test/o
go run main.go object-storage objects upload --src test.txt --dst s3://test/p
go run main.go -l "info:*" object-storage buckets delete --name test -o yaml
