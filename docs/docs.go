// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/api/v1/asset/despoit": {
            "post": {
                "description": "despoit an asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "asset despoit",
                "operationId": "v1.asset.despoit",
                "parameters": [
                    {
                        "description": "despoit request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.DespoitRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/asset/transfer/{symbol}": {
            "post": {
                "description": "transfer an asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "asset transfer",
                "operationId": "v1.asset.transfer",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "transfer request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.TransferRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/asset/withdraw": {
            "post": {
                "description": "withdraw an asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "asset withdraw",
                "operationId": "v1.asset.withdraw",
                "parameters": [
                    {
                        "description": "withdraw request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.WithdrawRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/asset/{symbol}/history": {
            "get": {
                "description": "get an asset history",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "get asset history",
                "operationId": "v1.asset.history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/market/depth": {
            "get": {
                "description": "get depth",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "market"
                ],
                "summary": "depth",
                "operationId": "v1.market.depth",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/market/klines": {
            "get": {
                "description": "获取K线数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "market"
                ],
                "summary": "klines",
                "operationId": "v1.market.klines",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "M1",
                            "M3",
                            "M5",
                            "M15",
                            "M30",
                            "H1",
                            "H2",
                            "H4",
                            "H6",
                            "H8",
                            "H12",
                            "D1",
                            "D3",
                            "W1",
                            "MN"
                        ],
                        "type": "string",
                        "description": "period",
                        "name": "period",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end",
                        "name": "end",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/market/trades": {
            "get": {
                "description": "获取近期成交记录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "market"
                ],
                "summary": "trades",
                "operationId": "v1.market.trades",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/order": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "创建订单",
                "operationId": "v1.order",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_modules_base_order.CreateOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/order/history": {
            "get": {
                "description": "history list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户中心"
                ],
                "summary": "历史订单",
                "operationId": "v1.user.order.history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "end",
                        "name": "end",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/ping": {
            "get": {
                "description": "test if the server is running",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "base"
                ],
                "summary": "ping",
                "operationId": "v1.ping",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/product": {
            "get": {
                "description": "get product list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "product list",
                "operationId": "v1.product.list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/product/:symbol": {
            "get": {
                "description": "get product detail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "product detail",
                "operationId": "v1.product",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/login": {
            "post": {
                "description": "user login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "user login",
                "operationId": "v1.user.login",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/order/trade/history": {
            "get": {
                "description": "trade history list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户中心"
                ],
                "summary": "成交历史",
                "operationId": "v1.order.trade_history",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/order/unfinished": {
            "get": {
                "description": "unfinished list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户中心"
                ],
                "summary": "未成交的订单",
                "operationId": "v1.user.order.unfinished",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/register": {
            "post": {
                "description": "user register",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "user register",
                "operationId": "v1.user.register",
                "parameters": [
                    {
                        "description": "args",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/version": {
            "get": {
                "description": "程序版本号和编译相关信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "base"
                ],
                "summary": "version",
                "operationId": "v1.version",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.DespoitRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "symbol": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "controllers.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "captcha": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "controllers.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "repeat_password",
                "username"
            ],
            "properties": {
                "captcha": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "repeat_password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "controllers.TransferRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "from": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        },
        "controllers.WithdrawRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "symbol": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "github_com_yzimhao_trading_engine_v2_pkg_matching_types.OrderSide": {
            "type": "string",
            "enum": [
                "bid",
                "ask"
            ],
            "x-enum-varnames": [
                "OrderSideBuy",
                "OrderSideSell"
            ]
        },
        "github_com_yzimhao_trading_engine_v2_pkg_matching_types.OrderType": {
            "type": "string",
            "enum": [
                "limit",
                "market",
                "marketQty",
                "marketAmount"
            ],
            "x-enum-varnames": [
                "OrderTypeLimit",
                "OrderTypeMarket",
                "OrderTypeMarketQuantity",
                "OrderTypeMarketAmount"
            ]
        },
        "internal_modules_base_order.CreateOrderRequest": {
            "type": "object",
            "required": [
                "order_type",
                "side",
                "symbol"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "order_type": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/github_com_yzimhao_trading_engine_v2_pkg_matching_types.OrderType"
                        }
                    ],
                    "example": "limit"
                },
                "price": {
                    "type": "number",
                    "example": 1
                },
                "qty": {
                    "type": "number",
                    "example": 12
                },
                "side": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/github_com_yzimhao_trading_engine_v2_pkg_matching_types.OrderSide"
                        }
                    ],
                    "example": "buy"
                },
                "symbol": {
                    "type": "string",
                    "example": "btcusdt"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
