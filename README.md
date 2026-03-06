# Library App

This is a Library Management System written in Go.

# Tech Stack
* [Go (Echo Framework)](https://echo.labstack.com/)
* [Supabase](https://supabase.com/)
* [Echo JWT Middleware](https://github.com/labstack/echo-jwt)
* [Golang JWT](https://github.com/golang-jwt/jwt)
* [GORM](https://gorm.io/)

# Deployment URL

This project is deployed here: https://library-production-76e6.up.railway.app

# API Documentation

This documentation describes the API endpoints for the **Library** management system. The API handles user authentication, book browsing, rentals, and administrative reporting.

The documentation can also be found via the root "/" path, as the app serves an [index.html](./index.html) containing the API Docs.

## Authentication

Private endpoints require a **Bearer Token**. Log in to receive a token, you need this script in the API Client:
```bash
if (res.status == 200) {
      bru.setVar("user_token", res.body.token) # I use bruno, make sure to adjust if you use other API Client
}
```

Then it should be included in the `Authorization` header for subsequent requests.
- **Type**: Bearer Token
- **Variable**: `{{user_token}}`

---

## Endpoints

Below are the endpoints in this app.

### User Management

#### 1. Register
Create a new user account.
* **Method:** `POST`
* **URL:** `/users/register`
* **Body (JSON):**
    ```json
    {
      "full_name": "string",
      "email": "user@example.com",
      "password": "password"
    }
    ```

#### 2. Login
Authenticate a user and receive a session token.
* **Method:** `POST`
* **URL:** `/users/login`
* **Body (JSON):**
    ```json
    {
      "email": "user@example.com",
      "password": "password"
    }
    ```
* **Note:** The `user_token` is automatically captured from the response for use in other requests.

---

### Book & Rental Operations

#### 3. Get Books
Retrieve a list of all books.
* **Method:** `GET`
* **URL:** `/books`

#### 4. Rent Book
Rent a specific book for a defined duration.
* **Method:** `POST`
* **URL:** `/books/rent`
* **Body (JSON):**
    ```json
    {
      "book_id": 19,
      "duration": 8
    }
    ```

#### 5. Get My Rents
Retrieve the rental history for the currently authenticated user.
* **Method:** `GET`
* **URL:** `/users/rent`

#### 6. Return Book
Mark a rented book as returned.
* **Method:** `POST`
* **URL:** `/books/return/:book_id`

---

### Administration

#### 7. Admin Get Rents Report
Get a report of users rent counts in descending order.
* **Method:** `GET`
* **URL:** `/admin/rent`

#### 8. Admin Get Author Report
Get a report of authors book counts in descending order.
* **Method:** `GET`
* **URL:** `/admin/authors`