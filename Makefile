
.PHONY: install-tools
install-tools:
	go install github.com/apache/skywalking-eyes/cmd/license-eye@latest

.PHONY: addlicense
addlicense: install-tools
	license-eye -c .licenserc.yaml header fix

.PHONY: checklicense
checklicense: install-tools
	license-eye -c .licenserc.yaml header check
