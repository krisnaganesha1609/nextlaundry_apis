
# NEXTLAUNDRY API

Simple RESTFUL API Endpoint By Golang, Gin, JWT, and GORM for NEXTLAUNDRY WEB APPS. As the owner Of NEXTLAUNDRY, you can do thing such as generate and get the reports of transaction of the outlet you owned, or try all the CRUD operations as admin instead. There is also a cashier role, if you interested in joining us and be our employee to input the transaction :)



## API Reference
Open and see the routes/routers.go to clarify all the available route endpoints! User localhost as startpoint. If you change the Router Run Port in the main.go file, you must also change the port on your localhost link.
```http
if (r.Run(":8000")) {
    http://127.0.0.1:8000
} else if (r.Run(":8080")) {
    http://127.0.0.1:8080
}
  
```
### API ENDPOINTS

```http
  POST /api/auth                             Login
```

#### MUST PASS THE MIDDLEWARE ADMIN
##### Logout Handler

```http
  POST /api/nextlaundry/admin/logout         Logout
```

##### Users Resource Data
```http
  GET /api/nextlaundry/admin/users           Get All Users
```
```http
  GET /api/nextlaundry/admin/users/:id       Get Detailed User
```
```http
  POST /api/nextlaundry/admin/user           Create New User
```
```http
  PUT /api/nextlaundry/admin/user/:id        Update User
```
```http
  DELETE /api/nextlaundry/admin/user         Delete User
```

##### Outlet Resource Data
```http
  GET /api/nextlaundry/admin/outlets           Get All Outlets
```
```http
  GET /api/nextlaundry/admin/outlets/:id       Get Detailed Outlet
```
```http
  POST /api/nextlaundry/admin/outlet           Create New Outlet
```
```http
  PUT /api/nextlaundry/admin/outlet/:id        Update Outlet
```
```http
  DELETE /api/nextlaundry/admin/outlet         Delete Outlet
```
##### Member Resource Data
```http
  GET /api/nextlaundry/admin/members           Get All Members
```
```http
  GET /api/nextlaundry/admin/members/:id       Get Detailed Member
```
```http
  POST /api/nextlaundry/admin/member           Create New Member
```
```http
  PUT /api/nextlaundry/admin/member/:id        Update Member
```
```http
  DELETE /api/nextlaundry/admin/member         Delete Member
```
##### Package Resource Data
```http
  GET /api/nextlaundry/admin/products           Get All Products
```
```http
  GET /api/nextlaundry/admin/products/:id       Get Detailed Product
```
```http
  POST /api/nextlaundry/admin/product           Create New Product
```
```http
  PUT /api/nextlaundry/admin/product/:id        Update Product
```
```http
  DELETE /api/nextlaundry/admin/product         Delete Product
```

Check the endpoints for other roles by visiting routes/routers.go file!

## Installation

Simply clone or pull this repository and create your very own database. Then follow these 5 easy step :

```bash
  1. cmd: cd nextlaundry_apis || cd <your-project>
  2. cmd: go get -u
  3. open models/setup/setup.go file
  4. rename the value of dbName
  5. cmd: go run main.go
  6. Have Fun!
```
    
## Why Create This Project?

To Fulfill The Software Engineering Competention Exam - Paket 3 (Sistem Pengelolaan Laundry) - Web-Based Application

