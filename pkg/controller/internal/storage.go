package internal

import (
	"strings"

	"fmt"
	"os"

	"path"

	v1 "k8s.io/api/core/v1"
)

const (
	defaultStorageClass = "cephfs"
	defaultStorageSize  = "1Gi"
)

func GetStorageClassName() string {
	storageClass := os.Getenv("STORAGE_CLASS")
	if len(storageClass) == 0 {
		return defaultStorageClass
	}
	return storageClass
}

type Volume struct {
	Name      string
	MountPath string
	ReadOnly  bool
}

func (v Volume) AddToPodSpec(spec *v1.PodSpec) {
	spec.Volumes = append(spec.Volumes, v1.Volume{
		Name: v.Name,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: v.Name,
			},
		},
	})
	if len(spec.Containers) > 0 {
		for i := 0; i < len(spec.Containers); i++ {
			spec.Containers[i].VolumeMounts = append(spec.Containers[i].VolumeMounts, v1.VolumeMount{
				MountPath: v.MountPath,
				Name:      v.Name,
				ReadOnly:  v.ReadOnly,
			})
		}
	}
}

func (v Volume) AddToPod(pod *v1.Pod) {
	v.AddToPodSpec(&pod.Spec)
}

type GitRepository struct {
	VolumeName   string
	Repo         string
	Name         string
	Owner        string
	MountPath    string
	GitUserName  string
	GitUserEmail string
}

const (
	gitSidecarImage         = "registry.njuics.cn/library/git:v2.7.2"
	gitSidecarContainerName = "git"

	gitHttpProxy     = "http://114.212.189.147:8118"
	gitHttpProxyKey  = "HTTP_PROXY"
	gitHttpsProxyKey = "HTTPS_PROXY"
)

func (g GitRepository) AddToPodSpec(spec *v1.PodSpec) {
	// gitSidecarCommand := fmt.Sprintf("git clone %s %s && git config --global user.name %s && git config --global user.email %s && git config --global http.proxy %s && git config --global https.proxy %s", g.Repo, path.Join(g.MountPath, g.Name),g.GitUserName,g.GitUserEmail,gitHttpProxy,gitHttpProxy)
	gitSidecarCommand := fmt.Sprintf("git clone %s %s", g.Repo, path.Join(g.MountPath, g.Name))
	//git username  git email. .  http_proxy
	spec.InitContainers = append(spec.InitContainers, v1.Container{
		Name:    gitSidecarContainerName,
		Image:   gitSidecarImage,
		Command: strings.Fields(gitSidecarCommand),
		Env: []v1.EnvVar{
			{
				Name:  gitHttpProxyKey,
				Value: gitHttpProxy,
			},
			{
				Name:  gitHttpsProxyKey,
				Value: gitHttpProxy,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      g.VolumeName,
				MountPath: g.MountPath,
			},
		},
	})
}
