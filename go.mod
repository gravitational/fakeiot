module github.com/gravitational/fakeiot

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gravitational/kingpin v2.1.10+incompatible
	github.com/gravitational/logrus v0.10.1-0.20180402202453-dcdb95d728db // indirect
	github.com/gravitational/roundtrip v1.0.0
	github.com/gravitational/trace v1.1.6
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/sirupsen/logrus v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20190228161510-8dd112bcdc25 // indirect
	golang.org/x/net v0.0.0-20190301231341-16b79f2e4e95 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/sirupsen/logrus => github.com/gravitational/logrus v1.4.3
