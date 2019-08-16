![](https://raw.githubusercontent.com/timClicks/juju/develop-docs--readme-upgrade/doc/juju-logo.png)


## Why Juju?

* Reduces complexity involved with operations by modelling your complex software deployments
* Encourages sharing of operations knowledge by encapsulating it reusable components, called charms
* Enables repeatability of complex deployments
* Excels at day two operations, such as upgrades, updates and configuration changes
* Provides portability across clouds, both public and private

Most devops tools work best when operators have fine-grained knowledge about every detail.
Skilled operators know every hostname, every machine, every subnet and every storage component.
That makes changes complicated, makes on-boarding difficult and tends to create knowledge silos.

Juju flips the default.
If your infrastructure can’t be understood by everyone in your organisation, that’s a bug.
Instead of reasoning about discrete machines and their IP addresses, Juju focuses on the applications that are provided by your software model.

Keeping a complex deployment up-to-date and functioning perfectly in spite of these changes represents an immense challenge.
And for operators maintaining these deployments, the cost in both time and cognitive load can be tremendous.
Juju is built for operators who want simplicity, security and stability.


## What is Juju?

Juju is a devops tool that reduces operational complexity through application modelling.
Once a model is described, Juju identifies the necessary steps to put that plan into action.
Juju has three core concepts: models, applications and relations.

Consider a whiteboard drawing of your service.
The whiteboard's border is the model, its boxes are applications and the lines between the boxes are relations. 

Juju uses an active agent deployed alongside your applications.
That agent orchestrates infrastructure, manages applications through the product life cycle.
Juju can manage bare metal, virtual hardware, containers and applications hosted by any substrate.
The agent’s capabilities are confined by user permissions and all communication is fully encrypted.

Juju allows you to deploy, configure, manage, maintain, and scale cloud applications quickly and efficiently.
on public clouds, as well as on prem , OpenStack, Kubernetes and containers.
You can use Juju from the command line or through its beautiful GUI.


## Next steps

Install Juju https://jaas.ai/docs/installing 

Walk through a tutorial https://jaas.ai/docs/client-usage-tutorial 

Introduce yourself or ask a question https://discourse.jujucharms.com/
