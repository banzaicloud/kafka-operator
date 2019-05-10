// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sutil

import (
	"context"
	"errors"
	"reflect"

	objectmatch "github.com/banzaicloud/k8s-objectmatcher"
	banzaicloudv1alpha1 "github.com/banzaicloud/kafka-operator/pkg/apis/banzaicloud/v1alpha1"
	"github.com/banzaicloud/kafka-operator/pkg/scale"
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(log logr.Logger, client runtimeClient.Client, desired runtime.Object, cr *banzaicloudv1alpha1.KafkaCluster) error {
	desiredType := reflect.TypeOf(desired)
	var current = desired.DeepCopyObject()
	var err error

	switch desired.(type) {
	default:
		var key runtimeClient.ObjectKey
		key, err = runtimeClient.ObjectKeyFromObject(current)
		if err != nil {
			return emperror.With(err, "kind", desiredType)
		}
		log = log.WithValues("kind", desiredType, "name", key.Name)

		err = client.Get(context.TODO(), key, current)
		if err != nil && !apierrors.IsNotFound(err) {
			return emperror.WrapWith(err, "getting resource failed", "kind", desiredType, "name", key.Name)
		}
		if apierrors.IsNotFound(err) {
			if err := client.Create(context.TODO(), desired); err != nil {
				return emperror.WrapWith(err, "creating resource failed", "kind", desiredType, "name", key.Name)
			}
			log.Info("resource created")
		}
	case *corev1.PersistentVolumeClaim:
		log = log.WithValues("kind", desiredType)
		log.Info("searching with label because name is empty")

		pvcList := &corev1.PersistentVolumeClaimList{}
		matchingLabels := map[string]string{
			"kafka_cr": cr.Name,
			"brokerId": desired.(*corev1.PersistentVolumeClaim).Labels["brokerId"],
		}
		err = client.List(context.TODO(), runtimeClient.InNamespace(current.(*corev1.PersistentVolumeClaim).Namespace).MatchingLabels(matchingLabels), pvcList)
		if err != nil && len(pvcList.Items) == 0 {
			return emperror.WrapWith(err, "getting resource failed", "kind", desiredType)
		}
		mountPath := current.(*corev1.PersistentVolumeClaim).Annotations["mountPath"]

		// Creating the first PersistentVolume For Pod
		if len(pvcList.Items) == 0 {
			err = apierrors.NewNotFound(corev1.Resource("PersistentVolumeClaim"), "kafkaBroker")
			if err := client.Create(context.TODO(), desired); err != nil {
				return emperror.WrapWith(err, "creating resource failed", "kind", desiredType)
			}
			log.Info("resource created")
			break
		}
		alreadyCreated := false
		for _, pvc := range pvcList.Items {
			if mountPath == pvc.Annotations["mountPath"] {
				current = pvc.DeepCopyObject()
				alreadyCreated = true
				break
			}
		}
		if !alreadyCreated {
			// Creating the 2+ PersistentVolumes for Pod
			err = apierrors.NewNotFound(corev1.Resource("PersistentVolumeClaim"), "kafkaBroker")
			if err := client.Create(context.TODO(), desired); err != nil {
				return emperror.WrapWith(err, "creating resource failed", "kind", desiredType)
			}
		}
	case *corev1.Pod:
		log = log.WithValues("kind", desiredType)
		log.Info("searching with label because name is empty")

		podList := &corev1.PodList{}
		matchingLabels := map[string]string{
			"kafka_cr": cr.Name,
			"brokerId": desired.(*corev1.Pod).Labels["brokerId"],
		}
		err = client.List(context.TODO(), runtimeClient.InNamespace(current.(*corev1.Pod).Namespace).MatchingLabels(matchingLabels), podList)
		if err != nil && len(podList.Items) == 0 {
			return emperror.WrapWith(err, "getting resource failed", "kind", desiredType)
		}
		if len(podList.Items) == 0 {
			err = apierrors.NewNotFound(corev1.Resource("Pod"), "kafkaBroker")
			if err := client.Create(context.TODO(), desired); err != nil {
				return emperror.WrapWith(err, "creating resource failed", "kind", desiredType)
			}
			scaleErr := scale.UpScaleCluster(desired.(*corev1.Pod).Labels["brokerId"], desired.(*corev1.Pod).Namespace)
			if scaleErr != nil {
				log.Error(err, "graceful upscale failed, or cluster just started")
			}
			log.Info("resource created")
		} else if len(podList.Items) == 1 {
			current = podList.Items[0].DeepCopyObject()
		} else {
			return emperror.WrapWith(errors.New("reconcile failed"), "more then one matching pod found", "labels", matchingLabels)
		}
	}
	if err == nil {
		objectsEquals, err := objectmatch.New(log).Match(current, desired)
		if err != nil {
			log.Error(err, "could not match objects", "kind", desiredType)
		} else if objectsEquals {
			log.V(1).Info("resource is in sync")
			return nil
		}

		switch desired.(type) {
		default:
			return emperror.With(errors.New("unexpected resource type"), "kind", desiredType)
		case *corev1.ConfigMap:
			cm := desired.(*corev1.ConfigMap)
			cm.ResourceVersion = current.(*corev1.ConfigMap).ResourceVersion
			desired = cm
		case *corev1.Service:
			svc := desired.(*corev1.Service)
			svc.ResourceVersion = current.(*corev1.Service).ResourceVersion
			svc.Spec.ClusterIP = current.(*corev1.Service).Spec.ClusterIP
			desired = svc
		case *corev1.Pod:
			err := updateCrWithNodeAffinity(current.(*corev1.Pod), cr, client)
			if err != nil {
				return emperror.WrapWith(err, "updating cr failed")
			}
			err = client.Delete(context.TODO(), current)
			if err != nil {
				return emperror.WrapWith(err, "deleting resource failed", "kind", desiredType)
			}
			return nil
		case *corev1.PersistentVolumeClaim:
			//TODO
			desired = current
		case *appsv1.Deployment:
			deploy := desired.(*appsv1.Deployment)
			deploy.ResourceVersion = current.(*appsv1.Deployment).ResourceVersion
			desired = deploy
		case *appsv1.StatefulSet:
			deploy := desired.(*appsv1.StatefulSet)
			deploy.ResourceVersion = current.(*appsv1.StatefulSet).ResourceVersion
			desired = deploy
		}
		if err := client.Update(context.TODO(), desired); err != nil {
			return emperror.WrapWith(err, "updating resource failed", "kind", desiredType)
		}
		log.Info("resource updated")
	}
	return nil
}
