package api

import (
	"maicare_go/async"
	"maicare_go/notification"
	"maicare_go/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// handleHealthCheck returns basic service health information
func (server *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, TestResponse{
		Status:    "healthy",
		Message:   "service is running",
		Timestamp: time.Now(),
	})
}

// handleEcho returns whatever JSON payload was sent
func (server *Server) handleEcho(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Status:    "error",
			Message:   "invalid JSON payload",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"echo":      requestBody,
		"timestamp": time.Now(),
	})
}

// handleLatency simulates a delayed response
func (server *Server) handleLatency(c *gin.Context) {
	ms := c.Param("ms")
	delay, err := time.ParseDuration(ms + "ms")
	if err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Status:    "error",
			Message:   "invalid latency value",
			Timestamp: time.Now(),
		})
		return
	}

	time.Sleep(delay)
	c.JSON(http.StatusOK, TestResponse{
		Status:    "success",
		Message:   "delayed response",
		Timestamp: time.Now(),
	})
}

func (server *Server) EmailAndAsynq(c *gin.Context) {
	server.asynqClient.EnqueueEmailDelivery(async.EmailDeliveryPayload{
		Name:         "Farjia Taha",
		To:           "farjiataha@gmail.com",
		UserEmail:    "farjiataha@gmail.com",
		UserPassword: "password",
	}, c)

	c.JSON(http.StatusOK, gin.H{
		"echo":      "Email sent",
		"timestamp": time.Now(),
	})
}

// NotificationResponse is the response structure for the notification endpoint
type NotificationResponse struct {
	Echo      string    `json:"echo"`
	Timestamp time.Time `json:"timestamp"`
}

// Notification to test Notification delivery .
// @Summary: Test Notification delivery
// @Description: Test Notification delivery
// @Tags: Test
// @Accept: json
// @Produce: json
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} Response[any]
// @Router /test/notification [get]
// @Security -
func (server *Server) Notification(c *gin.Context) {
	// payload, err := GetAuthPayload(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }
	// // Get the user ID from the payload
	// // _ := payload.UserId
	// // Enqueue the notification task

	server.asynqClient.EnqueueNotificationTask(c, notification.NotificationPayload{
		RecipientUserIDs: []int64{1},
		Type:             "new_client_assigned",
		Data: notification.NotificationData{
			NewClientAssignment: &notification.NewClientAssignmentData{
				ClientID:        12345,
				ClientFirstName: "John",
				ClientLastName:  "Doe",
				ClientLocation:  util.StringPtr("test Location"), // Assuming no location provided
			},
		},
		Message:   "You have been assigned a new client: John Doe",
		CreatedAt: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{
		"echo":      "Notification sent",
		"timestamp": time.Now(),
	})
}

// NotificationAppointement to test Notification delivery .
// @Summary: Test Notification delivery
// @Description: Test Notification delivery
// @Tags: Test
// @Accept: json
// @Produce: json
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} Response[any]
// @Router /test/notification_2 [get]
// @Security -
func (server *Server) NotificationAppointement(c *gin.Context) {
	// payload, err := GetAuthPayload(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }
	// // Get the user ID from the payload
	// // _ := payload.UserId
	// // Enqueue the notification task

	server.asynqClient.EnqueueNotificationTask(c, notification.NotificationPayload{
		RecipientUserIDs: []int64{1},
		Type:             "new_appointment",
		Data: notification.NotificationData{
			NewAppointment: &notification.NewAppointmentData{
				AppointmentID: uuid.New(),
				CreatedBy:     "Admin User",
				StartTime:     time.Now().Add(24 * time.Hour),
				EndTime:       time.Now().Add(25 * time.Hour),
				Location:      "Office",
			},
		},
		Message:   "You have been assigned a new appointment: John Doe",
		CreatedAt: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{
		"echo":      "Notification sent",
		"timestamp": time.Now(),
	})
}

// NotificationNewSchedule to test Notification delivery .
// @Summary: Test Notification delivery
// @Description: Test Notification delivery
// @Tags: Test
// @Accept: json
// @Produce: json
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} Response[any]
// @Router /test/notification_3 [get]
// @Security -
func (server *Server) NotificationNewSchedule(c *gin.Context) {
	// payload, err := GetAuthPayload(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }
	// // Get the user ID from the payload
	// // _ := payload.UserId
	// // Enqueue the notification task

	server.asynqClient.EnqueueNotificationTask(c, notification.NotificationPayload{
		RecipientUserIDs: []int64{1},
		Type:             "new_schedule_notification",
		Data: notification.NotificationData{
			NewScheduleNotification: &notification.NewScheduleNotificationData{
				ScheduleID: uuid.New(),
				CreatedBy:  2,
				StartTime:  time.Now().Add(24 * time.Hour),
				EndTime:    time.Now().Add(25 * time.Hour),
				Location:   "Office",
			},
		},
		Message:   "You have been assigned a new appointment: John Doe",
		CreatedAt: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{
		"echo":      "Notification sent",
		"timestamp": time.Now(),
	})
}

// NotificationNewIncident to test Notification delivery .
// @Summary: Test Notification delivery
// @Description: Test Notification delivery
// @Tags: Test
// @Accept: json
// @Produce: json
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} Response[any]
// @Router /test/notification_4 [get]
// @Security -
func (server *Server) NotificationNewIncident(c *gin.Context) {
	// payload, err := GetAuthPayload(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }
	// // Get the user ID from the payload
	// // _ := payload.UserId
	// // Enqueue the notification task

	server.asynqClient.EnqueueNotificationTask(c, notification.NotificationPayload{
		RecipientUserIDs: []int64{1},
		Type:             "incident_report",
		Data: notification.NotificationData{
			NewIncidentReport: &notification.NewIncidentReportData{
				ID:                 4,
				EmployeeID:         2,
				EmployeeFirstName:  "Jane",
				EmployeeLastName:   "Doe",
				LocationID:         3,
				LocationName:       "Main Office",
				ClientFirstName:    "John",
				ClientLastName:     "Doe",
				SeverityOfIncident: "High",
			},
		},
		Message:   "You have been assigned a new incident report: 4",
		CreatedAt: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{
		"echo":      "Notification sent",
		"timestamp": time.Now(),
	})
}
