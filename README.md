# ab-test-demo

In simplest terms, `A/B Testing`, in application development terms, implies an ability to compare two versions of the same application in a controlled environment to measure properties of each version (e.g. users' behavior or application performance).

In this example, I will overview how to configure a simple A/B testing experiment on Cloud Run.

## Scenario

Imagine you published a marketing site aiming to collect users' sign ups. After a couple of days, the marketing department tells you that they expected more sign ups and think that the low number of user registrations is related to the location of the sign up button (bottom of the page). They recommend you move the button to the top of the page to get more sign ups.

![](img/signup.png)

Now, you, as a developer, made data-driven decisions so you want to put together an experiment wherein subset of applications users will be seeing the new application layout and the rest will continue seeing the current version of the application.

## Pre-requirements

If you don't have one already, start by creating new project and configuring [Google Cloud SDK](https://cloud.google.com/sdk/docs/). To get access to the `traffic management` feature of Cloud Run you will also need to configure the `alpha` version of `gcloud`:

```shell
gcloud components install alpha
```

## Setup

> Throughout this sample I'm going to be using pre-build images of a demo application. These images are publicly accessible, but, if you want to build your own versions from this repository you can do so using the [bin/build-images](bin/build-images) script. As with any script, review it before executing it.

## Version "A"

To mimic the above scenario, I'm first going to deploy the original version of the application ("A").

```shell
gcloud beta run deploy signup \
	--image gcr.io/cloudylabs-public/ab-test-demo:a
```

> Note, depending on how you configured your `gcloud`, you may also have to append the `--cluster` and `--cluster-location` flags or run. More information on how to configure these flags by default see [Quickstart: Deploy to Cloud Run on GKE](https://cloud.google.com/run/docs/quickstarts/prebuilt-deploy-gke)

The `gcloud` command will return a URL to the active revision which is serving traffic, should look something like this:

```shell
Deploying container to Cloud Run on GKE service [signup] in namespace [demo] of cluster [cr-prod]
✓ Deploying... Done.
  ✓ Creating Revision...
  ✓ Routing traffic...
Done.
Service [signup] revision [signup-2zpw4] has been deployed and is serving traffic at http://signup.demo.cloudylabs.dev
```

In the above example, my `gcloud` configuration was already configured with defaults for name of the cluster (`cr-prod`), its location (`us-east1` region), and the target namespace (`demo`). The same deployment with all these flags in line would look like this:

```shell
gcloud beta run deploy signup \
    --image gcr.io/cloudylabs-public/ab-test-demo:a \
    --cluster cr-prod \
    --cluster-location us-east1 \
    --namespace demo
```

To see how the deployed version of the application should look like see https://signup-a.demo.cloudylabs.dev

## Version "B"

Cloud Run does a lot of things automatically so to implement my experiment I will first have to take it out of the "auto-pilot" mode so I can manage the revisions manually. To do that, I'll list the revisions:

```shell
gcloud alpha run revisions list --service signup
```

There should be only one revision at this point, but, in case you have run the deployment a few times just copy the revision ID from the top most row

```shell
For cluster [cr-prod] in [us-east1]:
   REVISION      ACTIVE  SERVICE  DEPLOYED
✔  signup-9p644  yes     signup   2019-08-26 22:44:10 UTC
✔  signup-2zpw4          signup   2019-08-26 22:12:13 UTC
✔  signup-m89nj          signup   2019-08-26 22:10:57 UTC
```

> Note, depending on how fast you are, the `yes` indicator for `ACTIVE` may not always be there as right now the revision doesn't retain that attribute when it scales down to 0

Once you have captured the revision ID, you can tell Cloud Run to that you want to manage the revisions manually but setting traffic 100% explicitly to that revision.

> Make sure to replace the revision ID in the following command before running it

```shell
gcloud alpha run services set-traffic signup \
    --to-revision signup-9p644=100
```

The result should look something like this:

```shell
Done.
Traffic set to signup-9p644=100.
```

Now we can deploy new revision to Cloud Run and will NOT take any traffic. 100% of the traffic will contuse to be routed to the revision we set above. We are ready to deploy the "B" version.

```shell
gcloud beta run deploy signup \
	--image gcr.io/cloudylabs-public/ab-test-demo:b
```

> The confirmation message will make it sound like the traffic actually is sent to the new version right now but that's just a CLI bug and should be fixed soon

Now if you list the service revisions you should see new one

```shell
gcloud alpha run revisions list --service signup
```

Our version "B" will be the top most on that returned list and the version "A" will be the third one from the top. There is a "side-effect" revision created right now. Fix for this is coming as well.

```shell
   REVISION      ACTIVE  SERVICE  DEPLOYED
✔  signup-xf95v          signup   2019-08-26 23:00:55 UTC
✔  signup-7cbsz          signup   2019-08-26 22:58:28 UTC
✔  signup-9p644          signup   2019-08-26 22:44:10 UTC
✔  signup-2zpw4          signup   2019-08-26 22:12:13 UTC
```

Now we are ready to execute our experiment. We will send 90% of the traffic to the original version ("A") and 10% to the update ("B") version.

```shell
gcloud alpha run services set-traffic signup \
    --to-revision signup-9p644=90,signup-xf95v=10
```

The response should be

```shell
Done.
Traffic set to signup-9p644=90, signup-xf95v=10.
```

To compare how the version "A" and "B" should look like take a look at these already deployed versions:

* http://signup-a.demo.cloudylabs.dev
* http://signup-b.demo.cloudylabs.dev



## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
