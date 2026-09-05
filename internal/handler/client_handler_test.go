package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterClientRoutesIncludesEmergencyContactDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	noop := func(ctx *gin.Context) { ctx.Next() }

	RegisterClientRoutes(router.Group(""), NewClientHandler(nil), noop, func(string) gin.HandlerFunc {
		return noop
	})

	for _, route := range router.Routes() {
		if route.Method == "DELETE" && route.Path == "/clients/:id/emergency_contacts/:contact_id" {
			return
		}
	}

	t.Fatal("emergency contact DELETE route is not registered")
}

func TestRegisterClientRoutesIncludesMainCoordinatorDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	noop := func(ctx *gin.Context) { ctx.Next() }

	RegisterClientRoutes(router.Group(""), NewClientHandler(nil), noop, func(string) gin.HandlerFunc {
		return noop
	})

	for _, route := range router.Routes() {
		if route.Method == "DELETE" && route.Path == "/clients/:id/coordinator" {
			return
		}
	}

	t.Fatal("main coordinator DELETE route is not registered")
}
