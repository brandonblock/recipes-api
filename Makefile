run-api:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	REDIS_URI=localhost:6379 \
	JWT_SECRET=eUbP9shywUygMx7u \
	go run api/main.go
run-web:
	INDEX_PATH=web/index.html \
	go run web/main.go
swagger:
	swagger generate spec -o ./api/swagger.json
	swagger serve -F swagger ./api/swagger.json