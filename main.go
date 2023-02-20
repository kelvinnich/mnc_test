package main

import (
	"fmt"
	"log"
	"mnc/controller"
	"mnc/model"
	"mnc/repository"
	"mnc/usecase"
	"net/http"
	"os"
)

var (
	//database temporary
	db = model.NewDatabase()


	//repository
	custrepo repository.CustomersRepository = repository.NewCustomerRepository(db)
	logRepo repository.LogTransactionRepository = repository.NewLogTransaction(db)
	merchantRepo repository.MerchantRepository = repository.NewMerchantRepository(db)

	//usecase
	dataUseCase usecase.DataUseCase = usecase.NewDataUseCase(db)
	merchantUseCase usecase.MerchantUsecase = usecase.NewMerchantUsecase(merchantRepo)


	//controller
	controllerLogin controller.LoginController = controller.NewLoginController(custrepo,dataUseCase)
	paymentController controller.PaymentController = controller.NewPaymentController(db, custrepo, merchantRepo, logRepo, dataUseCase)
	merchantController controller.MerchantController = controller.NewMerchantController(merchantUseCase)
)

func main() {
	// load data
	db, err := dataUseCase.LoadData("customers.json", "merchants.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	// set up handlers
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllerLogin.RegisterHandler(w,r,db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllerLogin.LoginHandler(w, r, db)
	})

	http.HandleFunc("/addMerchant", func(w http.ResponseWriter, r *http.Request) {
		merchantController.CreateMerchant(w,r)
	})


	http.HandleFunc("/payment", func(w http.ResponseWriter, r *http.Request) {
		 paymentController.PaymentHandler(w, r)
	})
	http.HandleFunc("/logout", controllerLogin.LogoutHandler)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
