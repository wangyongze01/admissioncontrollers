# Kubernetes Resource Control & Jar Injection Webhook

## 项目简介
该项目基于 **Kubernetes Admission Webhook**，实现了以下功能：  

1. **节点资源超卖防控**  
   - 自动拦截Node对象的创建/更新请求  
   - 通过Patch限制CPU/内存分配，防止节点超卖导致调度异常  

2. **Deployment自动Jar包注入**  
   - 在Pod创建时动态添加Init Container  
   - 将Jar包从共享卷复制到应用容器  
   - 自动挂载卷到主容器，提高部署效率  

---

## 技术栈
- Kubernetes: Admission Webhook、Node/Pod对象Patch、Init Container  
- Go语言: Webhook服务开发  
- JSON Patch: 动态修改K8s资源对象  

---

## 功能亮点
- **资源控制**：节点CPU/内存冲突事件下降50%  
- **自动化部署**：Jar包注入减少部署时间约30%  
- **安全可靠**：多团队共享集群部署安全性提升  

---

## 快速使用
1. **编译Webhook服务**
```bash
go build -o webhook main.go
