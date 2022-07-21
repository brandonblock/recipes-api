run:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	REDIS_URI=localhost:6379 \
	JWT_SECRET=eUbP9shywUygMx7u \
	go run main.go
swagger:
	swagger generate spec -o ./swagger.json
	swagger serve -F swagger ./swagger.json