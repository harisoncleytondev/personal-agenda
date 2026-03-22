package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harisoncleytondev/personal-agenda/internal/api/routes/dto"
	"github.com/harisoncleytondev/personal-agenda/internal/model"
	"github.com/harisoncleytondev/personal-agenda/internal/service"
)

type AppointmentHandle struct {
	serviceAppointment *service.AppointmentService
}

func NewAppointmentHandle(serviceAppointment *service.AppointmentService) *AppointmentHandle {
	return &AppointmentHandle{serviceAppointment: serviceAppointment}
}

func (s *AppointmentHandle) Create(c *gin.Context) {
	var body dto.CreateAppointmentRequest
	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userI, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	user := userI.(*model.User)

	err := s.serviceAppointment.CreateAppointment(c.Request.Context(), user.ID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Agendamento criado com sucesso!"})
}

func (s *AppointmentHandle) Update(c *gin.Context) {
	id := c.Param("id")
	var body dto.UpdateAppointmentRequest
	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userI, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	user := userI.(*model.User)

	err := s.serviceAppointment.UpdateAppointment(c.Request.Context(), id, user.ID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agendamento atualizado com sucesso!"})
}

func (s *AppointmentHandle) Delete(c *gin.Context) {
	id := c.Param("id")

	userI, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	user := userI.(*model.User)

	err := s.serviceAppointment.DeleteAppointment(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agendamento deletado com sucesso!"})
}

func (s *AppointmentHandle) GetAll(c *gin.Context) {
	userI, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	user := userI.(*model.User)

	appointments, err := s.serviceAppointment.GetAllByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (s *AppointmentHandle) GetByID(c *gin.Context) {
	id := c.Param("id")

	userI, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	user := userI.(*model.User)

	appointment, err := s.serviceAppointment.GetByID(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}