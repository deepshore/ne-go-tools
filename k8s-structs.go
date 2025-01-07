package negotools

// wrapper for creating commonly used k8s structs

import (
	externalsecretsv1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
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
	VolumeName                 string
	ImagePullSecretName        string
	ContainerName              string
	Image                      string
	PortName                   string
	EnvFromSecretName          string
	VolumeMountPath            string
	VolumeSubPath              string
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

	var configMapVolumeSource corev1.VolumeSource = corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			DefaultMode:          &config.DefaultConfigMapVolumeMode,
			LocalObjectReference: corev1.LocalObjectReference{Name: config.VolumeName},
		},
	}
	//
	var envVars []corev1.EnvVar = make([]corev1.EnvVar, len(config.EnvVarData))
	for key, value := range config.EnvVarData {
		envVars = append(envVars, corev1.EnvVar{Name: key, Value: value})
	}
	var secretEnvSrc corev1.SecretEnvSource = corev1.SecretEnvSource{
		LocalObjectReference: corev1.LocalObjectReference{Name: config.EnvFromSecretName},
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
	var volumeMounts []corev1.VolumeMount = []corev1.VolumeMount{{
		Name:      config.VolumeName,
		MountPath: config.VolumeMountPath,
		SubPath:   config.VolumeSubPath,
	}}
	//
	var volumes []corev1.Volume = []corev1.Volume{
		{
			Name:         config.VolumeName,
			VolumeSource: configMapVolumeSource,
		},
	}
	var containers []corev1.Container = []corev1.Container{
		{
			Name:            config.ContainerName,
			Image:           config.Image,
			ImagePullPolicy: config.ImagePullPolicy,
			EnvFrom: []corev1.EnvFromSource{{
				SecretRef: &secretEnvSrc,
			}},
			Env: envVars,
			Resources: corev1.ResourceRequirements{
				Limits:   resourceLimit,
				Requests: resourceRequest,
			},
			Ports: []corev1.ContainerPort{{
				Name:          config.PortName,
				ContainerPort: config.ContainerPort,
			}},
			LivenessProbe:  &readinessProbe,
			ReadinessProbe: &livenessProbe,
			VolumeMounts:   volumeMounts,
		},
	}
	//
	var podMeta metav1.ObjectMeta = metav1.ObjectMeta{
		Labels: config.PodLabels,
	}
	var podSpec corev1.PodSpec = corev1.PodSpec{
		Volumes:    volumes,
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
	port int32, pathType networking.PathType,
) networking.Ingress {

	var ingressHost string = dnsUri + "." + ingressBaseUrl
	var ingressService networking.IngressServiceBackend = networking.IngressServiceBackend{
		Name: serviceName,
		Port: networking.ServiceBackendPort{
			Name:   serviceName,
			Number: port,
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

func GeneratePushSecret(
	pushSecretName, targetSecretName, sourceSecretName, namespaceName string,
	secretStoreName, secretStoreKind string, keys []string, refreshInterval metav1.Duration,
) externalsecretsv1alpha1.PushSecret {

	var metadata metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      pushSecretName,
		Namespace: namespaceName,
	}

	var pushSecretData []externalsecretsv1alpha1.PushSecretData = make([]externalsecretsv1alpha1.PushSecretData, len(keys))
	for _, key := range keys {
		pushSecretData = append(pushSecretData, externalsecretsv1alpha1.PushSecretData{
			Match: externalsecretsv1alpha1.PushSecretMatch{
				SecretKey: key,
				RemoteRef: externalsecretsv1alpha1.PushSecretRemoteRef{
					RemoteKey: targetSecretName,
					Property:  key,
				},
			},
		})
	}

	var pushSecret externalsecretsv1alpha1.PushSecret = externalsecretsv1alpha1.PushSecret{
		ObjectMeta: metadata,
		Spec: externalsecretsv1alpha1.PushSecretSpec{
			RefreshInterval: &refreshInterval,
			SecretStoreRefs: []externalsecretsv1alpha1.PushSecretStoreRef{
				{
					Name: secretStoreName,
					Kind: secretStoreKind,
				},
			},
			Selector: externalsecretsv1alpha1.PushSecretSelector{
				Secret: &externalsecretsv1alpha1.PushSecretSecret{Name: sourceSecretName},
			},
			Data: pushSecretData,
		},
	}
	return pushSecret
}

func GenerateExternalSecret(
	name, namespace, secretStoreName, secretStoreKind, targetSecretName, remoteSecretName string,
	refreshInterval metav1.Duration,
	externalSecretKeyMapping map[string]string,
) externalsecretsv1alpha1.ExternalSecret {

	var externalSecretDataSpec []externalsecretsv1alpha1.ExternalSecretData = make([]externalsecretsv1alpha1.ExternalSecretData, len(externalSecretKeyMapping))
	for localKey, remoteKey := range externalSecretKeyMapping {
		externalSecretDataSpec = append(externalSecretDataSpec, externalsecretsv1alpha1.ExternalSecretData{
			SecretKey: localKey,
			RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
				Key:      remoteSecretName,
				Property: remoteKey,
			},
		})
	}
	var spec externalsecretsv1alpha1.ExternalSecret = externalsecretsv1alpha1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: externalsecretsv1alpha1.ExternalSecretSpec{
			SecretStoreRef: externalsecretsv1alpha1.SecretStoreRef{
				Name: secretStoreName,
				Kind: secretStoreKind,
			},
			Target: externalsecretsv1alpha1.ExternalSecretTarget{
				Name: targetSecretName,
			},
			RefreshInterval: &refreshInterval,
			Data:            externalSecretDataSpec,
		},
	}
	return spec
}