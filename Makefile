ARTIFACTS_DIR = artifacts
DEPLOYMENT_DIR = deployment
OS = linux
ARCH = arm64
GO_TAGS = lambda.norpc
DEPLOYMENT_HANDLER = bootstrap
DEPLOYMENT_PACKAGE = $(DEPLOYMENT_DIR).zip
SRC_DIR = src

# create final .zip
z: $(ARTIFACTS_DIR)/$(DEPLOYMENT_HANDLER)
	@echo "Compressing deployment package"
	cd $(ARTIFACTS_DIR) && zip -r ../$(DEPLOYMENT_DIR)/$(DEPLOYMENT_PACKAGE) ./*

# build the binaries whenever a .go file inside SRC_DIR changes
$(ARTIFACTS_DIR)/$(DEPLOYMENT_HANDLER): $(shell find $(SRC_DIR) -type f -name *.go)
# create artifacts and deployment dirs if not exists
	@echo "Creating artifacts directories"
	mkdir -p $(ARTIFACTS_DIR) $(DEPLOYMENT_DIR)
	@echo "\n"

# compile lambdas
	@echo "Compiling lambdas"
	GOOS=$(OS) GOARCH=$(ARCH) go build -tags $(GO_TAGS) -o $(DEPLOYMENT_HANDLER) $(SRC_DIR)/**.go
	@echo "\n"

	@echo "moving binaries to artifacts"
	mv $(DEPLOYMENT_HANDLER) $(ARTIFACTS_DIR)/
	@echo "\n"

.PHONY: deploy
deploy: $(DEPLOYMENT_DIR)/$(DEPLOYMENT_PACKAGE)
	cd cdk && cdk deploy -x "*.DS_Store"