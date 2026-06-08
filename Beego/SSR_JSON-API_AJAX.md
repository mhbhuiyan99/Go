# SSR vs JSON API vs AJAX

## Why This Topic Matters

Many beginners confuse:

```text
SSR
JSON API
AJAX
```

and often build applications where everything is an API or everything is server-rendered.

Modern web applications commonly combine all three.

Understanding the difference is a fundamental web development skill.

---

# 1. SSR (Server-Side Rendering)

Server-Side Rendering means the server generates a complete HTML page and sends it to the browser.

Example request:

```http
GET /countries
```

Server:

```go
func (c *CountryController) Get() {
    c.Data["Countries"] = countries
    c.TplName = "countries.tpl"
}
```

Response:

```html
<html>
<head>
    <title>Countries</title>
</head>
<body>
    <h1>Countries</h1>
    ...
</body>
</html>
```

The browser receives a complete page and displays it.

---

## SSR Flow

```text
Browser
    ↓
GET /countries
    ↓
Server
    ↓
Render Template
    ↓
HTML
    ↓
Browser
```

---

## Examples

```text
GET /
GET /countries
GET /countries/bangladesh
GET /dashboard
GET /wishlist
```

These routes usually return:

```text
HTML
```

not JSON.

---

## When to Use SSR

Use SSR when:

* User navigates to a page
* Browser needs a complete document
* SEO matters
* Initial page load

Example:

```text
User clicks:
    /countries/bangladesh

Server returns:
    destination.tpl
```

---

# 2. JSON API

An API route returns data instead of HTML.

Example:

```http
GET /api/countries
```

Response:

```json
[
    {
        "name": "Bangladesh",
        "capital": "Dhaka"
    }
]
```

No HTML.

Only data.

---

## API Flow

```text
Client
    ↓
GET /api/countries
    ↓
Server
    ↓
JSON
    ↓
Client
```

---

## Examples

```text
GET    /api/countries
GET    /api/countries/:slug
GET    /api/wishlist
POST   /api/wishlist
PUT    /api/wishlist/:id
DELETE /api/wishlist/:id
```

These routes usually return:

```text
JSON
```

not HTML.

---

## When to Use APIs

Use APIs when:

* JavaScript needs data
* Mobile apps need data
* Frontend frameworks need data
* AJAX requests need data

---

# 3. AJAX

AJAX is a technique that allows JavaScript to communicate with the server without reloading the page.

AJAX itself is not a route.

AJAX usually calls JSON API endpoints.

---

## Example

User is already on:

```text
/countries
```

and types:

```text
Bangladesh
```

JavaScript:

```javascript
fetch("/api/countries?search=bangladesh")
```

Server returns:

```json
[
    {
        "name": "Bangladesh"
    }
]
```

JavaScript updates:

```html
<div id="country-results">
```

without reloading the page.

---

## AJAX Flow

```text
Browser
    ↓
JavaScript
    ↓
GET /api/countries
    ↓
JSON
    ↓
Update Part of Page
```

Notice:

```text
No page reload
```

---

# SSR vs AJAX

## SSR

User visits:

```text
/countries
```

Browser receives:

```html
Entire Page
```

---

## AJAX

User searches:

```text
Bangladesh
```

Browser receives:

```json
Search Results
```

Only a small section updates.

---

# TravelSphere Example

## Step 1: SSR

User opens:

```text
/countries
```

Server returns:

```text
countries.tpl
```

Complete HTML page.

---

## Step 2: AJAX

User types:

```text
Bangladesh
```

JavaScript calls:

```http
GET /api/countries?search=bangladesh
```

Server returns:

```json
[
    {
        "name": "Bangladesh"
    }
]
```

Only:

```html
#country-results
```

changes.

---

## Step 3: SSR Again

User clicks:

```html
<a href="/countries/bangladesh">
```

Browser performs normal navigation.

Server returns:

```text
destination.tpl
```

Full page.

---

# Common Beginner Mistakes

## Mistake 1

Returning HTML from API routes.

Bad:

```text
GET /api/countries
```

returns:

```html
<html>...</html>
```

API routes should return:

```json
{}
```

or

```json
[]
```

---

## Mistake 2

Returning JSON from SSR routes.

Bad:

```text
GET /countries
```

returns:

```json
[]
```

SSR routes should return:

```html
<html>...</html>
```

---

## Mistake 3

Using AJAX for page navigation.

Bad:

```javascript
fetch("/countries/bangladesh")
```

to load an entire page.

Normal page navigation should use:

```html
<a href="/countries/bangladesh">
```

---

# Rule of Thumb

Ask:

### Is the user navigating to a page?

Use:

```text
SSR Route
```

Example:

```text
/
/countries
/countries/bangladesh
/wishlist
/dashboard
```

---

### Is JavaScript requesting data?

Use:

```text
JSON API
```

Example:

```text
/api/countries
/api/wishlist
/api/dashboard/summary
```

---

### Do I want to update only part of the page?

Use:

```text
AJAX
```

which usually calls:

```text
JSON API endpoints
```

---

# Mental Model

Think of it this way:

```text
SSR
    = Build a house

JSON API
    = Deliver materials

AJAX
    = Bring materials into one room
      without rebuilding the whole house
```

A modern web application typically works like this:

```text
SSR
    ↓
Initial Page Load

AJAX
    ↓
Calls JSON APIs

JSON APIs
    ↓
Provide Data

JavaScript
    ↓
Updates Part of Page
```

This combination is the foundation of most web applications built today.
