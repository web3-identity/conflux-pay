package cns_errors

import "net/http"

func init() {
	// AUTH ERRORS
	rainbowErrorInfos[ERR_AUTHORIZATION_COMMON] = RainbowErrorInfo{"Unauthorized", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_JWT] = RainbowErrorInfo{"Unauthorized, invalid JWT token", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_TOKEN_MISSING] = RainbowErrorInfo{"Authorization header is empty", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_TOKEN_INVALID] = RainbowErrorInfo{"Authorization token is invalid", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_TOKEN_EXPIRED] = RainbowErrorInfo{"Token is expired", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_NOT_KYC] = RainbowErrorInfo{"KYC required", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_NO_PERMISSION] = RainbowErrorInfo{"No permission to access", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_NO_AUTH_HEADER] = RainbowErrorInfo{"Auth header is empty", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_MISS_EXP_FIELD] = RainbowErrorInfo{"Missing exp field", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_EXP_FORMAT] = RainbowErrorInfo{"Exp must be float64 format", http.StatusUnauthorized}
	rainbowErrorInfos[ERR_AUTHORIZATION_JWT_PAYLOAD] = RainbowErrorInfo{"Jwt payload content uncorrect", http.StatusUnauthorized}

	// VALIDATION ERRORS
	rainbowErrorInfos[ERR_INVALID_REQUEST_COMMON] = RainbowErrorInfo{"Invalid request", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_APP_ID] = RainbowErrorInfo{"Invalid app id", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_ADDRESS] = RainbowErrorInfo{"Invalid address", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_CHAIN] = RainbowErrorInfo{"Chain is not supported", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_CONTRACT_TYPE] = RainbowErrorInfo{"Contract type is not supported", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_URL] = RainbowErrorInfo{"Invalid url", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_METADATA_ID] = RainbowErrorInfo{"Invalid metadataId", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_MINT_AMOUNT] = RainbowErrorInfo{"Invalid mint amount, mint amount could not be 0", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_MINT_AMOUNT_721] = RainbowErrorInfo{"Invalid mint amount, mint amount could not more than 1 for erc 721 contract", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_TOKEN_ID] = RainbowErrorInfo{"Invalid token ID", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_CONTRACT_TYPE_UNMATCH] = RainbowErrorInfo{"Contract type and contract address not match", http.StatusBadRequest}
	rainbowErrorInfos[ERR_INVALID_PAGINATION] = RainbowErrorInfo{"Invalid page or limit", http.StatusBadRequest}

	// CONFLICT ERRORS
	rainbowErrorInfos[ERR_CONFLICT_COMMON] = RainbowErrorInfo{"Conflict", http.StatusConflict}
	rainbowErrorInfos[ERR_CONFLICT_COMPANY_EXISTS] = RainbowErrorInfo{"Company already exists", http.StatusConflict}

	// RATELIMIT ERRORS
	rainbowErrorInfos[ERR_TOO_MANY_REQUEST_COMMON] = RainbowErrorInfo{"Too many requests", http.StatusTooManyRequests}

	// INTERNAL SERVER ERRORS
	rainbowErrorInfos[ERR_INTERNAL_SERVER_COMMON] = RainbowErrorInfo{"Internal Server error", http.StatusInternalServerError}
	rainbowErrorInfos[ERR_INTERNAL_SERVER_DB] = RainbowErrorInfo{"Database operation error", http.StatusInternalServerError}

	// BUSINESS ERRORS
	rainbowErrorInfos[ERR_BUSINESS_COMMON] = RainbowErrorInfo{"Business error", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_MINT_LIMIT_EXCEEDED] = RainbowErrorInfo{"Mint limit exceeded", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_DEPLOY_LIMIT_EXCEEDED] = RainbowErrorInfo{"Deploy limit exceeded", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_UPLOADE_FILE_LIMIT_EXCEEDED] = RainbowErrorInfo{"Uploade file limit exceeded", HTTP_STATUS_BUSINESS_ERROR}

	rainbowErrorInfos[ERR_NO_SPONSOR] = RainbowErrorInfo{"Contract has no sponsor", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_NO_SPONSOR_BALANCE] = RainbowErrorInfo{"Contract sponsor balance not enough", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_NO_SPONSOR_FOR_USER] = RainbowErrorInfo{"Contract has no sponsor for application admin", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_NO_PERMISSION_TO_UPDATE_ADMIN] = RainbowErrorInfo{"Only admin can reset admin", HTTP_STATUS_BUSINESS_ERROR}
	rainbowErrorInfos[ERR_CONTRACT_NOT_OWNED_BY_APP] = RainbowErrorInfo{"Contract is not belong to this application", HTTP_STATUS_BUSINESS_ERROR}
}

