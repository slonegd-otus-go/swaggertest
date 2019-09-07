# проба swagger

всё отсюда https://ops.tips/blog/a-swagger-golang-hello-world/

интереснее https://posener.github.io/openapi-intro/#example

## команды для проверки
список curl localhost:8080/pets/

один curl localhost:8080/pets/0

добавить curl -d '{"Kind":"dog", "Name":"Roma"}' -H "Content-Type: application/json" -X POST http://localhost:8080/pets/