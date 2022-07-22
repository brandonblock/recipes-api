run:
	make -j 2 run-api run-web-react
run-api:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	REDIS_URI=localhost:6379 \
	JWT_SECRET=eUbP9shywUygMx7u \
	./api/api
run-web-go:
	./web-go/web
run-web-react:
	cd web-react; npm start
build: build-api build-web-go
build-api:
	cd api; go build -o api
build-web-go:
	cd web-go; go build -o web
swagger:
	swagger generate spec -o ./api/swagger.json
	swagger serve -F swagger ./api/swagger.json