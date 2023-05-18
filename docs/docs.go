// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/": {
            "get": {
                "description": "get the status of server",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "main"
                ],
                "summary": "Show the status of server",
                "responses": {
                    "200": {
                        "description": "Status available",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/project/": {
            "post": {
                "description": "create a new project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Create project",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Project data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api-your-accounts_project_infrastructure_model.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/api-your-accounts_project_infrastructure_model.CreateResponse"
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
        "/api/v1/project/logs/{id}": {
            "get": {
                "description": "read logs associated to a project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Read logs by project",
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
                        "description": "Project ID",
                        "name": "id",
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
        "/api/v1/project/{id}": {
            "delete": {
                "description": "Delete an project by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Delete project",
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
                        "description": "Project ID",
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
        "/api/v1/project/{user}": {
            "get": {
                "description": "read projects associated to an user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "project"
                ],
                "summary": "Read projects by user",
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
                            "$ref": "#/definitions/api-your-accounts_user_infrastructure_model.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/api-your-accounts_user_infrastructure_model.CreateResponse"
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
        },
        "/user/auth": {
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
                            "$ref": "#/definitions/model.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.AuthResponse"
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
        "/user/refresh-token": {
            "put": {
                "description": "refresh token for access",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Refresh token of user",
                "parameters": [
                    {
                        "description": "Refresh token data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.RefreshTokenResponse"
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
        }
    },
    "definitions": {
        "api-your-accounts_project_infrastructure_model.CreateRequest": {
            "type": "object",
            "properties": {
                "cloneId": {
                    "type": "integer",
                    "minimum": 1
                },
                "name": {
                    "type": "string",
                    "maxLength": 40
                },
                "type": {
                    "enum": [
                        "budget"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.ProjectType"
                        }
                    ]
                },
                "userId": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "api-your-accounts_project_infrastructure_model.CreateResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "api-your-accounts_user_infrastructure_model.CreateRequest": {
            "type": "object",
            "required": [
                "email",
                "uuid"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "api-your-accounts_user_infrastructure_model.CreateResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "domain.ProjectType": {
            "type": "string",
            "enum": [
                "budget"
            ],
            "x-enum-varnames": [
                "Budget"
            ]
        },
        "model.AuthRequest": {
            "type": "object",
            "required": [
                "email",
                "uuid"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "model.AuthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.ReadLogsResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
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
                "name": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/domain.ProjectType"
                }
            }
        },
        "model.RefreshTokenRequest": {
            "type": "object",
            "required": [
                "email",
                "token",
                "uuid"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "model.RefreshTokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
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
