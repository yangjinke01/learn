# Remove kubernetes namespace stuck in terminating

1. First, dump the namespace spec in json format as seen below:

   ```shell
   kubectl get ns  -o json > namespace.json
   ```

2. Next, we edit the namespace.json and then remove the finalizer portion in the spec. So, we have to change to from:

   ```text
   ”spec”: {
   	“finalizers”:
     },
   ```

to:

   ```text
   ”spec”: {
     },
   ```

Alternatively, we can remove everything inside the spec to get the job done.

3. After that, we have to manually patch the namespace with the command below:

   ```shell
   kubectl replace --raw "/api/v1/namespaces/{namespace}/finalize" -f namespace.json
   ```