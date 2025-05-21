package negotools

// wrapper for creating commonly used k8s structs

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GenerateSecret(
	name, namespaceName string, data map[string]string,
) corev1.Secret {

	var metadata metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: namespaceName,
	}
	var secret corev1.Secret = corev1.Secret{
		ObjectMeta: metadata,
		Type:       "Opaque",
		StringData: data,
	}
	return secret
}

func GenerateConfigMap(name, namespaceName string, data map[string]string,
) corev1.ConfigMap {

	var metadata metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: namespaceName,
	}
	var configMap corev1.ConfigMap = corev1.ConfigMap{
		ObjectMeta: metadata,
		Data:       data,
	}
	return configMap
}

type ProbeSpec struct {
	HttpGetPath         string
	HttpGetPort         int32
	InitialDelaySeconds int32
	TimeoutSeconds      int32
	PeriodSeconds       int32
	FailureThreshold    int32
	SuccessThreshold    int32
}
type DeploymentConfig struct {
	Name                       string
	Namespace                  string
	Volumes                    []corev1.Volume
	ImagePullSecretName        string
	ContainerName              string
	Image                      string
	PortName                   string
	EnvFromSecretNames         []string
	EnvFromConfigMapNames      []string
	VolumeMounts               []corev1.VolumeMount
	ImagePullPolicy            corev1.PullPolicy
	ContainerPort              int32
	DefaultConfigMapVolumeMode int32
	Replicas                   int32
	EnvVarData                 map[string]string
	PodLabels                  map[string]string
	MatchLabels                map[string]string
	CpuRequestMilli            int64
	CpuLimitMilli              int64
	MemoryRequestMi            int64
	MemoryLimitMi              int64
	LivenessProbeSpec          ProbeSpec
	ReadinessProbeSpec         ProbeSpec
}

