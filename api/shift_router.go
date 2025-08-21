package api

import "github.com/gin-gonic/gin"

func (server *Server) setupShiftsRoutes(baseRouter *gin.RouterGroup) {
	shifts := baseRouter.Group("/locations")
	shifts.Use(server.AuthMiddleware())

	{
		shifts.POST("/:id/shifts", server.RBACMiddleware("SHIFT.CREATE"), server.CreateShiftApi)             // POST /locations/:id/shifts
		shifts.GET("/:id/shifts", server.RBACMiddleware("SHIFT.VIEW"), server.ListShiftByLocationID)         // GET /locations/:id/shifts
		shifts.PUT("/:id/shifts/:shift_id", server.RBACMiddleware("SHIFT.UPDATE"), server.UpdateShiftApi)    // PUT /locations/:id/shifts/:shift_id
		shifts.DELETE("/:id/shifts/:shift_id", server.RBACMiddleware("SHIFT.DELETE"), server.DeleteShiftApi) // DELETE /locations/:id/shifts/:shift_id
	}

}
