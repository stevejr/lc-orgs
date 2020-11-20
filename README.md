# lc-orgs


## TODOs

### Automate generation of faas secrets

faas-cli secret create mysql-host \
  --from-literal "lc-db-1.cluster-cboq7plsgav7.eu-west-2.rds.amazonaws.com"  --gateway https://openfaas.lc.lcmvp.dockerps.io --tls-no-verify

faas-cli secret create mysql-username \
  --from-literal "lc-orgs" --gateway https://openfaas.lc.lcmvp.dockerps.io --gateway https://openfaas.lc.lcmvp.dockerps.io --tls-no-verify

faas-cli secret create mysql-password \
  --from-literal "Ue8ieMaithevaix" --gateway https://openfaas.lc.lcmvp.dockerps.io --gateway https://openfaas.lc.lcmvp.dockerps.io --tls-no-verify

### Automate creation of mysql 'database' and tables

CREATE TABLE orgs
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    unique_id varchar(36) NOT NULL,
    short_name text NOT NULL,
    long_name text NOT NULL,
    owner_email text NOT NULL
);


INSERT INTO orgs(unique_id,short_name,long_name,owner_email)
  VALUES (uuid(), 'org1', 'I am org1', 'admin@org1.com');

INSERT INTO orgs(unique_id,short_name,long_name,owner_email)
  VALUES (uuid(), 'org2', 'I am org2', 'admin@org2.com');