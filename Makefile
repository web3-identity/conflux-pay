swag:
	swag fmt
	swag init --parseDependency --parseInternal

gosdk:
	make swag
	rm -rf $(ls -A ${CONFLUXPAY_SDK_DIR_GO}| grep -v .git)
	openapi-generator generate -c ./code_gen/go_config.yaml -o ${CONFLUXPAY_SDK_DIR_GO}
	cd ${CONFLUXPAY_SDK_DIR_GO} && rm -rf ./test && go mod tidy && go mod edit -module github.com/web3-identity/conflux-pay-sdk-go
