# msm
Multi-layer Synthetic Monitor

The goal is to write a synthetic monitor that extends synthetic monitoring from websites to the application and DB layer. A
single monitor instance acts as a website, application and DB client to ensure the target (website, application or DB) is 
actually servicing requests. The results are sent to some other system (a monitoring system) so that the output data can be
used to so something sensible (raise alerts, take action to restart unrepsonsive targets, etc.). The goal is not to write 
another monitoring an alerting system.

From a user perspective MSM is intended to be very easy to install and run, so plain text configuration files with a simple
syntax are used to configure the software and list the targets to be checked. The software should be highly tolerant of 
both configuration file content (spaces, commas, etc. in the wrong place), and the way all the varieties of websites, 
applications and DBs work.

Written in Go for portability and ease of deployment. Also as a way to find out more about Go. Tested on a Raspberry PI 3 
(Raspian) and x86 (Windows and Linux)
