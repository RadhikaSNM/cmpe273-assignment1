/*
Radhika Srirangam Nadhamuni Manohar
Sjsu id: 009426196
CMPE273:Assignment 1
*/
package main

import (
"log"
"fmt"
"net/http"
"github.com/gorilla/rpc"
"github.com/gorilla/rpc/json"
"github.com/gorilla/mux"
"strconv"
"strings"
"io/ioutil"
"errors"
"regexp"
json1 "encoding/json"  
"os"  
)

var indexAssigned int64
var Transactions map[int64]*Transaction

func init() {
    //Initiliazing count
    indexAssigned=0;
    //initialing map
    Transactions=make(map[int64]*Transaction)

    r := mux.NewRouter()    
    jsonRPC := rpc.NewServer()
    jsonCodec := json.NewCodec()
    jsonRPC.RegisterCodec(jsonCodec, "application/json")
//  jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
    jsonRPC.RegisterService(new(PurchaseStockService), "")
    r.Handle("/api/", jsonRPC)  
    http.ListenAndServe(":1111", r)
} 


type StockBuyingArgs struct
{
    StockSymbolAndPercentage string `json: stockSymbolAndPercentage`
    Budget string `json: budget`
}

type StockBuyingReply struct{
    TradeId string `json: tradeId`
    Stocks string `json: stocks`
    UnvestedAmount string `json: unvestedAmount`

}

type CheckPortfolioArgs struct{
    Id string
}


type CheckPortfolioReply struct{
    Stocks string `json: stocks`
    CurrentMarketValue string `json: currentMarketValue`
    UnvestedAmount string `json: unvestedAmount`
}



type Price struct {
    List struct {
        Resources []struct{
            Resource struct{
                Fields struct{
                    Price string `json: price`
                    } `json: fields`

                    } `json: resource`

                    } `json: resources`

                    } `json: list`
                }

