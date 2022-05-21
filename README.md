# Go Authentication Server

Simple authentication api written in go inspired by [gotrue](https://github.com/netlify/gotrue)

#### Requirements
- [Golang](https://go.dev/)

#### Libraries
- [Gin](https://gin-gonic.com) - Go Web Framework
- [GORM](https://gorm.io) - ORM library
- [Uber-Fx](https://github.com/uber-go/fx) - Dependency Injection Framework
- [Simple-Mail](https://github.com/xhit/go-simple-mail) - Mail Client

### Usage :
Copy `sample.env` file as `.env` in the root directory and edit all the values.
___

### **Debug build**

Runs server without building
``` console
user@main:~$ go run cmd/main.go
```
Build and run server
``` console
user@main:~$ go build cmd/main.go
user@main:~$ ./main
```
___
#### Release build
``` console
user@main:~$ go build cmd/main.go
user@main:~$ GIN_MODE=release ./main
```
After running the commands, Authentication server runs on port `8080`
___
### **Test API**
```console
user@main:~$ make test
```
`or`
```console
user@main:~$ go test ./...
```
## API Endpoints
---
## Sign up new user
#### POST /api/v1/auth/signup

**Parameters**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `name` | required | string |
|     `email` | required | string  |
|     `password` | required | string  |

**Response**
```json
{
    "error": false,
    "message": "account created"
}
```
---
## Sign in user
#### POST /api/v1/auth/signin

**Query**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `type` | required | string | email or magiclink |
|     `redirect_to` | required | string  | redirect url |

**Parameters**

|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `email` | required | string  | email|
|     `password` | optional | string  | required if login type is email |

**Response**
```json
{
    "error": false,
    "access_token": "JWT",
    "refresh_token": "JWT",
}
```
---
## Get user data
#### POST /api/v1/auth/me

**Headers**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `X-Access-Token` | required | string | access token |

**Response**
```json
{
    "error": false,
    "user": {
        "id": "0c57528e-4d4f-4737-b6e9-da1f251ce8b7",
        "name": "My Name",
        "email": "email123@gmail.com"
    }
}
```
---
## Refresh user token
#### POST /api/v1/auth/refresh

**Headers**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `X-Refresh-Token` | required | string | refresh token |

**Response**
```json
{
    "error": false,
    "access_token": "JWT"
}
```
---
## Verify magic link login
#### GET /api/v1/auth/verify
#### POST /api/v1/auth/verify

**Query**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `type` | required | string | type of verification ( magiclink ) |
|     `token` | required | string | token |
|     `redirect_to` | required | string | redirect url |

**Response**

- Redirects to `redirect_url` with `access_token` and `refresh_token` as query.
Example: http://localhost:3000?access_token=JWT&refresh_toke=JWT
---
## Sign out user
#### POST /api/v1/auth/signout

**Headers**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `X-Refresh-Token` | required | string | refresh token |

**Response**
```json
{
    "error": false,
    "message": "successfully signed out"
}
```
---
## Send password reset request
#### POST /api/v1/auth/request-password-reset

**Parameters**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `email` | required | string  |

**Response**
```json
{
    "error": false,
    "message": "password recovery email sent"
}
```
---
## Reset request
#### POST /api/v1/auth/reset-password

**Query**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `token` | required | string  |

**Parameters**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `password` | required | string  |

**Response**
```json
{
    "error": false,
    "message": "password reset successfully"
}
```
---
## Confirm user mail
#### GET /api/v1/auth/confirm
#### POST /api/v1/auth/confirm

**Query**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `token` | required | string  |

**Response**
```json
{
    "error": false,
    "message": "user mail confirmed"
}
```