Bitbucket Web Service hook server
=================================

A web server handling POST hooks from bitbucket.org.

The idea and code is heavily inspired by
http://github.com/derekpitt/githubservicehook

Installation
-------------

    go install github.com/doist/bitbuckethook


Usage
-----

As you start a hook server, you could define following options

- `-c`: configuration file for hook listener commands (see below for format,
  default value is "bitbuckethook.json" in current directory)
- `-p`: port to listen (default is 4007)
- `-t`: optional secret token

If you launch the command as

    bitbuckethook -p 4007 -c bitbuckethook.json -t foobar

Then your POST hook URL should look like this:

    http://yourdomain.tld/?token=foobar

See [POST hook management](https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management)
for more details.

Configuration file
------------------

The file `bitbuckethook.json` is used to map repository names to commands
to execute. Repository names don't contain user part, commands are lists
of strings.

    {
        "repo1": ["make", "-C", "/path/to/directory"],
        "repo2": ["bash", "-c", "cd /foo && git pull && sudo service nginx reload"],
    }
