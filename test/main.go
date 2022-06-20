package main

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
	cronjobv1 "tutorial.kubebuilder.io/project/api/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client

func main() {
	const (
		CronjobName      = "test-cronjob"
		CronjobNamespace = "default"
		JobName          = "test-job"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	ctx := context.Background()
	cronJob := &cronjobv1.CronJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch.tutorial.kubebuilder.io/v1",
			Kind:       "CronJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      CronjobName,
			Namespace: CronjobNamespace,
		},
		Spec: cronjobv1.CronJobSpec{
			Schedule: "1 * * * *",
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					// For simplicity, we only fill out the required fields.
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							// For simplicity, we only fill out the required fields.
							Containers: []v1.Container{
								{
									Name:  "hello",
									Image: "bazel/driver:image",
								},
							},
							RestartPolicy: v1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	scheme := runtime.NewScheme()
	cronjobv1.AddToScheme(scheme)
	k8sClient, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(k8sClient.Create(ctx, cronJob))
}
