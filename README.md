# Payment Gateway API

The Payment Gateway API enables merchants (clients) to easily manage payments. 

## User stories:
A client can request an authorisation call:

i.e.: 

        A merchant requests an authorisation call. 
            This call contains:
                    The customer credit card data as well as an amount and currency. 
            It will return:
                    A unique ID that will be used in all next API calls.
   
A client can request a void call:

i.e.: 

        A merchant can cancel a transaction without billing a cardholder. 
            This call contains:
                     The authorisation code generated previously using the authorise call.
            It will return:
                     Whether the cancellation was successfull or an error 
                     and the amount and currency available. 
               
               
A client can request a capture call:

i.e.: 

        A merchant requests a capture call which will capture the money from the issuer. 
            This call contains:
                     The authorisation code generated previously using the authorise call 
                     and the amount of monye to be captured.
            It will return:
                     Whether the capture was successfull or an error 
                     and the amount and currency available. 
               

A client can request a refund call:

i.e.: 

        A merchant request an authorisation call. 
            This call contains:
                     The authorisation code generated previously using the authorise call 
                     and the amount of monye to be returned.
            It will return:
                     Whether the capture was successfull or an error 
                     and the amount and currency available. 
               
## Assumptions
* Definitions about merchant, acquirer, issuer and cardholder can be retrieved from [Payments terminology](https://www.marqeta.com/payments-basics)
* Financial amounts will be managed as float32 values as proper implementation of monetary value data structure is out of scope.
* In a real world scenario "name on the card" and "billing address" string values might be useful in terms of fraud detection and troubleshooting,
 however, they have been considered as out of scope for this API and they won't be stored.
* Sensitive data such as card details should be stored in PCI DSS compliant way. Such implementation is our of scope. 
* I assume in both "capture" and "refund" endpoint, currency code will be the same as the authorisation call. Currency conversion is out of scope.
* Currency conversion will not be implemented, in case currencies don't match, an error will be returned back to the client.
* Client sends only positive values for amount. Hence during validation, the amount will be checked so that it will fail if negative.
* I assume to void a transaction, this cannot be in either capture or refund stages of its lifecycle, hence, the service will check that
in the db whether a refund or a capture operation has previously been done. In fact, the operation table will keep track of what operations (authorisation,
void, capture and refund) have been executed for each authID. So, if there is more than 1 operation (one being the authorization operation) then,
according to the requirements, the void cannot be executed.

## How to run: 
### Prerequisites: 
- Go 1.13
- Docker 19.03.5
- DockerHub access to pull image

### How run the api:
To pull the docker image from DockerHub and start the service on port 8080, please run the below command:
```
docker run --rm -p 8080:8080 sam195/paymentgatewayapi:latest
```
The API will run on ```http://localhost:8080/``` the usage is described in the "Usage" section of this README.

If that fails, open terminal in the root of the project and run the below command:

```
go run main.go
```

## Usage

This can be done using multiple tools such as Postman and Curl commands.

## API endpoints definition

### Authorisation call

Returns the authorisation unique ID.

<details>
  <summary>Call definition</summary>
  
* **URL**

  /authorize

* **Method:**

  `POST`
  
* **Data Params**

     **Required:**
   
    ```json
    {
      "card_details":{
        "card_number": "integer indicating the cardholder's card number",
        "expiry_date": "string indicating the date of expiration of the card in MM-YYYY format",
        "cvv": "integer indicating the card verification value"
      },
      "amount": "floating point (float32) value with the amount to be authorised",
      "currency": "string in three letter format indicating the currency of the amount to be authorised."
    }
    ```

* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Content:** 
    ```json
    {
     "id": "string indicating the authorisation unique id",
     "success": "boolean indicating whether the call was successful or not",
     "amount": "floating point (float32) value with the amount that has been authorised",
     "currency": "string in three letter format indicating the currency of the amount that has been authorised."
    }
    ```
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST <br />
  
      In case the required fields are wrong or invalid.
      
      **Content:** `{ "error": "string indicating the errors" }`
  
  OR  
     
  * **Code:** 401 UNAUTHORISED <br />
  
      In case the card number is rejected.
      
      **Content:** `{ "error": "string indicating the errors" }`
        
  OR

  * **Code:** 422 UNPROCESSABLE ENTITY <br />
  
      In case any of the fields are invalid. e.g. if the card is expired.
      
      **Content:** `{ "error": "string indicating the error" }`
  
  OR
  
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    
      In case there is no connection to the database or marshalling issues within the service.
        
      **Content:** `{ "error": "string indicating the error" }`
      
</details>

### Void call

Returns the amount and currency available after the avoid call has been processed.

<details>
  <summary>Call definition</summary>

* **URL**

  /void

* **Method:**

  `PATCH`
  
* **Data Params**

     **Required:**
   
    ```json
    {
       "id": "string indicating the authorisation unique id"
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 
    ```json
    {
     "success": "boolean indicating whether the authorisation call was successful",
     "amount": "floating point (float32) value with the amount that has been authorised",
     "currency": "string in three letter format indicating the currency of the amount that has been authorised."
    }
    ```
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
  
    In case the authorisation ID cannot be found.
  
    **Content:** `{ "error": "string indicating the error" }`
    
  OR
  
  * **Code:** 400 BAD REQUEST <br />
  
      In case the required fields are wrong or invalid.
      
      **Content:** `{ "error": "string indicating the errors" }`
    
  OR

  * **Code:** 422 UNPROCESSABLE ENTITY <br />
  
    In case any of the fields are invalid.    
  
    **Content:** `{ "error": "string indicating the error" }`
    
  OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    
      In case there is no connection to the database or marshalling issues within the service.
        
      **Content:** `{ "error": "string indicating the error" }`
    
</details>
    
### Capture call

Returns the amount and currency available after capturing some or all of the authorised money.

<details>
  <summary>Call definition</summary>

* **URL**

  /capture

* **Method:**

  `PATCH`
  
* **Data Params**

     **Required:**
   
    ```json
    {
     "id": "string indicating the authorisation unique id",
     "amount": "floating point (float32) value indicating the available authorised amount"
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 
    ```json
    {
     "success": "boolean indicating whether the authorisation call was successful",
     "amount": "floating point (float32) value with the amount that has been authorised",
     "currency": "string in three letter format indicating the currency of the amount that has been authorised."
    }
    ```
 
* **Error Response:**

  * **Code:** 204 No Content <br />
  
    In case the authorisation ID cannot be found.
  
    **Content:** `{ "error": "string indicating the error" }`
    
  OR
  
  * **Code:** 400 BAD REQUEST <br />
  
      In case the required fields are wrong or invalid.
      
      **Content:** `{ "error": "string indicating the errors" }`
  OR
      
  * **Code:** 401 UNAUTHORISED <br />
  
      In case the authorised card is now expired or the card number is rejected.
      
      **Content:** `{ "error": "string indicating the errors" }`
        
  OR

  * **Code:** 422 UNPROCESSABLE ENTITY <br />
  
      In case any of the fields are invalid.
      
      **Content:** `{ "error": "string indicating the error" }`
    
  OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    
      In case there is no connection to the database or marshalling issues within the service.
        
      **Content:** `{ "error": "string indicating the error" }`
  
</details>
  
### Refund call

Returns the amount and currency available after capturing some or all of the authorised money.

<details>
  <summary>Call definition</summary>

* **URL**

  /capture

* **Method:**

  `PATCH`
  
* **Data Params**

     **Required:**
   
    ```json
    {
     "id": "string indicating the authorisation unique id",
     "amount": "floating point (float32) value indicating the available authorised amount"
    }
    ```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 
    ```json
    {
     "success": "boolean indicating whether the authorisation call was successful",
     "amount": "floating point (float32) value with the amount that has been authorised",
     "currency": "string in three letter format indicating the currency of the amount that has been authorised."
    }
    ```
 
* **Error Response:**

  * **Code:** 204 No Content <br />
  
      In case the authorisation ID cannot be found.
    
      **Content:** `{ "error": "string indicating the error" }`
    
  OR

  * **Code:** 400 BAD REQUEST <br />
  
      In case the required fields are wrong or invalid.
      
      **Content:** `{ "error": "string indicating the errors" }`
  OR
      
  * **Code:** 401 UNAUTHORISED <br />
  
      In case the authorised card is now expired or the card number is rejected.
      
      **Content:** `{ "error": "string indicating the errors" }`
            
  OR
  
  * **Code:** 422 UNPROCESSABLE ENTITY <br />
  
      In case any of the fields are invalid.
      
      **Content:** `{ "error": "string indicating the error" }`
    
  OR
    
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    
      In case there is no connection to the database or marshalling issues within the service.
        
      **Content:** `{ "error": "string indicating the error" }`

</details>

## How to test
The project contains both Unit and Integration tests, below are steps to run them

### Unit tests
The unit tests mocks the call to the database in order to only check the functionality of the service.

To run the unit tests:

Open terminal in the root of the project

```
go test ./... -short
```

### Integration tests
The integration tests test calls between the service and its dependency i.e. database.

```
go test ./... -run Integration
```

### Future work
* Implementation of proper data structure for financial information such as money amounts. 
* Any sensitive card details storage should adhere to PCI data security standard requirements, in this solution, the CVV is 
not persisted into the db as only if needed, these information are required to be stored. 
* Currency conversion, [API](https://exchangeratesapi.io/) can be used to query foreign exchange rates for currency conversion.
By integrating with this API, the currency code check can also be improved as the proposed solution only
checks whether the code is a 3 letter string without checking if it is an actual currency code.
* The database store should be persisted using Docker Volumes so that even when the service is restarted, the transaction data
are kept safe and ready to be used once the service is up and running again.
