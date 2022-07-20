run:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	go run main.go
swagger:
	swagger generate spec -o ./swagger.json
	swagger serve -F swagger ./swagger.json