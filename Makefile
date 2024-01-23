mocks-clear:
	rm -r ./mocks && mkdir ./mocks
mocks: mocks-clear
	@mockery --all --keeptree