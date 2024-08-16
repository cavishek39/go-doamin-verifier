package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	/* scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("domain, hasMX, hasSPF, SPFRecord, hasDMARC, DMARCRecord\n")

	for scanner.Scan() {
		domain := scanner.Text()
		checkDomain(domain) 
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading input: ", err)
	} */

	r := Router() // call 	the router
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r)) // Run the server
}

type Email struct {
	Email string `json:"email"`
	HasMX bool `json:"hasMX"`
	HasSPF bool `json:"hasSPF"`
	SPFRecord string `json:"spfRecord"`
	HasDMARC bool `json:"hasDMARC"`
	DMARCRecord string `json:"dmarcRecord"`
}

func checkDomain(domain string) Email {
	if domain == "" {
		return Email{}
	}

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	hasMXRecord, err := net.LookupMX(domain)

	if err != nil {
		fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, false, false, "", false, "")
		return Email{}
	}

	if len(hasMXRecord) > 0 {
		hasMX = true
	}

	hasSPFRecord, err := net.LookupTXT(domain)

	if err != nil {
		fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMX, false, "", false, "")
		return Email{}
	}

	if len(hasSPFRecord) > 0 {
		hasSPF = true
		fmt.Println(hasSPFRecord)
		spfRecord = strings.Join(hasSPFRecord, ", ")
	}

	hasDMARCRecord, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMX, hasSPF, spfRecord, false, "")
		return Email{}
	}

	if len(hasDMARCRecord) > 0 {
		hasDMARC = true
		dmarcRecord = strings.Join(hasDMARCRecord, ", ")
	}

	fmt.Printf("END RESULT => ", "%s, %t, %t, %s, %t, %s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)

	return Email{
		Email: domain,
		HasMX: hasMX,
		HasSPF: hasSPF,
		SPFRecord: spfRecord,
		HasDMARC: hasDMARC,
		DMARCRecord: dmarcRecord,
	}
}

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

func GetEmailVerifier(w http.ResponseWriter, r *http.Request) {
	// Get the email from the request body
	var req EmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	var res response

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		res = response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:   nil,
		}
		json.NewEncoder(w).Encode(res)
		return 
	}
	
	email := req.Email
	fmt.Println("email => ", email)

	// Check if the email is valid
	if email == "" {
		fmt.Printf("Email is required")
		res = response{
			Status:  http.StatusBadRequest,
			Message: "Email is required",
			Data:   nil,
		}
		json.NewEncoder(w).Encode(res)
		return
		// panic("Email is required")
	}

	// Check the domain
     result :=	checkDomain(email)

	 if result.Email == "" {
		res = response{
			Status:  http.StatusBadRequest,
			Message: "Invalid email",
			Data:   nil,
		} 
	 } else {
			res = response{
				Status:  http.StatusOK,
				Message: "Email verified successfully",
				Data:    result,
			}
		}



	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)


	
}

func Router () * mux.Router{
	router := mux.NewRouter() // create a new router

	router.HandleFunc("/api/email-verifier", GetEmailVerifier).Methods("POST", "OPTIONS") // get all email-verifier

	return router
}