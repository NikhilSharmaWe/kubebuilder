# Deploy Image Plugin

The deploy-image plugin allows you to create [controllers][controller-runtime] and custom resources which will deploy and manage an image on the cluster following the guidelines and 
best practices. It abstracts the complexities to achieve this goal and still giving the possibility of users improve and 
customize their projects according to their requirements.
<aside class="note">
<h1>Examples</h1>

You can check samples using this plugin by looking at the "project-v3-with-deploy-image" sample inside the [testdata][testdata] directory of the Kubebuilder project.

</aside> 


## When to use it ?
- If you are looking to scaffold an operator solution which helps you in deploying an image of your application, so that you do not need to know exactly how to implement it.

## How to use it ?

After initializing the project with one of the available Golang plugins. 
```sh
kubebuilder init --plugins=go/v3 --domain example.org --repo example.org/guestbook-operator
```
You will need to use deploy-image/v1-alpha plugin for creating the
custom resource and controller which will deploy and manage your image on the cluster. Use --image flag (available only for this plugin) to state your image.

```sh
kubebuilder create api --plugins=deploy-image/v1-alpha --group ship --version v1beta1 --kind Frigate --image=example.com/frigate:v0.0.1 
```

## Subcommands

The declarative plugin implements the following subcommands:

* create api (`$ kubebuilder create api [OPTIONS]`)

## Affected files

The following scaffolds will be created or updated by this plugin:

* `controllers/*_controller.go`
* `controllers/*_suite_test.go`
* `api/*_types.go`
* `api/*_groupversion_info.go`
* `config/samples/*_.yaml`
* `hack/boilerplate.txt`

[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
