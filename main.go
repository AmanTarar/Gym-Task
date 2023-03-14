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
// var membershiptypeprice = map[string]float64{"Gold": 2000, "Silver": 1000}






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



//create membership price database
//working
func CreateMembershipPriceDB(w http.ResponseWriter,r *http.Request){

//post request

	var memtypePrice MembershipTypePrice
	_=json.NewDecoder(r.Body).Decode(&memtypePrice)

	result:=db.Create(&memtypePrice)
	if result.Error!=nil{
		fmt.Println("error in DB")
	}

	json.NewEncoder(w).Encode(&memtypePrice)
	fmt.Fprint(w,"membershipPrice has been successfully set")

}

//working
func SetMembershipPrice(w http.ResponseWriter,r *http.Request){

//function to change or update the price of membership types

	//fmt.Println("prices will update soon")
	fmt.Fprintln(w, "updated membership prices!!")
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

     params := mux.Vars(r)

	 var enrol enrolledInfo

	           db.Where("id=?", params["id"]).First(&enrol)
				fmt.Println(enrol.Name,":is deleted from database")
				//enrolmentdataBASE=append(enrolmentdataBASE[:i],enrolmentdataBASE[i+1:]... )
				enrol.MemberShip="INACTIVE"
			
			
				fmt.Println(enrol.Name," membership has been deactivated")
				fmt.Println("enrolmentdataBASE",enrolmentdataBASE)

			    //refund calculated price
				currentTime:=time.Now().Format("2006-01-02")
				enrol.EndDate=currentTime

				currentTimesliced:=strings.Split(currentTime,"-")
				monthinString:=currentTimesliced[1]
				intmonth, err := strconv.Atoi(monthinString)
				if err!=nil{
					fmt.Println("fatt gya")
				}
				//duartion -intmonth
				refundMonth:=enrol.Duration-intmonth
				enrol.Duration=enrol.Duration-refundMonth


				var memtypeprice MembershipTypePrice

				if enrol.MemberShipTYPE=="gold"{
					db.Where("name = ?", "Gold").First(&memtypeprice)
					enrol.Refund=(float64(refundMonth)*memtypeprice.Price/2)

				}else{

					db.Where("name = ?", "Silver").First(&memtypeprice)
					enrol.Refund=(float64(refundMonth)*memtypeprice.Price/2)
				}
				fmt.Println("details with refund added",enrol.Name)

				// enrol_in_Bytes,_:=json.Marshal(enrol)
				// w.Write(enrol_in_Bytes)
				// fmt.Println("",enrol)


				db.Where("id =?", params["id"]).Updates(&enrolledInfo{Refund: enrol.Refund, Duration: enrol.Duration, EndDate: time.Now().Format("2006-01-02"), MemberShip:enrol.MemberShip })
				db.Where("id=?", params["id"]).Delete(&enrol)
				fmt.Fprint(w,"this record has been deleted\n\n")
				json.NewEncoder(w).Encode(&enrol)
				
			
		
		
		


	
	fmt.Fprint(w,"Not present in database")
	
}
	


//Working
func CustomerEnrolmentsDatabyID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var enrollment enrolledInfo
	
	
	
	db.Where("id=?",params["id"]).First(&enrollment)
	json.NewEncoder(w).Encode(&enrollment)
	fmt.Fprintf(w,"here is your customer with id")
	
}

//working
func CustomerEnrolmentsData(w http.ResponseWriter,r * http.Request){


	// Reading from DB
	result := db.Unscoped().Find(&enrolmentdataBASE)
	if result.Error != nil {
		http.Error(w, "Reading users failed", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrolmentdataBASE) // this will send response as a json value
	fmt.Fprintf(w, "Reading successful")
}

//All DATABASES
//------------------------------------------>
//dataBASE in slice form, for enrolments
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

db.AutoMigrate(&enrolledInfo{},&MembershipTypePrice{})




router.HandleFunc("/gym/enrollmentData", CustomerEnrolmentsData)
router.HandleFunc("/gym/enrollmentData/{id}", CustomerEnrolmentsDatabyID)
router.HandleFunc("/gym/enrollment", CustomerEnrolments)
router.HandleFunc("/gym/deleteMember/{id}", DeleteMembership)//
router.HandleFunc("/gym/setPrice", SetMembershipPrice)
router.HandleFunc("/gym/createMembershipPrice",CreateMembershipPriceDB)//post method




log.Fatal(http.ListenAndServe(":8080", router))
	// fmt.Println("dataBASE",dataBASE)



}