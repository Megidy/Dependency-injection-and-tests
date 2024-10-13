package api

import (
	"database/sql"
	"log"

	"github.com/API/services/order"
	"github.com/API/services/product"
	"github.com/API/services/user"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := gin.Default()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(router)

	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(orderStore, productStore, userStore)
	orderHandler.RegisterRoutes(router)

	log.Println("started server on 8080 ")
	return router.Run()

}
