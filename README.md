homework
slice operation

#switch  the default namespace for k8s

➜  golang go run src/switchns/switchns.go -h

switchns               # Get the current namespace!
switchns -n namespace  # switch the namespace in k8s!
➜  golang go run src/switchns/switchns.go

The Current Namespace is jenkins
➜  golang go run src/switchns/switchns.go -n sonarqube

Successfully switch the namespaces: sonarqube!
