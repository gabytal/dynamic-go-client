package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)

	namespace := "argo"
	app, err := GetArgoApp(dynamic, ctx, "demo-multi-canary-gcp", namespace)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// change OwnerReferences
	ownerReferences := app.GetOwnerReferences()
	fmt.Printf("old ownerReferences: %s", ownerReferences[0].Name)

	ownerReferences[0].Name = "new-owner"
	ownerReferences[0].UID = "new-uid"
	ownerReferences[0].Controller = new(bool)
	ownerReferences[0].BlockOwnerDeletion = new(bool)
	ownerReferences[0].Kind = "new-kind"

	app.SetOwnerReferences(ownerReferences)

	err = updateResourceDynamically(dynamic, ctx, app)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("\nsuccess\n")

}

func GetArgoApp(dynamic dynamic.Interface, ctx context.Context, name string, namespace string) (
	*unstructured.Unstructured, error) {

	resourceId := schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}
	app, err := dynamic.Resource(resourceId).Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return &unstructured.Unstructured{}, err
	}
	return app, nil
}

func updateResourceDynamically(dynamic dynamic.Interface, ctx context.Context, obj *unstructured.Unstructured) (err error) {
	resourceId := schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}
	_, err = dynamic.Resource(resourceId).Namespace(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})

	if err != nil {
		return err
	}
	return nil
}
