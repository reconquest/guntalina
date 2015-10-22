Guntalina
=========

Guntalina is the utility for creating and executing command list basing on list
of modified/overwrited/touched/chmoded/chowned/created files, action
definitions and list of rules.


It's convinient for using in integration with **Guntalina**'s brother
**Gunter**.

For example, your configuration manager changed or created following list of
files:

```
/etc/nginx/conf.d/sites-available.conf
/etc/nginx/nginx.conf
/var/whatever
```

So, you should reload nginx, what you should do? Write configuration rule!

## Configuration

Guntalina configuration file should be written in YAML language, and consists
of two basic directives:

- `actions` - array of action definitions, for example, `nginx-reload`, or
    `nginx-restart`, if you want.
- `rules` - array of rule definitions, which should give an answer on the
    question like *When some action should be invoked?*

#### Actions

Action, in **guntalina** meaning, it's array of commands which should be
executed when action triggered, which should be declared in subdirective
`commands:`.

Let's write some typical actions for `nginx reload` and `nginx restart`:

```yaml
actions:
    - nginx-reload:
        commands:
            - systemctl reload nginx

    - nginx-restart:
        commands:
            - nginx -t # let's force check nginx config before real restart
            - systemctl restart nginx
```

#### Rules

Basically rules has two parts, list of actions, which named as `workflow` and
list of file path (which modified/created/etc.) glob patterns which named as
`masks`.

With mask you can handle files by patterns like `/etc/nginx/conf.d/*` or
`/etc/nginx/*.conf` and do specified workflow.


Let's proceed to write rules for our arbitary example with `nginx`:

```yaml
rules:
    - masks:
        - /etc/nginx/conf.d/*
        - /etc/nginx/nginx.conf
      workflow:
        - nginx-reload
```

So, action `nginx-reload` will be triggered if `/etc/nginx/nginx.conf` or some
file in `/etc/nginx/conf.d/` directive has been changed.

When **guntalina** passes
through the list of modified files and list of rules, she disables rules after
their using, so if you have two or more modified files which fit a rule,
guntalina will execute rule workflow once.

Actually, if rule contains the `group` directive, then guntalina will disable all
rules with same group.

Also, if you have some rules, which
has repeatable workflow actions, guntalina will execute actions only once.

Let's write more complex config, when **guntalina** should make a decision for
restart, but do not typical reload of software which will be named, for
example, **exampled**.

**exampled** should be restarted when files like as
`/etc/exampled/conf.d/*_cache_zone` changed, and should be reloaded when any
another file in the directory `/etc/exampled/conf.d/` has been changed.

For this case we should use `group` directive.

```yaml
actions:
    - exampled-restart:
        commands:
            - exampled --check-config
            - examplectl stop
            - examplectl start
    - exampled-reload:
        commands:
            - examplectl reload

rules:
    - group: exampled
      masks:
        - /etc/nginx/conf.d/*_cache_zone
      workflow:
        - exampled-restart
    - group: exampled
      masks:
        - /etc/nginx/conf.d/*
      workflow:
        -  exampled-reload
```

How you can see, both rules uses optional *group* directive with `exampled`
value, it's means that guntalina will not runs `exampled-reload`, if
`exampled-restart` already triggered, because after triggering
`exampled-restart` group `exampled` will be disabled.

#### Includes

Guntalina has own YAML with Blackjack and includes. Include can be invoked
using `!include path/to/another/config`.

I'm recommend file structure which looks like this:

```
.
├── conf.d
│   ├── haproxy
│   │   ├── actions
│   │   └── rules
│   └── nginx
│       ├── actions
│       └── rules
└── guntalina.conf
```

For these reasons default configuration file looks like this:

```yaml
actions:
    !include conf.d/*/actions

rules:
    !include conf.d/*/rules
```

## Usage

#### Synopsis

- `-s <source>` - Specify source file, which should consist of list of
     modified/overwrited/created files.
- `-c <config` - Specify configuration file
- `-r --dry-run` - Dry-run mode, in this mode commands will be not really
     executed.
- `-f --force` - Do not stop if any command has been failed.
- `-v --version` - Show **guntalina** version.
- `-h --help` - Show help message.

So, **guntalina** have one required argument it's `-s <source>`, configuration
data by default will be read from `/etc/guntalina/guntalina.conf` file.
