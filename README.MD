Teknologi Yang Di Gunakan :

- Go Gin Gonic
- GORM
- Mysql
- godotenv

Set up local :

- clone repo
- ubah nama file .env.sample jadi .env
- ubah isi file .env sesuai dengan nama database, username, password, host dan port masing masing
  ```bash
  DB_NAME=db_toko
  DB_USERNAME=root
  DB_PASSWORD=
  DB_PORT=3306
  DB_HOST=127.0.0.1
  PORT_APP=:8080
  ```
- jalankan di terminal go run main.go

  ```bash
  go run main.go

  ```

- URL Documentation API : https://documenter.getpostman.com/view/22397647/2sB3Hhrgqt