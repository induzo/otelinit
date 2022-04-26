.PHONY: test bench bench-compare lint sec-scan upgrade release release-tag changelog-gen changelog-commit

########
# test #
########

test:
	go test ./... -cover

test-leak:
	go test ./ -leak

bench:
	go test ./... -bench=. -benchmem | tee ./bench.txt

bench-compare:
	benchstat ./bench.txt

########
# lint #
########

lint:
	golangci-lint run ./... --config=./.golangci.toml

#######
# sec #
#######

sec-scan:
	trivy fs ./

############
# upgrades #
############

upgrade:
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy

###########
# release #
###########

MOD_VERSION = $(shell git describe --abbrev=0 --tags)

release: release-tag changelog-gen changelog-commit
	
release-tag:
	@printf "here is the latest tag present: "; \
	git describe --abbrev=0 --tags; \
	printf "what tag do you want to give? (use the form vX.X.X): "; \
	read -r TAG && \
	git tag $$TAG && \
	printf "\nrelease tagged $$TAG !\n"

#############
# changelog #
#############

MESSAGE_CHANGELOG_COMMIT="update CHANGELOG.md for $(MOD_VERSION)"

changelog-gen:
	@git cliff \
		-c ./cliff.toml \
		-o ./CHANGELOG.md  && \
	printf "\nchangelog generated!\n"

# keep this commit unconventional so it doesnt appear in the changelog
changelog-commit:
	git commit -m $(MESSAGE_CHANGELOG_COMMIT) ./CHANGELOG.md

