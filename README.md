# Go Authentication Server

Simple authentication api written in go inspired by [gotrue](https://github.com/netlify/gotrue)

**Supports email and password login, magic link, oauth**

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
* [/api/v1/auth](#sign-up-new-user)
    * [/signup](#sign-up-new-user) Sign Up New User
    * [/confirm](#confirm-user-mail) 

    * [/signin](#sign-in-user) Sign in User (email or magiclink)
    * [/oauth](#sign-in-with-oauth) Sign in with oauth
    * [/refresh](#refresh-user-token) Refresh user Token
    * [/verify](#verify-magic-link-login) Verify Magic Link Login
    * [/request-reset-password](#send-password-reset-request) Send password reset request
    * [/reset-password](#reset-password) Reset Password

    * [/signout](#sign-out-user) Sign out User

    * [/me](#get-user-data) Get Logged in User Data

    * [/admin](#get-all-users) Admin Endpoints
        * [/users](#get-all-users) Get all Users
        * [/user](#get-user-by-email) Get User By Email
        * [/user](#delete-user) Delete user

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
## Sign in With Oauth
#### GET /api/v1/auth/oauth

**Query**
|          Name | Required |  Type   | Description|
| -------------:|:--------:|:-------:| ----------:|
|     `oauth_provider` | required | string | github |
|     `redirect_to` | required | string  | redirect url |

`Redirects to redirect url with access_token and refresh_token`

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
## Reset password
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
---
---
# Admin
## Get All Users
#### GET /api/v1/auth/admin/users

**Headers**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `x-api-token` | require | string  |

**Query**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `page` | optional, default 1 | int  |
|     `limit` | optional, default 10 | int  |


**Response**
```json
{
    "error": false,
    "limit": 10,
    "page": 1,
    "totalPage": 1,
    "users": [{
        "id": "d2f25e7e-0e5d-49cd-b791-4d8fcabeb073",
        "Name": "Ram",
        "Email": "ram@gmail.com",
    }]
}
```
---
## Get User By Email
#### DELETE /api/v1/auth/admin/user?email=sad@gmail.com

**Headers**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `x-api-token` | require | string  |



**Response**
```json
{
    "error": false,
    "user": {
        "id": "d2f25e7e-0e5d-49cd-b791-4d8fcabeb073",
        "Name": "Ram",
        "Email": "ram@gmail.com",
    }
}
```
---
## Delete User
#### DELETE /api/v1/auth/admin/user/{user_id}

**Headers**

|          Name | Required |  Type   |
| -------------:|:--------:|:-------:| 
|     `x-api-token` | require | string  |



**Response**
```json
{
    "error": false,
    "message": "user deleted"
}
```
---