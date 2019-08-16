* Simplify deployment of complex workloads
* Automate updates, upgrades and other configuration changes
* Enable fine-grained access control to infrastructure and deployment changes 
* Ease the on-boarding and induction process for new staff
* Use identical hosting infrastructure within local, dev, staging, testing and production environments

Here's how to deploy a high-availability PostgreSQL cluster on AWS from scratch on any Linux computer:

```bash
sudo snap install --classic juju
juju autoload-credentials
juju bootstrap aws/us-west-1
juju deploy -n3 postgresql
```


## What is Juju?

Juju is a leading devops software project for people that don’t want to care about devops.
It reduces operational complexity by presenting a holistic view first.
We call this application modelling.
Once a model is described, Juju identifies the necessary steps to put that plan into action. 

Juju uses an active agent deployed alongside your applications.
That agent orchestrates infrastructure, manages applications through the product life cycle.
Juju can manage bare metal, virtual hardware, containers and applications hosted by any substrate.
The agent’s capabilities are confined by user permissions and all communication is fully encrypted.

## Why Juju?

* Simplicity
* Security
* Reliability
* Responsiveness

Most devops tools work best when operators have fine-grained knowledge about every detail.
Skilled operators know every hostname, every machine, every subnet and every storage component.
That makes changes complicated, excludes less technical people, makes on-boarding difficult and tends to create knowledge silos.

Juju flips the default.
If your infrastructure can’t be understood by everyone in your organisation, that’s a bug.

## Next steps

Install Juju https://jaas.ai/docs/installing 

Walk through a tutorial https://jaas.ai/docs/client-usage-tutorial 

Introduce yourself or ask a question https://discourse.jujucharms.com/
