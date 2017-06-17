# gin-swagger - dry templates for go-swagger

## Generate code

```
./gin-swagger -A cluster-registry -f swagger.yaml
```

## Features
* [ ] Validate + bind input to gin ctx.
  * [x] bind body input.
    * [ ] [Nice to have] custom input Models with required fields.
  * [x] bind params input.
  * [x] bind query params input.
  * [ ] consume more than `application/json`
* [ ] Security.
  * [ ] basic
  * [ ] apiKey
  * [ ] OAuth2
    * [x] `password` (Bearer token)
    * [ ] `accessCode`
    * [ ] `application`
    * [ ] `implicit`
  * [ ] Auth chain
  * [x] Custom authorize on individual routes.
* [x] Set custom middleware on each router Pre/post 'main' handler.
  * Use case pre: *custom authorization pre handler*.
  * Use case post: *audit report events post handler*.
* [ ] Ginize generated code.
* [ ] Set and get user info (uid, realm) from gin context.
* [ ] Response helper functions
