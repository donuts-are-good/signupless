
![signupless banner logo](https://user-images.githubusercontent.com/96031819/224571551-a577e46c-059a-41b1-a897-339bf7fc6d55.png)

![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors) ![swaggerdoc](https://img.shields.io/swagger/valid/3.0?color=555555&label=API%20&specUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fdonuts-are-good%2Fsignupless%2Fmaster%2Fopenapi.yml&style=for-the-badge)

# 🎟️ **SIGNUPLESS** v1.0.0
Session tokens as a Service

## 🎉 Getting Started

Running a release build doesn't require any build steps and is usually the easiest way.

### 📦 Run a release build


Windows, MacOS, Linux [Download](https://github.com/donuts-are-good/signupless/releases/latest)

### 👨‍💻 Run from source

```
git clone https://github.com/donuts-are-good/signupless.git
cd signupless
go get -d ./...
go run main.go
```

Signupless will run on port 8080 like this: http://localhost:8080/


## 💡 Usage

Signupless offers two endpoints: `/session/add` and `/session/check`.

### **`/session/add`**
The `/session/add` endpoint is used to add a new session. It requires a JSON payload with an `id` field. It returns a session token.

## 📤 Example Request:

```
{
    "id": "user-123"
}
```
## 📫 Example Response:


```
{
    "id": "user-123",
    "token": "f87d319fb7e58b0a439b8f9c17f37d76cf0e6591f27c8d33e45a2c96b6dd26e"
}
```

### **`/session/check`**
The `/session/check` endpoint is used to check the validity of a session token. It requires the session token to be passed in the `session-token` header. If the token is valid, a new token is returned.

## 📤 Example Request:

```
GET /session/check HTTP/1.1
Host: localhost:8080
session-token: f87d319fb7e58b0a439b8f9c17f37d76cf0e6591f27c8d33e45a2c96b6dd26e
```

## 📫 Example Response:

```
{
    "id": "user-123",
    "token": "759e19a7178d9771a6b52eabf1bcb9b8e365db0e325c305d227688a07e2f8f98"
}
```
