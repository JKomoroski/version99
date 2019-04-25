# Version 99 Does Not Exist

Version99 service in Go.

Please see [my blog](http://day-to-day-stuff.blogspot.com/2007/10/announcement-version-99-does-not-exist.html) to read why I created Version 99 Does Not Exist and its predecessor no-commons-logging.

Version 99 Does Not Exist emulates a Maven 2 repository and serves empty jars for any valid package that has version number *99.0-does-not-exist*. It also generates poms and of course the appropriate hashes.

For example the following links will give an [empty jar](http://version99.grons.nl/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.jar), its [pom](http://version99.grons.nl/mvn2/commons-logging/commons-logging/99.0-does-not-exist/commons-logging-99.0-does-not-exist.pom) for commons-logging.</p>

The original Ruby/Camping application was ported to Go by my colleague Frank Schroeders.

Build it as follows:

    go build

This service is running on http://version99.grons.nl. However, I encourage you to run this service yourself as bandwidth to my machine is limited.

Often you only need this sevice for just 1 package. In that case you can also download the jar and pom for that package and upload it to your local repository (that is, if you have one).

## Update April 2019

It has been 6 years since the Go version of this service was developed and deployed. Luckily, commons-logging is rarely encountered anymore and this service is less and less necessary. Nonetheless, the service is still running and I have found no reason yet to change that :)
