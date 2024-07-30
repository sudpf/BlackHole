// Code generated by swaggo/swag. DO NOT EDIT.

package voidengine

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "这里写联系人信息",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/v1//user": {
            "put": {
                "description": "Add a User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "modify user param",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/message.ModifyUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ApiResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "add user param",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/message.AddUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ApiResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "delete user param",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/message.DeleteUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ApiResponse"
                        }
                    }
                }
            }
        },
        "/v1/ping": {
            "get": {
                "description": "Ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ping"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ApiResponse"
                        }
                    }
                }
            }
        },
        "/v1/user": {
            "get": {
                "description": "Add a User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "list user param",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/message.ListUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ApiResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "message.AddUserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8
                },
                "username": {
                    "description": "用户名",
                    "type": "string",
                    "maxLength": 32
                }
            }
        },
        "message.DeleteUserRequest": {
            "type": "object",
            "required": [
                "username"
            ],
            "properties": {
                "username": {
                    "description": "用户名",
                    "type": "string",
                    "maxLength": 32
                }
            }
        },
        "message.ListUserRequest": {
            "type": "object",
            "properties": {
                "orderBy": {
                    "description": "排序方式[desc, asc]",
                    "type": "string",
                    "maxLength": 4,
                    "example": "desc"
                },
                "pageNo": {
                    "description": "当前页码",
                    "type": "integer",
                    "minimum": 1,
                    "example": 1
                },
                "pageSize": {
                    "description": "每页数量",
                    "type": "integer",
                    "maximum": 100,
                    "minimum": 1,
                    "example": 20
                }
            }
        },
        "message.ModifyUserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8
                },
                "username": {
                    "description": "用户名",
                    "type": "string",
                    "maxLength": 32
                }
            }
        },
        "response.ApiResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "这里写接口服务的host",
	BasePath:         "这里写base path",
	Schemes:          []string{},
	Title:            "\"API 文档\"",
	Description:      "这里写描述信息",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}