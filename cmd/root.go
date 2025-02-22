package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang-k8s-temp-access/internal/k8s"
)

var (
	kubeconfig string
	namespaces []string
	resources  []string
	expiration string
)

var rootCmd = &cobra.Command{
	Use:   "kube-temp-access",
	Short: "A tool to create temporary Kubernetes access",
	Long:  `Creates temporary service accounts and roles for Kubernetes dashboard login, auto-deleting after a configurable duration.`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create temporary Kubernetes access",
	Long:  `Creates a temporary service account, assigns specified permissions, generates a token, and schedules deletion.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := k8s.NewClient(kubeconfig)
		if err != nil {
			fmt.Printf("Error initializing Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		for _, ns := range namespaces {
			if err := k8s.CreateTemporaryAccess(clientset, ns, resources, expiration); err != nil {
				fmt.Printf("Error creating access in namespace %s: %v\n", ns, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", os.Getenv("HOME")+"/.kube/config", "Path to kubeconfig file")
	createCmd.Flags().StringSliceVarP(&namespaces, "namespace", "n", []string{"default"}, "Namespaces to create resources in (comma-separated or multiple flags)")
	createCmd.Flags().StringSliceVarP(&resources, "resources", "r", []string{"view"}, "Resources to grant access to (e.g., deployments, pods, statefulsets, view)")
	createCmd.Flags().StringVarP(&expiration, "expiration", "e", "15m", "Token expiration and cleanup duration (e.g., 15m, 1h)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}