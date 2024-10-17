#!/bin/bash

ACE_PLATFORM=$1
API_SECRET=$2
APPLICATION_NAME=$3
INSTALLER_URL=$4
LOCATION=$5
PUBLIC_IP=$6
RESOURCE_GROUP=$7

sudo su
HOME="/root"
cd $HOME
apt-get -y update
apt upgrade -y
set -xeo pipefail
exec >/root/userdata.log 2>&1

#constants (don't touch)
SKU="Standard_LRS"
STORAGE_ACCOUNT_NAME="ace"
CONTAINER_NAME="ace"
ACCESS_KEY=""
INSTALLER_ID=$(echo $INSTALLER_URL | awk -F '[/]' '{ print $8 }')

timestamp() {
    date +"%Y/%m/%d %T"
}
log() {
    local type="$1"
    local msg="$2"
    local script_name=${0##*/}
    echo "$(timestamp) [$script_name] [$type] $msg"
}
retry() {
    local retries="$1"
    shift
    local count=0
    local wait=5
    until "$@"; do
        exit="$?"
        if [ $count -lt $retries ]; then
            log "INFO" "Attempt $count/$retries. Command exited with exit_code: $exit. Retrying after $wait seconds..."
            sleep $wait
        else
            log "INFO" "Command failed in all $retries attempts with exit_code: $exit. Stopping trying any further...."
            return $exit
        fi
        count=$(($count + 1))
    done
    return 0
}

create_k3s() {
    echo 'fs.inotify.max_user_instances=100000' | sudo tee -a /etc/sysctl.conf
    echo 'fs.inotify.max_user_watches=100000' | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p

    # Create k3s cluster
    SERVER_IP=${PUBLIC_IP}
    cmd="curl -sfL https://get.k3s.io"
    retry 5 $cmd | INSTALL_K3S_EXEC="--disable=traefik --disable=metrics-server" sh -s - --tls-san "$SERVER_IP"

    echo 'alias k=kubectl' >> ${HOME}/.bashrc
    echo 'export KUBECONFIG=/etc/rancher/k3s/k3s.yaml' >> ${HOME}/.bashrc
    source "${HOME}/.bashrc"

    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

    # wait for 2 pods to become running
    cmd="kubectl wait --for=condition=ready pods --all -A --timeout=5m"
    retry 5 $cmd

    # Install helm
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
}

download_values(){
    mkdir old
    cd old
    curl -L "${INSTALLER_URL}" -o "archive.tar.gz"
    tar -xvzf archive.tar.gz

    #soruce azure credential file from archive.tar.gz
    source env.sh

    cd ..
}

###test az cli
az_cli() {
    curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

    #azure cli login
    az login \
    --service-principal \
    -t ${AZURE_TENANT_ID} \
    -u ${AZURE_CLIENT_ID} \
    -p ${AZURE_CLIENT_SECRET}

    #set subscription id
    az account set -s ${AZURE_SUBSCRIPTION_ID}

    #install jq
    apt-get install jq -y

    STORAGE_ACCOUNT_NAME=${STORAGE_ACCOUNT_NAME}$(head /dev/urandom | tr -dc 'a-z' | head -c 6)
    echo "storage account name: "${STORAGE_ACCOUNT_NAME}

    az storage account create --name ${STORAGE_ACCOUNT_NAME} --resource-group ${RESOURCE_GROUP} --location ${LOCATION} --sku ${SKU}

    #in the --assignee-object-id flag you have to give the objectId of the service account
    #not the appId. Currently the service principal that we are using has the contributor permission,
    #that's why it can't assign the role to himself
#    az ad signed-in-user show --query id -o tsv | az role assignment create \
#            --role "Storage Blob Data Contributor" \
#            --assignee-object-id "0000-000-000000-0000000" \
#            --scope "/subscriptions/0000-000-000000-0000000/resourceGroups/<resource group>/providers/Microsoft.Storage/storageAccounts/<storage account>"


    az storage container create \
        --account-name ${STORAGE_ACCOUNT_NAME} \
        --name ${CONTAINER_NAME}

    ACCESS_KEY=$(az storage account keys list --account-name ${STORAGE_ACCOUNT_NAME} | jq -r '.[0].value')

    #call the webhook here
    resp=$(curl -X POST https://appscode."$ACE_PLATFORM"/marketplace/api/v1/marketplaces/azure/notification/resource?secret=${API_SECRET} \
      -H "Content-Type: application/json" \
      -d '{
            "eventType": "BIND",
            "applicationId": "/subscriptions/'${AZURE_SUBSCRIPTION_ID}'/resourceGroups/'${RESOURCE_GROUP}'/providers/Microsoft.Solutions/applications/'${APPLICATION_NAME}'",
            "bindingInfo": {
              "installerID": "'${INSTALLER_ID}'",
              "options": {
                "infra": {
                  "dns": {
                    "provider": "none",
                    "targetIPs": ["'${PUBLIC_IP}'"]
                  },
                  "cloudServices": {
                    "objstore": {
                      "auth": {
                        "azure": {
                          "AZURE_ACCOUNT_KEY": "'${ACCESS_KEY}'",
                          "AZURE_ACCOUNT_NAME": "'${STORAGE_ACCOUNT_NAME}'"
                        }
                      },
                      "bucket": "azblob://'${CONTAINER_NAME}'"
                    },
                    "provider": "azure"
                  },
                  "kubestash": {
                    "backend": {
                      "azure": {
                        "container": "'${CONTAINER_NAME}'",
                        "prefix": "ace"
                      }
                    },
                    "retentionPolicy": "keep-1mo",
                    "schedule": "0 */2 * * *",
                    "storageSecret": {
                      "create": true
                    }
                  }
                },
                "initialSetup": {
                  "cluster": {
                    "region": "'${LOCATION}'"
                  },
                  "subscription": {
                    "azure": {
                      "applicationId": "/subscriptions/'${AZURE_SUBSCRIPTION_ID}'/resourceGroups/'${RESOURCE_GROUP}'/providers/Microsoft.Solutions/applications/'${APPLICATION_NAME}'"
                    }
                  }
                }
              }
            }
          }')
    link=$(echo ${resp} | jq -r '.link')
    if [ ${link} == "null" ]; then   exit ; fi

    mkdir new
    cd new
    curl -L "${link}" -o "archive.tar.gz"
    tar -xvzf archive.tar.gz
    cd ..
}

install_fluxcd() {
    helm upgrade -i flux2 \
    oci://ghcr.io/appscode-charts/flux2 \
    --version ${FLUXCD_CHART_VERSION} \
    --namespace flux-system --create-namespace \
    --set helmController.create=true \
    --set sourceController.create=true \
    --set imageAutomationController.create=false \
    --set imageReflectionController.create=false \
    --set kustomizeController.create=false \
    --set notificationController.create=false \
    --set-string helmController.labels."ace\.appscode\.com/managed=true" \
    --set-string sourceController.labels."ace\.appscode\.com/managed=true" \
    --wait --debug --burst-limit=10000
}
deploy_ace(){
    helm upgrade -i ace-installer \
    oci://ghcr.io/appscode-charts/ace-installer \
    --version ${ACE_INSTALLER_CHART_VERSION} \
    --namespace kubeops --create-namespace \
    --values=./new/values.yaml \
    --wait --debug --burst-limit=10000
    #--set helm.releases.ace.values.global.infra.dns.targetIPs={${PUBLIC_IP}}

}
init(){
    create_k3s
    download_values
    az_cli
    install_fluxcd
    deploy_ace
}
init
