/*
Radhika Srirangam Nadhamuni Manohar
Sjsu id: 009426196
CMPE273: Assignment 1
*/

package main

import (
"fmt"
"net/http"
"bytes"
"io/ioutil"
"encoding/json"
"strconv"
"os"
)

type StockBuyingMessage struct 
{
    Method string
    Params [1]StockBuyingParam
    Id string
}

type StockBuyingParam struct
{
    StockSymbolAndPercentage string `json: stockSymbolAndPercentage`
    Budget string `json: budget`
}


type StockBuyingResponse struct
{
    Error string `json: error`
    Id string `json: id`
    Result struct{
        TradeId string `json: tradeId`
        Stocks string `json: stocks`
        UnvestedAmount string `json: unvestedAmount`

        } `json: result`

    }


    type CheckPortfolioParam struct{
        Id string
    }


    type CheckPortfolioResponse struct
    {
        Error string `json: error`
        Id string `json: id`
        Result struct{
            Stocks string `json: stocks`
            CurrentMarketValue string `json: currentMarketValue`
            UnvestedAmount string `json: unvestedAmount`
            } `json: result`

        }


        type CheckPortfolioMessage struct 
        {
            Method string
            Params [1]CheckPortfolioParam
            Id string
        }

        func main() {
            url := "http://localhost:1111/api/"
            fmt.Println("Server URL is:>", url)



            fmt.Println("1 - Buy Stocks")
            fmt.Println("2 - Check Portfolio")
            fmt.Println("Please enter 1 or 2 to proceed.")
            var choice int
            var newline float64
            _,err:=fmt.Scanf("%d",&choice)

//checking if integer and if it is equal to 1 or 2
            if err != nil {
             fmt.Println("Invalid input!")
             return
         }
         if (!(choice==1 || choice ==2 )){
            fmt.Println("Invalid choice!")
            return
        }

        fmt.Scanln(&newline)
        if(choice==1){
            fmt.Println("============Buy Stocks============")
            fmt.Println("Please enter the stock buying string [For eg: GOOG:50%,YHOO:30%,AMKR:20%]" )
            var stockBuyingString string
            _,err1:=fmt.Scanf("%s",&stockBuyingString)
            if err1!= nil {
             fmt.Println("Invalid input!")
             return
         }
         fmt.Scanln(&newline)

         fmt.Println("Please enter the purchasing amount (numerical)" )
         var budget float64
         _,err2:=fmt.Scanf("%f",&budget)
         if err2!= nil {
             fmt.Println("Invalid input!")
             return
         }
         fmt.Scanln(&newline)

       //converting the float into string:
         budgetString:=strconv.FormatFloat(budget,'f',-1,64)


         paramVal:=[1]StockBuyingParam{StockBuyingParam{stockBuyingString,budgetString}}
         m := StockBuyingMessage{"PurchaseStockService.Buy",paramVal,"1"}
         b, err3:= json.Marshal(m)
         if err3 != nil {
            errorCheck(err3)
        }


        req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
        req.Header.Set("Content-Type", "application/json")

        //Display request
        fmt.Println(req)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errorCheck(err)
        }

        defer resp.Body.Close()

        body, _ := ioutil.ReadAll(resp.Body)

    //unmarshalling
        var res StockBuyingResponse
        err4 := json.Unmarshal(body, &res)
        if err4 != nil {
            errorCheck(err4)
        }

        if(len(res.Error)>0){
    
            fmt.Println("")
            fmt.Println("Response error: ",res.Error)

        }else{
           
            fmt.Println("")
            fmt.Println("*** Response Result ***")
            fmt.Println("")
            fmt.Println("Trade Id:",res.Result.TradeId)
            fmt.Println("Stocks:",res.Result.Stocks)
            fmt.Println("Unvested amount:","$"+res.Result.UnvestedAmount)
        }

    }else{

    //Checking portfolio:
        fmt.Println("============Checking Portfolio========")
        fmt.Println("Please enter the Trade ID (Interger) to check portfolio status:")
        var requestNum int64
        var newline float64
        _,err5:=fmt.Scanf("%d",&requestNum)
//checking if integer and if it is equal to 1 or 2
        if err5 != nil {
         fmt.Println("Invalid input!")
         return
     }
     fmt.Scanln(&newline) 

     requestNoString:=strconv.FormatInt(requestNum,10)

    //Marshalling

     paramVal1:=[1]CheckPortfolioParam{CheckPortfolioParam{requestNoString}}
     m1 := CheckPortfolioMessage{"PurchaseStockService.CheckPortfolio",paramVal1,"1"}
     b1, err6:= json.Marshal(m1)
     if err6 != nil {
        errorCheck(err6)
    }
    

    //POST
    req1, err7 := http.NewRequest("POST", url, bytes.NewBuffer(b1))
    if err7 != nil {
        errorCheck(err7)
    }

    req1.Header.Set("Content-Type", "application/json")

     fmt.Println(req1)

    client1 := &http.Client{}
    resp2, err8 := client1.Do(req1)
    if err8 != nil {
        errorCheck(err8)
    }
    defer resp2.Body.Close()


    body2, _ := ioutil.ReadAll(resp2.Body)
 

    //unmarshalling
    var res2 CheckPortfolioResponse
    err9 := json.Unmarshal(body2, &res2)
    if err9 != nil {
        errorCheck(err9)
    }

    if(len(res2.Error)>0){
        fmt.Println("Response error: ",res2.Error)

    }else{
        fmt.Println("")
        fmt.Println("*** Response Result ***")
        fmt.Println("")
        fmt.Println("Stocks:",res2.Result.Stocks)
        fmt.Println("CurrentMarketValue:","$"+res2.Result.CurrentMarketValue)
        fmt.Println("UnvestedAmount:","$"+res2.Result.UnvestedAmount)
    }

}
}

func errorCheck(err error){
            if err!=nil{
                    fmt.Println("Error :", err.Error())
                    os.Exit(1)
            }

        }


