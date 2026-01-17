POST {{BASE_URL}}/users/auth/login
{
  "email": "{{EMAIL}}",
  "password": "foobar123"
}
HTTP 200
[Asserts]
cookie "refreshToken" exists
jsonpath "$.data.access_token" matches "ey\\w.+"
