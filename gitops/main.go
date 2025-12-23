package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v3"
)

func main() {
	// CLI inputs
	serviceFlag := flag.String("service", "", "Service name (auth/order)")
	imageFlag := flag.String("image", "", "Image tag to deploy")
	replicasFlag := flag.Int("replicas", 0, "Number of replicas")
	flag.Parse()

	if *serviceFlag == "" || *imageFlag == "" || *replicasFlag == 0 {
		log.Fatal("Usage: --service <name> --image <tag> --replicas <num>")
	}

	// 1️⃣ Read config.yaml
	configPath := "config.yaml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var cfg map[string]interface{}
	yaml.Unmarshal(data, &cfg)

	services := cfg["services"].(map[string]interface{})
	svc, ok := services[*serviceFlag].(map[string]interface{})
	if !ok {
		log.Fatalf("Service %s not found in config.yaml", *serviceFlag)
	}

	// 2️⃣ Update config.yaml values
	svc["image"] = *imageFlag
	svc["replicas"] = *replicasFlag

	out, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(configPath, out, 0644)
	fmt.Println("✅ config.yaml updated")

	// 3️⃣ Update deployment.yaml
	deployPath := "../" + svc["path"].(string)
	deployData, err := os.ReadFile(deployPath)
	if err != nil {
		log.Fatal(err)
	}

	deployStr := string(deployData)
	deployStr = replaceYAML(deployStr, `image:.*`, fmt.Sprintf("image: %s", *imageFlag))
	deployStr = replaceYAML(deployStr, `replicas:.*`, fmt.Sprintf("replicas: %d", *replicasFlag))

	os.WriteFile(deployPath, []byte(deployStr), 0644)
	fmt.Println("✅ deployment.yaml updated")

	// 4️⃣ Apply deployment to cluster
	cmd := exec.Command("kubectl", "apply", "-f", deployPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Kubernetes deployment applied")

	// 5️⃣ Git commit & push
	exec.Command("git", "add", ".").Run()
	commitMsg := fmt.Sprintf("Update %s: image=%s replicas=%d", *serviceFlag, *imageFlag, *replicasFlag)
	exec.Command("git", "commit", "-m", commitMsg).Run()
	exec.Command("git", "push").Run()
	fmt.Println("✅ Changes pushed to GitHub")
}

// Helper for regex replacement
func replaceYAML(input, pattern, replace string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, replace)
}