//whole
                type Transaction struct {
                    ID int64
                    StockList []Stock
                    Budget float32
                    Unvested float32
                }

                type Stock struct {
                    Name string
                    NoOfShares int64
                    BuyingPrice float32
                }



                type PurchaseStockService struct {}



                func (p *PurchaseStockService) Buy(r *http.Request, args *StockBuyingArgs, reply *StockBuyingReply) error {
                    

    //call the buying function
                    tradeId,str,remaining,err:=buyStock(args.StockSymbolAndPercentage,args.Budget)      
                    reply.TradeId = tradeId
                    reply.Stocks =str
                    reply.UnvestedAmount=remaining
                    return err

                }


                func (p *PurchaseStockService) CheckPortfolio(r *http.Request, args *CheckPortfolioArgs, reply *CheckPortfolioReply) error {
                    log.Printf(args.Id)
    //calling the portfolio checking function:
                    currStr,currAmount,unvestedAmount,err:=getPorfolio(args.Id)

    //setting return values:
                    reply.Stocks=currStr
                    reply.CurrentMarketValue=currAmount
                    reply.UnvestedAmount=unvestedAmount

                    return err
                }



                func buyStock(stockParam string,budget string) (string,string,string,error) {

                    fmt.Println("=============================================")
                    var stockPer []string=strings.Split(stockParam,",")
                    var remainingAmount float32
                    remainingAmount=0.0

//Returned string:  
                    var returnedString string
                    returnedString=""



                    noOfStocks:=len(stockPer)
                    stocks:=make([]Stock,noOfStocks,noOfStocks)

                    var stringFormat bool
                    stringFormat=true;

                    var totalPercentageValue float32

    //Confirming the total %age is =100.

    //obtaining total %age
                    for i:=0;i<noOfStocks&&stringFormat;i++ {
                        wholeString:=stockPer[i]

    //Checking if the individual stock expression format is correct.
                        stringFormat,_:=regexp.MatchString("^[A-Za-z]+:[0-9]+%$", wholeString)
                        if(!stringFormat){
                            err_Format:=errors.New("Input String does not match the format.")
                            return "","","",err_Format 
                        }

    //Splitting on ":"
                        share:=strings.Split(wholeString,":")
    //Strip the % symbol
                        perString:=strings.TrimSuffix(share[1],"%")
                        percentage,_:=strconv.ParseFloat(perString,32)
                        totalPercentageValue=totalPercentageValue+float32(percentage)    

                    }

                    //Throw error for total percentage >100
                    if !(totalPercentageValue==float32(100)) {
                     err_100:=errors.New("Given percentage splits not equal to 100. Exiting")
                     return "","","",err_100 
                 }


                 //checking if budget is correct float32
                 budget1,errBudget:=strconv.ParseFloat(budget,32)
                 if errBudget!=nil{
                    err_Budget:=errors.New("Given Budget is not a valid number")
                    return "","","",err_Budget 

                }


                for i:=0;i<noOfStocks;i++ {
    //addding , except for zero index
                    if(!(i==0)){
                        returnedString =returnedString+","
                    }

    //Adding the charac "
                    returnedString+="\""

                    fmt.Println(stockPer[i])
                    wholeString:=stockPer[i]

    //Splitting on ":"
                    share:=strings.Split(wholeString,":")
                    name:=share[0]
    //Strip the % symbol
                    perString:=strings.TrimSuffix(share[1],"%")
                    fmt.Println("percentage1: ",perString)
                    percentage,_:=strconv.ParseFloat(perString,32)
                    fmt.Println("percentage: ",percentage)

                    amountToSpend:=(float32(percentage)/float32(100))*(float32(budget1))
                    fmt.Println("Amount to spend: ",amountToSpend)

    //call yahoo finance rest api and get the current price: 
                    link:="http://finance.yahoo.com/webservice/v1/symbols/"+name+"/quote?format=json"
                    resp, err := http.Get(link);
                    if err != nil {
                        err_yahoo:=errors.New("Yahoo didnt respond!")
                        return "","","",err_yahoo 
                    }



                    defer resp.Body.Close()
                    body, err1 := ioutil.ReadAll(resp.Body)
                    if err1 != nil {
                        errorCheck(err1)

                    }


                    var obtPrice Price
                    err2:=json1.Unmarshal(body,&obtPrice)
                    if err2 != nil {
                        errorCheck(err2)
                    }




    //Checking for an invalid stock symbol

                    if (len(obtPrice.List.Resources)==0){
                        err_noRes:=errors.New("One of the stock symbol is invalid. Please Check")
                        return "","","",err_noRes 

                    }


                    currentPriceString:=obtPrice.List.Resources[0].Resource.Fields.Price
                    currentPrice,_:=strconv.ParseFloat(currentPriceString,32)
                    fmt.Println("The price of ",name,": ",currentPrice)
                    var noOfSingleShares int64
                    noOfSingleShares=int64(amountToSpend/float32(currentPrice))

                    remainingAmount=remainingAmount+(amountToSpend-(float32(noOfSingleShares)*float32(currentPrice)))


    //Insert in the end
                    stocks[i] = Stock{name,noOfSingleShares,float32(currentPrice)}

    //returning the correct values
                    returnedString+=name+":"+ strconv.FormatInt(noOfSingleShares,10)+":"+"$"+ currentPriceString + "\""

                }
//complete the creation of the transaction:
                //Create new whole transaction:

                indexAssigned=indexAssigned+1;
                StringIndex:=strconv.FormatInt(indexAssigned,10)

                trans:= new(Transaction)

                trans.ID=indexAssigned
                trans.Budget=float32(budget1)
                trans.Unvested=remainingAmount
                trans.StockList=stocks


//Print the entire data structure for verification
                fmt.Printf("%+v\n",trans)

//adding the value to the map 
                Transactions[indexAssigned]=trans


                return StringIndex, returnedString, strconv.FormatFloat(float64(remainingAmount),'f',-1,32),nil

            }




            func getPorfolio(ID string) (string, string, string,error){

                id,err:=strconv.ParseInt(ID,10,64)
                if(err!=nil){
                    err_ID:=errors.New("Given Request number(trade ID) is not a valid number")
                    return "","","",err_ID 
                }

//Throwing an error if the map does not contain the key 
                if _, ok := Transactions[id]; !ok {
                   err_noKey:=errors.New("Supplied Request number(trade ID) is not found in the system. Please check.")
                   return "","","",err_noKey 

               }


               trans:=Transactions[id]
               remainingAmount:=trans.Unvested

    //total present value
               var totalMarketValue float32

    //returned string:
               var returnedString string
               stockList:=trans.StockList

               for i:=0;i<len(stockList);i++{
                stock:=stockList[i]

                if(!(i==0)){
                    returnedString =returnedString+","
                }


    //Adding the charac "
                returnedString+="\""

                fmt.Println(stock.Name,stock.NoOfShares,stock.BuyingPrice)
                name:=stock.Name
                oldPrice:=stock.BuyingPrice
                NoOfShares:=stock.NoOfShares

         //call yahoo finance rest api and get the current price:
                link:="http://finance.yahoo.com/webservice/v1/symbols/"+name+"/quote?format=json"
                resp, err := http.Get(link);
                if err != nil {
                    errorCheck(err)
                }

                defer resp.Body.Close()
                body, err1 := ioutil.ReadAll(resp.Body)
                if err1 != nil {
                    errorCheck(err1)
                }
 
                var obtPrice Price
                err2:=json1.Unmarshal(body,&obtPrice)
                if err2 != nil {
                    errorCheck(err2)
                }

    //Getting the current price:

                currentPriceString:=obtPrice.List.Resources[0].Resource.Fields.Price
                currentPrice,_:=strconv.ParseFloat(currentPriceString,32)
                fmt.Println("The price of ",name,": ",currentPrice)

                var profitIndicator string

                if(float32(currentPrice)<oldPrice){
                    profitIndicator="-"
                }else if float32(currentPrice)>oldPrice{
                    profitIndicator="+"
                }else{
                    profitIndicator=""
                }

                totalMarketValue+=float32(NoOfShares)*float32(currentPrice)
                fmt.Println("total market value for present",totalMarketValue)

    //Creating the return string
                returnedString+=name+":"+ strconv.FormatInt(NoOfShares,10)+":"+profitIndicator+"$"+ currentPriceString + "\""
            }

            fmt.Println("TOTAL market value:",totalMarketValue)

            return returnedString,strconv.FormatFloat(float64(totalMarketValue),'f',-1,32),strconv.FormatFloat(float64(remainingAmount),'f',-1,32),nil
        }

        func errorCheck(err error){
            if err!=nil{
                    fmt.Println("Error :", err.Error())
                    os.Exit(1)
            }

        }


        func main() {
            fmt.Println("The server has started!")
        }


