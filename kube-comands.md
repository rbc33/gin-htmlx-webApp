# kubernetes commands

kubectl apply -f gocms-deployment.yaml  
kubectl delete namespace gocms  
kubectl logs gocms-admin-5dc9b544d9-qw9n6    -n gocms  
kubectl get pods -n gocms -w 
kubectl get pods -n gocms -w
kubectl rollout restart deployment gocms gocms-admin -n gocms  
kubectl edit deployment gocms -n gocms
kubectl describe pod gocms-64789f7c8c-p2jxh -n gocms 
kubectl create namespace gocms   

## ip y port forwading

kubectl port-forward svc/gocms-admin-svc 8081:8081 -n gocms
kubectl port-forward mysql-6fc5db8cc8-nzfzm   3306:3306 -n gocms  

kubectl get svc gocms-admin-svc -n gocms 
kubectl get svc gocms-app-svc -n gocms 

kubectl get nodes -o wide 

kubectl create secret generic gocms-db-secret \
  --from-literal=uri="root:root@tcp(mysql.gocms.svc.cluster.local:3306)/gocms" \
  -n gocms
kubectl delete secret gocms-db-secret -n gocms
kubectl delete secret mysql-secret -n gocms

kubectl exec gocms-app-c4b7758d4-rct5k -n gocms -- sh -c 'goose -dir migrations mysql "$DOCKER_DB_URI" up'



kubectl apply -f mysql-deployment.yaml
kubectl wait --for=condition=ready pod -l app=mysql -n gocms --timeout=300s
kubectl apply -f gocms-deployment.yaml

minikube service gocms-app-svc -n gocms

