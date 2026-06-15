# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### تنظیمات عمومی ###
# تنظیم namespace برای Kubernetes
k8s_namespace('microservice-dev')

# تعیین contextهای مجاز
allow_k8s_contexts(['docker-desktop', 'minikube', 'docker-for-desktop'])

### k8s Config ###
k8s_yaml('./infra/development/k8s/configs/app-config.yaml')
#k8s_yaml('./infra/development/k8s/configs/secrets.yaml')
### End k8s Config ###

### RabbitMQ ###
k8s_yaml('./infra/local/k8s/rabbitmq-service/deployment.yaml')
k8s_yaml('./infra/local/k8s/rabbitmq-service/service.yaml')
k8s_resource('rabbitmq',
             port_forwards=[5672, 15672],
             labels=['tooling', 'infra'],
             extra_pod_selectors=[{'app': 'rabbitmq'}])
### End RabbitMQ ###

### PostgresDB ###
k8s_yaml('./infra/local/k8s/postgres-service/pvc.yaml')
k8s_yaml('./infra/local/k8s/postgres-service/deployment.yaml')
k8s_yaml('./infra/local/k8s/postgres-service/service.yaml')
k8s_resource('postgres',
             port_forwards=[5432],
             labels=['databases', 'infra'],
             extra_pod_selectors=[{'app': 'postgres'}])
### End PostgresDB ###

### API Gateway ###
# کامپایل API Gateway
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'

local_resource(
    name='api-gateway-compile',
    cmd=gateway_compile_cmd,
    deps=['./services/api-gateway'],
    ignore=['./services/api-gateway/vendor'],
    labels=['compiles'],
    trigger_mode=TRIGGER_MODE_AUTO)

# ساخت Docker image
docker_build(
    'microservice/api-gateway:dev',
    '.',
    dockerfile='./infra/local/docker/api-gateway/api-gateway.Dockerfile',
    only=[
        './services/api-gateway',
        './build/api-gateway',
        './env'
    ],
    live_update=[
        sync('./services/api-gateway', '/app'),
        sync('./build/api-gateway', '/app/build/api-gateway'),
        sync('./env', '/app/.env')
    ]
)

# کانفیگ Kubernetes مخصوص development
k8s_yaml('./infra/local/k8s/api-gateway/deployment.yaml')
k8s_yaml('./infra/local/k8s/api-gateway/service.yaml')

k8s_resource('api-gateway',
             port_forwards=[8085],
             labels=['services', 'gateway'],
             extra_pod_selectors=[{'app': 'api-gateway'}])
### End of API Gateway ###

### Auth Service ###
# کامپایل Auth Service
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/auth-service ./services/auth-service/cmd'

local_resource(
    name='auth-service-compile',
    cmd=auth_compile_cmd,
    deps=['./services/auth-service'],
    ignore=['./services/auth-service/vendor'],
    labels=['compiles'],
    trigger_mode=TRIGGER_MODE_AUTO)

# ساخت Docker image برای Auth Service
docker_build(
    'microservice/auth-service:dev',
    '.',
    dockerfile='./services/auth-service/docker/auth-service.Dockerfile',
    only=[
        './services/auth-service',
        './build/auth-service',
        './env'
    ],
    live_update=[
        sync('./services/auth-service', '/app'),
        sync('./build/auth-service', '/app/build/auth-service'),
        sync('./env', '/app/.env')
    ]
)

# کانفیگ Kubernetes مخصوص development
k8s_yaml('./infra/local/k8s/auth-service/deployment.yaml')
k8s_yaml('./infra/local/k8s/auth-service/service.yaml')

k8s_resource('auth-service',
             port_forwards=[9092],
             labels=['services', 'auth'],
             extra_pod_selectors=[{'app': 'auth-service'}])
### End Auth Service ###


### Notification Service ###
# کامپایل Notification Service
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/notification-service ./services/notification-service/cmd'

local_resource(
    name='notification-service-compile',
    cmd=auth_compile_cmd,
    deps=['./services/notification-service'],
    ignore=['./services/notification-service/vendor'],
    labels=['compiles'],
    trigger_mode=TRIGGER_MODE_AUTO)

# ساخت Docker image برای Auth Service
docker_build(
    'microservice/notification-service:dev',
    '.',
    dockerfile='./services/notification-service/docker/notification-service.Dockerfile',
    only=[
        './services/notification-service',
        './build/notification-service',
        './env'
    ],
    live_update=[
        sync('./services/notification-service', '/app'),
        sync('./build/notification-service', '/app/build/notification-service'),
        sync('./env', '/app/.env')
    ]
)

# کانفیگ Kubernetes مخصوص development
k8s_yaml('./infra/local/k8s/notification-service/deployment.yaml')

k8s_resource('notification-service',
             port_forwards=[9093],
             labels=['services', 'notification'],
             extra_pod_selectors=[{'app': 'notification-service'}])
### End Notification Service ###



### نمایش وضعیت ###
# غیرفعال کردن حذف خودکار images
docker_prune_settings(disable=True)