package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// type date struct{

// 	Day string `json:"day"`
// 	Month string `json:"month"`
// 	Year string `json:"year"`
// }


type enrolledInfo struct{
	gorm.Model

	Id             string  `json:"id" gorm:"default:uuid_generate_v4()"`
	Name string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	MoneySubmitted float64 `json:"moneySubmitted"`
	Duration int `json:"duration"`
	MemberShipTYPE string `json:"memberShip_Type"`
	Refund float64 `json:"refund"`
	MemberShip string `json:"memberShip"`

}


type MembershipTypePrice struct {
		Name  string `json:"name"`
		Price float64 `json:"price"`
	
}
	
//Global for current membership Price
var membershiptypeprice = map[string]float64{"Gold": 2000, "Silver": 1000}






func CustomerEnrolments(w http.ResponseWriter, r *http.Request) {

	//post request

	var newEnrolinput enrolledInfo
	_=json.NewDecoder(r.Body).Decode(&newEnrolinput)
	
	

	// _=json.NewDecoder(r.Body).Decode(&newEnrolinput)
	// fmt.Println("newEnrol input ",newEnrolinput)

	//activate the membership
	newEnrolinput.MemberShip="ACTIVE"
	
	//currentTIME
	currentTime:=time.Now().Format("2006-01-02")
	// newEnrolinput.StartDate=currentTime
	newEnrolinput.StartDate="2023-01-02"
	//calculate duration from money submitted
	//    moneysubmitted/membershiptypeprice =no. of months available



var memtypeprice MembershipTypePrice
 
// monthly price and duration set krna hai


// fmt.Println("price: ", price.Price)
// member.Duration = member.MoneySubmitted / price.Price
// member.MonthlyPrice = price.Price

	
	if newEnrolinput.MemberShipTYPE=="gold"{
		db.Where("name = ?", "Gold").First(&memtypeprice)
		newEnrolinput.Duration=int(newEnrolinput.MoneySubmitted/memtypeprice.Price)
	}else{
		db.Where("name = ?", "Silver").First(&memtypeprice)
		newEnrolinput.Duration=int(newEnrolinput.MoneySubmitted/memtypeprice.Price)

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


	result:=db.Create(&newEnrolinput)
	if result.Error!=nil{
		fmt.Println("error in DB")
	}

	
	// newEnrolinputinBytes,_:=json.Marshal(newEnrolinput)
	// fmt.Fprintf(w,"updated price of membership types\n")
	// w.Write(newEnrolinputinBytes)


	// enrolmentdataBASE=append(enrolmentdataBASE,newEnrolinput)
	// fmt.Println("new customer enrolled")
	// fmt.Println("enrolmentdataBASE",enrolmentdataBASE)
	// enrolmentdataBASEinBytes,_:=json.Marshal(enrolmentdataBASE)
	// w.Write(CustomerEnrolmentsData)

  


}

// func MembershipPriceUpdateHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Now you can update membership prices!!")

// 	var memPrice Price
// 	err := json.NewDecoder(r.Body).Decode(&memPrice)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	DB.Model(&Price{}).Where("name=?", memPrice.Name).Updates(&memPrice)
// 	json.NewEncoder(w).Encode(memPrice)

// 	fmt.Fprint(w, "Membership prices updated")

// }


func SetMembershipPrice(w http.ResponseWriter,r *http.Request){

//function to change or update the price of membership types

	//fmt.Println("prices will update soon")
	fmt.Fprintln(w, "Now you can update membership prices!!")
	var memtypePrice MembershipTypePrice
	_=json.NewDecoder(r.Body).Decode(&memtypePrice)
	

    db.Model(&MembershipTypePrice{}).Where("name=?", memtypePrice.Name).Updates(&memtypePrice)

	json.NewEncoder(w).Encode(&memtypePrice)


	// membershipPriceinBytes,_:=json.Marshal(membershipPrice)
	// fmt.Fprintf(w,"updated price of membership types\n")
	// w.Write(membershipPriceinBytes)

	// fmt.Println("membership prices has been changed",membershipPrice)

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
					v.Refund=(float64(refundMonth)*membershiptypeprice["Gold"])/2

				}else{

					v.Refund=(float64(refundMonth)*membershiptypeprice["Silver"])/2
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

func CustomerEnrolmentsDatabyID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var enrollment enrolledInfo
	db.First(&enrollment, params["id"])
	json.NewEncoder(w).Encode(&enrollment)
}

func CustomerEnrolmentsData(w http.ResponseWriter,r * http.Request){


	params := mux.Vars(r)
	var enrollmentrecord enrolledInfo
	db.First(&enrollmentrecord, params["id"])
	json.NewEncoder(w).Encode(&enrollmentrecord)
}

//All DATABASES
//------------------------------------------>
//dataBASE for enrolments
var enrolmentdataBASE []enrolledInfo





var db *gorm.DB
var err error



// ------------------------------>(MAIN FUNCTION)
func main() {

	//Default pricing
//------------------------------------------>	
	

	
//------------------------------------------->
	
router := mux.NewRouter()

dsn := "host=localhost port=5432 user=postgres password=6280912015 dbname=gorm_db sslmode=disable TimeZone=Asia/Shanghai"

db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})


if err != nil {
	panic("failed to connect database")
}

fmt.Println("DB connection established")

//defer db.Close()

db.AutoMigrate(&enrolledInfo{})



router.HandleFunc("/gym", CustomerEnrolmentsData).Methods("GET")
router.HandleFunc("/gym/{id}", CustomerEnrolmentsDatabyID)
router.HandleFunc("/gym", CustomerEnrolments)
router.HandleFunc("/gym/{id}", DeleteMembership).Methods("DELETE")
router.HandleFunc("/gym/setPrice", SetMembershipPrice)



log.Fatal(http.ListenAndServe(":8080", router))




	// fmt.Println("working...")
	// http.HandleFunc("/CustomerEnrolments", CustomerEnrolments)
	// http.HandleFunc("/SetMembershipPrice",SetMembershipPrice)
	// http.HandleFunc("/DeleteMembership", DeleteMembership)
	
	


	
	// fmt.Println("dataBASE",dataBASE)



}