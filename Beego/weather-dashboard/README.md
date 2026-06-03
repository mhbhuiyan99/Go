# Weather Dashboard (Beego + Goroutines + HTML Templates)

## Project Overview

This project is a simple Weather Dashboard built using the Beego framework in Go. It demonstrates how goroutines can be used to perform multiple tasks concurrently and how data can be rendered on a web page using Go's HTML templates.

The application simulates fetching weather information for multiple cities and displays the results in a table.

---

## Features

* Built with Beego Framework
* Uses Goroutines for concurrent processing
* Uses sync.WaitGroup for synchronization
* Uses HTML Templates for frontend rendering
* Demonstrates MVC architecture

---

## Technologies Used

* Go (Golang)
* Beego v2
* HTML Templates
* Goroutines
* WaitGroup

---

## Project Structure

```text
weather-dashboard/
│
├── controllers/
│   └── home.go
│
├── models/
│   └── weather.go
│
├── routers/
│   └── router.go
│
├── views/
│   └── index.tpl
│
├── conf/
│   └── app.conf
│
├── main.go
│
├── go.mod
└── README.md
```

---

## How It Works

1. The user visits the home page.
2. The controller creates a list of cities.
3. A separate goroutine is started for each city.
4. Each goroutine simulates fetching weather data.
5. A WaitGroup waits until all goroutines finish.
6. The collected weather data is sent to the HTML template.
7. The template displays the weather information in a table.

---

## Installation

### Clone the Repository

```bash
git clone <repository-url>
cd weather-dashboard
```

### Install Dependencies

```bash
go mod tidy
```

### Run the Application

```bash
bee run
```

or

```bash
go run main.go
```

---

## Open in Browser

Visit:

```text
http://localhost:8080
```

---

## Sample Output

```text
Weather Dashboard

City         Temperature
--------------------------------
Dhaka         31.52 °C
Khulna        28.91 °C
Rajshahi      34.17 °C
Tangail       29.44 °C
```

