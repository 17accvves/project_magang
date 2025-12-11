package routes

import (
    "backend/handlers"
    "github.com/gorilla/mux"
)

func SetupRoutes(auth *handlers.AuthHandler, cafe *handlers.CafeHandler) *mux.Router {
    r := mux.NewRouter()

    // ==============================
    // Auth routes
    // ==============================
    r.HandleFunc("/login", auth.Login).Methods("POST")
    r.HandleFunc("/register-cafe", auth.RegisterCafe).Methods("POST")

    // ==============================
    // Admin Cafe routes
    // ==============================
    r.HandleFunc("/approve-cafe", cafe.ApproveCafe).Methods("POST")      // Approve cafe
    r.HandleFunc("/reject-cafe", cafe.RejectCafe).Methods("POST")        // Tolak cafe
    r.HandleFunc("/all-cafes", cafe.ListAllCafes).Methods("GET")         // Ambil semua cafe beserta statusnya

    return r
}
