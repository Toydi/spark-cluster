package sparkcluster

import (
	sparkv1alpha1 "github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	gitHttpProxy     = "http://114.212.189.147:8118"
	gitHttpProxyKey  = "HTTP_PROXY"
	gitHttpsProxyKey = "HTTPS_PROXY"
)

func (r *ReconcileSparkCluster) newMasterPod(instance *sparkv1alpha1.SparkCluster) *corev1.Pod {
	var volumeMounts []corev1.VolumeMount
	var volumeMounts_master []corev1.VolumeMount
	var volumeMounts_vscode []corev1.VolumeMount
	var volumes []corev1.Volume

	if instance.Spec.NFS {
		// nfs := corev1.NFSVolumeSource{Server: ShareServer, Path: "/hadoop/share-dir"}
		// volumeMounts = append(volumeMounts, corev1.VolumeMount{Name: "share-dir", MountPath: nfs.Path})
		// volumes = append(volumes, corev1.Volume{Name: "share-dir", VolumeSource: corev1.VolumeSource{NFS: &nfs}})
		nfs := corev1.NFSVolumeSource{Server: ShareServer, Path: "/hadoop/share-data"}
		volumes = append(volumes, corev1.Volume{Name: "hadoop-share-volume", VolumeSource: corev1.VolumeSource{NFS: &nfs}})
		volumeMounts_master = append(volumeMounts_master, corev1.VolumeMount{Name: "hadoop-share-volume", MountPath: "/hadoop/share-data"})
		volumeMounts_vscode = append(volumeMounts_vscode, corev1.VolumeMount{Name: "hadoop-share-volume", MountPath: "/hadoop/share-data"})
	}

	if instance.Spec.PvcEnable {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{Name: "dfs", MountPath: "/root/hdfs/namenode"})
		pvc := corev1.PersistentVolumeClaimVolumeSource{ClaimName: masterPvc(instance)}
		volumes = append(volumes, corev1.Volume{Name: "dfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &pvc}})
	}
	volumeMounts_master = append(volumeMounts_master, corev1.VolumeMount{Name: "code", MountPath: "/root/code"})
	volumeMounts_vscode = append(volumeMounts_vscode, corev1.VolumeMount{Name: "code", MountPath: "/workspace"})
	// volumeMounts_git=append(volumeMounts_git,corev1.VolumeMount{Name:"code",MountPath:"/root/code"})
	pvc := corev1.PersistentVolumeClaimVolumeSource{ClaimName: instance.Spec.ClusterPrefix + "-" + "vscode-pvc"}
	volumes = append(volumes, corev1.Volume{Name: "code", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &pvc}})
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      masterName(instance),
			Namespace: instance.Namespace,
			Labels:    masterLabel(instance),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  masterName(instance),
					Image: MasterImage,
					// Command:         []string{"bash", "-c", "/etc/master-bootstrap.sh " + fmt.Sprintf("%d", instance.Spec.SlaveNum)},
					ImagePullPolicy: "Always",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8020,
						},
						{
							ContainerPort: 50070,
						},
						{
							ContainerPort: 50470,
						},
					},
					Env:          []corev1.EnvVar{{Name: "PREFIX", Value: instance.Spec.ClusterPrefix}},
					Resources:    instance.Spec.Resources,
					VolumeMounts: volumeMounts_master,
				},
				{
					Name:            masterName(instance) + "-" + "codeserver",
					Image:           VscodeImage,
					Args:            []string{"--allow-http", "--no-auth"},
					ImagePullPolicy: "Always",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8443,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "PREFIX",
							Value: instance.Spec.ClusterPrefix,
						},
						{
							Name:  gitHttpProxyKey,
							Value: gitHttpProxy,
						},
						{
							Name:  gitHttpsProxyKey,
							Value: gitHttpProxy,
						},
						{
							Name:  "GIT_USERNAME",
							Value: instance.Spec.GitUserName,
						},
						{
							Name:  "GIT_USEREMAIL",
							Value: instance.Spec.GitUserEmail,
						},
					},
					VolumeMounts: volumeMounts_vscode,
				},
			},
			Volumes: volumes,
		},
	}
}

//gitSidecarCommand := fmt.Sprintf("git clone %s %s && git config --global user.name %s && git config --global user.email %s && git config --global http.proxy %s && git config --global https.proxy %s", g.Repo, path.Join(g.MountPath, g.Name),g.GitUserName,g.GitUserEmail,gitHttpProxy,gitHttpProxy)
// gitHttpProxy     = "http://114.212.189.147:8118"
// gitHttpProxyKey  = "HTTP_PROXY"
// gitHttpsProxyKey = "HTTPS_PROXY"

