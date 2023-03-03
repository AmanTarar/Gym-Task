package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// type date struct{

// 	Day string `json:"day"`
// 	Month string `json:"month"`
// 	Year string `json:"year"`
// }


type enrolledInfo struct{

	Name string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	MoneySubmitted float64 `json:"moneySubmitted"`
	Duration int `json:"duration"`
	MemberShipTYPE string `json:"memberShip_Type"`
	Refund float64 `json:"refund"`
	MemberShip string `json:"memberShip"`

}

type MembershipTypePrice struct{

	Gold float64 `json:"gold"`
	Silver float64 `json:"silver"`
}



func CustomerEnrolments(w http.ResponseWriter, r *http.Request) {

		//post request
	

	var newEnrolinput enrolledInfo

	_=json.NewDecoder(r.Body).Decode(&newEnrolinput)
	fmt.Println("newEnrol input ",newEnrolinput)

	//activate the membership
	newEnrolinput.MemberShip="ACTIVE"
	
	//currentTIME
	currentTime:=time.Now().Format("2006-01-02")
	// newEnrolinput.StartDate=currentTime
	newEnrolinput.StartDate="2023-01-02"
	//calculate duration from money submitted
	//    moneysubmitted/membershiptypeprice =no. of months available

	membershiptyp:=newEnrolinput.MemberShipTYPE
	
	if membershiptyp=="gold"{
		newEnrolinput.Duration=int(newEnrolinput.MoneySubmitted/membershipPrice.Gold)
	}else{
		newEnrolinput.Duration=int(newEnrolinput.MoneySubmitted/membershipPrice.Silver)

	}
	//end date 
	//current time + no. of months in duration
	currentTimesliced:=strings.Split(currentTime,"-")
	monthinString:=currentTimesliced[1]
	intmonth, err := strconv.Atoi(monthinString)
	if err!=nil{
		fmt.Println("fatt gya")
	}
	intmonth=intmonth+newEnrolinput.Duration
	//intmonth back to string representation
	monthstr:= strconv.Itoa(intmonth)
	currentTimesliced[1]=monthstr
	var endDATEstring string
	//currentTimesliced back to string
	for i,_:=range currentTimesliced{

		endDATEstring=endDATEstring+currentTimesliced[i]+"-"
	}
	strings.TrimSuffix(endDATEstring,"-")
	

	//end date simplified
	newEnrolinput.EndDate=endDATEstring

	enrolmentdataBASE=append(enrolmentdataBASE,newEnrolinput)
	fmt.Println("new customer enrolled")
	fmt.Println("enrolmentdataBASE",enrolmentdataBASE)
	enrolmentdataBASEinBytes,_:=json.Marshal(enrolmentdataBASE)
	w.Write(enrolmentdataBASEinBytes)

  


}


func SetMembershipPrice(w http.ResponseWriter,r *http.Request){

//function to change or update the price of membership types

	_=json.NewDecoder(r.Body).Decode(&membershipPrice)

	membershipPriceinBytes,_:=json.Marshal(membershipPrice)
	fmt.Fprintf(w,"updated price of membership types\n")
	w.Write(membershipPriceinBytes)

	fmt.Println("membership prices has been changed",membershipPrice)

}


func DeleteMembership(w http.ResponseWriter,r *http.Request){

//delete the membership function
var enroldata enrolledInfo

_=json.NewDecoder(r.Body).Decode(&enroldata)

	for _,v:=range enrolmentdataBASE{

		if enroldata.Name==v.Name{
			//fmt.Println(v.Name,":is deleted from database")
			//enrolmentdataBASE=append(enrolmentdataBASE[:i],enrolmentdataBASE[i+1:]... )
			v.MemberShip="INACTIVE"
			
			
			fmt.Println(v.Name," membership has been deactivated")
			fmt.Println("enrolmentdataBASE",enrolmentdataBASE)

			    //refund calculated price
				currentTime:=time.Now().Format("2006-01-02")
				v.EndDate=currentTime

				currentTimesliced:=strings.Split(currentTime,"-")
				monthinString:=currentTimesliced[1]
				intmonth, err := strconv.Atoi(monthinString)
				if err!=nil{
					fmt.Println("fatt gya")
				}
				//duartion -intmonth
				refundMonth:=v.Duration-intmonth
				v.Duration=v.Duration-refundMonth

				if v.MemberShipTYPE=="gold"{
					v.Refund=(float64(refundMonth)*membershipPrice.Gold)/2

				}else{

					v.Refund=(float64(refundMonth)*membershipPrice.Silver)/2
				}
				fmt.Println("details with refund added",v.Name)

				v_in_Bytes,_:=json.Marshal(v)
				w.Write(v_in_Bytes)
				fmt.Println("",v)
			return
		}
		

	}
	fmt.Println("member not found in database")

}

//All DATABASES
//------------------------------------------>
//dataBASE for enrolments
var enrolmentdataBASE []enrolledInfo

//Global for current membership Price
var membershipPrice MembershipTypePrice


func main() {

	//Default pricing
//------------------------------------------>	
	

	membershipPrice.Gold=2000
	membershipPrice.Silver=1000
//------------------------------------------->
	
	fmt.Println("working...")
	http.HandleFunc("/CustomerEnrolments", CustomerEnrolments)
	http.HandleFunc("/SetMembershipPrice",SetMembershipPrice)
	http.HandleFunc("/DeleteMembership", DeleteMembership)
	http.ListenAndServe(":8080", nil)
	


	
	// fmt.Println("dataBASE",dataBASE)



}