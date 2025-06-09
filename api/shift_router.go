package api

import "github.com/gin-gonic/gin"

func (server *Server) setupShiftsRoutes(baseRouter *gin.RouterGroup) {
	shifts := baseRouter.Group("/locations")
	shifts.Use(AuthMiddleware(server.tokenMaker))

	{
		shifts.POST("/:id/shifts", RBACMiddleware(server.store, "SHIFT.CREATE"), server.CreateShiftApi)             // POST /locations/:id/shifts
		shifts.GET("/:id/shifts", RBACMiddleware(server.store, "SHIFT.VIEW"), server.ListShiftByLocationID)         // GET /locations/:id/shifts
		shifts.PUT("/:id/shifts/:shift_id", RBACMiddleware(server.store, "SHIFT.UPDATE"), server.UpdateShiftApi)    // PUT /locations/:id/shifts/:shift_id
		shifts.DELETE("/:id/shifts/:shift_id", RBACMiddleware(server.store, "SHIFT.DELETE"), server.DeleteShiftApi) // DELETE /locations/:id/shifts/:shift_id
	}

}
