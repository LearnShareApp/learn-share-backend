// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Ruslan's Support",
            "url": "https://t.me/Ruslan20007",
            "email": "ruslanrbb8@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Login with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login Credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/login.request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/login.response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Register a new user (student) in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "Registration Info",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/registration.request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/registration.response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/categories": {
            "get": {
                "description": "Get list of all categories",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Get categories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/get_categories.response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/teacher": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get user id by jwt token, and he became teach (if he was not be registrate himself as teacher)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "teacher"
                ],
                "summary": "User registrate also as teacher",
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/teacher/skill": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Registrate new skill for teacher (if he not exists create and registrate skill)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "teacher"
                ],
                "summary": "Registrate new skill",
                "parameters": [
                    {
                        "description": "Skill data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/add_skill.request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get info about user by jwt token (in Authorization enter: Bearer \u003cyour_jwt_token\u003e)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/get_profile.response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        },
        "/users/{id}/profile": {
            "get": {
                "description": "Get info about user by id in route (/api/users/{id}/profile)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user profile",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/get_profile.response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/jsonutils.ErrorStruct"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "add_skill.request": {
            "type": "object",
            "required": [
                "category_id"
            ],
            "properties": {
                "about": {
                    "type": "string",
                    "example": "I am Groot"
                },
                "category_id": {
                    "type": "integer",
                    "example": 1
                },
                "video_card_link": {
                    "type": "string",
                    "example": "https://youtu.be/HIcSWuKMwOw?si=FtxN1QJU9ZWnXy85"
                }
            }
        },
        "get_categories.category": {
            "description": "data of category",
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "min_age": {
                    "type": "integer",
                    "example": 12
                },
                "name": {
                    "type": "string",
                    "example": "Programing"
                }
            }
        },
        "get_categories.response": {
            "description": "get categories response",
            "type": "object",
            "properties": {
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/get_categories.category"
                    }
                }
            }
        },
        "get_profile.response": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string",
                    "example": "2002-09-09T10:10:10+09:00"
                },
                "email": {
                    "type": "string",
                    "example": "qwerty@example.com"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "is_teacher": {
                    "type": "boolean",
                    "example": false
                },
                "name": {
                    "type": "string",
                    "example": "John"
                },
                "registration_date": {
                    "type": "string",
                    "example": "2022-09-09T10:10:10+09:00"
                },
                "surname": {
                    "type": "string",
                    "example": "Smith"
                }
            }
        },
        "jsonutils.ErrorStruct": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "login.request": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "strongpass123"
                }
            }
        },
        "login.response": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        },
        "registration.request": {
            "description": "User registration request",
            "type": "object",
            "required": [
                "birthdate",
                "email",
                "name",
                "password",
                "surname"
            ],
            "properties": {
                "birthdate": {
                    "type": "string",
                    "example": "2000-01-01T00:00:00Z"
                },
                "email": {
                    "type": "string",
                    "example": "john@gmail.com"
                },
                "name": {
                    "type": "string",
                    "example": "John"
                },
                "password": {
                    "type": "string",
                    "example": "strongpass123"
                },
                "surname": {
                    "type": "string",
                    "example": "Smith"
                }
            }
        },
        "registration.response": {
            "description": "User registration response",
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:81",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Learn-Share API",
	Description:      "back-end part for mobile application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
