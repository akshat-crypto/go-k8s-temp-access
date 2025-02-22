package k8s

import (
	"context"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetServiceAccountToken(client *kubernetes.Clientset, namespace, name string, expirationSeconds int64) (string, error) {
	tokenRequest := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &expirationSeconds,
		},
	}
	tokenResponse, err := client.CoreV1().ServiceAccounts(namespace).CreateToken(context.TODO(), name, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return tokenResponse.Status.Token, nil
}