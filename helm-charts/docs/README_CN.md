# [星舰（Starship）](https://github.com/tricorder-observability/starship) 安装文档

[![Release Helm Charts](https://github.com/tricorder-observability/Starship/actions/workflows/release-chart.yaml/badge.svg)](https://github.com/tricorder-observability/Starship/actions/workflows/release-chart.yaml)

![image](https://user-images.githubusercontent.com/112656580/219543149-2e2bbebc-1891-4dcb-ba66-0f8b7f1bcd68.png)
![image](https://user-images.githubusercontent.com/112656580/219542981-5a4e5fb1-0603-4c0b-91e2-c94c36a92c0b.png)

本仓库是用于安装[星舰（Starship）](https://github.com/tricorder-observability/starship)
可观测平台的 Helm charts，
星舰是[三度观测云](https://tricorder.dev)开发的下一代可观测性平台。

🖖 [星舰](https://github.com/tricorder-observability/starship) 🪐 基于 🐝 eBPF➕WASM
![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)
打造，定义了零成本可观测性。

🚀 星舰重新定义可观测性，其意义类似 ChatGPT 对搜索的颠覆性意义；🐝 eBPF 支持无侵入式数据采集，同时 WASM
![image](https://user-images.githubusercontent.com/112656580/219543881-046af389-ca10-4dda-b79a-a60088a1220a.png)
弥补了 eBPF 无法实现复杂数据处理的缺陷。

目前，Starship 仅支持 Kubernetes 平台上运行。Starship 基于 eBPF 技术，无需更改应用程序代码，
即可立即获取云原生应用的服务运行拓扑图（Service Map）。
Starship 也支持收集 Prometheus、OpenTelemetry 数据。
Starship 使用 Grafana 对收集的数据进行可视化，使用下面的用户名、密码访问：

```text
username: admin
password: tricorder
```

## 系统要求

### AWS EKS

如果您使用AWS EKS，请在EKS群集上安装EBS CSI。这是必须的，因为Helm
charts需要为数据库pod创建持久卷PerstientVolume，这需要EBS CSI。

### 阿里云 ACK

- 如果你使用阿里云 ACK，阿里云会默认安装
  [ACK CSI插件](https://help.aliyun.com/document_detail/134722.html) 在你的集群中。
  你的集群中必须要要有 CSI 插件，因为 Helm Charts 需要为数据库 Pod 创建持久卷
  [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)，
  这需要 CSI 插件。
- 在安装 Starship 之前，你需要检查你的集群中是否有默认的 StorageClass，
  如果没有，你需要创建一个。
  你可以通过下面的命令来查看你的集群中是否有默认的 StorageClass。

  ```shell
  kubectl get storageclass | grep default
  ```
  
- 如果你的集群中没有默认的 StorageClass，你可以通过下面的命令来创建默认的 StorageClass。

  ```shell
  kubectl patch storageclass <you-storageclass-name> --patch '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  # 例如:
  kubectl patch storageclass alibabacloud-cnfs-nas --patch '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

**警告** 如果你的集群中没有默认的 StorageClass，你需要创建一个默认的 StorageClass，否则 Starship 无法正常安装在你的集群中。

## 安装

创建namespace，这里以tricorder为例。

```shell
helm repo add tricorder-stable \
    https://tricorder-observability.github.io/helm-charts
helm repo update
kubectl create namespace tricorder
helm install my-starship tricorder-stable/starship -n tricorder
```

## 通过 LoadBalancer External IP 访问 Starship Web 管理界面

Starship默认通过 LoadBalancer 来发布 Web UI
服务。如果您的集群配置了支持外部访问的LoadBalancer， 比如 AWS LoadBalancer
Controller 或 Ingress Controller，您可以直接通过 api-server 服务的 ExternIP
访问该服务:

```shell
kubectl get service -n tricorder
```

![image](https://user-images.githubusercontent.com/112656580/215043391-6c4cd4bd-3a58-472f-a688-b88f11ef90c1.png)

在你的浏览器输入`http://${EXTERNAL-IP}`，注意这里是**HTTP**，**不是HTTPS**。

## 使用kubectl port-forward转发访问Starship 管理界面

您也可以使用kubectl port-forward命令将服务发布给本地系统；下面的命令将 API
Server 的 80 端口暴露于本地的 18080 端口；可以通过 `http://localhost:18080/`
来访问。

```shell
kubectl -n tricorder port-forward service/my-starship-api-server
18080:80
```

## 使用 Ingress 转发访问Starship 管理界面

Ingress 用于公开从集群外部到集群内服务的 HTTP 和 HTTPS 路由。
使用以下命令创建 Ingress 规则使用 `starship.io` 作为 host，`80` 作为端口访问 Starship。

```shell
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tricorder
  namespace: tricorder
spec:
  rules:
  - host: starship.io
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: api-server
            port:
              number: 80
EOF
```

如果你想使用更多的 Ingress 功能，比如 TLS 访问，
请根据 Ingress
[文档](https://kubernetes.io/zh-cn/docs/concepts/services-networking/ingress/) 来编写对应的配置。

## 注意事项

**注意** 本项目正在积极开发中，当前版本为技术预览版。
**警告** 不要同时在多个名字空间安装星舰可观测性平台；如果你不小心在多个名字空间
  安装了星舰可观测性平台，按照卸载说明删除所有部件并重新安装。
**警告** 不要在同一个命名空间中安装多个发行版（Release）。

## 配置数据留存时间

指标和跟踪数据会被保留一段时间，然后被自动丢弃；默认数据保留周期为7天。
数据留存时间的配置对系统的稳定性至关重要。 留存时间必须足够小，
才能保证存储系统能够容纳所有的数据，否则会因为磁盘耗尽、或数据量过大导致查询速度
大幅下降， 导致系统不稳定、可用性下降，甚至直接崩溃：

```yaml
promscale:
  config:
    startup.dataset.config: |
      metrics:
        compress_data: true
        default_retention_period: 7d
      traces:
        default_retention_period: 7d
```

上面的保留时间可以通过--valus 参数自定义，我们将 default_retention_period
从7天改为30天:

- 为自定义保留时间创建补丁yaml文件

```shell
cat > rentention_patch.yaml << EOF
promscale:
  config:
    startup.dataset.config: |
      metrics:
        compress_data: true
        default_retention_period: 30d
      traces:
        default_retention_period: 30d
EOF
```

- 通过--values 参数覆盖默认设置:

```shell
helm install my-starship tricorder-stable/starship -n tricorder \
    --values rentention_patch.yaml
```

## 卸载

你可以运行以下命令来卸载startship。

```shell
helm uninstall my-starship -n tricorder
```

在卸载helm 版本后，一些对象会有遗留，可通过以下步骤删除遗留的对象。

### 删除密钥

在部署时创建的密钥默认不会被删除，这些密钥需要被手动删除：

```shell
kubectl delete -n tricorder \
    $(kubectl get secrets -n tricorder -l "app=timescaledb" -o name)
```

### 删除configmap

kubectl delete -n tricorder \
    $(kubectl get configmap -n tricorder -l "app=my-starship-promscale" -o name)

### 删除Kube-Prometheus secret

部署时创建的Kube-Prometheus secret没有被删除，需要被手动删除：

```shell
kubectl delete secret -n tricorder my-starship-kube-prometheus-stack-admission
```

### 删除 DB PVC和备份

删除部署并不会同时删除这个部署的持久卷（Persistent Volume Claims (pvc)），
可使用以下指令全部删除：

```shell
kubectl delete -n tricorder \
    $(kubectl get pvc -n tricorder -l release=my-starship -o name)
```

### 删除 Prometheus 相关组件

删除部署时，相关联的 Prometheus的持久卷不会被删除，可使用以下指令全部删除：

```shell
kubectl delete -n tricorder $(kubectl get pvc -n tricorder \
  -l operator.prometheus.io/name=my-starship-kube-prometheus-stack-prometheus \
  -o name)

kubectl delete crd alertmanagerconfigs.monitoring.coreos.com \
    alertmanagers.monitoring.coreos.com \
    probes.monitoring.coreos.com \
    prometheuses.monitoring.coreos.com \
    prometheusrules.monitoring.coreos.com \
    servicemonitors.monitoring.coreos.com \
    thanosrulers.monitoring.coreos.com \
    podmonitors.monitoring.coreos.com
kubectl delete MutatingWebhookConfiguration my-starship-kube-promethe-admission
kubectl delete ValidatingWebhookConfiguration my-starship-kube-prometheus-admission
```

### 删除namespace

最后将 namespace 删除，至此，星舰可观测性平台所有的组件都以从你的集群上移除。

```shell
kubectl delete namespace tricorder
```

## 高级主题

### 采集 OpenTelemetry 数据

[将 OTel 数据发送到星舰可观测性平台](./send-otlp-data-to-starship.md).

### 覆盖默认配置

您可以使用--set 参数来覆盖在 charts/starship/Values.yaml 的配置值。

```shell
# 将服务类型改成 ClusterIP
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set service.type=ClusterIP

# 使用特定版本的 Starship
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set images.tag=<a specific tag>

# 使用特定版本的 Starship apiSever
helm upgrade my-starship tricorder-stable/starship -n tricorder \
    --set apiServer.image.tag=<a specific tag>
```

### 从本地仓库中安装

这个存储库，并用本地路径charts/starship 取代 tricorder-stable/starship.

```shell
git clone https://github.com/tricorder-observability/helm-charts
cd helm-charts
# You'll need this step to fetch the dependent charts
helm dep update charts/starship
helm install my-starship charts/starship -n tricorder
```

在前面安装章节所列的所有指令，您需要将tricorder-stable/starship替换为charts/starship，并且从克隆的仓库的根目录运行前述的所有命令。