func (r *ReconcileSparkCluster) newSlavePod(instance *sparkv1alpha1.SparkCluster, index int) *corev1.Pod {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	if instance.Spec.PvcEnable {
		volumeMounts = []corev1.VolumeMount{{Name: "dfs", MountPath: "/root/hdfs/datanode"}}
		pvc := corev1.PersistentVolumeClaimVolumeSource{ClaimName: slavePvc(instance, index)}
		volumes = []corev1.Volume{{Name: "dfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &pvc}}}
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      slaveName(instance, index),
			Namespace: instance.Namespace,
			Labels:    slaveLabel(instance, index),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            Slave,
					Image:           SlaveImage,
					ImagePullPolicy: "Always",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 50010,
						},
						{
							ContainerPort: 50020,
						},
						{
							ContainerPort: 50075,
						},
						{
							ContainerPort: 50475,
						},
					},
					Env:          []corev1.EnvVar{{Name: "PREFIX", Value: instance.Spec.ClusterPrefix}},
					Resources:    instance.Spec.Resources,
					VolumeMounts: volumeMounts,
				},
			},
			Volumes: volumes,
		},
	}
}

func (r *ReconcileSparkCluster) newMasterService(instance *sparkv1alpha1.SparkCluster) *corev1.Service {
	labels := map[string]string{
		"app": masterName(instance),
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      masterName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "rpc",
					Port: 8020,
				},
				{
					Name: "p1",
					Port: 50070,
				},
				{
					Name: "p2",
					Port: 50470,
				},
			},
			Selector: labels,
		},
	}
}

func (r *ReconcileSparkCluster) newUIService(instance *sparkv1alpha1.SparkCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "ssh",
			Port: 22,
		},
		{
			Name: "hdfs-client",
			Port: 9000,
		},
		{
			Name: "resource-manager",
			Port: 8088,
		},
		{
			Name: "name-node",
			Port: 50070,
		},
		{
			Name: "code-server",
			Port: 8443,
		},
		{
			Name: "spark",
			Port: 8080,
		},
		{
			Name: "spark-shell",
			Port: 4040,
		}}
	ports = append(ports, instance.Spec.Ports...)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.ClusterPrefix + "-ui-service",
			Namespace: instance.Namespace,
			Labels:    masterLabel(instance),
		},
		Spec: corev1.ServiceSpec{
			Type:     "NodePort",
			Ports:    ports,
			Selector: masterLabel(instance),
		},
	}
}

func (r *ReconcileSparkCluster) newSlaveService(instance *sparkv1alpha1.SparkCluster, index int) *corev1.Service {
	serviceName := slaveName(instance, index)
	labels := map[string]string{
		"app": serviceName,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name: "rpc",
					Port: 8020,
				},
				{
					Name: "p1",
					Port: 50070,
				},
				{
					Name: "p2",
					Port: 50470,
				},
			},
			Selector: map[string]string{"app": serviceName},
		},
	}
}
func (r *ReconcileSparkCluster) newVscodePvc(instance *sparkv1alpha1.SparkCluster) *corev1.PersistentVolumeClaim {
	vscodePvcPrefix := instance.Spec.ClusterPrefix
	vscodePvcNamespace := instance.Namespace
	storageClassName := StorageClassName
	accessModes := make([]corev1.PersistentVolumeAccessMode, 1)
	accessModes[0] = corev1.ReadWriteMany
	resourceList := corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("20Gi")}
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vscodePvcPrefix + "-" + "vscode-pvc",
			Namespace: vscodePvcNamespace,
			Labels:    masterLabel(instance),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources:        corev1.ResourceRequirements{Requests: resourceList},
		},
	}
}

func (r *ReconcileSparkCluster) newMasterPvc(instance *sparkv1alpha1.SparkCluster) *corev1.PersistentVolumeClaim {
	storageClassName := StorageClassName
	accessModes := make([]corev1.PersistentVolumeAccessMode, 1)
	accessModes[0] = corev1.ReadWriteOnce
	resourceList := corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("5Gi")}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MasterPvc,
			Namespace: instance.Namespace,
			Labels:    masterLabel(instance),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources:        corev1.ResourceRequirements{Requests: resourceList},
		},
	}
}

func (r *ReconcileSparkCluster) newSlavePvc(instance *sparkv1alpha1.SparkCluster, index int) *corev1.PersistentVolumeClaim {
	storageClassName := StorageClassName
	accessModes := make([]corev1.PersistentVolumeAccessMode, 1)
	accessModes[0] = corev1.ReadWriteOnce
	q := resource.MustParse("5Gi")
	resourceList := corev1.ResourceList{corev1.ResourceStorage: q}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      slavePvc(instance, index),
			Namespace: instance.Namespace,
			Labels:    slaveLabel(instance, index),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources:        corev1.ResourceRequirements{Requests: resourceList},
		},
	}
}
