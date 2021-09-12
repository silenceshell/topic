# topic

top in container.

Running the original `top` command in a container will not get information of the container, many metrics like uptime, users, tasks, cpu, memory, is about the host in fact. `topic`(top in container) will retrieve those metrics from container instead.

Attention: `topic` have not supported load average metrics.
