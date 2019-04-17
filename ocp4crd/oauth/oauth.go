package oauth

import (
	configv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	// TODO: Is this line needed at all? It may be superflous to
	// ocp3.go/init()/configv1.InstallLegacy(scheme.Scheme)
	oauthv1.Install(scheme.Scheme)
}

// TODO: Generated yamls are results of pure imagination. Structure and consistency
// must be reviewed and fixed. I guess this code can be simplified once we know
// how the output should look exactly.

// reference:
//   [v3] OCPv3:
//   - [1] https://docs.openshift.com/container-platform/3.11/install_config/configuring_authentication.html#identity_providers_master_config
//   [v4] OCPv4:
//   - [2] htpasswd: https://docs.openshift.com/container-platform/4.0/authentication/understanding-identity-provider.html
//   - [3] github: https://docs.openshift.com/container-platform/4.0/authentication/identity_providers/configuring-github-identity-provider.html

// Structures defining custom resource definitions / manifests / yamls
// TODO: figure out the OKD terminology

// Shared CRD part, present in all types of OAuth CRDs
type v4OAuthCRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	MetaData   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		IdentityProviders []interface{} `yaml:"identityProviders"`
	} `yaml:"spec"`
}

// Generate converts OCPv3 OAuth to OCPv4 OAuth Custom Resources
func Generate(masterconfig configv1.MasterConfig) (*v4OAuthCRD, error) {
	var auth = masterconfig.OAuthConfig.DeepCopy()
	var err error

	var crd v4OAuthCRD
	crd.APIVersion = "config.openshift.io/v1"
	crd.Kind = "OAuth"
	crd.MetaData.Name = "cluster"

	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	for _, p := range auth.IdentityProviders {
		p.Provider.Object, _, err = serializer.Decode(p.Provider.Raw, nil, nil)
		if err != nil {
			return nil, err
		}

		switch kind := p.Provider.Object.GetObjectKind().GroupVersionKind().Kind; kind {
		case "HTPasswdPasswordIdentityProvider":
			idP := buildHTPasswdIP(serializer, p)
			crd.Spec.IdentityProviders = append(crd.Spec.IdentityProviders, idP)
		case "GitHubIdentityProvider":
			idP := buildGitHubIP(serializer, p)
			crd.Spec.IdentityProviders = append(crd.Spec.IdentityProviders, idP)
		default:
			logrus.Print("can't handle: ", kind)
		}
	}

	return &crd, nil
}

// PrintCRD Print generated CRD
func PrintCRD(crd *v4OAuthCRD) {
	yamlBytes, err := yaml.Marshal(&crd)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Print(string(yamlBytes))
}