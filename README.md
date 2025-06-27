# Ferroxide

A community fork of [emersion/ferroxide](https://github.com/emersion/ferroxide)

Primary changes:
- CalDAV
- Tor and proxies
- Custom config directory
- Systemd socket support

# Systemd socket usage

The systemd socket and service files can be found under the directory `dist`.
Using these you can basically run the services on demand.
If no one is connected to them, no resources are used but when someone connects to the service port it's automatically started.

You also gain a lot of security measure related to process isolation.
I still would not expose these services to the world.
Bind to localhost or an interface behind a firewall/VPN.

If you chose to use the system socket services (best option) you'll need to log into your protonmail account but under the `ferroxide` Linux user.

simplest way is:

    sudo -u ferroxide ferroxide auth email@protonmail.com

Or copy your auth data at `$HOME/.config/ferroxide/auth.json` to `/var/lib/ferroxide/auth.json`

The services can be started as usual with

    sudo systemctl start ferroxide-{imap,smtp,etc}.socket

Or you can just enable them by default.

    sudo systemctl enable ferroxide-{imap,smtp,etc}.socket

If you want to bind to different ports or other customizations use

    sudo systemctl edit ferroxide-{imap,smtp,etc}.socket

And change what you want.

Environment variables can be set in `/etc/ferroxide.conf`.

if you use the user systemd services don't forget to use the `--user` flag in the previous comments.

The service will then run with your own user, with all the potential problems from the lack of isolation.

# Original Hydroxide ReadMe

A third-party, open-source ProtonMail bridge. For power users only, designed to
run on a server.

ferroxide supports CardDAV, CalDAV, IMAP and SMTP.

Rationale:

* No GUI, only a CLI (so it runs in headless environments)
* Standard-compliant (we don't care about Microsoft Outlook)
* Fully open-source

Feel free to join the IRC channel: #emersion on Libera Chat.

## How does it work?

ferroxide is a server that translates standard protocols (SMTP, IMAP, CardDAV, CalDAV)
into ProtonMail API requests. It allows you to use your preferred e-mail clients
and `git-send-email` with ProtonMail.

    +-----------------+             +-------------+  ProtonMail  +--------------+
    |                 | IMAP, SMTP  |             |     API      |              |
    |  E-mail client  <------------->  ferroxide  <-------------->  ProtonMail  |
    |                 |             |             |              |              |
    +-----------------+             +-------------+              +--------------+

## Setup

### Go

ferroxide is implemented in Go. Head to [Go website](https://golang.org) for
setup information.

### Installing

Start by installing ferroxide:

```shell
go install github.com/vcalv/ferroxide-systemd/cmd/ferroxide@latest
```

Then you'll need to login to ProtonMail via ferroxide, so that ferroxide can
retrieve e-mails from ProtonMail. You can do so with this command:

```shell
ferroxide auth <username>
```

Once you're logged in, a "bridge password" will be printed. Don't close your
terminal yet, as this password is not stored anywhere by ferroxide and will be
needed when configuring your e-mail client.

Your ProtonMail credentials are stored on disk encrypted with this bridge
password (a 32-byte random password generated when logging in).

## Usage

ferroxide can be used in multiple modes.

> Don't start ferroxide multiple times, instead you can use `ferroxide serve`.
> This requires ports 1025 (smtp), 1143 (imap), 8080 (carddav) and 8081 (caldav).

### SMTP

To run ferroxide as an SMTP server:

```shell
ferroxide smtp
```

Once the bridge is started, you can configure your e-mail client with the
following settings:

* Hostname: `localhost`
* Port: 1025
* Security: none
* Username: your ProtonMail username
* Password: the bridge password (not your ProtonMail password)

### CardDAV

You must setup an HTTPS reverse proxy to forward requests to `ferroxide`.

```shell
ferroxide carddav
```

Tested on GNOME (Evolution) and Android (DAVDroid).

### CalDAV

```shell
ferroxide caldav
```

Tested on GNOME (Evolution), Thunderbird, KOrganizer.

### IMAP

⚠️  **Warning**: IMAP support is work-in-progress. Here be dragons.

For now, it only supports unencrypted local connections.

```shell
ferroxide imap
```

## License

MIT
