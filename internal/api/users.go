package api

import (
	"log"
	"net/http"
)

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("get user")
}
func (c *Controller) PostUser(w http.ResponseWriter, r *http.Request) {
	log.Println("post user")
}
func (c *Controller) GetUserId(w http.ResponseWriter, r *http.Request, id string) {
	log.Println("get user id")
}
func (c *Controller) PutUserId(w http.ResponseWriter, r *http.Request, id string) {
	log.Println("put user id")
}
func (c *Controller) DeleteUserId(w http.ResponseWriter, r *http.Request, id string) {
	log.Println("delete user id")
}
