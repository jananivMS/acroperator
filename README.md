# acroperator
Azure operator POC to create delete Azure Container Registry from a K8s cluster

- This uses Kubebuilder (Steps here - https://book.kubebuilder.io/quick-start.html) to setup the scaffolding.
- The actual operation of creation/deletion of the Azure Container Registry uses the Azure GO SDK (https://github.com/Azure/azure-sdk-for-go)
  The SDK repo also has samples for a lot of the services in the 'services' folder but doesn't have one for ACR.
  
## Steps to try this out (We will use kubernetes with docker-desktop for this for ease)
1. Clone the repo to your computer

2. Create a setenv.sh file in that folder with the below information
`

export AZURE_CLIENT_ID=<Client ID of your service principal>

export AZURE_CLIENT_SECRET=<Client Secret of your service principal>

export AZURE_TENANT_ID=<Azure tenant ID where you want to deploy the ACR>

export AZURE_SUBSCRIPTION_ID=<Azure Subscription ID where you want to deploy the ACR>`

The client ID and client secret is of the service principal in Azure that you have created and have given write access to your resource group where you want to create the ACR. You will need to do this in the Azure portal before running this.

3. Run the above file using `source setenv.sh` in the command prompt where you would be executing the project from. This sets these values as environment variables for the program to read from.

4. Update the values of the "resourcegroup", "location" in the spec section of the YAML file under config/samples. Also update "name" to a unique value.

5. Enable Kubernetes on Docker Desktop: Click on the Docker whale icon on the taskbar --> Preferences --> Kubernetes tab --> Enable Kubernetes

6. Once this is done, you can check that it was enabled by typing the below:
`
kubectl cluster-info

Kubernetes master is running at https://127.0.0.1:6443
KubeDNS is running at https://127.0.0.1:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy`

7. Now run "make install" from the folder. And then "make run". This starts the operator code which is ready to handle requests for creating/deleting a new ACR.

8. Now from another terminal run 'kubectl apply -f config/samples". This will instruct the operator to install a ACR.
