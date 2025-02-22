package k8s

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"golang-k8s-temp-access/internal/utils"
)

func CreateTemporaryAccess(client *kubernetes.Clientset, namespace string, resources []string, expiration string) error {
	uid := utils.GenerateUUID()
	saName := "temp-sa-" + uid
	crName := "temp-role-" + uid
	crbName := "temp-binding-" + uid
	jobName := "temp-deletion-job-" + uid

	// Parse expiration duration
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return fmt.Errorf("invalid expiration duration: %v", err)
	}

	// Create Service Account
	if err := createServiceAccount(client, namespace, saName); err != nil {
		return fmt.Errorf("creating service account: %v", err)
	}
	fmt.Printf("Created service account %s in namespace %s\n", saName, namespace)

	// Determine role to use
	roleName := "view"
	if !contains(resources, "view") {
		roleName = crName
		if err := createCustomClusterRole(client, crName, resources); err != nil {
			return fmt.Errorf("creating custom cluster role: %v", err)
		}
		fmt.Printf("Created custom cluster role %s with resources %v\n", crName, resources)
	}

	// Create ClusterRoleBinding
	if err := createClusterRoleBinding(client, crbName, saName, namespace, roleName); err != nil {
		return fmt.Errorf("creating cluster role binding: %v", err)
	}
	fmt.Printf("Created cluster role binding %s\n", crbName)

	// Get Token with specified expiration
	expirationSeconds := int64(duration.Seconds())
	token, err := GetServiceAccountToken(client, namespace, saName, expirationSeconds)
	if err != nil {
		return fmt.Errorf("getting token: %v", err)
	}
	fmt.Printf("Token for dashboard login in namespace %s: %s\n", namespace, token)

	// Create Deletion Job
	if err := createDeletionJob(client, namespace, jobName, saName, crbName, crName, expiration); err != nil {
		return fmt.Errorf("creating deletion job: %v", err)
	}
	fmt.Printf("Created deletion job %s in namespace %s, resources will be deleted in %s\n", jobName, namespace, expiration)

	return nil
}

func createServiceAccount(client *kubernetes.Clientset, namespace, name string) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := client.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), sa, metav1.CreateOptions{})
	return err
}

func createCustomClusterRole(client *kubernetes.Clientset, name string, resources []string) error {
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"", "apps"},
				Verbs:     []string{"get", "list", "watch"},
				Resources: resources,
			},
		},
	}
	_, err := client.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{})
	return err
}

func createClusterRoleBinding(client *kubernetes.Clientset, roleBindingName, saName, namespace, roleName string) error {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      saName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     roleName,
		},
	}
	_, err := client.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{})
	return err
}

func createDeletionJob(client *kubernetes.Clientset, namespace, jobName, saName, crbName, crName, expiration string) error {
	ttlSeconds := int32(30)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSeconds,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "deletion-container",
							Image: "bitnami/kubectl:latest",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("sleep %s && kubectl delete serviceaccount $SERVICE_ACCT_NAME -n $NAMESPACE && kubectl delete clusterrolebinding $CLUSTERROLEBINDING_NAME && kubectl delete clusterrole $CLUSTERROLE_NAME", expiration),
							},
							Env: []corev1.EnvVar{
								{Name: "SERVICE_ACCT_NAME", Value: saName},
								{Name: "NAMESPACE", Value: namespace},
								{Name: "CLUSTERROLEBINDING_NAME", Value: crbName},
								{Name: "CLUSTERROLE_NAME", Value: crName},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}