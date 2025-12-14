# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### k8s Config ###
k8s_yaml('./infra/development/k8s/app-config.yaml')
### End k8s Config ###


### API Gateway ###
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
    gateway_compile_cmd = './infra/development/docker/api-gateway/api-gateway-build.bat'

local_resource(
    'api-gateway-compile',
    gateway_compile_cmd,
    deps=['./services/api-gateway'], labels="compiles")

docker_build_with_restart(
    'microservice/api-gateway',
    '.',
    entrypoint=['/app/build/api-gateway'],
    dockerfile='./infra/development/docker/api-gateway/api-gateway.Dockerfile',
    only=[
        './build/api-gateway'
    ],
    live_update=[
        sync('./build', '/app/build')
    ]
)


k8s_yaml('./infra/development/k8s/api-gateway/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile'], labels="services")
### End of API Gateway ###







