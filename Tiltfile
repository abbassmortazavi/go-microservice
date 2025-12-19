# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### k8s Config ###
k8s_yaml('./infra/development/k8s/app-config.yaml')
### End k8s Config ###

### PostgresDB ###
k8s_yaml('./infra/development/k8s/postgres-service/postgres-deployment.yaml')
k8s_resource('postgres', port_forwards=5432, labels='databases')
### End PostgresDB ###


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



### Auth Service ###
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/auth-service ./services/auth-service/cmd/main.go'
if os.name == 'nt':
    auth_compile_cmd = './services/auth-service/docker/auth-service-build.bat'

local_resource(
    'auth-service-compile',
    auth_compile_cmd,
    deps=['./services/auth-service'], labels="compiles")

docker_build_with_restart(
    'microservice/auth-service',
    '.',
    entrypoint=['/app/build/auth-service'],
    dockerfile='./services/auth-service/docker/auth-service.Dockerfile',
    only=[
        './build/auth-service'
    ],
    live_update=[
        sync('./build', '/app/build')
    ]
)


k8s_yaml('./infra/development/k8s/auth-service/auth-service-deployment.yaml')
k8s_resource('auth-service',resource_deps=['auth-service-compile'], labels="services")
### End Auth Service ###







