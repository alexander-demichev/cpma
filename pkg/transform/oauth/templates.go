package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	notlegacy "github.com/openshift/api/config/v1"
	configv1 "github.com/openshift/api/legacyconfig/v1"
)

const (
	loginSecret             = "templates-login-secret"
	errorSecret             = "templates-error-secret"
	providerSelectionSecret = "templates-providerselect-secret"
)

func translateTemplates(templates configv1.OAuthTemplates) (*notlegacy.OAuthTemplates, []*secrets.Secret, error) {
	var templateSecrets []*secrets.Secret

	translatedTemplates := &notlegacy.OAuthTemplates{
		Login: notlegacy.SecretNameReference{
			Name: loginSecret,
		},
		Error: notlegacy.SecretNameReference{
			Name: errorSecret,
		},
		ProviderSelection: notlegacy.SecretNameReference{
			Name: providerSelectionSecret,
		},
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(templates.Login))
	secret, err := secrets.GenSecret(loginSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, err
	}
	templateSecrets = append(templateSecrets, secret)

	encoded = base64.StdEncoding.EncodeToString([]byte(templates.Error))
	secret, err = secrets.GenSecret(errorSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, err
	}
	templateSecrets = append(templateSecrets, secret)

	encoded = base64.StdEncoding.EncodeToString([]byte(templates.ProviderSelection))
	secret, err = secrets.GenSecret(providerSelectionSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, err
	}
	templateSecrets = append(templateSecrets, secret)

	return translatedTemplates, templateSecrets, nil
}