const (
	HTTP_STATUS_BUSINESS_ERROR = 599
)

// AUTH ERRORS
const (
	ERR_AUTHORIZATION_COMMON RainbowError = http.StatusUnauthorized*100 + iota //40100
	ERR_AUTHORIZATION_JWT
	ERR_AUTHORIZATION_TOKEN_MISSING
	ERR_AUTHORIZATION_TOKEN_INVALID
	ERR_AUTHORIZATION_TOKEN_EXPIRED
	ERR_AUTHORIZATION_NOT_KYC
	ERR_AUTHORIZATION_NO_PERMISSION
	ERR_AUTHORIZATION_NO_AUTH_HEADER
	ERR_AUTHORIZATION_MISS_EXP_FIELD
	ERR_AUTHORIZATION_EXP_FORMAT
	ERR_AUTHORIZATION_JWT_PAYLOAD
)

// VALIDATION ERRORS
const (
	ERR_INVALID_REQUEST_COMMON RainbowError = http.StatusBadRequest*100 + iota //40000
	ERR_INVALID_APP_ID
	ERR_INVALID_ADDRESS
	ERR_INVALID_CHAIN
	ERR_INVALID_CONTRACT_TYPE
	ERR_INVALID_URL
	ERR_INVALID_METADATA_ID
	ERR_INVALID_MINT_AMOUNT
	ERR_INVALID_MINT_AMOUNT_721
	ERR_INVALID_TOKEN_ID
	ERR_INVALID_CONTRACT_TYPE_UNMATCH
	ERR_INVALID_PAGINATION
)

// RESOURCE CONFLICT ERRORS
const (
	ERR_CONFLICT_COMMON RainbowError = http.StatusConflict*100 + iota //40900
	ERR_CONFLICT_COMPANY_EXISTS
)

// RATELIMIT ERRORS
const (
	ERR_TOO_MANY_REQUEST_COMMON RainbowError = http.StatusTooManyRequests*100 + iota //42900
)

// INTERNAL SERVER ERRORS
const (
	ERR_INTERNAL_SERVER_COMMON RainbowError = http.StatusInternalServerError*100 + iota //50000
	ERR_INTERNAL_SERVER_DB
	ERR_INTERNAL_SERVER_DB_NOT_FOUND
)

// BUSINESS ERRORS
const (
	ERR_BUSINESS_COMMON RainbowError = HTTP_STATUS_BUSINESS_ERROR*100 + iota //60000
	ERR_NO_SPONSOR
	ERR_NO_SPONSOR_BALANCE
	ERR_NO_SPONSOR_FOR_USER
	ERR_MINT_LIMIT_EXCEEDED
	ERR_DEPLOY_LIMIT_EXCEEDED
	ERR_UPLOADE_FILE_LIMIT_EXCEEDED
	ERR_NO_PERMISSION_TO_UPDATE_ADMIN
	ERR_CONTRACT_NOT_OWNED_BY_APP
)

func GetRainbowOthersErrCode(httpStatusCode int) int {
	return httpStatusCode * 100
}
