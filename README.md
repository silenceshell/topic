# topic

top in container.

Running the original `top` command in a container will not get information of the container, many metrics like uptime, users, load average, tasks, cpu, memory, are about the host in fact. 
`topic`(top in container) will retrieve those metrics from container instead, so it shows the status of the container, not the host.
