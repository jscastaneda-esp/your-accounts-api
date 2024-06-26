// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Your Accounts Support",
            "email": "jonathancastaneda@jsc-developer.me"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/budget/": {
            "get": {
                "description": "read budgets associated to an user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Read budgets by user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ReadResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "create a new budget",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Create budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Budget data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/your-accounts-api_budgets_infrastructure_model.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/your-accounts-api_budgets_infrastructure_model.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/budget/available/": {
            "post": {
                "description": "create a new available for budget",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Create available for budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Available data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateAvailableRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.CreateAvailableResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/budget/bill/": {
            "post": {
                "description": "create a new bill for budget",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Create bill for budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Bill data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateBillRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.CreateBillResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/budget/bill/transaction": {
            "put": {
                "description": "create a new transaction for bill of the budget",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Create transaction for bill of the budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Bill transaction data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateBillTransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/budget/{id}": {
            "get": {
                "description": "read budget by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Read budget by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Budget ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ReadByIDResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "receive changes associated to a budget",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Receive changes in budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Budget ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Changes data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ChangesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ChangesResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an budget by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "budget"
                ],
                "summary": "Delete budget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Budget ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/log/{id}/code/{code}": {
            "get": {
                "description": "read logs associated to a resource and code",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "log"
                ],
                "summary": "Read logs by resource and code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "budget",
                            "budget_bill"
                        ],
                        "type": "string",
                        "description": "Code",
                        "name": "code",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ReadLogsResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "create token for access",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Authenticate user",
                "parameters": [
                    {
                        "description": "Authentication data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "Create user in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/your-accounts-api_users_infrastructure_model.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/your-accounts-api_users_infrastructure_model.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Action": {
            "type": "string",
            "enum": [
                "update",
                "delete"
            ],
            "x-enum-varnames": [
                "Update",
                "Delete"
            ]
        },
        "domain.BudgetBillCategory": {
            "type": "string",
            "enum": [
                "house",
                "entertainment",
                "personal",
                "vehicle_transportation",
                "education",
                "services",
                "financial",
                "saving",
                "others"
            ],
            "x-enum-varnames": [
                "House",
                "Entertainment",
                "Personal",
                "Vehicle_Transportation",
                "Education",
                "Services",
                "Financial",
                "Saving",
                "Others"
            ]
        },
        "domain.BudgetSection": {
            "type": "string",
            "enum": [
                "main",
                "available",
                "bill"
            ],
            "x-enum-varnames": [
                "Main",
                "Available",
                "Bill"
            ]
        },
        "model.ChangeRequest": {
            "type": "object",
            "required": [
                "action",
                "id",
                "section"
            ],
            "properties": {
                "action": {
                    "enum": [
                        "update",
                        "delete"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.Action"
                        }
                    ]
                },
                "detail": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "id": {
                    "type": "integer",
                    "minimum": 1
                },
                "section": {
                    "enum": [
                        "main",
                        "available",
                        "bill"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.BudgetSection"
                        }
                    ]
                }
            }
        },
        "model.ChangeResponse": {
            "type": "object",
            "properties": {
                "change": {
                    "$ref": "#/definitions/model.ChangeRequest"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "model.ChangesRequest": {
            "type": "object",
            "required": [
                "changes"
            ],
            "properties": {
                "changes": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/model.ChangeRequest"
                    }
                }
            }
        },
        "model.ChangesResponse": {
            "type": "object",
            "properties": {
                "changes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ChangeResponse"
                    }
                }
            }
        },
        "model.CreateAvailableRequest": {
            "type": "object",
            "required": [
                "budgetId",
                "name"
            ],
            "properties": {
                "budgetId": {
                    "type": "integer",
                    "minimum": 1
                },
                "name": {
                    "type": "string",
                    "maxLength": 40
                }
            }
        },
        "model.CreateAvailableResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.CreateBillRequest": {
            "type": "object",
            "required": [
                "budgetId",
                "category",
                "description"
            ],
            "properties": {
                "budgetId": {
                    "type": "integer",
                    "minimum": 1
                },
                "category": {
                    "enum": [
                        "house",
                        "entertainment",
                        "personal",
                        "vehicle_transportation",
                        "education",
                        "services",
                        "financial",
                        "saving",
                        "others"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.BudgetBillCategory"
                        }
                    ]
                },
                "description": {
                    "type": "string",
                    "maxLength": 200
                }
            }
        },
        "model.CreateBillResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.CreateBillTransactionRequest": {
            "type": "object",
            "required": [
                "amount",
                "billId",
                "description"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "billId": {
                    "type": "integer",
                    "minimum": 1
                },
                "description": {
                    "type": "string",
                    "maxLength": 500
                }
            }
        },
        "model.LoginRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "model.LoginResponse": {
            "type": "object",
            "properties": {
                "expiresAt": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "model.ReadByIDResponse": {
            "type": "object",
            "properties": {
                "additionalIncome": {
                    "type": "number"
                },
                "availables": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ReadByIDResponseAvailable"
                    }
                },
                "bills": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ReadByIDResponseBill"
                    }
                },
                "fixedIncome": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "month": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "model.ReadByIDResponseAvailable": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.ReadByIDResponseBill": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "category": {
                    "$ref": "#/definitions/domain.BudgetBillCategory"
                },
                "complete": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "dueDate": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "payment": {
                    "type": "number"
                }
            }
        },
        "model.ReadLogsResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "detail": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.ReadResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "month": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "pendingBills": {
                    "type": "integer"
                },
                "totalAvailable": {
                    "type": "number"
                },
                "totalPending": {
                    "type": "number"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "your-accounts-api_budgets_infrastructure_model.CreateRequest": {
            "type": "object",
            "properties": {
                "cloneId": {
                    "type": "integer",
                    "minimum": 1
                },
                "name": {
                    "type": "string",
                    "maxLength": 40
                }
            }
        },
        "your-accounts-api_budgets_infrastructure_model.CreateResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "your-accounts-api_users_infrastructure_model.CreateRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "your-accounts-api_users_infrastructure_model.CreateResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Your Accounts API",
	Description:      "This is the API from project Your Accounts",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
