swagger:
	swagger generate spec -o ./swagger.json
	swagger serve -F swagger ./swagger.json