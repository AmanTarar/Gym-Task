package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// type date struct{

// 	Day string `json:"day"`
// 	Month string `json:"month"`
// 	Year string `json:"year"`
// }


type enrol struct{

	Name string `json:"name"`
	MembershipType string `json:"membership"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	MoneySubmitted float64 `json:"moneySubmitted"`
	Duration string `json:"duration"`


}



func CustomerEnrolments(w http.ResponseWriter, r *http.Request) {

		//post request


	var newEnrol enrol

	_=json.NewDecoder(r.Body).Decode(&newEnrol)
	fmt.Println("",newEnrol)
	newEnrolbyte,_:=json.Marshal(newEnrol)
	w.Write(newEnrolbyte)

}
//dataBASE
var dataBASE []enrol



func main() {


	
	
	
	http.HandleFunc("/CustomerEnrolments", CustomerEnrolments)
	http.ListenAndServe(":8080", nil)


	
	fmt.Println("dataBASE",dataBASE)



}