run-api:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	REDIS_URI=localhost:6379 \
	JWT_SECRET=eUbP9shywUygMx7u \
	./api/api
run-web:
	./web-go/web
build: build-api build-web
build-api:
	cd api; go build -o api
build-web:
	cd web-go; go build -o web
swagger:
	swagger generate spec -o ./api/swagger.json
	swagger serve -F swagger ./api/swagger.json