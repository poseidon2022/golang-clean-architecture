package router

import (
	"golang-clean-architecture/delivery/controllers"
	"golang-clean-architecture/infrastructure"
	"golang-clean-architecture/repository"
	usecase "golang-clean-architecture/use_cases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(db *mongo.Database, router *gin.Engine) {
	publicRouter := router.Group("")

	NewSignUpRouter(db, publicRouter)
	NewLoginRouter(db, publicRouter)

	privateRouter := router.Group("")
	privateRouter.Use(infrastructure.AuthMiddleWare())
	NewTaskRouter(db, privateRouter)
	EscalatePrevilige(db, privateRouter)
}

func EscalatePrevilige(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, "users")
	uc := &controllers.UserController{
		UserUseCase : usecase.NewUserUseCase(ur),
	}

	group.PUT("/promote/:id", uc.PromoteUser())
}

func NewLoginRouter(db *mongo.Database, group *gin.RouterGroup) {
	//here we should make the appropriate invocations to the controller function and
	//instantiate the userUseCase usecase and pass it as an argument. uc.register => uc.login
	//but before that we have to assign somethings to the uc struct
	//usercontroller.somestruct.taskRepository setup the db and context here
	ur := repository.NewUserRepository(db, "users")
	uc := &controllers.UserController {
		UserUseCase : usecase.NewUserUseCase(ur),
	}
	group.POST("/login", uc.Login())
}

func NewSignUpRouter(db *mongo.Database, group *gin.RouterGroup) {

	ur := repository.NewUserRepository(db, "users")
	uc := &controllers.UserController{
		UserUseCase : usecase.NewUserUseCase(ur),
	}
	group.POST("/register", uc.Register())
}

func NewTaskRouter(db *mongo.Database, group *gin.RouterGroup) {
	//now we prepare a task controller function that returns a handler when it is called
	tr := repository.NewTaskRepository(db, "tasks")
	tc := &controllers.TaskController{
		TaskUseCase: usecase.NewTaskUseCase(tr),
	}
	group.POST("/tasks", tc.PostTask())
	group.GET("/tasks", tc.GetTasks())
	group.GET("/tasks/:id", tc.GetTask())
	group.PUT("/tasks/:id", tc.UpdateTask())
	group.DELETE("/tasks/:id", tc.DeleteTask())
}