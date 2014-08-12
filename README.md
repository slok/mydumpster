mydumpster
==========

Mysql dumps based on a config file

NOT READY FOR USE
-----------------
Some features work, for now don't use it, we are working hard to get an
stable version, be patient :)

|Workflow | CI | 
|---------|----|
| [![Stories in Ready](https://badge.waffle.io/slok/mydumpster.png?label=ready)](http://waffle.io/sharestack/sharestack-api) | |

Description
-----------

We always want dumps from our databases, but sometimes we don't want all the database,
imagine a dev team that has a 2GB database and downloads these dumps every few days
to develop, but at the end of the day this team only used the latest registries 
of teh database (imagine a 3% of the database), that sucks

This project takes this problem and tries to solve it. The dumps are configured
in simple json files so you can do easily different custom dumps.

Features (planned)
------------------

* Censore database
* Dump N tables
* Filter tables (a.k.a get me only the rows that I need)
* Referential integrity
* dump in parallel (warning!)
* Foreign key automatic discovery when dumping
* JSON config files

License
-------
MIT

Author
------
Xabier Larrakoetxea

