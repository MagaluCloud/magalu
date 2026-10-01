---
sidebar_position: 2
---
# Create

Create Load Balancer

## Usage:
```
mgc beta load-balancer application create [flags]
```

## Examples:
```
mgc beta load-balancer application create --backends='[{"balance_algorithm":"round_robin","description":"Some optional backend description 1","health_check_name":"alb-health-check-1","name":"alb-backend-1","panic_threshold":50,"port":80,"targets":[{"instance_id":"00000000-0000-0000-0000-000000000000"},{"instance_id":"00000000-0000-0000-0000-000000000001"}],"targets_type":"instance"}]' --health-checks='[{"description":"Some optional health-check description 1","healthy_status_code":200,"healthy_threshold_count":8,"initial_delay_seconds":30,"interval_seconds":30,"name":"alb-health-check-1","path":"/health-check","port":5000,"protocol":"http","timeout_seconds":10,"unhealthy_threshold_count":3}]' --listeners='[{"backend_name":"alb-backend-1","description":"Some optional listener description 1","name":"alb-listener-1","port":443,"protocol":"https","tls_certificate_name":"alb-tls-certificate-1"}]' --routes='[{"description":"Routes /api and / to distinct backends","listener":"alb-listener-1","name":"general-route","rules":[{"backend":"alb-backend-1","headers":[{"name":"X-Env","value":"prod"}],"method":"GET","prefix":"/api"}]}]' --tls-certificates='[{"certificate":"SGVsbG8sIFdvcmxkIQ==","description":"Some optional tls-certificate description 1","expiration_date":"2027-12-31T23:59:59Z","name":"alb-tls-certificate-1","private_key":"SGVsbG8sIFdvcmxkIQ=="}]'
```

## Flags:
```
    --backends array(object)           Backends
                                       Use --backends=help for more details (required)
    --description string               Description
    --health-checks array(object)      Health Checks
                                       Use --health-checks=help for more details (required)
-h, --help                             help for create
    --listeners array(object)          Listeners
                                       Use --listeners=help for more details (required)
    --name string                      Name (required)
    --parent-id uuid                   Parent Id
    --public-ip-id string              Public Ip Id
    --routes array(object)             Routes
                                       Use --routes=help for more details (required)
    --subnet-id string                 Subnet Id
    --tls-certificates array(object)   Tls Certificates
                                       Use --tls-certificates=help for more details (required)
    --type enum                        ApplicationLoadBalancerType (must be "application") (required)
    --visibility enum                  LoadBalancerVisibility (one of "external" or "internal") (required)
    --vpc-id string                    Vpc Id (required)
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. You can use 'yaml', 'json' or 'table'.
    --project-type enum        ProjectType: Specifies the project type to which the load balancer belongs (one of "dbaas", "default" or "k8saas")
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
    --zone enum                Availability zone where the resource will be created (one of "br-se1-a", "br-se1-b" or "br-se1-c") (required)
```

