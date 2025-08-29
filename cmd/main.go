package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models"
	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/pb"
	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/server"
	repository_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/repository"
	usecase_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/usecase"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=root dbname=auth_service port=5432 sslmode=disable"

	//connect to DB

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error Connecting to DB: %v", err)
	}
	fmt.Println("Database connected")
		// Auto-migrate tables
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

	// Initialize repository
	userRepo :=repository_authSvc.NewUserRepository(db)

	// Initialize usecase
	userUseCase:=usecase_authSvc.NewUserUseCase(userRepo)

	//Create GRPC Server
	grpcServer := grpc.NewServer()

	// Initialize auth service with all dependencies
	authService := server.NewAuthService(userUseCase)

	// Register the service correctly
	pb.RegisterAuthServiceServer(grpcServer,authService)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("grpc Auth Service running on %s\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc %v", err)
	}
}
