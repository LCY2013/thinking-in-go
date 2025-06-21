#!/bin/bash
kubectl apply -f nfs-pv.yaml
kubectl apply -f nfs-pvc.yaml
kubectl apply -f pod-with-nfs.yaml

#在 Pod 中写入文件：
#kubectl exec -it nfs-pod -- sh -c "echo 'Hello K8s NFS' > /mnt/nfs/test.txt"

#在 NFS 服务器上检查
#cat /nfs/test.txt  # 应输出 "Hello K8s NFS"
