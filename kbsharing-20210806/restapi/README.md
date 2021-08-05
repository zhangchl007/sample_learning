#restapi by python3

➜  kubectl apply -f admin-role.yaml

➜  kubectl -n kube-system get secret|grep admin-token

➜  oc describe secret/admin-token-j5vjc -n kube-system 

replace the token in main.py with mytoken = os.getenv("mytoken", "key")

```
$ python3 main.py
the pods list is:
demo-deployment-6b4d4fbcdb-db652
demo-deployment-6b4d4fbcdb-nm5w8
```


