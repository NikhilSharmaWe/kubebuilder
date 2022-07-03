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

package controllers

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &SuiteTest{}
var _ machinery.Inserter = &SuiteTest{}

// SuiteTest scaffolds the file that sets up the controller tests
// nolint:maligned
type SuiteTest struct {
	machinery.TemplateMixin
	machinery.MultiGroupMixin
	machinery.BoilerplateMixin
	machinery.ResourceMixin

	// CRDDirectoryRelativePath define the Path for the CRD
	CRDDirectoryRelativePath string

	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *SuiteTest) SetTemplateDefaults() error {
	if f.Path == "" {
		if f.MultiGroup && f.Resource.Group != "" {
			f.Path = filepath.Join("controllers", "%[group]", "suite_test.go")
		} else {
			f.Path = filepath.Join("controllers", "suite_test.go")
		}
	}
	f.Path = f.Resource.Replacer().Replace(f.Path)

	f.TemplateBody = fmt.Sprintf(controllerSuiteTestTemplate,
		machinery.NewMarkerFor(f.Path, importMarker),
		machinery.NewMarkerFor(f.Path, addSchemeMarker),
	)

	// If is multigroup the path needs to be ../../ since it has
	// the group dir.
	f.CRDDirectoryRelativePath = `".."`
	if f.MultiGroup && f.Resource.Group != "" {
		f.CRDDirectoryRelativePath = `"..", ".."`
	}

	// This one is to overwrite the controller_suitetest.go if it exist
	f.IfExistsAction = machinery.OverwriteFile

	return nil
}

const (
	importMarker    = "imports"
	addSchemeMarker = "scheme"
)

// GetMarkers implements file.Inserter
func (f *SuiteTest) GetMarkers() []machinery.Marker {
	return []machinery.Marker{
		machinery.NewMarkerFor(f.Path, importMarker),
		machinery.NewMarkerFor(f.Path, addSchemeMarker),
	}
}

const (
	apiImportCodeFragment = `%s "%s"
`
	addschemeCodeFragment = `err = %s.AddToScheme(scheme.Scheme)
Expect(err).NotTo(HaveOccurred())

`
)

// GetCodeFragments implements file.Inserter
func (f *SuiteTest) GetCodeFragments() machinery.CodeFragmentsMap {
	fragments := make(machinery.CodeFragmentsMap, 2)

	// Generate import code fragments
	imports := make([]string, 0)
	if f.Resource.Path != "" {
		imports = append(imports, fmt.Sprintf(apiImportCodeFragment, f.Resource.ImportAlias(), f.Resource.Path))
	}

	// Generate add scheme code fragments
	addScheme := make([]string, 0)
	if f.Resource.Path != "" {
		addScheme = append(addScheme, fmt.Sprintf(addschemeCodeFragment, f.Resource.ImportAlias()))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(imports) != 0 {
		fragments[machinery.NewMarkerFor(f.Path, importMarker)] = imports
	}
	if len(addScheme) != 0 {
		fragments[machinery.NewMarkerFor(f.Path, addSchemeMarker)] = addScheme
	}
	f.IfExistsAction = machinery.OverwriteFile
	return fragments
}

const controllerSuiteTestTemplate = `{{ .Boilerplate }}

{{if and .MultiGroup .Resource.Group }}
package {{ .Resource.PackageName }}
{{else}}
package controllers
{{end}}

import (
	"path/filepath"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	%s
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx       context.Context
var cancel    context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
	"Controller Suite",
	[]Reporter{printer.NewlineReporter{}})
}

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj client.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join({{ .CRDDirectoryRelativePath }}, "config", "crd", "bases")},
		ErrorIfCRDPathMissing: {{ .Resource.HasAPI }},
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	%s

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&{{ .Resource.Kind }}Reconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
`
