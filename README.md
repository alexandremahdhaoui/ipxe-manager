# ipxe-api

## iPXE booting workflow

| Phase             | Action                          | Description                                                                        |
|-------------------|---------------------------------|------------------------------------------------------------------------------------|
| `[BOOTSTRAPPING]` | Call `/boot.ipxe`               | Machine starts and DHCP rule 67 specifies next server as this ipxe-api.            |
| `[ASSIGNMENT]`    | Chainload `/ipxe?labels=values` | Machine chainload to this endpoint specifying labels for scheduling/assignment.    |
| `[BOOT]`          | Run `#ipxe...`                  | Contains an optional uuid reference to a config such as an ignition or cloud-init. |
| `[OPTIONAL]`      | Fetch `/config/{uuid}`          | Fetch the optional config identified by a UUID.                                    |

## Custom Resource Definitions

We designed the Profile and Assignment CRDs in way that let 

### Profile

```yaml
apiVersion: ipxe.cloud.alexandre.mahdhaoui.com/v1alpha1
kind: Profile
metadata:
  name: your-profile
  labels:
    assignment/ipxe/buildarch: aarch64
    assignment/extrinsic/region: us-cal
spec:
  # ipxe string
  ipxe: |
    # ipxe
    command ... --with-config {{ config0 }} --ignition-url {{ ignitionFile }} --or-cloud-init {{ cloudInit }}
  # additionalConfig map[string]string
  additionalConfig:
    config0: |
      YOUR CONFIG HERE
    ignitionFile: |
      YOUR IGNITION CONFIG HERE
    cloudInit: |
      YOUR CLOUD INIT CONFIG HERE
status:
  # UUIDs that are used to fetch
  additionalConfig:
    config0: 89952e35-2a85-4f03-a6b2-7f9526bfafc0
    ignitionFile: 445a4753-3d59-4429-8cea-7db9febdecad
```

### Assignment

Because the `ipxe-api` should not endorse any `scheduler` or `assigner` role, but serve the purpose of other processes,
assignments should be authored by them.

This purpose is served by the `Assignment` CRD.

TODO: How do we avoid issues w/ conflicting subjectSelectors. Should we rank labels?

```yaml
apiVersion: ipxe.cloud.alexandre.mahdhaoui.com/v1alpha1
kind: Assignment
metadata:
  name: your-assignment
spec:
  # subjectSelector map[string]string
  # the specified labels selects a subject that iPXE boots.
  subjectSelector:
    serialNumber: c4a94672-05a1-4eda-a186-b4aa4544b146
    uuid: 47c6da67-7477-4970-aa03-84e48ff4f6ad
  # profileSelectors map[string]string
  # the specified labels selects which profile should be used.
  profileSelectors:
    assigment/ipxe/buildarch: aarch64
status:
  conditions: []
```

## Architecture

We have controllers, admission webhooks and a REST API.

The **REST API** is an iPXE server that only serves GET requests. The API endpoints are described below:
- `/boot.ipxe` to chainload into `/ipxe` endpoint.
- `/ipxe?key=value` to load the iPXE manifest assigned to the booting machine.
- `/config/{config-uuid}?key=value` to dynamically load any arbitrary configuration files.

**Admission webhooks** ensures Assignment & Profile custom resources are conform, and optionally enriched them with more
information.

**Controllers** maintain datastructures queried by the REST API.

## Deployment

Replicas of the REST API queries the datastructures maintained by the controllers. These communications are performed
via mTLS. Hence, cert-manager is required for a production deployment.

## Acknowledgement

This project was inspired by [poseidon/matchbox](https://github.com/poseidon/matchbox).
