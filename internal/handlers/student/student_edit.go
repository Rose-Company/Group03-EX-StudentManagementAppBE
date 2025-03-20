package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudent(c *gin.Context) {
	log.Println("Handling request: CreateStudent - Creating a new student")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in CreateStudent")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		log.Printf("Error binding request body in CreateStudent: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.Service.Student.CreateAStudent(c.Request.Context(), &student)
	if err != nil {
		log.Printf("Error creating student: %v", err)
		common.AbortWithError(c, err)
		return
	}

	log.Println("Student created successfully")
	c.JSON(http.StatusCreated, common.Response{
		Code:    200,
		Message: "Student created successfully",
	})
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	log.Println("Handling request: UpdateStudent - Updating student")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in UpdateStudent")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		log.Printf("Error binding request body in UpdateStudent: %v", err)
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	updatedStudent, err := h.Service.Student.UpdateStudent(c.Request.Context(), id, &student)
	if err != nil {
		log.Printf("Error updating student with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Student updated successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student updated successfully",
		Data:    updatedStudent,
	})
}

func (h *Handler) DeleteStudentByID(c *gin.Context) {
	log.Println("Handling request: DeleteStudentByID - Deleting student")

	ok, _ := common.ProfileFromJwt(c)
	if !ok {
		log.Println("Unauthorized access attempt in DeleteStudentByID")
		common.AbortWithError(c, common.ErrInvalidToken)
		return
	}

	id := c.Param("id")
	err := h.Service.Student.DeleteByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error deleting student with ID %s: %v", id, err)
		common.AbortWithError(c, err)
		return
	}

	log.Printf("Student deleted successfully with ID: %s", id)
	c.JSON(http.StatusOK, common.Response{
		Code:    200,
		Message: "Student deleted successfully",
	})
}
