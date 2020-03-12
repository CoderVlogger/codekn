.PHONY: deploy/bin
deploy/bin:
	sup --debug -f Supfile development deploy-bin

.PHONY: deploy/docker
deploy/docker:
	sup --debug -f Supfile development deploy-docker

.PHONY: st/run
st/%/run:
	cat ./target | vegeta attack -rate=50/1s -duration=30s > results.$*.bin

.PHONY: st/report
st/%/report:
	vegeta report results.$*.bin

.PHONY: st/plot
st/%/plot:
	cat results.$*.bin | vegeta plot > results.$*.html
	open results.$*.html

.PHONY: st/%
st/%: st/run st/report st/plot

.PHONY: test/mk
test/mk/%:
	echo 'test/mk $*'
