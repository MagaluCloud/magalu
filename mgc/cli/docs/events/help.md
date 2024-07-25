# # Consult API

## Usage:
```bash
With the continued growth of the distributed architecture in Magalu
Cloud, everyone should have the ability to audit all (or important)
events triggered by users or internal systems. This capability is
essential not only for visualizing and tracking changes requested by
internal developers but also for those made by tenant users. Our
initial approach was to create a "broker" (i.e., an event system and a
consultation API) that could store all these events. All events are
exposed using the CloudEvent specification.
```

## Product catalog:
- That API is a internal and External (filtered by Tenant ID) usage.

## Other commands:
- ## About CloudSpec
- Info about the spec: https://cloudevents.io/

## Flags:
```bash
### The cloud spec fields used today ([Spec](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md)):
- id: The request id to track.
- source: Identifies the context of this event.
- type: The “action” or type of the event. SHOULD be prefixed with a reverse-DNS name.
- specversion: The version of the CloudEvents specification which the event uses. Today is 1.0.
- subject: Identifies the subject of the event in the context of the producer. e.g: source
  is https://virtual-machine.pre-prod.br-se-1.jaxyendy.com/v1/instances/cddd202a-7017-4d80-b702-6c49ad89ae99/start and
  the subject is the instance id.
- data: Open field with the raw event.
- time: ISO8601 date of the event.
- authid: An unique identifier of the principal that triggered the occurrence.
- authtype: An enum representing the type of principal that triggered the occurrence, e.g: tenant, api_key or
  unauthenticated, system_api_key.
```

