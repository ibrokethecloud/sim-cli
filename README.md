# sim-cli

`sim-cli` is a utility to help manage multiple instances of the rancher support-bundle-kit simulator.

Currently, the easiest way to run simulator instances is to run the binary locally. Due to port mapping requirements
it is hard to run multiple instances of the simulator on the same host.

`sim-cli` attempts to bridge this gap by allowing a user to run and load multiple support bundles in containers.

The cli does this by building a custom image based on `rancher/support-bundle-kit:master-head`, which has been layered
with the contents of the support bundle itself. It then runs this custom image with specific command arguments to trigger
the processing of the packaged support bundle.

## To build
`make` will leverage dapper and build and drop the binary in `bin` directory of the project.

## Usage
```markdown
sim-cli is a utility to help create and manage multiple support bundle kid instances in a docker container. 
This allows users to have multiple copies of support bundle kit running on your desktop to allow debugging of harvester issues

Usage:
  sim-cli [flags]
  sim-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      create a support bundle kit simulator instance
  delete      delete a support bundle kit simulator instance
  export      export kubeconfig for an existing simulator instance
  help        Help about any command
  list        list existing simulator instances

Flags:
  -h, --help      help for sim-cli
      --verbose   verbose output


```

### Creating a new instance
```
sim-cli create --name issue-7007 --bundle-path $HOME/Downloads/supportbundle_207d0deb-1cf3-46c8-aedb-fd3d28d04530_2024-09-04T07-00-02Z.zip
```
will create a new simulator instance by packaging the zip file of the support bundle into a base image of support-bundle-kit.
It will run a new instance using the newly create image, and export the kubeconfig from the running instance and merge
it in to the default simulator config file `$HOME/.sim/admin.kubeconfig`
```markdown
sim-cli create --name issue-7007 --bundle-path $HOME/Downloads/supportbundle_207d0deb-1cf3-46c8-aedb-fd3d28d04530_2024-09-04T07-00-02Z.zip
INFO[0001] Step 1/4 : FROM rancher/support-bundle-kit:dev 
INFO[0001]  ---> d58fca6009e7                           
INFO[0001] Step 2/4 : EXPOSE 6443/tcp                   
INFO[0001]  ---> Using cache                            
INFO[0001]  ---> 7ef1ef326246                           
INFO[0001] Step 3/4 : COPY bundle /bundle               
INFO[0001]  ---> Using cache                            
INFO[0001]  ---> 7a3519649d88                           
INFO[0001] Step 4/4 : LABEL harvesterhci.io/bundle-name=issue-7007 
INFO[0001]  ---> Running in 5bd40db14d27                
INFO[0001]  ---> 61082bce11ad                           
INFO[0001] Successfully built 61082bce11ad              
INFO[0001] Successfully tagged sim-cli-managed:issue-7007 
INFO[0001] simulator instance exposed on port 32773      name=issue-7007
INFO[0011] exporting kubeconfig for instance issue-7007 
INFO[0011] exported kubeconfig to context issue-7007  
```

Users can use the newly added context to access the simulator instance using any tooling used to access a k8s cluster.

### Listing instances
`sim-cli list` will list all running instances of simulator along with details of related image, support bundle file
and port this instance is exposed on
```markdown
sim-cli list
+---------------+---------------------------------------------+-------------------------------+------------------+-----------------+
|          name |                                  bundlePath |                         image |           status |    exposed port |
+===============+=============================================+===============================+==================+=================+
|    issue-7007 |    /home/random/Downloads/supportbundle_207 |    sim-cli-managed:issue-7007 |    Up 40 minutes |           32770 |
|               |    d0deb-1cf3-46c8-aedb-fd3d28d04530_2024-0 |                               |                  |                 |
|               |                          9-04T07-00-02Z.zip |                               |                  |                 |
+---------------+---------------------------------------------+-------------------------------+------------------+-----------------+
```


### Deleting an instance
`sim-cli delete --name issue-7007` will find the associated container and image for instance, stop the container, 
remove the container and associated image. It will also remove the context associated for the instance from the kubeconfig

```markdown
sim-cli delete --name issue-7007
INFO[0000] removing instance issue-7007                 
INFO[0000] removing image for instance issue-7007       
INFO[0000] removed image: { sim-cli-managed:issue-7007} 
INFO[0000] removed image: {sha256:4cbca2ba8ec2626ba904b6d3b6c9570e3734a74f4d41f6457f14778497f9efe9 } 
INFO[0000] removing context for instance issue-7007     
```

### Export kubeconfig for an instance
`sim-cli export --name issue-7007` can be used to export the kubeconfig for an already running instance.
The export config will be added as a new context into `$HOME/.sim/admin.kubeconfig`
```markdown
sim-cli export --name issue-7007
INFO[0000] exporting kubeconfig for instance issue-7007 
INFO[0000] exported kubeconfig to context issue-7007    
```