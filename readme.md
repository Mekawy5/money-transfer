
# Money Transfer RESTful API

  

This is a Go project that emulates a RESTful API for money transfers between accounts. The system is designed to run an in-memory datastore that is not persistent over restarts.

  

#### Data Ingestion

Accounts are parsed from a file and ingested into the system. Once all accounts have been consumed, the system will indicate that it is ready to make a transfer through a console output message.

  

# Endpoints

` GET /accounts `

  

Returns a list of all accounts in the system.

  

### Request

  

1- ` GET /accounts `

  

### Response

  

    {
	    "result": [
		    {
			    "id": "3d253e29-8785-464f-8fa0-9e4b57699db9",
			    "name": "Trupe",
			    "balance": 82.11
		    },
		    {
			    "id": "17f904c1-806f-4252-9103-74e7a5d3e340",
			    "name": "Fivespan",
			    "balance": 951.15
		    }
	    ],
	    "success": true
    }

  

------------

  
  

2- ` POST /accounts/transfer`

  

Transfers an amount of money from one account to another.

  

### Request

` POST /accounts/transfer `


### Form Inputs

| key | value | type | 
|--|--|--|
| sender_id | 3d253e29-8785-464f-8fa0-9e4b57699db9 | string |
| receiver_id | 17f904c1-806f-4252-9103-74e7a5d3e340 | string |
| amount | 10.0 | float |
 

` Response `

  

    {
	    "result": {
		    "from": {
			    "id": "3d253e29-8785-464f-8fa0-9e4b57699db9",
			    "name": "Trupe",
			    "balance": 77.11
		    },
		    "to": {
			    "id": "17f904c1-806f-4252-9103-74e7a5d3e340",
			    "name": "Fivespan",
			    "balance": 956.15
		    }
	    },
	    "success": true
    }

  

------------

  
  
  

# Error Handling

  

The API returns appropriate error messages and HTTP status codes for invalid requests, such as attempting to transfer more money than is available in an account.

The following response returns for this case with status code  `400`

    {
	    "error": "no enough balance",
	    "success": false
    }

  
  

------------

  

# Testing the Application

To run the test cases, use go test command :

  

run ` go test ./... `

  

------------

  
  
  

# Running the Application

To run the application, use go build to generate the binary and simply you can execute it:

  

run ` go build . `

run ` ./money-transfer `

# Design suggestions


- use a persisted transactional database like Postgres
- store transaction audit
- authentication/authorization implementation
- idempotency for transfer requests to prevent redundant transactions
- use caching mechanism to improve the application performance with proper cache invalidation/updates to maintain consistency
- use log monitoring service for fast debugging
- use sentry for better error tracking
