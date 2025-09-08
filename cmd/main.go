package main

import (
	"fmt"
	"log"
	"net"

	domain_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/domain"
	config_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/config"
	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/pb"
	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/server"
	repository_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/repository"
	usecase_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/usecase"
	hashpassword_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/hash_password"
	jwttoken_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/jwt.go"
	randnumgene_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/random_number"
	sendgrid_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/send_grid"
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
	err = db.AutoMigrate(&domain_authSvc.User{},domain_authSvc.OtpInfo{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

	config, err := config_authSvc.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load Config:", err)
	}

	// Initialize repository
	userRepo := repository_authSvc.NewUserRepository(db)
	hashUtils := hashpassword_authSvc.NewHashUtil()
	jwtUtil := jwttoken_authSvc.NewJwtUtil()
	randNumUtil:=randnumgene_authSvc.NewRandomNumUtil()
	sendGrid := sendgrid_authSvc.NewSendGrid(&config.SendGrid)

	// Initialize usecase
	userUseCase := usecase_authSvc.NewUserUseCase(userRepo, hashUtils, jwtUtil, &config.Token, randNumUtil,sendGrid)

	//Create GRPC Server
	grpcServer := grpc.NewServer()

	// Initialize auth service with all dependencies
	authService := server.NewAuthService(userUseCase)

	// Register the service correctly
	pb.RegisterAuthServiceServer(grpcServer, authService)

	listener, err := net.Listen("tcp", config.PortMngr.RunnerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("grpc Auth Service running on %s\n", config.PortMngr.RunnerPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc %v", err)
	}
}
