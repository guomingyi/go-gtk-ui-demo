OUT := usb-hidapi

all:
	@[ -d build ] || mkdir build
	@./mk ${PWD} ${OUT} amd64 all

linux:
	@[ -d build ] || mkdir build
	@./mk ${PWD} ${OUT} amd64 l

win:
	@[ -d build ] || mkdir build
	@./mk ${PWD} ${OUT} amd64 w

clean:
	@rm -rf build


