// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/CA": {
            "get": {
                "description": "Retrieve a list of all CAs stored in the in-memory store.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CAs"
                ],
                "summary": "List all Certificate Authorities (CAs)",
                "responses": {
                    "200": {
                        "description": "List of all CAs",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/cert.CaResponse"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new SSH CA and store it in the applications store.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CAs"
                ],
                "summary": "Create a new SSH Certificate Authority (CA)",
                "parameters": [
                    {
                        "description": "New CA",
                        "name": "CA",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cert.CaRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "The newly created CA",
                        "schema": {
                            "$ref": "#/definitions/cert.CaResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Could not create CA",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/CA/{id}": {
            "get": {
                "description": "Retrieve a CA by its ID from the applications store.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CAs"
                ],
                "summary": "Get a SSH Certificate Authority (CA) by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "CA ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/cert.CaResponse"
                        }
                    },
                    "404": {
                        "description": "CA not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/CA/{id}/Sign": {
            "post": {
                "description": "Use the specified CA to sign a provided public key and return the signed key.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CAs"
                ],
                "summary": "Sign a public key with a specific CA",
                "parameters": [
                    {
                        "type": "string",
                        "description": "CA ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Public key to be signed",
                        "name": "public_key",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cert.SignRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "The signed public key will be returned under the 'signed_key' field",
                        "schema": {
                            "$ref": "#/definitions/cert.SignResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request or failed to parse public key",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "CA not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to sign public key",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "cert.CaRequest": {
            "type": "object",
            "properties": {
                "bits": {
                    "description": "Key length",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of CA",
                    "type": "string"
                },
                "type": {
                    "description": "Type of ca, rsa, ed25519",
                    "type": "string"
                }
            }
        },
        "cert.CaResponse": {
            "type": "object",
            "properties": {
                "bits": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "public_key": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "cert.SignRequest": {
            "type": "object",
            "properties": {
                "public_key": {
                    "type": "string"
                }
            }
        },
        "cert.SignResponse": {
            "type": "object",
            "properties": {
                "signed_key": {
                    "type": "string"
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
