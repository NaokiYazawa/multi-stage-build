## exec create user
.PHONY: create_user
create_user:
	curl -X POST "http://localhost:8080/users" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{\"name\":\"Jack\"}"

## exec get users
.PHONY: get_users
get_users:
	curl -X GET "http://localhost:8080/users" -H  "accept: application/json"

## exec get user
.PHONY: get_user
get_user:
	curl -X GET "http://localhost:8080/users/1" -H  "accept: application/json"

## exec update user by name
.PHONY: update_user
update_user:
	curl -X PUT "http://localhost:8080/users/1" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{\"name\":\"updated-Jack\"}"

## exec delete user
.PHONY: delete_user
delete_user:
	curl -X DELETE "http://localhost:8080/users/1" -H  "accept: application/json" -H  "Content-Type: application/json"
