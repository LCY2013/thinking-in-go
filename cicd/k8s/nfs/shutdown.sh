#!/bin/bash
kubectl delete -f pod-with-nfs.yaml
kubectl delete -f nfs-pvc.yaml
kubectl delete -f nfs-pv.yaml
