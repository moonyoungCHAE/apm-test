* 8080 : java spring application
* 8081 : golang gin application

---

GET localhost:8080/ping, localhost:8081/ping
* function trace 확인 (+argument)

---
GET localhost:8080/golang -> GET localhost:8081/golang

---

GET localhost:8080/golang/db -> GET localhost:8081/golang/db -> mysql data 조회

