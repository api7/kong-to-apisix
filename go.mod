module github.com/api7/kongtoapisix

go 1.14

require (
	github.com/Azure/go-ansiterm v0.0.0-20210608223527-2377c96fe795 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/apache/apisix-ingress-controller v0.0.0-20210614074814-7e146b66f26c
	github.com/globocom/gokong v1.9.1-0.20200127185249-0b630f045649
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/icza/dyno v0.0.0-20200205103839-49cb13720835
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible // indirect
)

// use personal branch for now
replace github.com/apache/apisix-ingress-controller v0.0.0-20210614074814-7e146b66f26c => github.com/yiyiyimu/apisix-ingress-controller v0.0.0-20210618042149-49bc57f52079