// use a struct to avoid mistakes in the order of arguments and keep things
// read- and debugable
func GenerateDeployment(config DeploymentConfig) appsv1.Deployment {

	var envVars []corev1.EnvVar = []corev1.EnvVar{}
	for key, value := range config.EnvVarData {
		LogTrace(fmt.Sprintf("Adding ENV %q=%q to DeploymentSpec", key, value))
		envVars = append(envVars, corev1.EnvVar{Name: key, Value: value})
	}
	var envFromSources []corev1.EnvFromSource = []corev1.EnvFromSource{}
	for _, configMapName := range config.EnvFromConfigMapNames {
		var ref corev1.EnvFromSource = corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: string(configMapName)},
			},
		}
		envFromSources = append(envFromSources, ref)
	}
	for _, secretName := range config.EnvFromSecretNames {
		var ref corev1.EnvFromSource = corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: string(secretName)},
			},
		}
		envFromSources = append(envFromSources, ref)
	}
	var cpuRequest *resource.Quantity = resource.NewMilliQuantity(config.CpuRequestMilli, resource.DecimalSI)
	var memoryRequest *resource.Quantity = resource.NewQuantity(config.MemoryRequestMi*1024*1024, resource.BinarySI)
	var resourceRequest corev1.ResourceList = corev1.ResourceList{
		corev1.ResourceName("Cpu"):    *cpuRequest,
		corev1.ResourceName("Memory"): *memoryRequest,
	}
	var cpuLimit *resource.Quantity = resource.NewMilliQuantity(config.CpuLimitMilli, resource.DecimalSI)
	var memoryLimit *resource.Quantity = resource.NewQuantity(config.MemoryLimitMi*1024*1024, resource.BinarySI)
	var resourceLimit corev1.ResourceList = corev1.ResourceList{
		corev1.ResourceName("Cpu"):    *cpuLimit,
		corev1.ResourceName("Memory"): *memoryLimit,
	}
	var livenessProbe corev1.Probe = corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: config.LivenessProbeSpec.HttpGetPath,
				Port: intstr.FromInt32(config.LivenessProbeSpec.HttpGetPort),
			},
		},
		InitialDelaySeconds: config.LivenessProbeSpec.InitialDelaySeconds,
		TimeoutSeconds:      config.LivenessProbeSpec.TimeoutSeconds,
		PeriodSeconds:       config.LivenessProbeSpec.PeriodSeconds,
		FailureThreshold:    config.LivenessProbeSpec.FailureThreshold,
		SuccessThreshold:    config.LivenessProbeSpec.SuccessThreshold,
	}
	var readinessProbe corev1.Probe = corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: config.ReadinessProbeSpec.HttpGetPath,
				Port: intstr.FromInt32(config.ReadinessProbeSpec.HttpGetPort),
			},
		},
		InitialDelaySeconds: config.ReadinessProbeSpec.InitialDelaySeconds,
		TimeoutSeconds:      config.ReadinessProbeSpec.TimeoutSeconds,
		PeriodSeconds:       config.ReadinessProbeSpec.PeriodSeconds,
		FailureThreshold:    config.ReadinessProbeSpec.FailureThreshold,
		SuccessThreshold:    config.ReadinessProbeSpec.SuccessThreshold,
	}
	//
	var containers []corev1.Container = []corev1.Container{
		{
			Name:            config.ContainerName,
			Image:           config.Image,
			ImagePullPolicy: config.ImagePullPolicy,
			EnvFrom:         envFromSources,
			Env:             envVars,
			Resources: corev1.ResourceRequirements{
				Limits:   resourceLimit,
				Requests: resourceRequest,
			},
			Ports: []corev1.ContainerPort{{
				Name:          config.PortName,
				ContainerPort: config.ContainerPort,
			}},
			LivenessProbe:  &livenessProbe,
			ReadinessProbe: &readinessProbe,
			VolumeMounts:   config.VolumeMounts,
		},
	}
	//
	var podMeta metav1.ObjectMeta = metav1.ObjectMeta{
		Labels: config.PodLabels,
	}
	var podSpec corev1.PodSpec = corev1.PodSpec{
		Volumes:    config.Volumes,
		Containers: containers,
		ImagePullSecrets: []corev1.LocalObjectReference{{
			Name: config.ImagePullSecretName,
		}},
	}
	//
	var selector metav1.LabelSelector = metav1.LabelSelector{
		MatchLabels: config.MatchLabels,
	}
	var podTemplate corev1.PodTemplateSpec = corev1.PodTemplateSpec{
		ObjectMeta: podMeta,
		Spec:       podSpec,
	}
	//
	var meta metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      config.Name,
		Namespace: config.Namespace,
	}
	var spec appsv1.DeploymentSpec = appsv1.DeploymentSpec{
		Replicas: &config.Replicas,
		Selector: &selector,
		Template: podTemplate,
	}
	//
	var deployment appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: meta,
		Spec:       spec,
	}
	return deployment
}

func GenerateIngress(
	name, namespace, dnsUri, ingressBaseUrl, serviceName, path, ingressClassName string,
	k8sServiceName string, pathType networking.PathType,
) networking.Ingress {

	var ingressHost string = dnsUri + "." + ingressBaseUrl
	var ingressService networking.IngressServiceBackend = networking.IngressServiceBackend{
		Name: k8sServiceName,
		Port: networking.ServiceBackendPort{
			Name: serviceName,
		},
	}
	var ingressPath networking.HTTPIngressPath = networking.HTTPIngressPath{
		Path:     path,
		PathType: &pathType,
		Backend: networking.IngressBackend{
			Service: &ingressService,
		},
	}
	var paths []networking.HTTPIngressPath = []networking.HTTPIngressPath{ingressPath}
	var httpIngressRuleValue networking.HTTPIngressRuleValue = networking.HTTPIngressRuleValue{
		Paths: paths,
	}
	var ingressRuleValue networking.IngressRuleValue = networking.IngressRuleValue{
		HTTP: &httpIngressRuleValue,
	}
	var ingressRule networking.IngressRule = networking.IngressRule{
		Host:             ingressHost,
		IngressRuleValue: ingressRuleValue,
	}
	var ingressRules []networking.IngressRule = []networking.IngressRule{
		ingressRule,
	}
	var ingressSpec networking.Ingress = networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networking.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules:            ingressRules,
		},
	}
	return ingressSpec
}
