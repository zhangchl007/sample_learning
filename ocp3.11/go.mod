module ocp3.11

go 1.14

require (
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.5 // indirect
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/apiserver => k8s.io/apiserver v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
)
