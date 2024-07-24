/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kdefault

import (
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &Kustomization{}

// Kustomization scaffolds a file that defines the kustomization scheme for the default overlay folder
type Kustomization struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
}

// SetTemplateDefaults implements file.Template
func (f *Kustomization) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", "default", "kustomization.yaml")
	}

	f.TemplateBody = kustomizeTemplate

	f.IfExistsAction = machinery.Error

	return nil
}

const kustomizeTemplate = `# Adds namespace to all resources.
namespace: {{ .ProjectName }}-system

# Value of this field is prepended to the
# names of all resources, e.g. a deployment named
# "wordpress" becomes "alices-wordpress".
# Note that it should also match with the prefix (text before '-') of the namespace
# field above.
namePrefix: {{ .ProjectName }}-

# Labels to add to all resources and selectors.
#labels:
#- includeSelectors: true
#  pairs:
#    someName: someValue

resources:
#- ../crd
- ../rbac
- ../manager
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
#- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
#- ../certmanager
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'.
#- ../prometheus
# [METRICS] Expose the controller manager metrics service.
- metrics_service.yaml
# [NETWORK POLICY] Protect the /metrics endpoint and Webhook Server with NetworkPolicy.
# Only Pod(s) running a namespace labeled with 'metrics: enabled' will be able to gather the metrics.
# Only CR(s) which requires webhooks and are applied on namespaces labeled with 'webhooks: enabled' will
# be able to communicate with the Webhook Server.
#- ../network-policy

# Uncomment the patches line if you enable Metrics, and/or are using webhooks and cert-manager
patches:
# [METRICS] The following patch will enable the metrics endpoint using HTTPS and the port :8443.
# More info: https://book.kubebuilder.io/reference/metrics
- path: manager_metrics_patch.yaml
  target:
    kind: Deployment

# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
#- path: manager_webhook_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
# Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
# 'CERTMANAGER' needs to be enabled to use ca injection
#- path: webhookcainjection_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
# Uncomment the following replacements to add the cert-manager CA injection annotations
#replacements:
#  - source: # Add cert-manager annotation to ValidatingWebhookConfiguration, MutatingWebhookConfiguration and CRDs
#      kind: Certificate
#      group: cert-manager.io
#      version: v1
#      name: serving-cert # this name should match the one in certificate.yaml
#      fieldPath: .metadata.namespace # namespace of the certificate CR
#    targets:
#      - select:
#          kind: ValidatingWebhookConfiguration
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 0
#          create: true
#      - select:
#          kind: MutatingWebhookConfiguration
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 0
#          create: true
#      - select:
#          kind: CustomResourceDefinition
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 0
#          create: true
#  - source:
#      kind: Certificate
#      group: cert-manager.io
#      version: v1
#      name: serving-cert # this name should match the one in certificate.yaml
#      fieldPath: .metadata.name
#    targets:
#      - select:
#          kind: ValidatingWebhookConfiguration
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 1
#          create: true
#      - select:
#          kind: MutatingWebhookConfiguration
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 1
#          create: true
#      - select:
#          kind: CustomResourceDefinition
#        fieldPaths:
#          - .metadata.annotations.[cert-manager.io/inject-ca-from]
#        options:
#          delimiter: '/'
#          index: 1
#          create: true
#  - source: # Add cert-manager annotation to the webhook Service
#      kind: Service
#      version: v1
#      name: webhook-service
#      fieldPath: .metadata.name # namespace of the service
#    targets:
#      - select:
#          kind: Certificate
#          group: cert-manager.io
#          version: v1
#        fieldPaths:
#          - .spec.dnsNames.0
#          - .spec.dnsNames.1
#        options:
#          delimiter: '.'
#          index: 0
#          create: true
#  - source:
#      kind: Service
#      version: v1
#      name: webhook-service
#      fieldPath: .metadata.namespace # namespace of the service
#    targets:
#      - select:
#          kind: Certificate
#          group: cert-manager.io
#          version: v1
#        fieldPaths:
#          - .spec.dnsNames.0
#          - .spec.dnsNames.1
#        options:
#          delimiter: '.'
#          index: 1
#          create: true
`
