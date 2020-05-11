# REALMS

REALMS Establishes A Library Management System, written in Go, using a MySQL database.

## Getting started

### 0. Prerequisities

To set up the environment, you need to have the following dependencies installed.

- [Go](https://golang.org/dl) 1.14 or above
- [GNU make](https://www.gnu.org/software/make) 4.0 or above

For Windows, try [MinGW-w64](https://sourceforge.net/projects/mingw-w64).

### 1. Installation

First, you need to obtain the REALMS package.

```bash {.line-numbers}
git clone https://github.com/hakula139/REALMS.git
cd REALMS
```

Then you can build the project using Make.

```bash {.line-numbers}
make build
```

For MinGW-w64 on Windows, use the command below.

```bash {.line-numbers}
mingw32-make build
```

You should see the following output, which indicates a successful installation.

```text {.line-numbers}
build: realms done.
build: realmsd done.
```

### 2. Usage

#### 2.1 realmsd

Run `realmsd` using the command below, and the server will listen to port `7274` by default.

```bash {.line-numbers}
./bin/realmsd
```

#### 2.2 realms

To interact with the back end, here's a simple CLI tool, namely, `realms`. Though, it's not necessarily required, since you can easily build another front end with the RESTful APIs, a guide to which will be provided later.

Run `realms` using the command below.

```bash {.line-numbers}
./bin/realms
```

You should see the welcome message.

```text {.line-numbers}
Welcome to REALMS! Check the manual using the command 'help'.
>
```

To get started with REALMS, use the command `help` to show all available commands.

```text {.line-numbers}
> help
```

```text {.line-numbers}
COMMANDS:
   Public:
      help           Shows a list of commands
      exit           Quit

      login          Log in to your library account
      logout         Log out of your library account
      status         Shows the current login status

      show books     Shows all books in the library
      show book      Shows the book of given ID
      find books     Finds books by title / author / ISBN

   Admin privilege required:
      add book       Adds a new book to the library
      update book    Updates data of a book
      remove book    Removes a book from the library

      add user       Adds a new user to the database
      update user    Updates data of a user
      remove user    Removes a user from the database
      show users     Shows all users in the library
      show user      Shows the user of given ID

   User privilege required:
      me             Shows the current logged-in user

      borrow book    Borrows a book from the library
      return book    Returns a book to the library
      check ddl      Checks the deadline to return a book
      extend ddl     Extends the deadline to return a book
      show list      Shows all books that you've borrowed
      show overdue   Shows all overdue books that you've borrowed
      show history   Shows all records
```

It's quite easy to understand how these commands work, nevertheless we're going to talk about them in the next chapter.

### 3. REST API

Here we'll demonstrate the usage of these RESTful APIs by example.

#### 3.1 Log in

##### 3.1.1 Request

Method: `POST /login`  
Content-Type: `multipart/form-data`  
CLI command: `login`

You'll be required to enter your username and password (FYI, the password is invisible while typing). There's no `signup` in REALMS, while a user account can only be acquired from an admin.

To authenticate a user's credentials, REALMS uses session. In the implementation of the CLI tool, a session cookie is used to store the essential information, and the cookies are handled by [cookiejar](https://golang.org/pkg/net/http/cookiejar).

On the server-side, the password will be hashed using [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) before save.

##### 3.1.2 Response

Status: `200 OK`  
Content-Type: `application/json`  

```json {.line-numbers}
{"data": true}
```

In case of a successful login, you'll receive a welcome message.

```text {.line-numbers}
Welcome Hakula!
```

Otherwise, an error will be returned. Possible error messages are shown below.

```text {.line-numbers}
cli: failed to make an http request, did you start realmsd?
auth: user not exist
auth: incorrect password
auth: already logged in
auth: failed to save session
```

#### 3.2 Log out

##### 3.2.1 Request

Method: `GET /logout`  
CLI command: `logout`

##### 3.2.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{"data": true}
```

In case of a successful logout, you'll receive a success message.

```text {.line-numbers}
Successfully logged out!
```

Otherwise, an error will be returned. Possible error messages are shown below.

```text {.line-numbers}
cli: failed to make an http request, did you start realmsd?
auth: invalid session token, have you logged in?
auth: failed to save session
```

#### 3.3 Show current logged-in user

##### 3.3.1 Request

Method: `GET /user/me`  
CLI command: `me`

**User** privilege is required, which means you have to login before doing this operation.

##### 3.3.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{"data": 3}
```

Normally, your user ID will be returned.

```text {.line-numbers}
Current user ID: 3
```

If you're not logged in, you'll receive an error message below.

```text {.line-numbers}
auth: unauthorized
```

#### 3.4 Show the current login status

##### 3.4.1 Request

Method: `GET /status`  
CLI command: `status`

##### 3.4.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{"data": true}
```

We expect the following output.

```text {.line-numbers}
Online
Offline
```

#### 3.5 Add a new book

##### 3.5.1 Request

Method: `POST /admin/books`  
Content-Type: `application/json`  
CLI command: `add book`

```json {.line-numbers}
{
  "title": "CS:APP",
  "author": "Randal E. Bryant",
  "publisher": "Pearson",
  "isbn": "978-0134092669"
}
```

**Admin** privilege is required. In REALMS, we use `level` to indicate a user's privilege, which is a property of the user model. When a user makes a request, the server will check if he/she has admin privilege. If not, an Unauthorized Error will be returned.

You'll be required to input the necessary information of the book, and the `title` field should not be blank, or an error will be returned.

On the server-side, the following message will be written to log using [zap](https://pkg.go.dev/mod/go.uber.org/zap). The default path to the log file is `./logs/realmsd.log`.

```json {.line-numbers}
{"level":"info","time":"2020-05-04T01:19:11.206+0800","msg":"Added book 20"}
```

##### 3.5.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{
  "data": {
    "id": 20,
    "title": "CS:APP",
    "author": "Randal E. Bryant",
    "publisher": "Pearson",
    "isbn": "978-0134092669"
  }
}
```

The complete information of the added book will be returned, since you may want to display it in your front-end application. For the sake of simplicity, here the CLI tool will just print the book ID.

```text {.line-numbers}
Successfully added book 20
```

If not authorized (i.e. you're not an admin), you'll receive an error message below.

```text {.line-numbers}
auth: unauthorized
```

#### 3.6 Update data of a book

##### 3.6.1 Request

Method: `PATCH /admin/books/:id`  
Content-Type: `application/json`  
CLI command: `update book`

```json {.line-numbers}
{
  "title": "Computer Systems",
  "author": "Randal E. Bryant, David R. O'Hallaron",
}
```

**Admin** privilege is required.

Here `:id` refers to the book ID. Simply sending a request including the fields that you want to update is fine, and empty values will be omitted. Still, there's an input checker for all inputs on the server-side, which will validate your request body to prevent invalid requests.

The following message will be written to log.

```json {.line-numbers}
{"level":"info","time":"2020-05-04T01:41:57.908+0800","msg":"Updated book 20"}
```

##### 3.6.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{
  "data": {
    "id": 20,
    "title": "Computer Systems",
    "author": "Randal E. Bryant, David R. O'Hallaron",
    "publisher": "Pearson",
    "isbn": "978-0134092669"
  }
}
```

The updated data will be returned.

```text {.line-numbers}
Successfully updated book 20
```

Possible error messages are shown below.

```text {.line-numbers}
auth: unauthorized
database: book not found
```

#### 3.7 Remove a book

##### 3.7.1 Request

Method: `DELETE /admin/books/:id`  
Content-Type: `application/json`  
CLI command: `remove book`

```json {.line-numbers}
{"message": "Book lost"}
```

**Admin** privilege is required.

The `message` field is optional, which is the explanation why you remove the book.

The following message will be written to log.

```json {.line-numbers}
{"level":"info","time":"2020-05-04T02:13:00.956+0800","msg":"Removed book 20 with explanation: Book lost"}
```

Or if there's no explanation:

```json {.line-numbers}
{"level":"info","time":"2020-05-04T02:13:00.956+0800","msg":"Removed book 20"}
```

##### 3.7.2 Response

Status: `200 OK`  
Content-Type: `application/json`

```json {.line-numbers}
{"data": true}
```

Since the book has already been removed, there's no need to return its information.

```text {.line-numbers}
Successfully removed book 20
```

Possible error messages are shown below.

```text {.line-numbers}
auth: unauthorized
database: book not found
```

## TODO

- [x] Add a simple CLI front-end
- [ ] Add unit tests
- [ ] Add a detailed document

## Contributors

- [**Hakula Chen**](https://github.com/hakula139)<[i@hakula.xyz](mailto:i@hakula.xyz)> - Fudan University

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](./LICENSE) file for details.
