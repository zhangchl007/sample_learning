# -*- coding: utf-8 -*-
import requests
from requests.packages.urllib3.exceptions import InsecureRequestWarning
import os
import json
import logging
import sys

log = logging.getLogger(__name__)
out_hdlr = logging.StreamHandler(sys.stdout)
out_hdlr.setFormatter(logging.Formatter('%(asctime)s %(message)s'))
out_hdlr.setLevel(logging.INFO)
log.addHandler(out_hdlr)
log.setLevel(logging.INFO)
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)


# env variable
namespace = os.getenv("res_namespace", "default")
base_url = os.getenv("base_url", "https://127.0.0.1:56733")
mytoken = os.getenv("mytoken", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImJ6THNOaHhIWUVBRUduS0FKUXVzdE1TYWNXOGFSeGlvbGxSc05CTU5YNEkifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi1qNXZqYyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6Ijc5ZTFlNjRhLTRmM2YtNGFmNS05YmY4LTYzNDJjZWJiMzQwMCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.WqxjT7HKQRlJFUn3cM7Emhxo8nXk4I84YEcObjfuDIybfQcyC9mYa0BI-3TFoVKOPjfLVrw97o18NO_l-t5F4MCfmqUPyEpnenOzkHwPfZbZJg7vq60oicPhIQZQdYC9hDN7K9TCIaWLweumETslW-ZToDGvzIjVL26K5KIoEbE-IjLsSQXKz0uCy_CezOXBSVQo-0age3RM8QF9tOj1iC3rqUlXh2P7RwEPZzGPk50wmhLU3c1IC4tvAPXni67_ZgdOA-h_dlpZEygvKw_XYFAZASFU3cDggnpzsslx70SGfazEFxoZZfHqf0Zlu6347hwHzH9o8V_7YZjHZEp_IA")
mytoken = mytoken.rstrip()
headers = {"Authorization": "Bearer {}".format(mytoken)}
os.environ['NO_PROXY'] = '127.0.0.1,192.168.49.2'

# This is the function that searches for and kills Pods by searching for them by label
def api_return(url):
   r = requests.get(url, headers=headers, verify=False)
   s=json.loads(r.content)
   return s

def list_pods():
        url = "{}/api/v1/namespaces/{}/pods".format(
            base_url, namespace)
        response = api_return(url)
        # Extract the Pod name from the list
        pods = [p['metadata']['name'] for p in response['items']]
        # For each Pod, issue an HTTP DELETE request
        print("the pods list is:")
        for p in pods:
            print("%s" % p)
def main():
    list_pods()

if __name__ == '__main__':
    main()
