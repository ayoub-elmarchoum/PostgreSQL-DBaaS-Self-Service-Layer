# PostgreSQL-DBaaS-Self-Service-Layer
# 📘 PostgreSQL DBaaS Self-Service Layer

This project explores the development of an **internal Database-as-a-Service (DBaaS) platform** at **Worldline**, addressing the need for **automation, scalability, and reliability** in managing **PostgreSQL instances** across sandbox, preproduction, and production environments.  

---

## 🚀 Project Overview

Traditional database provisioning often suffers from:
- Human intervention and delays  
- Inconsistent configurations  
- Lack of traceability  

In modern **DevOps** environments, these challenges limit agility.  
This project introduces a **DBaaS self-service layer** that enables:
- Automated provisioning  
- Consistent configuration  
- On-demand database lifecycle management  

---

## 🏗️ Architecture

The technical solution is based on an integration of modern DevOps tools and practices:

- **Custom Terraform Provider (Go)** → Interface between Terraform and API  
- **Reusable Terraform Modules** → Encapsulation of database configurations  
- **API Layer (Gin + Echo)** → Orchestration and workflow management  
- **Ansible Playbooks** → Automated provisioning & configuration of PostgreSQL  
- **CI/CD Integration (GitLab)** → Automated pipeline for delivery  
- **Monitoring & Maintenance**:  
  - `pg_repack` for automated cleanup  
  - Prometheus & Grafana hooks for observability  

📌 **Project Tree (main components):**
POSTGRESQL-DBAAS-SELF-SERVICE-LAYER/
│── module_pg_db-master/ # Core DB Terraform module
│── module_pg_exploit-master/ # Exploit module for PostgreSQL
│── postgres-exploit-sandbox/ # Sandbox exploitation
│── postgres-sandbox/ # PostgreSQL sandbox environment
│── terraform-provider-dbaas-postgres-exploit-master/ # Custom TF provider for exploit
│── terraform-provider-dbaas-postgres-master/ # Core TF provider in Go
│── README.md # Project documentation

---

## ⚙️ Features

✅ **Self-Service PostgreSQL DBaaS** (create, update, delete instances)  
✅ **Terraform-based automation** with custom provider  
✅ **Ansible-driven provisioning** with consistent configurations  
✅ **CI/CD with GitLab pipelines**  
✅ **Automated cleanup with `pg_repack`**  
✅ **Monitoring with Prometheus & Grafana**  
✅ **Multi-environment support**: Sandbox, Preproduction, Production  

---

## 🛠️ Technologies Used

- **PostgreSQL** – Target database engine  
- **Terraform** – Infrastructure as Code  
- **Go** – Custom Terraform provider  
- **Gin / Echo (Go frameworks)** – API orchestration layer  
- **Ansible** – Provisioning & configuration management  
- **GitLab CI/CD** – Automation pipeline  
- **Prometheus & Grafana** – Monitoring & observability  

---

## 🎯 Methodology

The project followed an **Agile approach**:
- Weekly sprints with clear objectives  
- Iterative design, implementation, and testing  
- Continuous feedback from SRE team  
- Documentation and alignment with industry standards  

---

## 📊 Impact & Benefits

- 🚀 **Reduced manual effort** → streamlined workflows  
- 🔄 **Consistent deployments** → no configuration drift  
- ⚡ **Faster delivery** → on-demand provisioning  
- 💰 **Lower operational costs**  
- 🔧 **Foundation for reuse** across other business units  

---

## 📸 Images & Diagrams

### Example: High-Level Architecture
![Architecture Example](https://upload.wikimedia.org/wikipedia/commons/4/48/Markdown-mark.svg)

---

## 🎥 Demo / Video

👉 [Project Demo Video](https://www.youtube.com/watch?v=dQw4w9WgXcQ)  
*(replace with actual demo link)*  

---

## 📌 Future Improvements

- Multi-database support (MySQL, Oracle, etc.)  
- Advanced monitoring & alerting  
- Extended self-service portal with UI  

---

## 🧑‍💻 Authors

- **ELMARCHOUM Ayoub**  

---

## 📄 License

This project is part of a **PFE (Projet de Fin d’Études)** at **Worldline**.  
Usage restricted to academic and internal purposes.
