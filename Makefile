.PHONY: radarctl clean-radarctl

radarctl:
	cd cli && go build -o bin/radarctl .

clean-radarctl:
	rm -f cli/bin/radarctl
