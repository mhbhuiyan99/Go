
# Learning Database with Golang
Learning in this [tutotial](https://youtu.be/P6AUDnQe360):
1. Connecting to Postgres SQL in Go
2. Running SQL migration from Go
3. Creating Records in Postgres 
4. Fetching Records
5. Paginating your records


## Step - 1 : Connect to Database

### Need:
1. data source name (here 'dsn')
2. Database Library (To connect database) <br>
   We are using Postgres. <br>
   So,<br>
   we need postgress implementation of that single api 
   that go provides a single interface for RDBMS and that
   would enable us communicate with postgres <br>
   So,<br>
   we need to install that library: ```go get github.com/lib/pq``` <br>
```
func connectToDB(dsn string) (*sql.DB, error) {

	// Open connection to database
	db, err := sql.Open("postgres", dsn) // Open(driver name, data source name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
```

## Step - 2 : Connect with Postgres Database

```
func main() {

	dsn := "user=postgres dbname=GoDB password=101010 sslmode=disable"

	_, err := connectToDB(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to DB")
}
```
```go run main.go``` result:
```
sql: unknown driver "postgres" (forgotten import?)
exit status 1
```
**Import:** ```github.com/lib/pq``` <br>

### Create user table

**Method-1:** Create on Postgres
```
CREATE TABLE IF NOT EXISTS users(
	id bigserial primary key,
	name text not null,
	email text not null,
	created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone default current_timestamp
)
```
**Method-2:** using Go
```
func main() {
	...

	q := `CREATE TABLE IF NOT EXISTS users(
		id bigserial primary key,
		name text not null,
		email text not null,
		created_at timestamp with time zone default current_timestamp,
		updated_at timestamp with time zone default current_timestamp
	)`

	_, err = db.Exec(q)
	if err != nil {
		log.Fatalln(err)
	}
}
```
refresh you database on Postgres to check. <br>

Looking bad to see all things are in ```main.go```. <br>

## Structure models


### Step - 1: Create ```struct```
```models/users.go```:
```
type Users struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

### Step - 2 : Create Users Model
To pass the Dependencies <br>
```models/users.go```:
```
type UsersModel struct {
	DB *sql.DB
}
```
## Inserting Data
```
func main() {

	dsn := "user=postgres dbname=GoDB password=101010 sslmode=disable"

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	um := models.UsersModel{DB: db}
	user := models.User{Name: "Mojammel", Email: "mhbhuiyan10023@gmail.com"}

	err = um.Insert(&user)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Created user with ID %d", user.ID)
}
```
Need to add this Insert method. <br>
```models/users.go```:
```
func (m UsersModel) Insert(u *User) error {
	q := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`

	return m.DB.QueryRow(q, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt)
}
```

Run ```go run main.go ```. Result: <br>
```
Created user with ID 1
```
Using Select Query on Postgres database you can see the result.
```select * from users```
<br>

## Fetching Data
If need to see all user data:
```
func main() {

	...

	users, err := um.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	for _, user := range users {
		fmt.Printf("%d: %s %s\n", user.ID, user.Name, user.Email)
	}
}
```
```models/users.go```:
```
func (m UsersModel) GetAll() ([]User, error) {
	var users []User
	q := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
```
Add/insert more data:
It is not the way. But need more data to learn something.
```
func main() {

	...
	um := models.UsersModel{DB: db}
	for i := 0; i < 100; i++ {
		user := models.User{
			Name: fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@temp.com", i),
		}

		err = um.Insert(&user)
		if err != nil {
			log.Fatalln(err)
		}
	}
	...
}
```

## Paginating

### Create model for Filtering
```models/filter.go```:
```
package models

type Filter struct {
	Page 	 int
	PageSize int
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	LastPage     int
	TotalRecords int
}

func ComputeMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	
	return Metadata {
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
```

Updates on ```models/users.go```:
```
func (m UsersModel) GetAll(filter Filter) ([]User, Metadata, error) {
	var users []User
	q := `SELECT COUNT(*) OVER(), id, name, email, created_at, updated_at 
		  FROM users 
		  LIMIT $1 OFFSET $2`

	rows, err := m.DB.Query(q, filter.Limit(), filter.Offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	var totalRec int
	for rows.Next() {
		var user User
		err = rows.Scan(&totalRec, &user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	return users, ComputeMetadata(totalRec, filter.Page, filter.PageSize), nil
}
```

Updates on ```main.go```:
```
func main() {

	...
	um := models.UsersModel{DB: db}

	f := models.Filter{
		PageSize: 5,
		Page:     1,
	}
	users, _, err := um.GetAll(f)
	...
}
```
Now you can see specific pages info.<br>
Try by changing the page number and page size.

## Structuring Models
Replace new model creating code from main().<br>
New file ```models/models.go```:
```
type Models struct {
	Users UsersModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users: UsersModel{DB: db},
	}
}
```
```server.go```:
```
func (app *application) serve() error {
	srv := http.Server{
		Handler: app.handlers(),
		Addr:    ":4000",
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) handlers() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homePage)

	return mux
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request) {		
	w.Write([]byte("Hello"))
}
```
```main.go``` update:
```
type application struct {
	Models models.Models
}

func main() {

	...

	app := application{
		Models: models.NewModel(db),
	}

	fmt.Println("Starting application...")
	err = app.serve()
	if err != nil {
		log.Fatalln(err)
	}
}
```

Using ```go run .``` you can see "Hello"  in http://localhost:4000

### Data in JSON

Updated ```server.go```:
```
func (app *application) homePage(w http.ResponseWriter, r *http.Request) {		
	
	f := models.Filter{
		Page:     1,
		PageSize: 20,
	}

	users, metadata, err := app.Models.Users.GetAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	res := struct {
		Users   []models.User
		Meta    models.Metadata
	}{
		Users: users,
		Meta:  metadata,
	}

	js, err := json.Marshal(res)
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}
```
Now using ```go run .``` you can see data in http://localhost:4000

<br>

Update User struct on ```models/users.go```:
```
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```
